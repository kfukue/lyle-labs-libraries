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
		token0_id=$6,
		token1_id=$7,
		url=$8,
		start_block=$9,
		latest_block_synced=$10,
		created_txn_hash=$11,
		IsActive=$12,
		description=$13,
		updated_by=$14, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$15`,
		liquidityPool.Name,              //1
		liquidityPool.AlternateName,     //2
		liquidityPool.PairAddress,       //3
		liquidityPool.ChainID,           //4
		liquidityPool.ExchangeID,        //5
		liquidityPool.Token0ID,          //6
		liquidityPool.Token1ID,          //7
		liquidityPool.Url,               //8
		liquidityPool.StartBlock,        //9
		liquidityPool.LatestBlockSynced, //10
		liquidityPool.CreatedTxnHash,    //11
		liquidityPool.IsActive,          //12
		liquidityPool.Description,       //13
		liquidityPool.UpdatedBy,         //14
		liquidityPool.ID)                //15
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
			current_timestamp at time zone 'UTC',
			$14,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		liquidityPool.Name,               //1
		liquidityPool.AlternateName,      //2
		&liquidityPool.PairAddress,       //3
		&liquidityPool.ChainID,           //4
		&liquidityPool.ExchangeID,        //5
		&liquidityPool.Token0ID,          //6
		&liquidityPool.Token1ID,          //7
		&liquidityPool.Url,               //8
		&liquidityPool.StartBlock,        //9
		&liquidityPool.LatestBlockSynced, //10
		&liquidityPool.CreatedTxnHash,    //11
		&liquidityPool.IsActive,          //12
		liquidityPool.Description,        //13
		liquidityPool.CreatedBy,          //14
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
			uuidString,                       //1
			liquidityPool.Name,               //2
			liquidityPool.AlternateName,      //3
			liquidityPool.PairAddress,        //4
			*liquidityPool.ChainID,           //5
			*liquidityPool.ExchangeID,        //6
			*liquidityPool.Token0ID,          //7
			*liquidityPool.Token1ID,          //8
			liquidityPool.Url,                //9
			*liquidityPool.StartBlock,        //10
			*liquidityPool.LatestBlockSynced, //11
			liquidityPool.CreatedTxnHash,     //12
			liquidityPool.IsActive,           //13
			liquidityPool.Description,        //14
			liquidityPool.CreatedBy,          //15
			&now,                             //16
			liquidityPool.CreatedBy,          //17
			&now,                             //18
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"liquidityPools"},
		[]string{
			"uuid",                //1
			"name",                //2
			"alternate_name",      //3
			"pair_address",        //4
			"chain_id",            //5
			"exchange_id",         //6
			"token0_id",           //7
			"token1_id",           //8
			"url",                 //9
			"start_block",         //10
			"latest_block_synced", //11
			"created_txn_hash",    //12
			"IsActive",            //13
			"description",         //14
			"created_by",          //15
			"created_at",          //16
			"updated_by",          //17
			"updated_at",          //18
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
	if liquidityPoolChain.LiquidityPoolID == nil || *liquidityPoolChain.LiquidityPoolID == 0 || liquidityPoolChain.ChainID == nil || *liquidityPoolChain.ChainID == 0 {
		return errors.New("liquidityPoolChain has invalid IDs")
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
		liquidityPoolChain.LiquidityPoolID, //1
		liquidityPoolChain.AssetID,         //2
		liquidityPoolChain.TokenNumber,     //3
		liquidityPoolChain.Name,            //4
		liquidityPoolChain.AlternateName,   //5
		liquidityPoolChain.Description,     //6
		liquidityPoolChain.UpdatedBy,       //7
		liquidityPoolChain.UUID)            //8
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertLiquidityPoolChain(liquidityPoolChain LiquidityPoolChain) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	// layoutPostgres := utils.LayoutPostgres
	_, err := database.DbConn.ExecContext(ctx, `INSERT INTO liquidityPool_chains 
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
		liquidityPoolChain.LiquidityPoolID, //1
		liquidityPoolChain.AssetID,         //2
		liquidityPoolChain.TokenNumber,     //3
		liquidityPoolChain.Name,            //4
		liquidityPoolChain.AlternateName,   //5
		liquidityPoolChain.Description,     //6
		liquidityPoolChain.CreatedBy,       //7
	)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertLiquidityPoolChains(liquidityPoolChains []LiquidityPoolChain) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range liquidityPoolChains {
		liquidityPoolChain := liquidityPoolChains[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(liquidityPoolChain.UUID)
		row := []interface{}{
			uuidString,                         //1
			liquidityPoolChain.LiquidityPoolID, //2
			liquidityPoolChain.AssetID,         //3
			liquidityPoolChain.TokenNumber,     //4
			liquidityPoolChain.Name,            //5
			liquidityPoolChain.AlternateName,   //6
			liquidityPoolChain.Description,     //7
			liquidityPoolChain.CreatedBy,       //8
			&now,                               //9
			liquidityPoolChain.CreatedBy,       //10
			&now,                               //11
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
