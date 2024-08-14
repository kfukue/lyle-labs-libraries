package dexjob

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

func GetDexTxnJob(dbConnPgx utils.PgxIface, dexTxnID *int) (*DexTxnJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	job_id,
	uuid, 
	name,
	alternate_name,
	start_date,
	end_date,
	description,
	status_id,
	chain_id,
	exchange_id,
	transaction_hashes,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM dex_txn_jobs 
	WHERE id = $1
	`, *dexTxnID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	defer row.Close()
	dexTxnJob, err := pgx.CollectOneRow(row, pgx.RowToStructByName[DexTxnJob])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &dexTxnJob, nil
}

func GetDexTxnJobByJobId(dbConnPgx utils.PgxIface, jobID *int) ([]DexTxnJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	job_id,
	uuid, 
	name,
	alternate_name,
	start_date,
	end_date,
	description,
	status_id,
	chain_id,
	exchange_id,
	transaction_hashes,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM dex_txn_jobs
	job_id = $1`, *jobID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	dexTxnJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[DexTxnJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return dexTxnJobs, nil
}

func GetDexTxnJobList(dbConnPgx utils.PgxIface) ([]DexTxnJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	job_id,
	uuid, 
	name,
	alternate_name,
	start_date,
	end_date,
	description,
	status_id,
	chain_id,
	exchange_id,
	transaction_hashes,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM dex_txn_jobs`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	dexTxnJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[DexTxnJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return dexTxnJobs, nil
}

func RemoveDexTxnJob(dbConnPgx utils.PgxIface, dexTxnID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveDexTxnJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM dex_txn_jobs WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *dexTxnID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func UpdateDexTxnJob(dbConnPgx utils.PgxIface, dexTxnJob *DexTxnJob) error {
	// if the dexTxnJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if dexTxnJob.ID == nil || *dexTxnJob.ID == 0 {
		return errors.New("dexTxnJob has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateDexTxnJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE dex_txn_jobs SET 
		name=$1,
		alternate_name=$2,
		start_date=$3,
		end_date=$4,
		description=$5,
		status_id=$6,
		chain_id=$7,
		exchange_id=$8,
		transaction_hashes=$9,
		updated_by=$10, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$11 `
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		dexTxnJob.Name,                        //1
		dexTxnJob.AlternateName,               //2
		dexTxnJob.StartDate,                   //3
		dexTxnJob.EndDate,                     //4
		dexTxnJob.Description,                 //5
		dexTxnJob.StatusID,                    //6
		dexTxnJob.ChainID,                     //7
		dexTxnJob.ExchangeID,                  //8
		pq.Array(dexTxnJob.TransactionHashes), //9
		dexTxnJob.UpdatedBy,                   //10
		dexTxnJob.ID,                          //11
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertDexTxnJob(dbConnPgx utils.PgxIface, dexTxnJob *DexTxnJob) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertDexTxnJob DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var ID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO dex_txn_jobs  
	(
		job_id,
		uuid, 
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		status_id,
		chain_id,
		exchange_id,
		transaction_hashes,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES ($1,
			uuid_generate_v4(),
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
			current_timestamp at time zone 'UTC',
			$11,
			current_timestamp at time zone 'UTC'
		)
		RETURNING dex_swap_id, job_id`,
		dexTxnJob.JobID,                       //1
		dexTxnJob.Name,                        //2
		dexTxnJob.AlternateName,               //3
		dexTxnJob.StartDate,                   //4
		dexTxnJob.EndDate,                     //5
		dexTxnJob.Description,                 //6
		dexTxnJob.StatusID,                    //7
		dexTxnJob.ChainID,                     //8
		dexTxnJob.ExchangeID,                  //9
		pq.Array(dexTxnJob.TransactionHashes), //10
		dexTxnJob.CreatedBy,                   //11
	).Scan(&ID)
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
	return int(ID), nil
}

func InsertDexTxnJobList(dbConnPgx utils.PgxIface, dexTxnJobList []DexTxnJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range dexTxnJobList {
		dexTxnJob := dexTxnJobList[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(dexTxnJob.UUID)
		row := []interface{}{
			dexTxnJob.JobID,                       //1
			uuidString,                            //2
			dexTxnJob.Name,                        //3
			dexTxnJob.AlternateName,               //4
			dexTxnJob.StartDate,                   //5
			dexTxnJob.EndDate,                     //6
			dexTxnJob.Description,                 //7
			dexTxnJob.StatusID,                    //8
			dexTxnJob.ChainID,                     //9
			dexTxnJob.ExchangeID,                  //10
			pq.Array(dexTxnJob.TransactionHashes), //11
			dexTxnJob.CreatedBy,                   //12
			dexTxnJob.CreatedAt,                   //13
			dexTxnJob.CreatedBy,                   //14
			now,                                   //15
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"dex_txn_jobs"},
		[]string{
			"job_id",             //1
			"uuid",               //2
			"name",               //3
			"alternate_name",     //4
			"start_date",         //5
			"end_date",           //6
			"description",        //7
			"status_id",          //8
			"chain_id",           //9
			"exchange_id",        //10
			"transaction_hashes", //11
			"created_by",         //12
			"created_at",         //13
			"updated_by",         //14
			"updated_at",         //15
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertDexTxnJobList : copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetDexTxnJobListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]DexTxnJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
	id,  
	job_id,
	uuid, 
	name,
	alternate_name,
	start_date,
	end_date,
	description,
	status_id,
	chain_id,
	exchange_id,
	transaction_hashes,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM dex_txn_jobs
	`
	if len(_filters) > 0 {
		sql += "WHERE "
		for i, filter := range _filters {
			sql += filter
			if i < len(_filters)-1 {
				sql += " AND "
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
	dexTxnJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[DexTxnJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return dexTxnJobs, nil
}

func GetTotalDexTxnJobCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM dex_txn_jobs`)
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
