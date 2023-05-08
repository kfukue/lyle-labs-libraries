package liquiditypool

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetLiquidityPool(liquidityPoolID int) (*LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
		updated_at
	FROM liquidityPools 
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
	)
	if err == sql.ErrNoRows {
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
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM liquidityPools WHERE id = $1`, liquidityPoolID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetLiquidityPools() ([]LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
		updated_at
	FROM liquidityPools`)
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
		updated_at
	FROM liquidityPools`
	if len(ids) > 0 {
		strIds := utils.SplitToString(ids, ",")
		additionalQuery := fmt.Sprintf(` WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := database.DbConn.QueryContext(ctx, sql)
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
		)
		liquidityPools = append(liquidityPools, liquidityPool)
	}
	return liquidityPools, nil
}

func GetLiquidityPoolsByUUIDs(UUIDList []string) ([]LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
		updated_at
	FROM liquidityPools
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
		)

		liquidityPools = append(liquidityPools, liquidityPool)
	}
	return liquidityPools, nil
}

func GetStartAndEndDateDiffLiquidityPools(diffInDate int) ([]LiquidityPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
		updated_at
	FROM liquidityPools
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
	_, err := database.DbConn.ExecContext(ctx, `UPDATE liquidityPools SET 
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
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$16`,
		liquidityPool.Name,                //1
		liquidityPool.AlternateName,       //2
		liquidityPool.PairAddress,         //3
		liquidityPool.ChainID,             //4
		liquidityPool.ExchangeID,          //5
		liquidityPool.LiquidityPoolTypeID, //6
		liquidityPool.Token0ID,            //7
		liquidityPool.Token1ID,            //8
		liquidityPool.Url,                 //9
		liquidityPool.StartBlock,          //10
		liquidityPool.LatestBlockSynced,   //11
		liquidityPool.CreatedTxnHash,      //12
		liquidityPool.IsActive,            //13
		liquidityPool.Description,         //14
		liquidityPool.UpdatedBy,           //15
		liquidityPool.ID)                  //16
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
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO liquidityPools 
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
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		liquidityPool.Name,                 //1
		liquidityPool.AlternateName,        //2
		&liquidityPool.PairAddress,         //3
		&liquidityPool.ChainID,             //4
		&liquidityPool.ExchangeID,          //5
		&liquidityPool.LiquidityPoolTypeID, //6
		&liquidityPool.Token0ID,            //7
		&liquidityPool.Token1ID,            //8
		&liquidityPool.Url,                 //9
		&liquidityPool.StartBlock,          //10
		&liquidityPool.LatestBlockSynced,   //11
		&liquidityPool.CreatedTxnHash,      //12
		&liquidityPool.IsActive,            //13
		liquidityPool.Description,          //14
		liquidityPool.CreatedBy,            //15
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
			uuidString,                         //1
			liquidityPool.Name,                 //2
			liquidityPool.AlternateName,        //3
			liquidityPool.PairAddress,          //4
			*liquidityPool.ChainID,             //5
			*liquidityPool.ExchangeID,          //6
			*liquidityPool.LiquidityPoolTypeID, //7
			*liquidityPool.Token0ID,            //8
			*liquidityPool.Token1ID,            //9
			liquidityPool.Url,                  //10
			*liquidityPool.StartBlock,          //11
			*liquidityPool.LatestBlockSynced,   //12
			liquidityPool.CreatedTxnHash,       //13
			liquidityPool.IsActive,             //14
			liquidityPool.Description,          //15
			liquidityPool.CreatedBy,            //16
			&now,                               //17
			liquidityPool.CreatedBy,            //18
			&now,                               //19
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"liquidityPools"},
		[]string{
			"uuid",                   //1
			"name",                   //2
			"alternate_name",         //3
			"pair_address",           //4
			"chain_id",               //5
			"exchange_id",            //6
			"liquidity_pool_type_id", //7
			"token0_id",              //8
			"token1_id",              //9
			"url",                    //10
			"start_block",            //11
			"latest_block_synced",    //12
			"created_txn_hash",       //13
			"IsActive",               //14
			"description",            //15
			"created_by",             //16
			"created_at",             //17
			"updated_by",             //18
			"updated_at",             //19
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
	_, err := database.DbConn.ExecContext(ctx, `UPDATE liquidity_pool_assets SET 
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
	_, err := database.DbConn.ExecContext(ctx, `INSERT INTO liquidity_pool_assets 
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
		pgx.Identifier{"liquidityPool_chains"},
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