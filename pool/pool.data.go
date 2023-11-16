package pool

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

func GetPool(poolID int) (*Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
		id,
		target_asset_id,
		strategy_id,
		account_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		description,
		chain_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM pools 
	WHERE id = $1`, poolID)

	pool := &Pool{}
	err := row.Scan(
		&pool.ID,
		&pool.TargetAssetID,
		&pool.StrategyID,
		&pool.AccountID,
		&pool.Name,
		&pool.UUID,
		&pool.AlternateName,
		&pool.StartDate,
		&pool.EndDate,
		&pool.Description,
		&pool.ChainID,
		&pool.FrequencyID,
		&pool.CreatedBy,
		&pool.CreatedAt,
		&pool.UpdatedBy,
		&pool.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return pool, nil
}

func GetTopTenStrategies() ([]Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,
	target_asset_id,
	strategy_id,
	account_id,
	name,
	uuid,
	alternate_name,
	start_date,
	end_date,
	description,
	chain_id ,
	frequency_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM pools 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	pools := make([]Pool, 0)
	for results.Next() {
		var pool Pool
		results.Scan(
			&pool.ID,
			&pool.TargetAssetID,
			&pool.StrategyID,
			&pool.AccountID,
			&pool.Name,
			&pool.UUID,
			&pool.AlternateName,
			&pool.StartDate,
			&pool.EndDate,
			&pool.Description,
			&pool.ChainID,
			&pool.FrequencyID,
			&pool.CreatedBy,
			&pool.CreatedAt,
			&pool.UpdatedBy,
			&pool.UpdatedAt,
		)

		pools = append(pools, pool)
	}
	return pools, nil
}

func RemovePool(poolID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM pools WHERE id = $1`, poolID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetPools(ids []int) ([]Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		target_asset_id,
		strategy_id,
		account_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		description,
		chain_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM pools`
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
	pools := make([]Pool, 0)
	for results.Next() {
		var pool Pool
		results.Scan(
			&pool.ID,
			&pool.TargetAssetID,
			&pool.StrategyID,
			&pool.AccountID,
			&pool.Name,
			&pool.UUID,
			&pool.AlternateName,
			&pool.StartDate,
			&pool.EndDate,
			&pool.Description,
			&pool.ChainID,
			&pool.FrequencyID,
			&pool.CreatedBy,
			&pool.CreatedAt,
			&pool.UpdatedBy,
			&pool.UpdatedAt,
		)

		pools = append(pools, pool)
	}
	return pools, nil
}

func GetPoolsByStrategyID(strategyID *int) ([]Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
		id,
		target_asset_id,
		strategy_id,
		account_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		description,
		chain_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM pools
	WHERE strategy_id = $1`, *strategyID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	pools := make([]Pool, 0)
	for results.Next() {
		var pool Pool
		results.Scan(
			&pool.ID,
			&pool.TargetAssetID,
			&pool.StrategyID,
			&pool.AccountID,
			&pool.Name,
			&pool.UUID,
			&pool.AlternateName,
			&pool.StartDate,
			&pool.EndDate,
			&pool.Description,
			&pool.ChainID,
			&pool.FrequencyID,
			&pool.CreatedBy,
			&pool.CreatedAt,
			&pool.UpdatedBy,
			&pool.UpdatedAt,
		)

		pools = append(pools, pool)
	}
	return pools, nil
}

func GetPoolsByUUIDs(UUIDList []string) ([]Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
		id,
		target_asset_id,
		strategy_id,
		account_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		description,
		chain_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM pools
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	pools := make([]Pool, 0)
	for results.Next() {
		var pool Pool
		results.Scan(
			&pool.ID,
			&pool.TargetAssetID,
			&pool.StrategyID,
			&pool.AccountID,
			&pool.Name,
			&pool.UUID,
			&pool.AlternateName,
			&pool.StartDate,
			&pool.EndDate,
			&pool.Description,
			&pool.ChainID,
			&pool.FrequencyID,
			&pool.CreatedBy,
			&pool.CreatedAt,
			&pool.UpdatedBy,
			&pool.UpdatedAt,
		)

		pools = append(pools, pool)
	}
	return pools, nil
}

func GetStartAndEndDateDiffPools(diffInDate int) ([]Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
		id,
		target_asset_id,
		strategy_id,
		account_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		description,
		chain_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM pools
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	pools := make([]Pool, 0)
	for results.Next() {
		var pool Pool
		results.Scan(
			&pool.ID,
			&pool.TargetAssetID,
			&pool.StrategyID,
			&pool.AccountID,
			&pool.Name,
			&pool.UUID,
			&pool.AlternateName,
			&pool.StartDate,
			&pool.EndDate,
			&pool.Description,
			&pool.ChainID,
			&pool.FrequencyID,
			&pool.CreatedBy,
			&pool.CreatedAt,
			&pool.UpdatedBy,
			&pool.UpdatedAt,
		)

		pools = append(pools, pool)
	}
	return pools, nil
}

func UpdatePool(pool Pool) error {
	// if the pool id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if pool.ID == nil || *pool.ID == 0 {
		return errors.New("pool has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE pools SET 
		target_asset_id=$1, 
		strategy_id=$2, 
		account_id=$3, 
		name=$4, 
		uuid=$5, 
		alternate_name=$6, 
		start_date=$7, 
		end_date=$8, 
		description=$9, 
		chain_id=$10, 
		frequency_id=$11, 
		updated_by=$12, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$13`,
		pool.TargetAssetID, //1
		pool.StrategyID,    //2
		pool.AccountID,     //3
		pool.Name,          //4
		pool.UUID,          //5
		pool.AlternateName, //6
		pool.StartDate,     //7
		pool.EndDate,       //8
		pool.Description,   //9
		pool.ChainID,       //10
		pool.FrequencyID,   //11
		pool.UpdatedBy,     //12
		pool.ID)            //13
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertPool(pool Pool) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	// layoutPostgres := utils.LayoutPostgres
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO pools 
	(
		target_asset_id,
		strategy_id,
		account_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		description,
		chain_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
		) VALUES (
			$1,
			$2, 
			$3, 
			$4, 
			uuid_generate_v4(), 
			$5, 
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			current_timestamp at time zone 'UTC',
			$11,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		&pool.TargetAssetID, //1
		&pool.StrategyID,    //2
		&pool.AccountID,     //3
		&pool.Name,          //4
		// &pool.UUID, //
		&pool.AlternateName, //5
		&pool.StartDate,     //6
		&pool.EndDate,       //7
		&pool.Description,   //8
		&pool.ChainID,       //9
		&pool.FrequencyID,   //10
		&pool.CreatedBy,     //11
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

func InsertPools(pools []Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range pools {
		pool := pools[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(pool.UUID)
		row := []interface{}{
			*pool.TargetAssetID, //1
			*pool.StrategyID,    //2
			*pool.AccountID,     //3
			pool.Name,           //4
			uuidString,          //5
			pool.AlternateName,  //6
			&pool.StartDate,     //7
			&pool.EndDate,       //8
			pool.Description,    //9
			*pool.ChainID,       //10
			*pool.FrequencyID,   //11
			pool.CreatedBy,      //12
			&now,                //13
			pool.CreatedBy,      //14
			&now,                //15
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"pools"},
		[]string{
			"target_asset_id", //1
			"strategy_id",     //2
			"account_id",      //3
			"name",            //4
			"uuid",            //5
			"alternate_name",  //6
			"start_date",      //7
			"end_date",        //8
			"description",     //9
			"chain_id",        //10
			"frequency_id",    //11
			"created_by",      //12
			"created_at",      //13
			"updated_by",      //14
			"updated_at",      //15
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
