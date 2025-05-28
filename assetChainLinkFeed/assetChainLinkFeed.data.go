package assetchainlinkfeed

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
)

func GetAssetChainLinkFeed(dbConnPgx utils.PgxIface, assetID, chainID *int) (*AssetChainLinkFeed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
		asset_id,
		chain_id,
		chainlink_data_feed_contract_address,
		created_by,
		created_at,
		updated_by,
		updated_at
		FROM asset_chain_link_data_feed 
		WHERE asset_id = $1 AND chain_id = $2`, *assetID, *chainID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	feed, err := pgx.CollectOneRow(row, pgx.RowToStructByName[AssetChainLinkFeed])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &feed, nil
}

func GetAssetChainLinkFeedList(dbConnPgx utils.PgxIface, assetIDs, chainIDs []int) ([]AssetChainLinkFeed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		asset_id,
		chain_id,
		chainlink_data_feed_contract_address,
		created_by,
		created_at,
		updated_by,
		updated_at
		FROM asset_chain_link_data_feed`
	if len(assetIDs) > 0 || len(chainIDs) > 0 {
		additionalQuery := ` WHERE`
		if len(assetIDs) > 0 {
			assetStrIds := utils.SplitToString(assetIDs, ",")
			additionalQuery += fmt.Sprintf(` asset_id IN (%s)`, assetStrIds)
		}
		if len(chainIDs) > 0 {
			if len(assetIDs) > 0 {
				additionalQuery += ` AND `
			}
			chainStrIds := utils.SplitToString(chainIDs, ",")
			additionalQuery += fmt.Sprintf(` chain_id IN (%s)`, chainStrIds)
		}
		sql += additionalQuery
	}
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	feeds, err := pgx.CollectRows(results, pgx.RowToStructByName[AssetChainLinkFeed])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return feeds, nil
}

func InsertAssetChainLinkFeed(dbConnPgx utils.PgxIface, feed *AssetChainLinkFeed) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertAssetChainLinkFeed DbConn.Begin   %s", err.Error())
		return err
	}
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO asset_chain_link_data_feed  
		(asset_id, chain_id, chainlink_data_feed_contract_address, created_by, created_at, updated_by, updated_at)
		VALUES ($1, $2, $3, $4, current_timestamp at time zone 'UTC', $5, current_timestamp at time zone 'UTC')`,
		feed.AssetID,
		feed.ChainID,
		feed.ChainlinkDataFeedContractAddress,
		feed.CreatedBy,
		feed.UpdatedBy,
	).Scan()
	if err != nil && err != pgx.ErrNoRows {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateAssetChainLinkFeed(dbConnPgx utils.PgxIface, feed *AssetChainLinkFeed) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if feed.AssetID == nil || feed.ChainID == nil {
		return errors.New("AssetChainLinkFeed has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateAssetChainLinkFeed DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE asset_chain_link_data_feed SET 
		chainlink_data_feed_contract_address=$1,
		updated_by=$2,
		updated_at=current_timestamp at time zone 'UTC'
		WHERE asset_id=$3 AND chain_id=$4`
	if _, err := dbConnPgx.Exec(ctx, sql,
		feed.ChainlinkDataFeedContractAddress, //1
		feed.UpdatedBy,                        //2
		feed.AssetID,                          //3
		feed.ChainID,                          //4
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemoveAssetChainLinkFeed(dbConnPgx utils.PgxIface, assetID, chainID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveAssetChainLinkFeed DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM asset_chain_link_data_feed WHERE asset_id = $1 AND chain_id = $2`
	if _, err := dbConnPgx.Exec(ctx, sql, *assetID, *chainID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertAssetChainLinkFeeds(dbConnPgx utils.PgxIface, feeds []AssetChainLinkFeed) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range feeds {
		feed := feeds[i]
		row := []interface{}{
			feed.AssetID,                          //1
			feed.ChainID,                          //2
			feed.ChainlinkDataFeedContractAddress, //3
			feed.CreatedBy,                        //4
			&now,                                  //5
			feed.CreatedBy,                        //6
			&now,                                  //7
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"asset_chain_link_data_feed"},
		[]string{
			"asset_id",                             //1
			"chain_id",                             //2
			"chainlink_data_feed_contract_address", //3
			"created_by",                           //4
			"created_at",                           //5
			"updated_by",                           //6
			"updated_at",                           //7
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertAssetChainLinkFeeds: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetAssetChainLinkFeedListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]AssetChainLinkFeed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
		asset_id,
		chain_id,
		chainlink_data_feed_contract_address,
		created_by,
		created_at,
		updated_by,
		updated_at
		FROM asset_chain_link_data_feed 
	`
	if len(_filters) > 0 {
		sql += "WHERE "
		for i, filter := range _filters {
			sql += filter
			if i < len(_filters)-1 {
				sql += " OR "
			}
		}
	}
	if _order != "" && _sort != "" {
		sql += fmt.Sprintf(" ORDER BY %s %s ", _sort, _order)
	}
	if (_start != nil && *_start > 0) && (_end != nil && *_end > 0) {
		pageSize := *_end - *_start
		sql += fmt.Sprintf(" OFFSET %d LIMIT %d ", *_start, pageSize)
	}

	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	feeds, err := pgx.CollectRows(results, pgx.RowToStructByName[AssetChainLinkFeed])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return feeds, nil
}
