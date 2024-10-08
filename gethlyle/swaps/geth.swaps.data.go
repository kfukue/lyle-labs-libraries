package gethlyleswaps

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

func GetGethSwapByBlockChain(dbConnPgx utils.PgxIface, txnHash string, blockNumber *uint64, indexNumber *uint, makerAddressID, liquidityPoolID *int) (*GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		swap_date,
		trade_type_id,
		txn_hash,
		maker_address,
		maker_address_id,
		is_buy,
		price,
		price_usd,
		token1_price_usd,
		total_amount_usd,
		pair_address,
		liquidity_pool_id,
		token0_asset_id,
		token1_asset_id,
		token0_amount,
		token1_Amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		oracle_price_usd,
		oracle_price_asset_id
	FROM geth_swaps
	WHERE txn_hash= $1
		AND block_number = $2
		AND index_number = $3
		AND maker_address_id = $4
		AND liquidity_pool_id = $5
	`, txnHash, *blockNumber, *indexNumber, *makerAddressID, *liquidityPoolID)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethSwap, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethSwap])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethSwap, nil
}

func GetGethSwap(dbConnPgx utils.PgxIface, gethSwapID *int) (*GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		swap_date,
		trade_type_id,
		txn_hash,
		maker_address,
		maker_address_id,
		is_buy,
		price,
		price_usd,
		token1_price_usd,
		total_amount_usd,
		pair_address,
		liquidity_pool_id,
		token0_asset_id,
		token1_asset_id,
		token0_amount,
		token1_Amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		oracle_price_usd,
		oracle_price_asset_id
	FROM geth_swaps
	WHERE id = $1
	`, *gethSwapID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethSwap, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethSwap])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethSwap, nil
}

func GetGethSwapByStartAndEndDates(dbConnPgx utils.PgxIface, startDate, endDate time.Time) ([]GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		swap_date,
		trade_type_id,
		txn_hash,
		maker_address,
		maker_address_id,
		is_buy,
		price,
		price_usd,
		token1_price_usd,
		total_amount_usd,
		pair_address,
		liquidity_pool_id,
		token0_asset_id,
		token1_asset_id,
		token0_amount,
		token1_Amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		oracle_price_usd,
  		oracle_price_asset_id
		FROM geth_swaps
		WHERE swap_date BETWEEN $1 AND $2
		`,
		startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethSwap])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethSwaps, nil
}

func GetGethSwapByFromMakerAddress(dbConnPgx utils.PgxIface, makerAddress string) ([]GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		swap_date,
		trade_type_id,
		txn_hash,
		maker_address,
		maker_address_id,
		is_buy,
		price,
		price_usd,
		token1_price_usd,
		total_amount_usd,
		pair_address,
		liquidity_pool_id,
		token0_asset_id,
		token1_asset_id,
		token0_amount,
		token1_Amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		oracle_price_usd,
  		oracle_price_asset_id
		FROM geth_swaps
		WHERE
		maker_address = $1
		ORDER BY gethSwap_date asc`,
		makerAddress,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethSwap])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethSwaps, nil
}

func GetGethSwapByFromMakerAddressId(dbConnPgx utils.PgxIface, makerAddressID *int) ([]GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		swap_date,
		trade_type_id,
		txn_hash,
		maker_address,
		maker_address_id,
		is_buy,
		price,
		price_usd,
		token1_price_usd,
		total_amount_usd,
		pair_address,
		liquidity_pool_id,
		token0_asset_id,
		token1_asset_id,
		token0_amount,
		token1_Amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		oracle_price_usd,
  		oracle_price_asset_id
		FROM geth_swaps
		WHERE
		maker_address_id = $1
		ORDER BY swap_date asc`,
		*makerAddressID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethSwap])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethSwaps, nil
}

func GetGethSwapByFromMakerAddressIdAndBeforeBlockNumber(dbConnPgx utils.PgxIface, baseAssetID, makerAddressID *int, blockNumber *uint64) ([]GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		swap_date,
		trade_type_id,
		txn_hash,
		maker_address,
		maker_address_id,
		is_buy,
		price,
		price_usd,
		token1_price_usd,
		total_amount_usd,
		pair_address,
		liquidity_pool_id,
		token0_asset_id,
		token1_asset_id,
		token0_amount,
		token1_Amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		oracle_price_usd,
		oracle_price_asset_id
		FROM geth_swaps
		WHERE
		base_asset_id = $1
		AND maker_address_id = $2
		AND block_number <= $3
		ORDER BY swap_date asc`,
		*baseAssetID, *makerAddressID, *blockNumber,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethSwap])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethSwaps, nil
}

func GetGethSwapByFromBaseAssetAndBeforeBlockNumber(dbConnPgx utils.PgxIface, baseAssetID *int, blockNumber *uint64) ([]GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		swap_date,
		trade_type_id,
		txn_hash,
		maker_address,
		maker_address_id,
		is_buy,
		price,
		price_usd,
		token1_price_usd,
		total_amount_usd,
		pair_address,
		liquidity_pool_id,
		token0_asset_id,
		token1_asset_id,
		token0_amount,
		token1_Amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		oracle_price_usd,
		oracle_price_asset_id
		FROM geth_swaps
		WHERE
		base_asset_id = $1
		AND block_number <= $2
		AND maker_address_id IS NOT NULL
		ORDER BY swap_date asc`,
		*baseAssetID, *blockNumber,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethSwap])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethSwaps, nil
}

func GetGethSwapByTxnHash(dbConnPgx utils.PgxIface, txnHash string, baseAssetID *int) ([]GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
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
		LEFT JOIN geth_addresses addresses ON gs.maker_address_id = addresses.id
		WHERE
		gs.txn_hash = $1 AND
		addresses.address_type_id = $2
		AND gs.base_asset_id  = $3
		ORDER BY gs.swap_date, gs.index_number asc`,
		txnHash, utils.EOA_ADDRESS_TYPE_STRUCTURED_VALUE_ID, *baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethSwap])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethSwaps, nil
}

// bulk swap methods
func GetGethSwapsByTxnHashes(dbConnPgx utils.PgxIface, txnHashes []string, baseAssetID *int) ([]GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
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
		LEFT JOIN geth_addresses addresses ON gs.maker_address_id = addresses.id
		WHERE
		txn_hash = ANY($1)
		AND addresses.address_type_id = $2
		AND gs.base_asset_id  = $3
		ORDER BY gs.swap_date, gs.index_number asc`,
		pq.Array(txnHashes), utils.EOA_ADDRESS_TYPE_STRUCTURED_VALUE_ID, *baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethSwap])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethSwaps, nil
}

func GetDistinctTransactionHashesFromAssetIdAndStartingBlock(dbConnPgx utils.PgxIface, assetID *int, startingBlock *uint64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
		
	SELECT DISTINCT txn_hash FROM geth_swaps
	WHERE block_number >= $1
	AND base_asset_id = $2
	`,
		*startingBlock, *assetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

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

func GetHighestBlockFromBaseAssetId(dbConnPgx utils.PgxIface, assetID *int) (*uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := dbConnPgx.QueryRow(ctx, `SELECT COALESCE (MAX(block_number), 0) FROM geth_swaps
	WHERE base_asset_id=$1
		`,
		*assetID)
	var maxBlockNumber uint64
	err := row.Scan(
		&maxBlockNumber)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &maxBlockNumber, nil
}

func GetDistinctMakerAddressesFromBaseTokenAssetID(dbConnPgx utils.PgxIface, baseAssetID *int) ([]GethSwapAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx,
		`SELECT	DISTINCT 
			maker_address_id as maker_address_id,
			maker_address as maker_address
		FROM geth_swaps
		WHERE
		base_asset_id = $1
		AND maker_address_id IS NOT NULL
		`,
		*baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	makerAddresses := make([]GethSwapAddress, 0)
	for results.Next() {
		var makerAddress GethSwapAddress
		results.Scan(
			&makerAddress.MakerAddressID,
			&makerAddress.MakerAddress,
		)
		makerAddresses = append(makerAddresses, makerAddress)
	}
	return makerAddresses, nil
}

func RemoveGethSwap(dbConnPgx utils.PgxIface, gethSwapID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethSwap DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_swaps WHERE id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *gethSwapID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemoveGethSwapsFromAssetIDAndStartBlockNumber(dbConnPgx utils.PgxIface, baseAssetID *int, startBlockNumber *uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethSwap DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_swaps WHERE base_asset_id = $1 AND block_number >= $2`

	if _, err := dbConnPgx.Exec(ctx, sql, *baseAssetID, *startBlockNumber); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func DeleteGethSwapsByBaseAssetId(dbConnPgx utils.PgxIface, baseAssetID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethSwap DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_swaps WHERE base_asset_id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *baseAssetID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetGethSwapList(dbConnPgx utils.PgxIface) ([]GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		swap_date,
		trade_type_id,
		txn_hash,
		maker_address,
		maker_address_id,
		is_buy,
		price,
		price_usd,
		token1_price_usd,
		total_amount_usd,
		pair_address,
		liquidity_pool_id,
		token0_asset_id,
		token1_asset_id,
		token0_amount,
		token1_Amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		oracle_price_usd,
		oracle_price_asset_id
	FROM geth_swaps `)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethSwap])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethSwaps, nil
}

func UpdateGethSwap(dbConnPgx utils.PgxIface, gethSwap *GethSwap) error {
	// if the gethSwap id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethSwap.ID == nil {
		return errors.New("gethSwap has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethSwap DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_swaps SET 
		chain_id = $1,
		exchange_id = $2,
		block_number= $3,
		index_number= $4,
		swap_date= $5,
		trade_type_id= $6,
		txn_hash= $7,
		maker_address= $8,
		maker_address_id =$9,
		is_buy= $10,
		price=$11,
		price_usd=$12,
		token1_price_usd=$13,
		total_amount_usd=$14,
		pair_address=$15,
		liquidity_pool_id=$16,
		token0_asset_id=$17,
		token1_asset_id=$18,
		token0_amount=$19,
		token1_Amount=$20,
		description=$21,
		updated_by=$22,
		updated_at=current_timestamp at time zone 'UTC',
		geth_process_job_id=$23,
		topics_str$=24,
		status_id=$25,
		base_asset_id=$26,
		oracle_price_usd$27,
  		oracle_price_asset_id$28
		WHERE id=$29`

	if _, err := dbConnPgx.Exec(ctx, sql,
		gethSwap.ChainID,             //1
		gethSwap.ExchangeID,          //2
		gethSwap.BlockNumber,         //3
		gethSwap.IndexNumber,         //4
		gethSwap.SwapDate,            //5
		gethSwap.TradeTypeID,         //6
		gethSwap.TxnHash,             //7
		gethSwap.MakerAddress,        //8
		gethSwap.MakerAddressID,      //9
		gethSwap.IsBuy,               //10
		gethSwap.Price,               //11
		gethSwap.PriceUSD,            //12
		gethSwap.Token1PriceUSD,      //13
		gethSwap.TotalAmountUSD,      //14
		gethSwap.PairAddress,         //15
		gethSwap.LiquidityPoolID,     //16
		gethSwap.Token0AssetId,       //17
		gethSwap.Token1AssetId,       //18
		gethSwap.Token0Amount,        //19
		gethSwap.Token1Amount,        //20
		gethSwap.Description,         //21
		gethSwap.UpdatedBy,           //22
		gethSwap.GethProcessJobID,    //23
		pq.Array(gethSwap.TopicsStr), //24
		gethSwap.StatusID,            //25
		gethSwap.BaseAssetID,         //26
		gethSwap.OraclePriceUSD,      //27
		gethSwap.OraclePriceAssetID,  //28
		gethSwap.ID,                  //29
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethSwap(dbConnPgx utils.PgxIface, gethSwap *GethSwap) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethSwap DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var gethSwapID int
	var gethSwapUUID string
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_swaps
	(
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		swap_date,
		trade_type_id,
		txn_hash,
		maker_address,
		maker_address_id,
		is_buy,
		price,
		price_usd,
		token1_price_usd,
		total_amount_usd,
		pair_address,
		liquidity_pool_id,
		token0_asset_id,
		token1_asset_id,
		token0_amount,
		token1_Amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		oracle_price_usd,
  		oracle_price_asset_id
		) VALUES (
		uuid_generate_v4(),
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15,
		$16,
		$17,
		$18,
		$19,
		$20,
		$21,
		$22,
		current_timestamp at time zone 'UTC',
		$22,
		current_timestamp at time zone 'UTC',
		$23,
		$24,
		$25,
		$26,
		$27,
		$28
		)
		RETURNING id, uuid`,
		gethSwap.ChainID,             //1
		gethSwap.ExchangeID,          //2
		gethSwap.BlockNumber,         //3
		gethSwap.IndexNumber,         //4
		gethSwap.SwapDate,            //5
		gethSwap.TradeTypeID,         //6
		gethSwap.TxnHash,             //7
		gethSwap.MakerAddress,        //8
		gethSwap.MakerAddressID,      //9
		gethSwap.IsBuy,               //10
		gethSwap.Price,               //11
		gethSwap.PriceUSD,            //12
		gethSwap.Token1PriceUSD,      //13
		gethSwap.TotalAmountUSD,      //14
		gethSwap.PairAddress,         //15
		gethSwap.LiquidityPoolID,     //16
		gethSwap.Token0AssetId,       //17
		gethSwap.Token1AssetId,       //18
		gethSwap.Token0Amount,        //19
		gethSwap.Token1Amount,        //20
		gethSwap.Description,         //21
		gethSwap.CreatedBy,           //22
		gethSwap.GethProcessJobID,    //23
		pq.Array(gethSwap.TopicsStr), //24
		gethSwap.StatusID,            //25
		gethSwap.BaseAssetID,         //26
		gethSwap.OraclePriceUSD,      //27
		gethSwap.OraclePriceAssetID,  //28
	).Scan(&gethSwapID, &gethSwapUUID)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, "", err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, "", err
	}
	return int(gethSwapID), gethSwapUUID, nil
}
func InsertGethSwaps(dbConnPgx utils.PgxIface, gethSwaps []GethSwap) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethSwaps {
		gethSwap := gethSwaps[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(gethSwap.UUID)
		row := []interface{}{
			uuidString,                   //1
			gethSwap.ChainID,             //2
			gethSwap.ExchangeID,          //3
			gethSwap.BlockNumber,         //4
			gethSwap.IndexNumber,         //5
			gethSwap.SwapDate,            //6
			gethSwap.TradeTypeID,         //7
			gethSwap.TxnHash,             //8
			gethSwap.MakerAddress,        //9
			gethSwap.MakerAddressID,      //10
			gethSwap.IsBuy,               //11
			gethSwap.Price,               //12
			gethSwap.PriceUSD,            //13
			gethSwap.Token1PriceUSD,      //14
			gethSwap.TotalAmountUSD,      //15
			gethSwap.PairAddress,         //16
			gethSwap.LiquidityPoolID,     //17
			gethSwap.Token0AssetId,       //18
			gethSwap.Token1AssetId,       //19
			gethSwap.Token0Amount,        //20
			gethSwap.Token1Amount,        //21
			gethSwap.Description,         //22
			gethSwap.CreatedBy,           //23
			&now,                         //24
			gethSwap.CreatedBy,           //25
			&now,                         //26
			gethSwap.GethProcessJobID,    //27
			pq.Array(gethSwap.TopicsStr), //28
			gethSwap.StatusID,            //29
			gethSwap.BaseAssetID,         //30
			gethSwap.OraclePriceUSD,      //31
			gethSwap.OraclePriceAssetID,  //32
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_swaps"},
		[]string{
			"uuid",                  //1
			"chain_id",              //2
			"exchange_id",           //3
			"block_number",          //4
			"index_number",          //5
			"swap_date",             //6
			"trade_type_id",         //7
			"txn_hash",              //8
			"maker_address",         //9
			"maker_address_id",      //10
			"is_buy",                //11
			"price",                 //12
			"price_usd",             //13
			"token1_price_usd",      //14
			"total_amount_usd",      //15
			"pair_address",          //16
			"liquidity_pool_id",     //17
			"token0_asset_id",       //18
			"token1_asset_id",       //19
			"token0_amount",         //20
			"token1_amount",         //21
			"description",           //22
			"created_by",            //23
			"created_at",            //24
			"updated_by",            //25
			"updated_at",            //26
			"geth_process_job_id",   //27
			"topics_str",            //28
			"status_id",             //29
			"base_asset_id",         //30
			"oracle_price_usd",      //31
			"oracle_price_asset_id", //32
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethSwaps: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetNullAddressStrsFromSwaps(dbConnPgx utils.PgxIface, assetID *int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
		SELECT DISTINCT gs.maker_address as address  
		FROM geth_swaps gs
		LEFT JOIN geth_addresses as ga
			ON LOWER(gs.maker_address) = LOWER(ga.address_str)
		WHERE gs.maker_address_id IS NULL
		AND gs.base_asset_id = $1
		AND ga.id IS NULL
		`, *assetID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethNullAddressStrs := make([]string, 0)
	for results.Next() {
		var gethNullAddressStr string
		results.Scan(
			&gethNullAddressStr,
		)
		gethNullAddressStrs = append(gethNullAddressStrs, gethNullAddressStr)
	}
	return gethNullAddressStrs, nil
}

func UpdateGethSwapAddresses(dbConnPgx utils.PgxIface, baseAssetID *int) error {
	// update address ids from existing addresses in geth_addresses
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethSwapAddresses DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `
		UPDATE geth_swaps as gs SET
		maker_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gs.maker_address) = LOWER(ga.address_str)
			AND gs.maker_address_id IS NULL
			AND gs.base_asset_id = $1;
	`

	if _, err := dbConnPgx.Exec(ctx, sql, *baseAssetID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

// for refinedev
func GetGethSwapListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		swap_date,
		trade_type_id,
		txn_hash,
		maker_address,
		maker_address_id,
		is_buy,
		price,
		price_usd,
		token1_price_usd,
		total_amount_usd,
		pair_address,
		liquidity_pool_id,
		token0_asset_id,
		token1_asset_id,
		token0_amount,
		token1_Amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		oracle_price_usd,
		oracle_price_asset_id
	FROM geth_swaps 
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

	gethSwaps, err := pgx.CollectRows(results, pgx.RowToStructByName[GethSwap])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return gethSwaps, nil
}

func GetTotalGethSwapsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_swaps
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
