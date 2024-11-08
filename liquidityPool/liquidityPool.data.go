package liquiditypool

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

func GetLiquidityPool(dbConnPgx utils.PgxIface, liquidityPoolID *int) (*LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		name,
		alternate_name,
		pair_address,
		chain_id,
		exchange_id,
		liquidity_pool_type_id,
		token0_id,
		token1_id,
		url,
		start_block,
		latest_block_synced,
		created_txn_hash,
		IsActive,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
		base_asset_id,
		quote_asset_id,
		quote_asset_chainlink_address_usd

	FROM liquidity_pools 
	WHERE id = $1`, *liquidityPoolID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	liquidityPool, err := pgx.CollectOneRow(row, pgx.RowToStructByName[LiquidityPool])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &liquidityPool, nil
}

func RemoveLiquidityPool(dbConnPgx utils.PgxIface, liquidityPoolID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM liquidity_pools WHERE id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *liquidityPoolID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetLiquidityPoolList(dbConnPgx utils.PgxIface, ids []int) ([]LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		pair_address,
		chain_id,
		exchange_id,
		liquidity_pool_type_id,
		token0_id,
		token1_id,
		url,
		start_block,
		latest_block_synced,
		created_txn_hash,
		IsActive,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
		base_asset_id,
		quote_asset_id,
		quote_asset_chainlink_address_usd
	FROM liquidity_pools`
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

	liquidityPools, err := pgx.CollectRows(results, pgx.RowToStructByName[LiquidityPool])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return liquidityPools, nil
}

func GetLiquidityPoolListByBaseAssetID(dbConnPgx utils.PgxIface, baseAssetID *int) ([]LiquidityPoolWithTokens, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	lp.id,
	lp.uuid,
	lp.name,
	lp.alternate_name,
	lp.pair_address,
	lp.chain_id,
	lp.exchange_id,
	lp.liquidity_pool_type_id,
	lp.token0_id,
	lp.token1_id,
	lp.url,
	lp.start_block,
	lp.latest_block_synced,
	lp.created_txn_hash,
	lp.IsActive,
	lp.description,
	lp.created_by, 
	lp.created_at, 
	lp.updated_by, 
	lp.updated_at,
	lp.base_asset_id,
	lp.quote_asset_id,
	lp.quote_asset_chainlink_address_usd,
	-- asset0
	token0.id as token0_id,
	token0.uuid as token0_uuid, 
	token0.name as token0_name, 
	token0.alternate_name as token0_alternate_name, 
	token0.cusip as token0_cusip,
	token0.ticker as token0_ticker,
	token0.base_asset_id as token0_base_asset_id,
	token0.quote_asset_id as token0_quote_asset_id,
	token0.description as token0_description,
	token0.asset_type_id as token0_asset_type_id,
	token0.created_by as token0_created_by, 
	token0.created_at as token0_created_at, 
	token0.updated_by as token0_updated_by, 
	token0.updated_at as token0_updated_at,
	token0.chain_id as token0_chain_id,
	token0.category_id as token0_category_id,
	token0.sub_category_id as token0_sub_category_id,
	token0.is_default_quote as token0_is_default_quote,
	token0.ignore_market_data as token0_ignore_market_data,
	token0.decimals as token0_decimals,
	token0.contract_address as token0_contract_address,
	token0.starting_block_number as token0_block_number,
	token0.import_geth as token0_import_geth,
	token0.import_geth_initial as token0_import_geth_initial,
	token0.chainlink_usd_address as token0_chainlink_usd_address,
	token0.chainlink_usd_chain_id as token0_chainlink_usd_chain_id,
	token0.total_supply as token0_total_supply,
	--asset 1
	token1.id as token1_id,
	token1.uuid as token1_uuid, 
	token1.name as token1_name, 
	token1.alternate_name as token1_alternate_name, 
	token1.cusip as token1_cusip,
	token1.ticker as token1_ticker,
	token1.base_asset_id as token1_base_asset_id,
	token1.quote_asset_id as token1_quote_asset_id,
	token1.description as token1_description,
	token1.asset_type_id as token1_asset_type_id,
	token1.created_by as token1_created_by, 
	token1.created_at as token1_created_at, 
	token1.updated_by as token1_updated_by, 
	token1.updated_at as token1_updated_at,
	token1.chain_id as token1_chain_id,
	token1.category_id as token1_category_id,
	token1.sub_category_id as token1_sub_category_id,
	token1.is_default_quote as token1_is_default_quote,
	token1.ignore_market_data as token1_ignore_market_data,
	token1.decimals as token1_decimals,
	token1.contract_address as token1_contract_address,
	token1.starting_block_number as token1_block_number,
	token1.import_geth as token1_import_geth,
	token1.import_geth_initial as token1_import_geth_initial,
	token1.chainlink_usd_address as token1_chainlink_usd_address,
	token1.chainlink_usd_chain_id as token1_chainlink_usd_chain_id,
	token1.total_supply as token1_total_supply
	FROM liquidity_pools lp
	LEFT JOIN assets token0 ON lp.token0_id = token0.id
	LEFT JOIN assets token1 ON lp.token1_id = token1.id
	WHERE lp.base_asset_id = $1
	`
	results, err := dbConnPgx.Query(ctx, sql, *baseAssetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	liquidityPoolsWithTokens := make([]LiquidityPoolWithTokens, 0)
	for results.Next() {
		var liquidityPoolWithTokens LiquidityPoolWithTokens
		results.Scan(
			&liquidityPoolWithTokens.ID,
			&liquidityPoolWithTokens.UUID,
			&liquidityPoolWithTokens.Name,
			&liquidityPoolWithTokens.AlternateName,
			&liquidityPoolWithTokens.PairAddress,
			&liquidityPoolWithTokens.ChainID,
			&liquidityPoolWithTokens.ExchangeID,
			&liquidityPoolWithTokens.LiquidityPoolTypeID,
			&liquidityPoolWithTokens.Token0ID,
			&liquidityPoolWithTokens.Token1ID,
			&liquidityPoolWithTokens.Url,
			&liquidityPoolWithTokens.StartBlock,
			&liquidityPoolWithTokens.LatestBlockSynced,
			&liquidityPoolWithTokens.CreatedTxnHash,
			&liquidityPoolWithTokens.IsActive,
			&liquidityPoolWithTokens.Description,
			&liquidityPoolWithTokens.CreatedBy,
			&liquidityPoolWithTokens.CreatedAt,
			&liquidityPoolWithTokens.UpdatedBy,
			&liquidityPoolWithTokens.UpdatedAt,
			&liquidityPoolWithTokens.BaseAssetID,
			&liquidityPoolWithTokens.QuoteAssetID,
			&liquidityPoolWithTokens.QuoteAssetChainlinkAddress,
			// token 0
			&liquidityPoolWithTokens.Token0.ID,
			&liquidityPoolWithTokens.Token0.UUID,
			&liquidityPoolWithTokens.Token0.Name,
			&liquidityPoolWithTokens.Token0.AlternateName,
			&liquidityPoolWithTokens.Token0.Cusip,
			&liquidityPoolWithTokens.Token0.Ticker,
			&liquidityPoolWithTokens.Token0.BaseAssetID,
			&liquidityPoolWithTokens.Token0.QuoteAssetID,
			&liquidityPoolWithTokens.Token0.Description,
			&liquidityPoolWithTokens.Token0.AssetTypeID,
			&liquidityPoolWithTokens.Token0.CreatedBy,
			&liquidityPoolWithTokens.Token0.CreatedAt,
			&liquidityPoolWithTokens.Token0.UpdatedBy,
			&liquidityPoolWithTokens.Token0.UpdatedAt,
			&liquidityPoolWithTokens.Token0.ChainID,
			&liquidityPoolWithTokens.Token0.CategoryID,
			&liquidityPoolWithTokens.Token0.SubCategoryID,
			&liquidityPoolWithTokens.Token0.IsDefaultQuote,
			&liquidityPoolWithTokens.Token0.IgnoreMarketData,
			&liquidityPoolWithTokens.Token0.Decimals,
			&liquidityPoolWithTokens.Token0.ContractAddress,
			&liquidityPoolWithTokens.Token0.StartingBlockNumber,
			&liquidityPoolWithTokens.Token0.ImportGeth,
			&liquidityPoolWithTokens.Token0.ImportGethInitial,
			&liquidityPoolWithTokens.Token0.ChainlinkUSDAddress,
			&liquidityPoolWithTokens.Token0.ChainlinkUSDChainID,
			&liquidityPoolWithTokens.Token0.TotalSupply,
			//token 1
			&liquidityPoolWithTokens.Token1.ID,
			&liquidityPoolWithTokens.Token1.UUID,
			&liquidityPoolWithTokens.Token1.Name,
			&liquidityPoolWithTokens.Token1.AlternateName,
			&liquidityPoolWithTokens.Token1.Cusip,
			&liquidityPoolWithTokens.Token1.Ticker,
			&liquidityPoolWithTokens.Token1.BaseAssetID,
			&liquidityPoolWithTokens.Token1.QuoteAssetID,
			&liquidityPoolWithTokens.Token1.Description,
			&liquidityPoolWithTokens.Token1.AssetTypeID,
			&liquidityPoolWithTokens.Token1.CreatedBy,
			&liquidityPoolWithTokens.Token1.CreatedAt,
			&liquidityPoolWithTokens.Token1.UpdatedBy,
			&liquidityPoolWithTokens.Token1.UpdatedAt,
			&liquidityPoolWithTokens.Token1.ChainID,
			&liquidityPoolWithTokens.Token1.CategoryID,
			&liquidityPoolWithTokens.Token1.SubCategoryID,
			&liquidityPoolWithTokens.Token1.IsDefaultQuote,
			&liquidityPoolWithTokens.Token1.IgnoreMarketData,
			&liquidityPoolWithTokens.Token1.Decimals,
			&liquidityPoolWithTokens.Token1.ContractAddress,
			&liquidityPoolWithTokens.Token1.StartingBlockNumber,
			&liquidityPoolWithTokens.Token1.ImportGeth,
			&liquidityPoolWithTokens.Token1.ImportGethInitial,
			&liquidityPoolWithTokens.Token1.ChainlinkUSDAddress,
			&liquidityPoolWithTokens.Token1.ChainlinkUSDChainID,
			&liquidityPoolWithTokens.Token1.TotalSupply,
		)
		liquidityPoolsWithTokens = append(liquidityPoolsWithTokens, liquidityPoolWithTokens)
	}
	return liquidityPoolsWithTokens, nil
}

func GetLiquidityPoolsByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		pair_address,
		chain_id,
		exchange_id,
		liquidity_pool_type_id,
		token0_id,
		token1_id,
		url,
		start_block,
		latest_block_synced,
		created_txn_hash,
		IsActive,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
		base_asset_id,
		quote_asset_id,
		quote_asset_chainlink_address_usd
	FROM liquidity_pools
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	liquidityPools, err := pgx.CollectRows(results, pgx.RowToStructByName[LiquidityPool])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return liquidityPools, nil
}
func UpdateLiquidityPool(dbConnPgx utils.PgxIface, liquidityPool *LiquidityPool) error {
	// if the liquidityPool id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if liquidityPool.ID == nil || *liquidityPool.ID == 0 {
		return errors.New("liquidityPool has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateLiquidityPool DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE liquidity_pools SET 
		name=$1,
		alternate_name=$2,
		pair_address=$3,
		chain_id=$4,
		exchange_id=$5,
		liquidity_pool_type_id=$6,
		token0_id=$7,
		token1_id=$8,
		url=$9,
		start_block=$10,
		latest_block_synced=$11,
		created_txn_hash=$12,
		IsActive=$13,
		description=$14,
		updated_by=$15, 
		updated_at=current_timestamp at time zone 'UTC',
		base_asset_id=$16,
		quote_asset_id=$17,
		quote_asset_chainlink_address_usd=$18
		WHERE id=$19`

	if _, err := dbConnPgx.Exec(ctx, sql,
		liquidityPool.Name,                       //1
		liquidityPool.AlternateName,              //2
		liquidityPool.PairAddress,                //3
		liquidityPool.ChainID,                    //4
		liquidityPool.ExchangeID,                 //5
		liquidityPool.LiquidityPoolTypeID,        //6
		liquidityPool.Token0ID,                   //7
		liquidityPool.Token1ID,                   //8
		liquidityPool.Url,                        //9
		liquidityPool.StartBlock,                 //10
		liquidityPool.LatestBlockSynced,          //11
		liquidityPool.CreatedTxnHash,             //12
		liquidityPool.IsActive,                   //13
		liquidityPool.Description,                //14
		liquidityPool.UpdatedBy,                  //15
		liquidityPool.BaseAssetID,                //16
		liquidityPool.QuoteAssetID,               //17
		liquidityPool.QuoteAssetChainlinkAddress, //18
		liquidityPool.ID,                         //19
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertLiquidityPool(dbConnPgx utils.PgxIface, liquidityPool *LiquidityPool) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertLiquidityPool DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var insertID int
	var insertUUID string
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO liquidity_pools 
	(
		uuid,
		name,
		alternate_name,
		pair_address,
		chain_id,
		exchange_id,
		liquidity_pool_type_id,
		token0_id,
		token1_id,
		url,
		start_block,
		latest_block_synced,
		created_txn_hash,
		IsActive,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
		base_asset_id,
		quote_asset_id,
		quote_asset_chainlink_address_usd
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
			current_timestamp at time zone 'UTC',
			$15,
			current_timestamp at time zone 'UTC',
			$16,
			$17,
			$18
		)
		RETURNING id, uuid`,
		liquidityPool.Name,                       //1
		liquidityPool.AlternateName,              //2
		liquidityPool.PairAddress,                //3
		liquidityPool.ChainID,                    //4
		liquidityPool.ExchangeID,                 //5
		liquidityPool.LiquidityPoolTypeID,        //6
		liquidityPool.Token0ID,                   //7
		liquidityPool.Token1ID,                   //8
		liquidityPool.Url,                        //9
		liquidityPool.StartBlock,                 //10
		liquidityPool.LatestBlockSynced,          //11
		liquidityPool.CreatedTxnHash,             //12
		liquidityPool.IsActive,                   //13
		liquidityPool.Description,                //14
		liquidityPool.CreatedBy,                  //15
		liquidityPool.BaseAssetID,                //16
		liquidityPool.QuoteAssetID,               //17
		liquidityPool.QuoteAssetChainlinkAddress, //18
	).Scan(&insertID, &insertUUID)
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
	return int(insertID), insertUUID, nil
}
func InsertLiquidityPools(dbConnPgx utils.PgxIface, liquidityPools []LiquidityPool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range liquidityPools {
		liquidityPool := liquidityPools[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(liquidityPool.UUID)
		row := []interface{}{
			uuidString,                               //1
			liquidityPool.Name,                       //2
			liquidityPool.AlternateName,              //3
			liquidityPool.PairAddress,                //4
			liquidityPool.ChainID,                    //5
			liquidityPool.ExchangeID,                 //6
			liquidityPool.LiquidityPoolTypeID,        //7
			liquidityPool.Token0ID,                   //8
			liquidityPool.Token1ID,                   //9
			liquidityPool.Url,                        //10
			liquidityPool.StartBlock,                 //11
			liquidityPool.LatestBlockSynced,          //12
			liquidityPool.CreatedTxnHash,             //13
			liquidityPool.IsActive,                   //14
			liquidityPool.Description,                //15
			liquidityPool.CreatedBy,                  //16
			&now,                                     //17
			liquidityPool.CreatedBy,                  //18
			&now,                                     //19
			liquidityPool.BaseAssetID,                //20
			liquidityPool.QuoteAssetID,               //21
			liquidityPool.QuoteAssetChainlinkAddress, //22
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"liquidity_pools"},
		[]string{
			"uuid",                              //1
			"name",                              //2
			"alternate_name",                    //3
			"pair_address",                      //4
			"chain_id",                          //5
			"exchange_id",                       //6
			"liquidity_pool_type_id",            //7
			"token0_id",                         //8
			"token1_id",                         //9
			"url",                               //10
			"start_block",                       //11
			"latest_block_synced",               //12
			"created_txn_hash",                  //13
			"IsActive",                          //14
			"description",                       //15
			"created_by",                        //16
			"created_at",                        //17
			"updated_by",                        //18
			"updated_at",                        //19
			"base_asset_id",                     //20
			"quote_asset_id",                    //21
			"quote_asset_chainlink_address_usd", //22
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertLiquidityPools: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetLiquidityPoolListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		pair_address,
		chain_id,
		exchange_id,
		liquidity_pool_type_id,
		token0_id,
		token1_id,
		url,
		start_block,
		latest_block_synced,
		created_txn_hash,
		IsActive,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
		base_asset_id,
		quote_asset_id,
		quote_asset_chainlink_address_usd
	FROM liquidity_pools
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

	liquidityPools, err := pgx.CollectRows(results, pgx.RowToStructByName[LiquidityPool])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return liquidityPools, nil
}

func GetTotalLiquidityPoolCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM liquidity_pools
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
