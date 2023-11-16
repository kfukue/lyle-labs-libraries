package gethlylejobs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/lib/pq"
)

func GetGethProcessVlogJob(gethProcessJobID int) (*GethProcessVlogJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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
	`, gethProcessJobID)

	gethProcessJob := &GethProcessVlogJob{}
	err := row.Scan(
		&gethProcessJob.ID,
		&gethProcessJob.GethProcessJobID,
		&gethProcessJob.UUID,
		&gethProcessJob.Name,
		&gethProcessJob.AlternateName,
		&gethProcessJob.StartDate,
		&gethProcessJob.EndDate,
		&gethProcessJob.Description,
		&gethProcessJob.StatusID,
		&gethProcessJob.JobCategoryID,
		&gethProcessJob.AssetID,
		&gethProcessJob.ChainID,
		&gethProcessJob.TxnHash,
		&gethProcessJob.AddressID,
		&gethProcessJob.BlockNumber,
		&gethProcessJob.IndexNumber,
		&gethProcessJob.TopicsStr,
		&gethProcessJob.CreatedBy,
		&gethProcessJob.CreatedAt,
		&gethProcessJob.UpdatedBy,
		&gethProcessJob.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethProcessJob, nil
}

func GetGethProcessVlogJobList() ([]GethProcessVlogJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	FROM geth_process_vlog_jobs`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethProcessJobs := make([]GethProcessVlogJob, 0)
	for results.Next() {
		var gethProcessJob GethProcessVlogJob
		results.Scan(
			&gethProcessJob.ID,
			&gethProcessJob.GethProcessJobID,
			&gethProcessJob.UUID,
			&gethProcessJob.Name,
			&gethProcessJob.AlternateName,
			&gethProcessJob.StartDate,
			&gethProcessJob.EndDate,
			&gethProcessJob.Description,
			&gethProcessJob.StatusID,
			&gethProcessJob.JobCategoryID,
			&gethProcessJob.AssetID,
			&gethProcessJob.ChainID,
			&gethProcessJob.TxnHash,
			&gethProcessJob.AddressID,
			&gethProcessJob.BlockNumber,
			&gethProcessJob.IndexNumber,
			&gethProcessJob.TopicsStr,
			&gethProcessJob.CreatedBy,
			&gethProcessJob.CreatedAt,
			&gethProcessJob.UpdatedBy,
			&gethProcessJob.UpdatedAt,
		)

		gethProcessJobs = append(gethProcessJobs, gethProcessJob)
	}
	return gethProcessJobs, nil
}

func RemoveGethProcessVlogJob(gethProcessJobID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_process_vlog_jobs WHERE 
	id = $1`, *gethProcessJobID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateGethProcessVlogJob(gethProcessJob *GethProcessVlogJob) error {
	// if the gethProcessJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethProcessJob.ID == nil || *gethProcessJob.ID == 0 {
		return errors.New("gethProcessJob has invalid ID")
	}
	_, err := database.DbConnPgx.Exec(ctx, `UPDATE geth_process_vlog_jobs SET 
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
		WHERE id=$17 `,
		gethProcessJob.GethProcessJobID,    //1
		gethProcessJob.Name,                //2
		gethProcessJob.AlternateName,       //3
		gethProcessJob.StartDate,           //4
		gethProcessJob.EndDate,             //5
		gethProcessJob.Description,         //6
		gethProcessJob.StatusID,            //7
		gethProcessJob.JobCategoryID,       //8
		gethProcessJob.AssetID,             //9
		gethProcessJob.ChainID,             //10
		gethProcessJob.TxnHash,             //11
		gethProcessJob.AddressID,           //12
		gethProcessJob.BlockNumber,         //13
		gethProcessJob.IndexNumber,         //14
		pq.Array(gethProcessJob.TopicsStr), //15
		gethProcessJob.UpdatedBy,           //16
		gethProcessJob.ID,                  //17
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertGethProcessVlogJob(gethProcessJob GethProcessVlogJob) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var ID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_process_vlog_jobs  
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
		&gethProcessJob.GethProcessJobID, //1
		gethProcessJob.Name,              //2
		gethProcessJob.AlternateName,     //3
		gethProcessJob.StartDate,         //4
		gethProcessJob.EndDate,           //5
		gethProcessJob.Description,       //6
		gethProcessJob.StatusID,          //7
		gethProcessJob.JobCategoryID,     //8
		gethProcessJob.AssetID,           //9
		gethProcessJob.ChainID,           //10
		gethProcessJob.TxnHash,           //11
		gethProcessJob.AddressID,         //12
		gethProcessJob.BlockNumber,       //13
		&gethProcessJob.IndexNumber,      //14
		gethProcessJob.TopicsStr,         //15
		gethProcessJob.CreatedBy,         //16
	).Scan(&ID)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(ID), nil
}

func InsertGethProcessVlogJobList(gethProcessJobList []GethProcessVlogJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethProcessJobList {
		gethProcessJob := gethProcessJobList[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(gethProcessJob.UUID)
		row := []interface{}{

			*gethProcessJob.ID,                 //1
			*gethProcessJob.GethProcessJobID,   //2
			uuidString,                         //3
			gethProcessJob.Name,                //4
			gethProcessJob.AlternateName,       //5
			&gethProcessJob.StartDate,          //6
			&gethProcessJob.EndDate,            //7
			gethProcessJob.Description,         //8
			*gethProcessJob.StatusID,           //9
			*gethProcessJob.JobCategoryID,      //10
			gethProcessJob.AssetID,             //11
			gethProcessJob.ChainID,             //12
			gethProcessJob.TxnHash,             //13
			gethProcessJob.AddressID,           //14
			gethProcessJob.BlockNumber,         //15
			gethProcessJob.IndexNumber,         //16
			pq.Array(gethProcessJob.TopicsStr), //17
			gethProcessJob.CreatedBy,           //18
			&gethProcessJob.CreatedAt,          //19
			gethProcessJob.CreatedBy,           //20
			&now,                               //21
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_process_vlog_jobs"},
		[]string{
			"id",                  //1
			"geth_process_job_id", //2
			"uuid",                //3
			"name",                //4
			"alternate_name",      //5
			"start_date",          //6
			"end_date",            //7
			"description",         //8
			"status_id",           //9
			"job_category_id",     //10
			"asset_id",            //11
			"chain_id",            //12
			"txn_hash",            //13
			"address_id",          //14
			"block_number",        //15
			"index_number",        //16
			"topics_str",          //17
			"created_by",          //18
			"created_at",          //19
			"updated_by",          //20
			"updated_at",          //21
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
