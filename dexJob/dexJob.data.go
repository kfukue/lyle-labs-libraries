package dexjob

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

func GetDexTxnJob(dexTxnID int) (*DexTxnJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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
	WHERE dex_swap_id = $1
	`, dexTxnID)

	dexTxnJob := &DexTxnJob{}
	err := row.Scan(
		&dexTxnJob.ID,
		&dexTxnJob.JobID,
		&dexTxnJob.UUID,
		&dexTxnJob.Name,
		&dexTxnJob.AlternateName,
		&dexTxnJob.StartDate,
		&dexTxnJob.EndDate,
		&dexTxnJob.Description,
		&dexTxnJob.StatusID,
		&dexTxnJob.ChainID,
		&dexTxnJob.ExchangeID,
		&dexTxnJob.TransactionHashes,
		&dexTxnJob.CreatedBy,
		&dexTxnJob.CreatedAt,
		&dexTxnJob.UpdatedBy,
		&dexTxnJob.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return dexTxnJob, nil
}

func GetDexTxnJobByJobId(jobID int) ([]DexTxnJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	job_id = $1	`, jobID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	dexTxnJobs := make([]DexTxnJob, 0)
	for results.Next() {
		var dexTxnJob DexTxnJob
		results.Scan(
			&dexTxnJob.ID,
			&dexTxnJob.JobID,
			&dexTxnJob.UUID,
			&dexTxnJob.Name,
			&dexTxnJob.AlternateName,
			&dexTxnJob.StartDate,
			&dexTxnJob.EndDate,
			&dexTxnJob.Description,
			&dexTxnJob.StatusID,
			&dexTxnJob.ChainID,
			&dexTxnJob.ExchangeID,
			&dexTxnJob.TransactionHashes,
			&dexTxnJob.CreatedBy,
			&dexTxnJob.CreatedAt,
			&dexTxnJob.UpdatedBy,
			&dexTxnJob.UpdatedAt,
		)
		dexTxnJobs = append(dexTxnJobs, dexTxnJob)
	}
	return dexTxnJobs, nil
}

func GetDexTxnJobList() ([]DexTxnJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	dexTxnJobs := make([]DexTxnJob, 0)
	for results.Next() {
		var dexTxnJob DexTxnJob
		results.Scan(
			&dexTxnJob.ID,
			&dexTxnJob.JobID,
			&dexTxnJob.UUID,
			&dexTxnJob.Name,
			&dexTxnJob.AlternateName,
			&dexTxnJob.StartDate,
			&dexTxnJob.EndDate,
			&dexTxnJob.Description,
			&dexTxnJob.StatusID,
			&dexTxnJob.ChainID,
			&dexTxnJob.ExchangeID,
			&dexTxnJob.TransactionHashes,
			&dexTxnJob.CreatedBy,
			&dexTxnJob.CreatedAt,
			&dexTxnJob.UpdatedBy,
			&dexTxnJob.UpdatedAt,
		)

		dexTxnJobs = append(dexTxnJobs, dexTxnJob)
	}
	return dexTxnJobs, nil
}

func RemoveDexTxnJob(dexTxnID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM dex_txn_jobs WHERE 
	id = $1`, dexTxnID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateDexTxnJob(dexTxnJob DexTxnJob) error {
	// if the dexTxnJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if dexTxnJob.ID == nil || *dexTxnJob.ID == 0 {
		return errors.New("dexTxnJob has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE dex_txn_jobs SET 

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
		WHERE id=$11 `,
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
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertDexTxnJob(dexTxnJob DexTxnJob) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var ID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO dex_txn_jobs  
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
		log.Println(err.Error())
		return 0, err
	}
	return int(ID), nil
}

func InsertDexTxnJobListManual(dexTxnJobList []DexTxnJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range dexTxnJobList {
		dexTxnJob := dexTxnJobList[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(dexTxnJob.UUID)
		row := []interface{}{

			*dexTxnJob.ID,                         //1
			*dexTxnJob.JobID,                      //2
			uuidString,                            //3
			dexTxnJob.Name,                        //4
			dexTxnJob.AlternateName,               //5
			&dexTxnJob.StartDate,                  //6
			&dexTxnJob.EndDate,                    //7
			dexTxnJob.Description,                 //8
			*dexTxnJob.StatusID,                   //9
			*dexTxnJob.ChainID,                    //10
			*dexTxnJob.ExchangeID,                 //11
			pq.Array(dexTxnJob.TransactionHashes), //12
			dexTxnJob.CreatedBy,                   //13
			&dexTxnJob.CreatedAt,                  //14
			dexTxnJob.CreatedBy,                   //15
			&now,                                  //16
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"dex_txn_jobs"},
		[]string{
			"id",                 //1
			"job_id",             //2
			"uuid",               //3
			"name",               //4
			"alternate_name",     //5
			"start_date",         //6
			"end_date",           //7
			"description",        //8
			"status_id",          //9
			"chain_id",           //10
			"exchange_id",        //11
			"transaction_hashes", //12
			"created_by",         //13
			"created_at",         //14
			"updated_by",         //15
			"updated_at",         //16
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

func InsertDexTxnJobList(dexTxnJobList []DexTxnJob) error {
	txn, err := database.DbConn.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	layoutPostgres := utils.LayoutPostgres
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn(
		"dex_txn_jobs",
		"id",                 //1
		"job_id",             //2
		"uuid",               //3
		"name",               //4
		"alternate_name",     //5
		"start_date",         //6
		"end_date",           //7
		"description",        //8
		"status_id",          //9
		"chain_id",           //10
		"exchange_id",        //11
		"transaction_hashes", //12
		"created_by",         //13
		"created_at",         //14
		"updated_by",         //15
		"updated_at",         //16

	))
	if err != nil {
		log.Fatal(err)
	}
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	for _, dexTxnJob := range dexTxnJobList {
		_, err = stmt.Exec(
			*dexTxnJob.ID,                         //1
			*dexTxnJob.JobID,                      //2
			dexTxnJob.UUID,                        //3
			dexTxnJob.Name,                        //4
			dexTxnJob.AlternateName,               //5
			dexTxnJob.StartDate,                   //6
			dexTxnJob.EndDate,                     //7
			dexTxnJob.Description,                 //8
			*dexTxnJob.StatusID,                   //9
			*dexTxnJob.ChainID,                    //10
			*dexTxnJob.ExchangeID,                 //11
			pq.Array(dexTxnJob.TransactionHashes), //12
			dexTxnJob.CreatedBy,                   //13
			dexTxnJob.CreatedAt,                   //14
			dexTxnJob.CreatedBy,                   //15
			now.Format(layoutPostgres),            //16
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
