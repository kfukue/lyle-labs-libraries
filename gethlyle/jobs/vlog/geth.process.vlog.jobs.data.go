package gethlylevlogjobs

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

func GetGethProcessVlogJob(dbConnPgx utils.PgxIface, gethProcessVlogJobID *int) (*GethProcessVlogJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	geth_process_job_id,
	uuid, 
	name,
	alternate_name,
	start_date,
	end_date,
	description,
	status_id,
	job_category_id,
	asset_id,
	chain_id,
	txn_hash,
	address_id,
	block_number,
	index_number,
	topics_str,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_process_vlog_jobs 
	WHERE id = $1
	`, *gethProcessVlogJobID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx

	gethProcessVlogJob, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethProcessVlogJob])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethProcessVlogJob, nil
}

func GetGethProcessVlogJobList(dbConnPgx utils.PgxIface) ([]GethProcessVlogJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,  
	geth_process_job_id,
	uuid, 
	name,
	alternate_name,
	start_date,
	end_date,
	description,
	status_id,
	job_category_id,
	asset_id,
	chain_id,
	txn_hash,
	address_id,
	block_number,
	index_number,
	topics_str,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_process_vlog_jobs`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethProcessVlogJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[GethProcessVlogJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethProcessVlogJobs, nil
}

func RemoveGethProcessVlogJob(dbConnPgx utils.PgxIface, gethProcessVlogJobID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethProcessVlogJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_process_vlog_jobs WHERE id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *gethProcessVlogJobID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func UpdateGethProcessVlogJob(dbConnPgx utils.PgxIface, gethProcessVlogJob *GethProcessVlogJob) error {
	// if the gethProcessVlogJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethProcessVlogJob.ID == nil || *gethProcessVlogJob.ID == 0 {
		return errors.New("gethProcessVlogJob has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethProcessVlogJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_process_vlog_jobs SET 
		geth_process_job_id = $1,
		name=$2,
		alternate_name=$3,
		start_date=$4,
		end_date=$5,
		description=$6,
		status_id=$7,
		job_category_id=$8,
		asset_id=$9,
		chain_id=$10,
		txn_hash=$11,
		address_id=$12,
		block_number=$13,
		index_number= $14, 
		topics_str=$15,
		updated_by=$16, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$17 `

	if _, err := dbConnPgx.Exec(ctx, sql,
		gethProcessVlogJob.GethProcessJobID,         //1
		gethProcessVlogJob.Name,                     //2
		gethProcessVlogJob.AlternateName,            //3
		gethProcessVlogJob.StartDate,                //4
		gethProcessVlogJob.EndDate,                  //5
		gethProcessVlogJob.Description,              //6
		gethProcessVlogJob.StatusID,                 //7
		gethProcessVlogJob.JobCategoryID,            //8
		gethProcessVlogJob.AssetID,                  //9
		gethProcessVlogJob.ChainID,                  //10
		gethProcessVlogJob.TxnHash,                  //11
		gethProcessVlogJob.AddressID,                //12
		gethProcessVlogJob.BlockNumber,              //13
		gethProcessVlogJob.IndexNumber,              //14
		pq.Array(gethProcessVlogJob.TopicsStrArray), //15
		gethProcessVlogJob.UpdatedBy,                //16
		gethProcessVlogJob.ID,                       //17
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethProcessVlogJob(dbConnPgx utils.PgxIface, gethProcessVlogJob *GethProcessVlogJob) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethProcessVlogJob DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var ID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_process_vlog_jobs  
	(
		geth_process_job_id,
		uuid, 
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		status_id,
		job_category_id,
		asset_id,
		chain_id,
		txn_hash,
		address_id,
		block_number,
		index_number,
		topics_str,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES (
			$1,
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
			$12,
			$13,
			$14,
			$15,
			$16,
			current_timestamp at time zone 'UTC',
			$16,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		gethProcessVlogJob.GethProcessJobID, //1
		gethProcessVlogJob.Name,             //2
		gethProcessVlogJob.AlternateName,    //3
		gethProcessVlogJob.StartDate,        //4
		gethProcessVlogJob.EndDate,          //5
		gethProcessVlogJob.Description,      //6
		gethProcessVlogJob.StatusID,         //7
		gethProcessVlogJob.JobCategoryID,    //8
		gethProcessVlogJob.AssetID,          //9
		gethProcessVlogJob.ChainID,          //10
		gethProcessVlogJob.TxnHash,          //11
		gethProcessVlogJob.AddressID,        //12
		gethProcessVlogJob.BlockNumber,      //13
		gethProcessVlogJob.IndexNumber,      //14
		gethProcessVlogJob.TopicsStrArray,   //15
		gethProcessVlogJob.CreatedBy,        //16
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

func InsertGethProcessVlogJobList(dbConnPgx utils.PgxIface, gethProcessVlogJobList []GethProcessVlogJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethProcessVlogJobList {
		gethProcessVlogJob := gethProcessVlogJobList[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(gethProcessVlogJob.UUID)
		row := []interface{}{

			gethProcessVlogJob.GethProcessJobID,         //1
			uuidString,                                  //2
			gethProcessVlogJob.Name,                     //3
			gethProcessVlogJob.AlternateName,            //4
			gethProcessVlogJob.StartDate,                //5
			gethProcessVlogJob.EndDate,                  //6
			gethProcessVlogJob.Description,              //7
			gethProcessVlogJob.StatusID,                 //8
			gethProcessVlogJob.JobCategoryID,            //9
			gethProcessVlogJob.AssetID,                  //10
			gethProcessVlogJob.ChainID,                  //11
			gethProcessVlogJob.TxnHash,                  //12
			gethProcessVlogJob.AddressID,                //13
			gethProcessVlogJob.BlockNumber,              //14
			gethProcessVlogJob.IndexNumber,              //15
			pq.Array(gethProcessVlogJob.TopicsStrArray), //16
			gethProcessVlogJob.CreatedBy,                //17
			&gethProcessVlogJob.CreatedAt,               //18
			gethProcessVlogJob.CreatedBy,                //19
			&now,                                        //20
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_process_vlog_jobs"},
		[]string{
			"geth_process_job_id", //1
			"uuid",                //2
			"name",                //3
			"alternate_name",      //4
			"start_date",          //5
			"end_date",            //6
			"description",         //7
			"status_id",           //8
			"job_category_id",     //9
			"asset_id",            //10
			"chain_id",            //11
			"txn_hash",            //12
			"address_id",          //13
			"block_number",        //14
			"index_number",        //15
			"topics_str",          //16
			"created_by",          //17
			"created_at",          //18
			"updated_by",          //19
			"updated_at",          //20
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethProcessVlogJobList: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetGethProcessVlogJobListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethProcessVlogJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,  
		geth_process_job_id,
		uuid, 
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		status_id,
		job_category_id,
		asset_id,
		chain_id,
		txn_hash,
		address_id,
		block_number,
		index_number,
		topics_str,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		FROM geth_process_vlog_jobs 
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

	gethProcessVlogJobList, err := pgx.CollectRows(results, pgx.RowToStructByName[GethProcessVlogJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethProcessVlogJobList, nil
}

func GetTotalGethProcessVlogJobCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM geth_process_vlog_jobs`)
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
