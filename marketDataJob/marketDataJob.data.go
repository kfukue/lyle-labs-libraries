package marketdatajob

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

func GetMarketDataJob(dbConnPgx utils.PgxIface, marketDataID, jobID *int) (*MarketDataJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
	market_data_id,  
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
	FROM market_data_jobs 
	WHERE market_data_id = $1
	AND job_id = $2
	`, *marketDataID, *jobID)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	marketDataJob, err := pgx.CollectOneRow(row, pgx.RowToStructByName[MarketDataJob])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &marketDataJob, nil
}

func GetMarketDataJobByMarketDataID(dbConnPgx utils.PgxIface, marketDataID *int) (*MarketDataJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
	market_data_id,
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
	FROM market_data_jobs 
	WHERE market_data_id = $1
	`, *marketDataID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	marketDataJob, err := pgx.CollectOneRow(row, pgx.RowToStructByName[MarketDataJob])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &marketDataJob, nil
}

func RemoveMarketDataJob(dbConnPgx utils.PgxIface, marketDataID, jobID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM market_data_jobs WHERE market_data_id = $1 AND job_id =$2`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *marketDataID, *jobID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetMarketDataJobList(dbConnPgx utils.PgxIface) ([]MarketDataJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	market_data_id,
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
	FROM market_data_jobs`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	marketDataJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[MarketDataJob])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return marketDataJobs, nil
}

func UpdateMarketDataJob(dbConnPgx utils.PgxIface, marketDataJob *MarketDataJob) error {
	// if the marketDataJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (marketDataJob.MarketDataID == nil || *marketDataJob.MarketDataID == 0) || (marketDataJob.JobID == nil || *marketDataJob.JobID == 0) {
		return errors.New("marketDataJob has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateMarketDataJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE market_data_jobs SET 
		name=$1,
		alternate_name=$2,
		start_date=$3,
		end_date=$4,
		description=$5,
		status_id=$6,
		response_status=$7,
		request_url=$8,
		request_body=$9,
		request_method=$10,
		response_data=$11,
		updated_by=$12, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE market_data_id=$13 AND job_id=$14`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		marketDataJob.Name,           //1
		marketDataJob.AlternateName,  //2
		marketDataJob.StartDate,      //3
		marketDataJob.EndDate,        //4
		marketDataJob.Description,    //5
		marketDataJob.StatusID,       //6
		marketDataJob.ResponseStatus, //7
		marketDataJob.RequestUrl,     //8
		marketDataJob.RequestBody,    //9
		marketDataJob.RequestMethod,  //10
		marketDataJob.ResponseData,   //11
		marketDataJob.UpdatedBy,      //12
		marketDataJob.MarketDataID,   //13
		marketDataJob.JobID,          //14
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertMarketDataJob(dbConnPgx utils.PgxIface, marketDataJob *MarketDataJob) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertMarketDataJob DbConn.Begin   %s", err.Error())
		return -1, -1, err
	}
	var marketDataID int
	var jobID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO market_data_jobs(
		market_data_id,
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
		) VALUES (
		 	$1,
			$2,
			uuid_generate_v4(),
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
			current_timestamp at time zone 'UTC',
			$15,
			current_timestamp at time zone 'UTC'
		)
		RETURNING market_data_id, job_id`,
		marketDataJob.MarketDataID,                           //1
		marketDataJob.JobID,                                  //2
		marketDataJob.Name,                                   //3
		marketDataJob.AlternateName,                          //4
		marketDataJob.StartDate.Format(utils.LayoutPostgres), //5
		marketDataJob.EndDate.Format(utils.LayoutPostgres),   //6
		marketDataJob.Description,                            //7
		marketDataJob.StatusID,                               //8
		marketDataJob.ResponseStatus,                         //9
		marketDataJob.RequestUrl,                             //10
		marketDataJob.RequestBody,                            //11
		marketDataJob.RequestMethod,                          //12
		marketDataJob.ResponseData,                           //13
		marketDataJob.ResponseDataJson,                       //14
		marketDataJob.CreatedBy,                              //15
	).Scan(&marketDataID, &jobID)
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
	return int(marketDataID), int(jobID), nil
}

func InsertMarketDataJobList(dbConnPgx utils.PgxIface, marketDataJobList []MarketDataJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range marketDataJobList {
		marketDataJob := marketDataJobList[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(marketDataJob.UUID)
		row := []interface{}{

			marketDataJob.MarketDataID,     //1
			marketDataJob.JobID,            //2
			uuidString,                     //3
			marketDataJob.Name,             //4
			marketDataJob.AlternateName,    //5
			marketDataJob.StartDate,        //6
			marketDataJob.EndDate,          //7
			marketDataJob.Description,      //8
			marketDataJob.StatusID,         //9
			marketDataJob.ResponseStatus,   //10
			marketDataJob.RequestUrl,       //11
			marketDataJob.RequestBody,      //12
			marketDataJob.RequestMethod,    //13
			marketDataJob.ResponseData,     //14
			marketDataJob.ResponseDataJson, //15
			marketDataJob.CreatedBy,        //16
			&now,                           //17
			marketDataJob.CreatedBy,        //18
			&now,                           //19
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"market_data_jobs"},
		[]string{
			"market_data_id",     //1
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
	log.Println(fmt.Printf("InsertMarketDataJobList: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetMarketDataJobListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]MarketDataJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
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
		job_category_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
	FROM market_data_jobs 
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
	marketDataJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[MarketDataJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return marketDataJobs, nil
}

func GetTotalMarketDataJobsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM market_data_jobs
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
