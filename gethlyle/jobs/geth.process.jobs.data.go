package gethlylejobs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetGethProcessJob(dbConnPgx utils.PgxIface, gethProcessJobID *int) (*GethProcessJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
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
	`, *gethProcessJobID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	defer row.Close()
	gethProcessJob, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethProcessJob])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethProcessJob, nil
}

func GetLatestGethProcessJobByImportTypeIDAndAssetID(dbConnPgx utils.PgxIface, importTypeID, assetID *int) (*GethProcessJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
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
	WHERE import_type_id = $1
	AND asset_id = $2
	-- needs to be success
	AND status_id =$3
	ORDER BY id desc
	LIMIT 1
	`, *importTypeID, *assetID, utils.SUCCESS_STRUCTURED_VALUE_ID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()
	gethProcessJob, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethProcessJob])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethProcessJob, nil
}

func GetGethProcessJobList(dbConnPgx utils.PgxIface) ([]GethProcessJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
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
	gethProcessJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[GethProcessJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethProcessJobs, nil
}

func RemoveGethProcessJob(dbConnPgx utils.PgxIface, gethProcessJobID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethProcessJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_process_jobs WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *gethProcessJobID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func UpdateGethProcessJob(dbConnPgx utils.PgxIface, gethProcessJob *GethProcessJob) error {
	// if the gethProcessJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethProcessJob.ID == nil || *gethProcessJob.ID == 0 {
		return errors.New("gethProcessJob has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethProcessJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_process_jobs SET 
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
		WHERE id=$14 `
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
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
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)

}

func InsertGethProcessJob(dbConnPgx utils.PgxIface, gethProcessJob *GethProcessJob) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethProcessJob DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var ID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_process_jobs  
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

func InsertGethProcessJobList(dbConnPgx utils.PgxIface, gethProcessJobList []GethProcessJob) error {
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
			uuidString,                       //1
			gethProcessJob.Name,              //2
			gethProcessJob.AlternateName,     //3
			&gethProcessJob.StartDate,        //4
			&gethProcessJob.EndDate,          //5
			gethProcessJob.Description,       //6
			*gethProcessJob.StatusID,         //7
			*gethProcessJob.JobCategoryID,    //8
			*gethProcessJob.ImportTypeID,     //9
			*gethProcessJob.ChainID,          //10
			*gethProcessJob.StartBlockNumber, //11
			*gethProcessJob.EndBlockNumber,   //12
			gethProcessJob.CreatedBy,         //13
			&gethProcessJob.CreatedAt,        //14
			gethProcessJob.CreatedBy,         //15
			&now,                             //16
			*gethProcessJob.AssetID,          //17
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_process_jobs"},
		[]string{
			"uuid",               //1
			"name",               //2
			"alternate_name",     //3
			"start_date",         //4
			"end_date",           //5
			"description",        //6
			"status_id",          //7
			"job_category_id",    //8
			"import_type_id",     //9
			"chain_id",           //10
			"start_block_number", //11
			"end_block_number",   //12
			"created_by",         //13
			"created_at",         //14
			"updated_by",         //15
			"updated_at",         //16
			"asset_id",           //17
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethProcessJobList: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetGethProcessJobListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethProcessJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
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
	gethProcessJobList, err := pgx.CollectRows(results, pgx.RowToStructByName[GethProcessJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethProcessJobList, nil
}

func GetTotalGethProcessJobCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM geth_process_jobs`)
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
