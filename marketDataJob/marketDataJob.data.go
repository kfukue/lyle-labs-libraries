package marketdatajob

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
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetMarketDataJob(marketDataID int, jobID int) (*MarketDataJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	`, marketDataID, jobID)

	marketDataJob := &MarketDataJob{}
	err := row.Scan(
		&marketDataJob.MarketDataID,
		&marketDataJob.JobID,
		&marketDataJob.UUID,
		&marketDataJob.Name,
		&marketDataJob.AlternateName,
		&marketDataJob.StartDate,
		&marketDataJob.EndDate,
		&marketDataJob.Description,
		&marketDataJob.StatusID,
		&marketDataJob.ResponseStatus,
		&marketDataJob.RequestUrl,
		&marketDataJob.RequestBody,
		&marketDataJob.RequestMethod,
		&marketDataJob.ResponseData,
		&marketDataJob.ResponseDataJson,
		&marketDataJob.CreatedBy,
		&marketDataJob.CreatedAt,
		&marketDataJob.UpdatedBy,
		&marketDataJob.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return marketDataJob, nil
}

func GetMarketDataJobByMarketDataID(marketDataID int) (*MarketDataJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	`, marketDataID)

	marketDataJob := &MarketDataJob{}
	err := row.Scan(
		&marketDataJob.MarketDataID,
		&marketDataJob.JobID,
		&marketDataJob.UUID,
		&marketDataJob.Name,
		&marketDataJob.AlternateName,
		&marketDataJob.StartDate,
		&marketDataJob.EndDate,
		&marketDataJob.Description,
		&marketDataJob.StatusID,
		&marketDataJob.ResponseStatus,
		&marketDataJob.RequestUrl,
		&marketDataJob.RequestBody,
		&marketDataJob.RequestMethod,
		&marketDataJob.ResponseData,
		&marketDataJob.ResponseDataJson,
		&marketDataJob.CreatedBy,
		&marketDataJob.CreatedAt,
		&marketDataJob.UpdatedBy,
		&marketDataJob.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return marketDataJob, nil
}

func GetTopTenMarketDataJobs() ([]MarketDataJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	marketDataJobs := make([]MarketDataJob, 0)
	for results.Next() {
		var marketDataJob MarketDataJob
		results.Scan(
			&marketDataJob.MarketDataID,
			&marketDataJob.JobID,
			&marketDataJob.UUID,
			&marketDataJob.Name,
			&marketDataJob.AlternateName,
			&marketDataJob.StartDate,
			&marketDataJob.EndDate,
			&marketDataJob.Description,
			&marketDataJob.StatusID,
			&marketDataJob.ResponseStatus,
			&marketDataJob.RequestUrl,
			&marketDataJob.RequestBody,
			&marketDataJob.RequestMethod,
			&marketDataJob.ResponseData,
			&marketDataJob.ResponseDataJson,
			&marketDataJob.CreatedBy,
			&marketDataJob.CreatedAt,
			&marketDataJob.UpdatedBy,
			&marketDataJob.UpdatedAt,
		)

		marketDataJobs = append(marketDataJobs, marketDataJob)
	}
	return marketDataJobs, nil
}

func RemoveMarketDataJob(marketDataID int, jobID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM market_data_jobs WHERE 
	market_data_id = $1 AND job_id =$2`, marketDataID, jobID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetMarketDataJobList() ([]MarketDataJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	marketDataJobs := make([]MarketDataJob, 0)
	for results.Next() {
		var marketDataJob MarketDataJob
		results.Scan(
			&marketDataJob.MarketDataID,
			&marketDataJob.JobID,
			&marketDataJob.UUID,
			&marketDataJob.Name,
			&marketDataJob.AlternateName,
			&marketDataJob.StartDate,
			&marketDataJob.EndDate,
			&marketDataJob.Description,
			&marketDataJob.StatusID,
			&marketDataJob.ResponseStatus,
			&marketDataJob.RequestUrl,
			&marketDataJob.RequestBody,
			&marketDataJob.RequestMethod,
			&marketDataJob.ResponseData,
			&marketDataJob.ResponseDataJson,
			&marketDataJob.CreatedBy,
			&marketDataJob.CreatedAt,
			&marketDataJob.UpdatedBy,
			&marketDataJob.UpdatedAt,
		)

		marketDataJobs = append(marketDataJobs, marketDataJob)
	}
	return marketDataJobs, nil
}

func UpdateMarketDataJob(marketDataJob MarketDataJob) error {
	// if the marketDataJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (marketDataJob.MarketDataID == nil || *marketDataJob.MarketDataID == 0) || (marketDataJob.JobID == nil || *marketDataJob.JobID == 0) {
		return errors.New("marketDataJob has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE market_data_jobs SET 

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
		WHERE market_data_id=$13 AND job_id=$14`,
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
		// marketDataJob.ResponseDataJson,//
		marketDataJob.UpdatedBy,    //12
		marketDataJob.MarketDataID, //13
		marketDataJob.JobID,        //14
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertMarketDataJob(marketDataJob MarketDataJob) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var MarketDataID int
	var JobID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO market_data_jobs  
	(
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
		) VALUES ($1,
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
		marketDataJob.MarketDataID,     //1
		marketDataJob.JobID,            //2
		marketDataJob.Name,             //3
		marketDataJob.AlternateName,    //4
		marketDataJob.StartDate,        //5
		marketDataJob.EndDate,          //6
		marketDataJob.Description,      //7
		marketDataJob.StatusID,         //8
		marketDataJob.ResponseStatus,   //9
		marketDataJob.RequestUrl,       //10
		marketDataJob.RequestBody,      //11
		marketDataJob.RequestMethod,    //12
		marketDataJob.ResponseData,     //13
		marketDataJob.ResponseDataJson, //14
		marketDataJob.CreatedBy,        //15
	).Scan(&MarketDataID, &JobID)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(MarketDataID), int(JobID), nil
}

func InsertMarketDataJobListManual(marketDataJobList []MarketDataJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range marketDataJobList {
		marketDataJob := marketDataJobList[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(marketDataJob.UUID)
		row := []interface{}{

			*marketDataJob.MarketDataID,  //1
			*marketDataJob.JobID,         //2
			uuidString,                   //3
			marketDataJob.Name,           //4
			marketDataJob.AlternateName,  //5
			&marketDataJob.StartDate,     //6
			&marketDataJob.EndDate,       //7
			marketDataJob.Description,    //8
			*marketDataJob.StatusID,      //9
			marketDataJob.ResponseStatus, //10
			marketDataJob.RequestUrl,     //11
			marketDataJob.RequestBody,    //12
			marketDataJob.RequestMethod,  //13
			marketDataJob.ResponseData,   //14
			// TODO: erroring out in json insert look into it later TAT-27
			// marketDataJob.ResponseDataJson //15
			marketDataJob.CreatedBy,  //16
			&marketDataJob.CreatedAt, //17
			marketDataJob.CreatedBy,  //18
			&now,                     //19
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"market_data_jobs"},
		[]string{
			"market_data_id",  //1
			"job_id",          //2
			"uuid",            //3
			"name",            //4
			"alternate_name",  //5
			"start_date",      //6
			"end_date",        //7
			"description",     //8
			"status_id",       //9
			"response_status", //10
			"request_url",     //11
			"request_body",    //12
			"request_method",  //13
			"response_data",   //14
			// "response_data_json", //15
			"created_by", //16
			"created_at", //17
			"updated_by", //18
			"updated_at", //19
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

func InsertMarketDataJobList(marketDataJobList []MarketDataJob) error {
	txn, err := database.DbConn.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	layoutPostgres := utils.LayoutPostgres
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn(
		"market_data_jobs",
		"market_data_id",  //1
		"job_id",          //2
		"uuid",            //3
		"name",            //4
		"alternate_name",  //5
		"start_date",      //6
		"end_date",        //7
		"description",     //8
		"status_id",       //9
		"response_status", //10
		"request_url",     //11
		"request_body",    //12
		"request_method",  //13
		"response_data",   //14
		// "response_data_json", //15
		"created_by", //16
		"created_at", //17
		"updated_by", //18
		"updated_at", //19

	))
	if err != nil {
		log.Fatal(err)
	}
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	for _, marketDataJob := range marketDataJobList {
		_, err = stmt.Exec(
			marketDataJob.MarketDataID,   //1
			marketDataJob.JobID,          //2
			marketDataJob.UUID,           //3
			marketDataJob.Name,           //4
			marketDataJob.AlternateName,  //5
			marketDataJob.StartDate,      //6
			marketDataJob.EndDate,        //7
			marketDataJob.Description,    //8
			marketDataJob.StatusID,       //9
			marketDataJob.ResponseStatus, //10
			marketDataJob.RequestUrl,     //11
			marketDataJob.RequestBody,    //12
			marketDataJob.RequestMethod,  //13
			marketDataJob.ResponseData,   //14
			// marketDataJob.ResponseDataJson, //15
			marketDataJob.CreatedBy,    //16
			marketDataJob.CreatedAt,    //17
			marketDataJob.CreatedBy,    //18
			now.Format(layoutPostgres), //19
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
