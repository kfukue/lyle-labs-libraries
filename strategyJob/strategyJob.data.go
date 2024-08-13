package strategyjob

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
)

func GetStrategyJob(dbConnPgx utils.PgxIface, strategyID, jobID *int) (*StrategyJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		strategy_id,  
		job_id,
		uuid, 
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		status_id,
		response_status,
		request_url,
		request_body,
		request_method,
		response_data,
		response_data_json,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM strategy_jobs 
	WHERE strategy_id = $1
	AND job_id = $2
	`, *strategyID, *jobID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	strategyJob, err := pgx.CollectOneRow(row, pgx.RowToStructByName[StrategyJob])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &strategyJob, nil
}

func GetStrategyJobByStrategyID(dbConnPgx utils.PgxIface, strategyID *int) (*StrategyJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		strategy_id,
		job_id,
		uuid, 
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		status_id,
		response_status,
		request_url,
		request_body,
		request_method,
		response_data,
		response_data_json,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM strategy_jobs 
	WHERE strategy_id = $1
	`, *strategyID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	strategyJob, err := pgx.CollectOneRow(row, pgx.RowToStructByName[StrategyJob])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &strategyJob, nil
}

func RemoveStrategyJob(dbConnPgx utils.PgxIface, strategyID, jobID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveStrategyJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM strategy_jobs WHERE 
		strategy_id = $1 AND job_id =$2`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *strategyID, *jobID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetStrategyJobList(dbConnPgx utils.PgxIface) ([]StrategyJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		strategy_id,
		job_id,
		uuid, 
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		status_id,
		response_status,
		request_url,
		request_body,
		request_method,
		response_data,
		response_data_json,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM strategy_jobs`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategyJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[StrategyJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyJobs, nil
}

func UpdateStrategyJob(dbConnPgx utils.PgxIface, strategyJob *StrategyJob) error {
	// if the strategyJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (strategyJob.StrategyID == nil || *strategyJob.StrategyID == 0) || (strategyJob.JobID == nil || *strategyJob.JobID == 0) {
		return errors.New("strategyJob has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateStrategyJob DbConn.Begin   %s", err.Error())
		return err
	}
	layoutPostgres := utils.LayoutPostgres
	sql := `UPDATE strategy_jobs SET 
		name=$1,  
		alternate_name=$2, 
		start_date =$3,
		end_date =$4,
		description=$5, 
		status_id=$6, 
		response_status=$7,
		request_url=$8, 
		request_body=$9, 
		request_method=$10, 
		response_data=$11, 
		response_data_json =$12,
		updated_by=$13, 
		updated_at=current_timestamp at time zone 'UTC'	
		WHERE strategy_id=$14 AND job_id=$15`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		strategyJob.Name,                             //1
		strategyJob.AlternateName,                    //2
		strategyJob.StartDate.Format(layoutPostgres), //3
		strategyJob.EndDate.Format(layoutPostgres),   //4
		strategyJob.Description,                      //5
		strategyJob.StatusID,                         //6
		strategyJob.ResponseStatus,                   //7
		strategyJob.RequestUrl,                       //8
		strategyJob.RequestBody,                      //9
		strategyJob.RequestMethod,                    //10
		strategyJob.ResponseData,                     //11
		strategyJob.ResponseDataJson,                 //12
		strategyJob.UpdatedBy,                        //13
		strategyJob.StrategyID,                       //14
		strategyJob.JobID,                            //15
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertStrategyJob(dbConnPgx utils.PgxIface, strategyJob *StrategyJob) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertStrategyJob DbConn.Begin   %s", err.Error())
		return -1, -1, err
	}
	var StrategyID int
	var JobID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO strategy_jobs  
	(
		strategy_id,
		job_id,
		uuid, 
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		status_id,
		response_status,
		request_url,
		request_body,
		request_method,
		response_data,
		response_data_json,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES ($1,
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
		RETURNING strategy_id, job_id`,
		strategyJob.StrategyID,       //1
		strategyJob.JobID,            //2
		strategyJob.UUID,             //3
		strategyJob.Name,             //4
		strategyJob.AlternateName,    //5
		strategyJob.StartDate,        //6
		strategyJob.EndDate,          //7
		strategyJob.Description,      //8
		strategyJob.StatusID,         //9
		strategyJob.ResponseStatus,   //10
		strategyJob.RequestUrl,       //11
		strategyJob.RequestBody,      //12
		strategyJob.RequestMethod,    //13
		strategyJob.ResponseData,     //14
		strategyJob.ResponseDataJson, //15
		strategyJob.CreatedBy,        //16
	).Scan(&StrategyID, &JobID)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, -1, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, -1, err
	}
	return int(StrategyID), int(JobID), nil
}

func InsertStrategyJobList(dbConnPgx utils.PgxIface, strategyJobList []StrategyJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range strategyJobList {
		strategyJob := strategyJobList[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(strategyJob.UUID)
		row := []interface{}{
			strategyJob.StrategyID,     //1
			strategyJob.JobID,          //2
			uuidString,                 //3
			strategyJob.Name,           //4
			strategyJob.AlternateName,  //5
			&strategyJob.StartDate,     //6
			&strategyJob.EndDate,       //7
			strategyJob.Description,    //8
			strategyJob.StatusID,       //9
			strategyJob.ResponseStatus, //10
			strategyJob.RequestUrl,     //11
			strategyJob.RequestBody,    //12
			strategyJob.RequestMethod,  //13
			strategyJob.ResponseData,   //14
			// TODO: erroring out in json insert look into it later TAT-27
			strategyJob.ResponseDataJson, //15
			strategyJob.CreatedBy,        //16
			&strategyJob.CreatedAt,       //17
			strategyJob.CreatedBy,        //18
			&now,                         //19
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"strategy_jobs"},
		[]string{
			"strategy_id",        //1
			"job_id",             //2
			"uuid",               //3
			"name",               //4
			"alternate_name",     //5
			"start_date",         //6
			"end_date",           //7
			"description",        //8
			"status_id",          //9
			"response_status",    //10
			"request_url",        //11
			"request_body",       //12
			"request_method",     //13
			"response_data",      //14
			"response_data_json", //15
			"created_by",         //16
			"created_at",         //17
			"updated_by",         //18
			"updated_at",         //19
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertStrategyJobList: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetStrategyJobListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]StrategyJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		strategy_id,  
		job_id,
		uuid, 
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		status_id,
		response_status,
		request_url,
		request_body,
		request_method,
		response_data,
		response_data_json,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM strategy_jobs 
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
	strategyJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[StrategyJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyJobs, nil
}

func GetTotalStrategyJobsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM strategy_jobs
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
