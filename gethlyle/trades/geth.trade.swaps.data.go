package gethlyletrades

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	gethlyleswaps "github.com/kfukue/lyle-labs-libraries/gethlyle/swaps"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetAllGethTradeSwapsByTradeID(dbConnPgx utils.PgxIface, gethTradeID *int) ([]GethTradeSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT 
		geth_trade_swaps.geth_trade_id,
		geth_trade_swaps.geth_swap_id,
		geth_trade_swaps.uuid, 
		geth_trade_swaps.name, 
		geth_trade_swaps.alternate_name, 
		geth_trade_swaps.description,
		geth_trade_swaps.created_by, 
		geth_trade_swaps.created_at, 
		geth_trade_swaps.updated_by, 
		geth_trade_swaps.updated_at 
	FROM geth_trades 
	LEFT JOIN geth_trade_swaps
	ON geth_trade_swaps.geth_trade_id = geth_trades.id 
	WHERE 
	geth_trades.id = $1
	`, *gethTradeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTradeSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTradeSwap])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTradeSwaps, nil
}
func GetGethTradeSwap(dbConnPgx utils.PgxIface, gethGethSwapID, gethTradeID *int) (*GethTradeSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `
	SELECT 
		geth_trade_id,
		geth_swap_id,
		uuid, 
		name, 
		alternate_name, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM geth_trade_swaps 
	WHERE 
	geth_trade_id = $1
	AND geth_swap_id = $2
	`, *gethTradeID, *gethGethSwapID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethTradeSwap, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethTradeSwap])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethTradeSwap, nil
}

func RemoveGethTradeSwap(dbConnPgx utils.PgxIface, gethTradeID, gethGethSwapID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethTradeSwap DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_trade_swaps WHERE geth_trade_id =$1 AND geth_swap_id = $2`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *gethTradeID, *gethGethSwapID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetGethTradeSwapList(dbConnPgx utils.PgxIface, gethTradeIds []int, swapIds []int) ([]GethTradeSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT 
		geth_trade_id,
		geth_swap_id,
		uuid, 
		name, 
		alternate_name, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM geth_trade_swaps`
	if len(gethTradeIds) > 0 || len(swapIds) > 0 {
		additionalQuery := ` WHERE`
		if len(gethTradeIds) > 0 {
			gethTradeStrIds := utils.SplitToString(gethTradeIds, ",")
			additionalQuery += fmt.Sprintf(`geth_trade_id IN (%s)`, gethTradeStrIds)
		}
		if len(swapIds) > 0 {
			if len(gethTradeIds) > 0 {
				additionalQuery += `AND `
			}
			swapStrIds := utils.SplitToString(swapIds, ",")
			additionalQuery += fmt.Sprintf(`geth_swap_id IN (%s)`, swapStrIds)
		}
		sql += additionalQuery
	}
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTradeSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTradeSwap])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTradeSwaps, nil
}

func UpdateGethTradeSwap(dbConnPgx utils.PgxIface, gethTradeSwap *GethTradeSwap) error {
	// if the gethTradeSwap id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (gethTradeSwap.GethSwapID == nil || *gethTradeSwap.GethSwapID == 0) || (gethTradeSwap.GethTradeID == nil || *gethTradeSwap.GethTradeID == 0) {
		return errors.New("gethTradeSwap has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethTradeSwap DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_trade_swaps SET 
		name=$1,  
		alternate_name=$2, 
		description=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE 
			geth_trade_id=$5 AND
			geth_swap_id=$6 
		`
	if _, err := dbConnPgx.Exec(ctx, sql,
		gethTradeSwap.Name,          //1
		gethTradeSwap.AlternateName, //2
		gethTradeSwap.Description,   //3
		gethTradeSwap.UpdatedBy,     //4
		gethTradeSwap.GethTradeID,   //5
		gethTradeSwap.GethSwapID,    //6
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethTradeSwap(dbConnPgx utils.PgxIface, gethTradeSwap *GethTradeSwap) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethTradeSwap DbConn.Begin   %s", err.Error())
		return -1, -1, err
	}
	var GethSwapID int
	var GethTradeID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_trade_swaps  
	(
		geth_trade_id,
		geth_swap_id,
		uuid,	 
		name, 
		alternate_name,  
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES ($1,
			$2,
			uuid_generate_v4(),
			$3,
			$4,
			$5,
			$6,
			current_timestamp at time zone 'UTC',
			$6,
			current_timestamp at time zone 'UTC'
		)
		RETURNING geth_trade_id, geth_swap_id`,
		gethTradeSwap.GethTradeID,   //1
		gethTradeSwap.GethSwapID,    //2
		gethTradeSwap.Name,          //3
		gethTradeSwap.AlternateName, //4
		gethTradeSwap.Description,   //5
		gethTradeSwap.CreatedBy,     //6
	).Scan(&GethTradeID, &GethSwapID)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, -1, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, -1, err
	}
	return int(GethTradeID), int(GethSwapID), nil
}

func InsertGethTradeSwaps(dbConnPgx utils.PgxIface, gethTradeSwaps []GethTradeSwap) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethTradeSwaps {
		gethTradeSwap := gethTradeSwaps[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(gethTradeSwap.UUID)
		row := []interface{}{
			gethTradeSwap.GethTradeID,   //1
			gethTradeSwap.GethSwapID,    //2
			uuidString,                  //3
			gethTradeSwap.Name,          //4
			gethTradeSwap.AlternateName, //5
			gethTradeSwap.Description,   //6
			gethTradeSwap.CreatedBy,     //7
			&now,                        //8
			gethTradeSwap.CreatedBy,     //9
			&now,                        //10
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_trade_swaps"},
		[]string{
			"geth_trade_id",  //1
			"geth_swap_id",   //2
			"uuid",           //3
			"name",           //4
			"alternate_name", //5
			"description",    //6
			"created_by",     //7
			"created_at",     //8
			"updated_by",     //9
			"updated_at",     //10
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethTradeSwaps: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetMissingTradesFromSwapsByBaseAssetID(dbConnPgx utils.PgxIface, baseAssetID *int) ([]gethlyleswaps.GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
		SELECT
			gs.id,
			gs.uuid,
			gs.chain_id,
			gs.exchange_id,
			gs.block_number,
			gs.index_number,
			gs.swap_date,
			gs.trade_type_id,
			gs.txn_hash,
			gs.maker_address,
			gs.maker_address_id,
			gs.is_buy,
			gs.price,
			gs.price_usd,
			gs.token1_price_usd,
			gs.total_amount_usd,
			gs.pair_address,
			gs.liquidity_pool_id,
			gs.token0_asset_id,
			gs.token1_asset_id,
			gs.token0_amount,
			gs.token1_Amount,
			gs.description,
			gs.created_by,
			gs.created_at,
			gs.updated_by,
			gs.updated_at,
			gs.geth_process_job_id,
			gs.topics_str,
			gs.status_id,
			gs.base_asset_id,
			gs.oracle_price_usd,
			gs.oracle_price_asset_id
		FROM geth_swaps gs
		LEFT JOIN geth_trade_swaps gts
			ON gs.id = gts.geth_swap_id
		WHERE
		gts.geth_swap_id IS NULL
		AND 
		gs.base_asset_id = $1
		ORDER BY gs.block_number asc
		`,
		*baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[gethlyleswaps.GethSwap])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethSwaps, nil
}

func GetMissingTxnHashesFromSwapsByBaseAssetID(dbConnPgx utils.PgxIface, baseAssetID *int, maxBlockNumber *uint64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
		SELECT
			DISTINCT gs.txn_hash
		FROM geth_swaps gs
		LEFT JOIN geth_trade_swaps gts
			ON gs.id = gts.geth_swap_id
		WHERE
			gts.geth_swap_id IS NULL
				AND 
			gs.base_asset_id = $1
				AND
			gs.status_id = $2
			AND
			gs.block_number > $3
		`,
		*baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID, *maxBlockNumber,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	txnHashes := make([]string, 0)
	for results.Next() {
		var txnHash string
		results.Scan(
			&txnHash,
		)
		txnHashes = append(txnHashes, txnHash)
	}
	return txnHashes, nil
}

func GetMinMaxBlocksOfMissingSwapByBaseAssetID(dbConnPgx utils.PgxIface, baseAssetID *int) (*uint64, *uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var minBlock, maxBlock uint64
	err := dbConnPgx.QueryRow(ctx, `SELECT
			MIN(gs.block_number),
			MAX(gs.block_number)
		FROM geth_swaps gs
		LEFT JOIN geth_trade_swaps gts
			ON gs.id = gts.geth_swap_id
		WHERE
			gts.geth_swap_id IS NULL
				AND 
			gs.base_asset_id = $1
				AND
			gs.status_id = $2
		`,
		*baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID,
	).Scan(&minBlock, &maxBlock)
	if err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}
	return &minBlock, &maxBlock, nil
}

func GetFirstNonProcessedSwapBlockNumberForTrades(dbConnPgx utils.PgxIface, baseAssetID *int) (*uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var startingBlock uint64
	err := dbConnPgx.QueryRow(ctx, `
	WITH max_existing_block_swaps as (
		SELECT COALESCE(MAX(block_number),0) as min_block_number from  geth_swaps gs
		LEFT JOIN geth_trade_swaps gts
			ON gs.id = gts.geth_swap_id
		WHERE
			gts.geth_swap_id IS NOT NULL
				AND 
			gs.base_asset_id = $1
				AND
			gs.status_id = $2
		),
		min_missing_block_swaps as (
		SELECT COALESCE(MIN(block_number),0) as min_block_number from  geth_swaps gs
		LEFT JOIN geth_trade_swaps gts
			ON gs.id = gts.geth_swap_id
		WHERE
			gts.geth_swap_id IS NULL
				AND 
			gs.base_asset_id = $1
				AND
			gs.status_id = $2
			)
		SELECT MAX(min_block_number) 
		FROM (
			SELECT min_block_number FROM max_existing_block_swaps
			UNION 
			SELECT min_block_number FROM min_missing_block_swaps
				)   missing_and_existing_swap
		`, *baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID,
	).Scan(&startingBlock)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &startingBlock, nil
}

// for refinedev
func GetGethTradeSwapListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethTradeSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
	geth_trade_id,
		geth_swap_id,
		uuid, 
		name, 
		alternate_name, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM geth_trade_swaps 
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
	minerTransactionInputs, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTradeSwap])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return minerTransactionInputs, nil
}

func GetTotalGethTradeSwapCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_trade_swaps
	`)
	totalCount := 0
	err := row.Scan(
		&totalCount,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &totalCount, nil
}
