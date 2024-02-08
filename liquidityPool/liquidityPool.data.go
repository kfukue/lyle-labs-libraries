package liquiditypool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetLiquidityPool(liquidityPoolID int) (*LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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
	WHERE id = $1`, liquidityPoolID)

	liquidityPool := &LiquidityPool{}
	err := row.Scan(
		&liquidityPool.ID,
		&liquidityPool.UUID,
		&liquidityPool.Name,
		&liquidityPool.AlternateName,
		&liquidityPool.PairAddress,
		&liquidityPool.ChainID,
		&liquidityPool.ExchangeID,
		&liquidityPool.LiquidityPoolTypeID,
		&liquidityPool.Token0ID,
		&liquidityPool.Token1ID,
		&liquidityPool.Url,
		&liquidityPool.StartBlock,
		&liquidityPool.LatestBlockSynced,
		&liquidityPool.CreatedTxnHash,
		&liquidityPool.IsActive,
		&liquidityPool.Description,
		&liquidityPool.CreatedBy,
		&liquidityPool.CreatedAt,
		&liquidityPool.UpdatedBy,
		&liquidityPool.UpdatedAt,
		&liquidityPool.BaseAssetID,
		&liquidityPool.QuoteAssetID,
		&liquidityPool.QuoteAssetChainlinkAddress,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return liquidityPool, nil
}

func RemoveLiquidityPool(liquidityPoolID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM liquidity_pools WHERE id = $1`, liquidityPoolID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetLiquidityPools() ([]LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	FROM liquidity_pools`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	liquidityPools := make([]LiquidityPool, 0)
	for results.Next() {
		var liquidityPool LiquidityPool
		results.Scan(
			&liquidityPool.ID,
			&liquidityPool.UUID,
			&liquidityPool.Name,
			&liquidityPool.AlternateName,
			&liquidityPool.PairAddress,
			&liquidityPool.ChainID,
			&liquidityPool.ExchangeID,
			&liquidityPool.LiquidityPoolTypeID,
			&liquidityPool.Token0ID,
			&liquidityPool.Token1ID,
			&liquidityPool.Url,
			&liquidityPool.StartBlock,
			&liquidityPool.LatestBlockSynced,
			&liquidityPool.CreatedTxnHash,
			&liquidityPool.IsActive,
			&liquidityPool.Description,
			&liquidityPool.CreatedBy,
			&liquidityPool.CreatedAt,
			&liquidityPool.UpdatedBy,
			&liquidityPool.UpdatedAt,
			&liquidityPool.BaseAssetID,
			&liquidityPool.QuoteAssetID,
			&liquidityPool.QuoteAssetChainlinkAddress,
		)

		liquidityPools = append(liquidityPools, liquidityPool)
	}
	return liquidityPools, nil
}

func GetLiquidityPoolList(ids []int) ([]LiquidityPool, error) {
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
	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	liquidityPools := make([]LiquidityPool, 0)
	for results.Next() {
		var liquidityPool LiquidityPool
		results.Scan(
			&liquidityPool.ID,
			&liquidityPool.UUID,
			&liquidityPool.Name,
			&liquidityPool.AlternateName,
			&liquidityPool.PairAddress,
			&liquidityPool.ChainID,
			&liquidityPool.ExchangeID,
			&liquidityPool.LiquidityPoolTypeID,
			&liquidityPool.Token0ID,
			&liquidityPool.Token1ID,
			&liquidityPool.Url,
			&liquidityPool.StartBlock,
			&liquidityPool.LatestBlockSynced,
			&liquidityPool.CreatedTxnHash,
			&liquidityPool.IsActive,
			&liquidityPool.Description,
			&liquidityPool.CreatedBy,
			&liquidityPool.CreatedAt,
			&liquidityPool.UpdatedBy,
			&liquidityPool.UpdatedAt,
			&liquidityPool.BaseAssetID,
			&liquidityPool.QuoteAssetID,
			&liquidityPool.QuoteAssetChainlinkAddress,
		)
		liquidityPools = append(liquidityPools, liquidityPool)
	}
	return liquidityPools, nil
}

func GetLiquidityPoolListByToken0(asset0ID *int) ([]LiquidityPoolWithTokens, error) {
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
	token1.import_geth_initial as token1_import_geth_initial
	FROM liquidity_pools lp
	LEFT JOIN assets token0 ON lp.token0_id = token0.id
	LEFT JOIN assets token1 ON lp.token1_id = token1.id
	WHERE lp.base_asset_id = $1
	`
	results, err := database.DbConnPgx.Query(ctx, sql, asset0ID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
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
		)
		liquidityPoolsWithTokens = append(liquidityPoolsWithTokens, liquidityPoolWithTokens)
	}
	return liquidityPoolsWithTokens, nil
}

func GetLiquidityPoolsByUUIDs(UUIDList []string) ([]LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	defer results.Close()
	liquidityPools := make([]LiquidityPool, 0)
	for results.Next() {
		var liquidityPool LiquidityPool
		results.Scan(
			&liquidityPool.ID,
			&liquidityPool.UUID,
			&liquidityPool.Name,
			&liquidityPool.AlternateName,
			&liquidityPool.PairAddress,
			&liquidityPool.ChainID,
			&liquidityPool.ExchangeID,
			&liquidityPool.LiquidityPoolTypeID,
			&liquidityPool.Token0ID,
			&liquidityPool.Token1ID,
			&liquidityPool.Url,
			&liquidityPool.StartBlock,
			&liquidityPool.LatestBlockSynced,
			&liquidityPool.CreatedTxnHash,
			&liquidityPool.IsActive,
			&liquidityPool.Description,
			&liquidityPool.CreatedBy,
			&liquidityPool.CreatedAt,
			&liquidityPool.UpdatedBy,
			&liquidityPool.UpdatedAt,
			&liquidityPool.BaseAssetID,
			&liquidityPool.QuoteAssetID,
			&liquidityPool.QuoteAssetChainlinkAddress,
		)

		liquidityPools = append(liquidityPools, liquidityPool)
	}
	return liquidityPools, nil
}

func GetStartAndEndDateDiffLiquidityPools(diffInDate int) ([]LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	liquidityPools := make([]LiquidityPool, 0)
	for results.Next() {
		var liquidityPool LiquidityPool
		results.Scan(
			&liquidityPool.ID,
			&liquidityPool.UUID,
			&liquidityPool.Name,
			&liquidityPool.AlternateName,
			&liquidityPool.PairAddress,
			&liquidityPool.ChainID,
			&liquidityPool.ExchangeID,
			&liquidityPool.LiquidityPoolTypeID,
			&liquidityPool.Token0ID,
			&liquidityPool.Token1ID,
			&liquidityPool.Url,
			&liquidityPool.StartBlock,
			&liquidityPool.LatestBlockSynced,
			&liquidityPool.CreatedTxnHash,
			&liquidityPool.IsActive,
			&liquidityPool.Description,
			&liquidityPool.CreatedBy,
			&liquidityPool.CreatedAt,
			&liquidityPool.UpdatedBy,
			&liquidityPool.UpdatedAt,
			&liquidityPool.BaseAssetID,
			&liquidityPool.QuoteAssetID,
			&liquidityPool.QuoteAssetChainlinkAddress,
		)

		liquidityPools = append(liquidityPools, liquidityPool)
	}
	return liquidityPools, nil
}

func UpdateLiquidityPool(liquidityPool LiquidityPool) error {
	// if the liquidityPool id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if liquidityPool.ID == nil || *liquidityPool.ID == 0 {
		return errors.New("liquidityPool has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE liquidity_pools SET 
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
		WHERE id=$19`,
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
		liquidityPool.ID)                         //19

	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertLiquidityPool(liquidityPool LiquidityPool) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	// layoutPostgres := utils.LayoutPostgres
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO liquidity_pools 
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
		RETURNING id`,
		liquidityPool.Name,                       //1
		liquidityPool.AlternateName,              //2
		&liquidityPool.PairAddress,               //3
		&liquidityPool.ChainID,                   //4
		&liquidityPool.ExchangeID,                //5
		&liquidityPool.LiquidityPoolTypeID,       //6
		&liquidityPool.Token0ID,                  //7
		&liquidityPool.Token1ID,                  //8
		&liquidityPool.Url,                       //9
		&liquidityPool.StartBlock,                //10
		&liquidityPool.LatestBlockSynced,         //11
		&liquidityPool.CreatedTxnHash,            //12
		&liquidityPool.IsActive,                  //13
		liquidityPool.Description,                //14
		liquidityPool.CreatedBy,                  //15
		liquidityPool.BaseAssetID,                //16
		liquidityPool.QuoteAssetID,               //17
		liquidityPool.QuoteAssetChainlinkAddress, //18
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertLiquidityPools(liquidityPools []LiquidityPool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range liquidityPools {
		liquidityPool := liquidityPools[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(liquidityPool.UUID)
		row := []interface{}{
			uuidString,                               //1
			liquidityPool.Name,                       //2
			liquidityPool.AlternateName,              //3
			liquidityPool.PairAddress,                //4
			*liquidityPool.ChainID,                   //5
			*liquidityPool.ExchangeID,                //6
			*liquidityPool.LiquidityPoolTypeID,       //7
			*liquidityPool.Token0ID,                  //8
			*liquidityPool.Token1ID,                  //9
			liquidityPool.Url,                        //10
			*liquidityPool.StartBlock,                //11
			*liquidityPool.LatestBlockSynced,         //12
			liquidityPool.CreatedTxnHash,             //13
			liquidityPool.IsActive,                   //14
			liquidityPool.Description,                //15
			liquidityPool.CreatedBy,                  //16
			&now,                                     //17
			liquidityPool.CreatedBy,                  //18
			&now,                                     //19
			*liquidityPool.BaseAssetID,               //20
			liquidityPool.QuoteAssetID,               //21
			liquidityPool.QuoteAssetChainlinkAddress, //22
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}
	return nil
}

// liquidityPool chain methods

func UpdateLiquidityPoolAssetByUUID(liquidityPoolAsset LiquidityPoolAsset) error {
	// if the liquidityPool id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if liquidityPoolAsset.LiquidityPoolID == nil || *liquidityPoolAsset.LiquidityPoolID == 0 || liquidityPoolAsset.UUID == "" {
		return errors.New("liquidityPoolAsset has invalid IDs")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE liquidity_pool_assets SET 
		liquidity_pool_id=$1,
		asset_id=$2,
		token_number=$3,
		name=$4,
		alternate_name=$5,
		description=$6,
		updated_by=$7, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE 
		uuid=$8,`,
		liquidityPoolAsset.LiquidityPoolID, //1
		liquidityPoolAsset.AssetID,         //2
		liquidityPoolAsset.TokenNumber,     //3
		liquidityPoolAsset.Name,            //4
		liquidityPoolAsset.AlternateName,   //5
		liquidityPoolAsset.Description,     //6
		liquidityPoolAsset.UpdatedBy,       //7
		liquidityPoolAsset.UUID)            //8
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertLiquidityPoolAsset(liquidityPoolAsset LiquidityPoolAsset) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	// layoutPostgres := utils.LayoutPostgres
	_, err := database.DbConnPgx.Query(ctx, `INSERT INTO liquidity_pool_assets 
	(
		uuid,
		liquidity_pool_id,
		asset_id,
		token_number,
		name,
		alternate_name,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
		) VALUES (
			uuid_generate_v4(),
			$1, 
			$2, 
			$3, 
			$4, 
			$5, 
			$6, 
			$7, 
			current_timestamp at time zone 'UTC',
			$7,
			current_timestamp at time zone 'UTC'
		)
		`,
		liquidityPoolAsset.LiquidityPoolID, //1
		liquidityPoolAsset.AssetID,         //2
		liquidityPoolAsset.TokenNumber,     //3
		liquidityPoolAsset.Name,            //4
		liquidityPoolAsset.AlternateName,   //5
		liquidityPoolAsset.Description,     //6
		liquidityPoolAsset.CreatedBy,       //7
	)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertLiquidityPoolAssets(liquidityPoolAssets []LiquidityPoolAsset) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range liquidityPoolAssets {
		liquidityPoolAssets := liquidityPoolAssets[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(liquidityPoolAssets.UUID)
		row := []interface{}{
			uuidString,                          //1
			liquidityPoolAssets.LiquidityPoolID, //2
			liquidityPoolAssets.AssetID,         //3
			liquidityPoolAssets.TokenNumber,     //4
			liquidityPoolAssets.Name,            //5
			liquidityPoolAssets.AlternateName,   //6
			liquidityPoolAssets.Description,     //7
			liquidityPoolAssets.CreatedBy,       //8
			&now,                                //9
			liquidityPoolAssets.CreatedBy,       //10
			&now,                                //11
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"liquidity_pool_assets"},
		[]string{
			"uuid",              //1
			"liquidity_pool_id", //2
			"asset_id",          //3
			"token_number",      //4
			"name",              //5
			"alternate_name",    //6
			"description",       //7
			"created_by",        //8
			"created_at",        //9
			"updated_by",        //10
			"updated_at",        //11
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}
	return nil
}

// for refinedev
func GetLiquidityPoolListByPagination(_start, _end *int, _order, _sort string, _filters []string) ([]LiquidityPool, error) {
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

	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	liquidityPools := make([]LiquidityPool, 0)
	for results.Next() {
		var liquidityPool LiquidityPool
		results.Scan(
			&liquidityPool.ID,
			&liquidityPool.UUID,
			&liquidityPool.Name,
			&liquidityPool.AlternateName,
			&liquidityPool.PairAddress,
			&liquidityPool.ChainID,
			&liquidityPool.ExchangeID,
			&liquidityPool.LiquidityPoolTypeID,
			&liquidityPool.Token0ID,
			&liquidityPool.Token1ID,
			&liquidityPool.Url,
			&liquidityPool.StartBlock,
			&liquidityPool.LatestBlockSynced,
			&liquidityPool.CreatedTxnHash,
			&liquidityPool.IsActive,
			&liquidityPool.Description,
			&liquidityPool.CreatedBy,
			&liquidityPool.CreatedAt,
			&liquidityPool.UpdatedBy,
			&liquidityPool.UpdatedAt,
			&liquidityPool.BaseAssetID,
			&liquidityPool.QuoteAssetID,
			&liquidityPool.QuoteAssetChainlinkAddress,
		)

		liquidityPools = append(liquidityPools, liquidityPool)
	}
	return liquidityPools, nil
}

func GetTotalLiquidityPoolCount() (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM liquidity_pools
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
