package gethlylejobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
)

func GetGethProcessJob(gethProcessJobID int) (*GethProcessJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	id,  
	uuid, 
	name,
	alternate_name,
	start_date,
	end_date,
	description,
	status_id,
	job_category_id,
	import_type_id,
	chain_id,
	start_block_number,
	end_block_number,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	asset_id
	FROM geth_process_jobs 
	WHERE id = $1
	`, gethProcessJobID)

	gethProcessJob := &GethProcessJob{}
	err := row.Scan(
		&gethProcessJob.ID,
		&gethProcessJob.UUID,
		&gethProcessJob.Name,
		&gethProcessJob.AlternateName,
		&gethProcessJob.StartDate,
		&gethProcessJob.EndDate,
		&gethProcessJob.Description,
		&gethProcessJob.StatusID,
		&gethProcessJob.JobCategoryID,
		&gethProcessJob.ImportTypeID,
		&gethProcessJob.ChainID,
		&gethProcessJob.StartBlockNumber,
		&gethProcessJob.EndBlockNumber,
		&gethProcessJob.CreatedBy,
		&gethProcessJob.CreatedAt,
		&gethProcessJob.UpdatedBy,
		&gethProcessJob.UpdatedAt,
		&gethProcessJob.AssetID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethProcessJob, nil
}

func GetGethProcessJobList() ([]GethProcessJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,  
	uuid, 
	name,
	alternate_name,
	start_date,
	end_date,
	description,
	status_id,
	job_category_id,
	import_type_id,
	chain_id,
	start_block_number,
	end_block_number,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	asset_id
	FROM geth_process_jobs`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethProcessJobs := make([]GethProcessJob, 0)
	for results.Next() {
		var gethProcessJob GethProcessJob
		results.Scan(
			&gethProcessJob.ID,
			&gethProcessJob.UUID,
			&gethProcessJob.Name,
			&gethProcessJob.AlternateName,
			&gethProcessJob.StartDate,
			&gethProcessJob.EndDate,
			&gethProcessJob.Description,
			&gethProcessJob.StatusID,
			&gethProcessJob.JobCategoryID,
			&gethProcessJob.ImportTypeID,
			&gethProcessJob.ChainID,
			&gethProcessJob.StartBlockNumber,
			&gethProcessJob.EndBlockNumber,
			&gethProcessJob.CreatedBy,
			&gethProcessJob.CreatedAt,
			&gethProcessJob.UpdatedBy,
			&gethProcessJob.UpdatedAt,
			&gethProcessJob.AssetID,
		)

		gethProcessJobs = append(gethProcessJobs, gethProcessJob)
	}
	return gethProcessJobs, nil
}

func RemoveGethProcessJob(gethProcessJobID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_process_jobs WHERE 
	id = $1`, *gethProcessJobID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateGethProcessJob(gethProcessJob GethProcessJob) error {
	// if the gethProcessJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethProcessJob.ID == nil || *gethProcessJob.ID == 0 {
		return errors.New("gethProcessJob has invalid ID")
	}
	_, err := database.DbConnPgx.Exec(ctx, `UPDATE geth_process_jobs SET 

		name=$1,
		alternate_name=$2,
		start_date=$3,
		end_date=$4,
		description=$5,
		status_id=$6,
		job_category_id=$7,
		import_type_id=$8,
		chain_id=$9,
		start_block_number=$10,
		end_block_number=$11,
		updated_by=$12, 
		updated_at=current_timestamp at time zone 'UTC',
		asset_id =$13
		WHERE id=$14 `,
		gethProcessJob.Name,             //1
		gethProcessJob.AlternateName,    //2
		gethProcessJob.StartDate,        //3
		gethProcessJob.EndDate,          //4
		gethProcessJob.Description,      //5
		gethProcessJob.StatusID,         //6
		gethProcessJob.JobCategoryID,    //7
		gethProcessJob.ImportTypeID,     //8
		gethProcessJob.ChainID,          //9
		gethProcessJob.StartBlockNumber, //10
		gethProcessJob.EndBlockNumber,   //11
		gethProcessJob.UpdatedBy,        //12
		gethProcessJob.AssetID,          //13
		gethProcessJob.ID,               //14
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertGethProcessJob(gethProcessJob GethProcessJob) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var ID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_process_jobs  
	(
		uuid, 
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		status_id,
		job_category_id,
		import_type_id,
		chain_id,
		start_block_number,
		end_block_number,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
		asset_id
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
			current_timestamp at time zone 'UTC',
			$12,
			current_timestamp at time zone 'UTC',
			$13
		)
		RETURNING id`,
		gethProcessJob.Name,             //1
		gethProcessJob.AlternateName,    //2
		gethProcessJob.StartDate,        //3
		gethProcessJob.EndDate,          //4
		gethProcessJob.Description,      //5
		gethProcessJob.StatusID,         //6
		gethProcessJob.JobCategoryID,    //7
		gethProcessJob.ImportTypeID,     //8
		gethProcessJob.ChainID,          //9
		gethProcessJob.StartBlockNumber, //10
		gethProcessJob.EndBlockNumber,   //11
		gethProcessJob.CreatedBy,        //12
		gethProcessJob.AssetID,          //13
	).Scan(&ID)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(ID), nil
}

func InsertGethProcessJobList(gethProcessJobList []GethProcessJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range gethProcessJobList {
		gethProcessJob := gethProcessJobList[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(gethProcessJob.UUID)
		row := []interface{}{

			*gethProcessJob.ID,               //1
			uuidString,                       //2
			gethProcessJob.Name,              //3
			gethProcessJob.AlternateName,     //4
			&gethProcessJob.StartDate,        //5
			&gethProcessJob.EndDate,          //6
			gethProcessJob.Description,       //7
			*gethProcessJob.StatusID,         //8
			*gethProcessJob.JobCategoryID,    //9
			*gethProcessJob.ImportTypeID,     //10
			*gethProcessJob.ChainID,          //11
			*gethProcessJob.StartBlockNumber, //12
			*gethProcessJob.EndBlockNumber,   //13
			gethProcessJob.CreatedBy,         //14
			&gethProcessJob.CreatedAt,        //15
			gethProcessJob.CreatedBy,         //16
			&now,                             //17
			*gethProcessJob.AssetID,          //18
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_process_jobs"},
		[]string{
			"id",                 //1
			"uuid",               //2
			"name",               //3
			"alternate_name",     //4
			"start_date",         //5
			"end_date",           //6
			"description",        //7
			"status_id",          //8
			"job_category_id",    //9
			"import_type_id",     //10
			"chain_id",           //11
			"start_block_number", //12
			"end_block_number",   //13
			"created_by",         //14
			"created_at",         //15
			"updated_by",         //16
			"updated_at",         //17
			"asset_id",           //18
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
