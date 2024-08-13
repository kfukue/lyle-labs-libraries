package strategymarketdataasset

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
)

func GetStrategyMarketDataAsset(dbConnPgx utils.PgxIface, strategyMarketDataAssetID *int) (*StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM strategy_market_data_assets 
	WHERE id = $1`, *strategyMarketDataAssetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	strategyMarketDataAsset, err := pgx.CollectOneRow(row, pgx.RowToStructByName[StrategyMarketDataAsset])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &strategyMarketDataAsset, nil
}

func RemoveStrategyMarketDataAsset(dbConnPgx utils.PgxIface, strategyMarketDataAssetID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveStrategyMarketDataAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM strategy_market_data_assets WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *strategyMarketDataAssetID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetStrategyMarketDataAssets(dbConnPgx utils.PgxIface, ids []int) ([]StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM strategy_market_data_assets`
	if len(ids) > 0 {
		strIds := utils.SplitToString(ids, ",")
		additionalQuery := fmt.Sprintf(` WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategyMarketDataAssets, err := pgx.CollectRows(results, pgx.RowToStructByName[StrategyMarketDataAsset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyMarketDataAssets, nil
}

func GetStrategyMarketDataAssetsByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
	FROM strategy_market_data_assets
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategyMarketDataAssets, err := pgx.CollectRows(results, pgx.RowToStructByName[StrategyMarketDataAsset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyMarketDataAssets, nil
}

func GetStrategyMarketDataAssetsByStrategyID(dbConnPgx utils.PgxIface, strategyID *int) ([]StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM strategy_market_data_assets
	WHERE strategy_id = $1
	`, *strategyID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategyMarketDataAssets, err := pgx.CollectRows(results, pgx.RowToStructByName[StrategyMarketDataAsset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyMarketDataAssets, nil
}

func GetStartAndEndDateDiffStrategyMarketDataAssets(dbConnPgx utils.PgxIface, diffInDate *int) ([]StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
	FROM strategy_market_data_assets
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, *diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategyMarketDataAssets, err := pgx.CollectRows(results, pgx.RowToStructByName[StrategyMarketDataAsset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyMarketDataAssets, nil
}

func UpdateStrategyMarketDataAsset(dbConnPgx utils.PgxIface, strategyMarketDataAsset *StrategyMarketDataAsset) error {
	// if the strategyMarketDataAsset id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if strategyMarketDataAsset.ID == nil || *strategyMarketDataAsset.ID == 0 {
		return errors.New("strategyMarketDataAsset has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateStrategyMarketDataAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE strategy_market_data_assets SET 
		strategy_id=$1, 
		base_asset_id=$2, 
		quote_asset_id=$3, 
		uuid=$4, 
		name=$5, 
		alternate_name=$6, 
		start_date=$7, 
		end_date=$8, 
		ticker=$9, 
		description=$10, 
		source_id=$11, 
		frequency_id=$12, 
		updated_by=$13, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$14`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		strategyMarketDataAsset.StrategyID,    //1
		strategyMarketDataAsset.BaseAssetID,   //2
		strategyMarketDataAsset.QuoteAssetID,  //3
		strategyMarketDataAsset.UUID,          //4
		strategyMarketDataAsset.Name,          //5
		strategyMarketDataAsset.AlternateName, //6
		strategyMarketDataAsset.StartDate,     //7
		strategyMarketDataAsset.EndDate,       //8
		strategyMarketDataAsset.Ticker,        //9
		strategyMarketDataAsset.Description,   //10
		strategyMarketDataAsset.SourceID,      //11
		strategyMarketDataAsset.FrequencyID,   //12
		strategyMarketDataAsset.UpdatedBy,     //13
		strategyMarketDataAsset.ID,            //14
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertStrategyMarketDataAsset(dbConnPgx utils.PgxIface, strategyMarketDataAsset *StrategyMarketDataAsset) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertStrategyMarketDataAsset DbConn.Begin   %s", err.Error())
		return -1, err
	}
	defer cancel()
	var insertID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO strategy_market_data_assets 
	(
		strategy_id,
		base_asset_id,
		quote_asset_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
		) VALUES (
			$1,
			$2, 
			$3, 
			uuid_generate_v4(), 
			$4, 
			$5, 
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			current_timestamp at time zone 'UTC',
			$12,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		strategyMarketDataAsset.StrategyID,    //1
		strategyMarketDataAsset.BaseAssetID,   //2
		strategyMarketDataAsset.QuoteAssetID,  //3
		strategyMarketDataAsset.Name,          //4
		strategyMarketDataAsset.AlternateName, //5
		strategyMarketDataAsset.StartDate,     //6
		strategyMarketDataAsset.EndDate,       //7
		strategyMarketDataAsset.Ticker,        //8
		strategyMarketDataAsset.Description,   //9
		strategyMarketDataAsset.SourceID,      //10
		strategyMarketDataAsset.FrequencyID,   //11
		strategyMarketDataAsset.CreatedBy,     //12
	).Scan(&insertID)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, err
	}
	return int(insertID), nil
}
func InsertStrategyMarketDataAssets(dbConnPgx utils.PgxIface, strategyMarketDataAssets []StrategyMarketDataAsset) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range strategyMarketDataAssets {
		strategyMarketDataAsset := strategyMarketDataAssets[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(strategyMarketDataAsset.UUID)
		row := []interface{}{
			strategyMarketDataAsset.StrategyID,    //1
			strategyMarketDataAsset.BaseAssetID,   //2
			strategyMarketDataAsset.QuoteAssetID,  //3
			uuidString,                            //4
			strategyMarketDataAsset.Name,          //5
			strategyMarketDataAsset.AlternateName, //6
			&strategyMarketDataAsset.StartDate,    //7
			&strategyMarketDataAsset.EndDate,      //8
			strategyMarketDataAsset.Ticker,        //9
			strategyMarketDataAsset.Description,   //10
			strategyMarketDataAsset.SourceID,      //11
			strategyMarketDataAsset.FrequencyID,   //12
			strategyMarketDataAsset.CreatedBy,     //13
			&now,                                  //14
			strategyMarketDataAsset.CreatedBy,     //15
			&now,                                  //16
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"strategy_market_data_assets"},
		[]string{
			"strategy_id",    //1
			"base_asset_id",  //2
			"quote_asset_id", //3
			"uuid",           //4
			"name",           //5
			"alternate_name", //6
			"start_date",     //7
			"end_date",       //8
			"ticker",         //9
			"description",    //10
			"source_id",      //11
			"frequency_id",   //12
			"created_by",     //13
			"created_at",     //14
			"updated_by",     //15
			"updated_at",     //16
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertStrategyMarketDataAssets: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetStrategyMarketDataAssetListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM strategy_market_data_assets 
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
	defer results.Close()
	strategyList, err := pgx.CollectRows(results, pgx.RowToStructByName[StrategyMarketDataAsset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyList, nil
}

func GetTotalStrategyMarketDataAssetsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM strategy_market_data_assets
	`)
	totalCount := 0
	err := row.Scan(
		&totalCount,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &totalCount, nil
}
