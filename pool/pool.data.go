package pool

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

func GetPool(dbConnPgx utils.PgxIface, poolID *int) (*Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		target_asset_id,
		strategy_id,
		account_id,
		uuid,
		name,
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
	WHERE id = $1`, *poolID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	pool, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Pool])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pool, nil
}

func RemovePool(dbConnPgx utils.PgxIface, poolID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM pools WHERE id = $1`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *poolID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetPools(dbConnPgx utils.PgxIface, ids []int) ([]Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		target_asset_id,
		strategy_id,
		account_id,
		uuid,
		name,
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
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	pools, err := pgx.CollectRows(results, pgx.RowToStructByName[Pool])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return pools, nil
}

func GetPoolsByStrategyID(dbConnPgx utils.PgxIface, strategyID *int) ([]Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		target_asset_id,
		strategy_id,
		account_id,
		uuid,
		name,
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
	pools, err := pgx.CollectRows(results, pgx.RowToStructByName[Pool])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return pools, nil
}

func GetPoolsByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		target_asset_id,
		strategy_id,
		account_id,
		uuid,
		name,
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
	pools, err := pgx.CollectRows(results, pgx.RowToStructByName[Pool])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return pools, nil
}

func GetStartAndEndDateDiffPools(dbConnPgx utils.PgxIface, diffInDate *int) ([]Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		target_asset_id,
		strategy_id,
		account_id,
		uuid,
		name,
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
	`, *diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	pools, err := pgx.CollectRows(results, pgx.RowToStructByName[Pool])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return pools, nil
}

func UpdatePool(dbConnPgx utils.PgxIface, pool *Pool) error {
	// if the pool id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if pool.ID == nil || *pool.ID == 0 {
		return errors.New("pool has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateMarketDataJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE pools SET 
		target_asset_id=$1, 
		strategy_id=$2, 
		account_id=$3, 
		uuid=$4, 
		name=$5, 
		alternate_name=$6, 
		start_date=$7, 
		end_date=$8, 
		description=$9, 
		chain_id=$10, 
		frequency_id=$11, 
		updated_by=$12, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$13`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		pool.TargetAssetID, //1
		pool.StrategyID,    //2
		pool.AccountID,     //3
		pool.UUID,          //4
		pool.Name,          //5
		pool.AlternateName, //6
		pool.StartDate,     //7
		pool.EndDate,       //8
		pool.Description,   //9
		pool.ChainID,       //10
		pool.FrequencyID,   //11
		pool.UpdatedBy,     //12
		pool.ID,            //13
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertPool(dbConnPgx utils.PgxIface, pool *Pool) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertMarketDataJob DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var insertID int
	var uuid string
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO pools 
	(
		target_asset_id,
		strategy_id,
		account_id,
		uuid,
		name,
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
			uuid_generate_v4(), 
			$4, 
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
		pool.TargetAssetID, //1
		pool.StrategyID,    //2
		pool.AccountID,     //3
		// &pool.UUID, //
		pool.Name,          //4
		pool.AlternateName, //5
		pool.StartDate,     //6
		pool.EndDate,       //7
		pool.Description,   //8
		pool.ChainID,       //9
		pool.FrequencyID,   //10
		pool.CreatedBy,     //11
	).Scan(&insertID, &uuid)
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
	return int(insertID), uuid, nil

}

func InsertPools(dbConnPgx utils.PgxIface, pools []Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range pools {
		pool := pools[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(pool.UUID)
		row := []interface{}{
			pool.TargetAssetID, //1
			pool.StrategyID,    //2
			pool.AccountID,     //3
			uuidString,         //4
			pool.Name,          //5
			pool.AlternateName, //6
			pool.StartDate,     //7
			pool.EndDate,       //8
			pool.Description,   //9
			pool.ChainID,       //10
			pool.FrequencyID,   //11
			pool.CreatedBy,     //12
			&now,               //13
			pool.CreatedBy,     //14
			&now,               //15
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"pools"},
		[]string{
			"target_asset_id", //1
			"strategy_id",     //2
			"account_id",      //3
			"uuid",            //4
			"name",            //5
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
	log.Println(fmt.Printf("InsertPools: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetPoolListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		target_asset_id,
		strategy_id,
		account_id,
		uuid,
		name,
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
	pools, err := pgx.CollectRows(results, pgx.RowToStructByName[Pool])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return pools, nil
}

func GetTotalPoolsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM pools
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
