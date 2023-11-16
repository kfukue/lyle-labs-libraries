package strategyjob

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
)

func GetStrategyJob(strategyID int, jobID int) (*StrategyJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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
	`, strategyID, jobID)

	strategyJob := &StrategyJob{}
	err := row.Scan(
		&strategyJob.StrategyID,
		&strategyJob.JobID,
		&strategyJob.UUID,
		&strategyJob.Name,
		&strategyJob.AlternateName,
		&strategyJob.StartDate,
		&strategyJob.EndDate,
		&strategyJob.Description,
		&strategyJob.StatusID,
		&strategyJob.ResponseStatus,
		&strategyJob.RequestUrl,
		&strategyJob.RequestBody,
		&strategyJob.RequestMethod,
		&strategyJob.ResponseData,
		&strategyJob.ResponseDataJson,
		&strategyJob.CreatedBy,
		&strategyJob.CreatedAt,
		&strategyJob.UpdatedBy,
		&strategyJob.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyJob, nil
}

func GetStrategyJobByStrategyID(strategyID int) (*StrategyJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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
	`, strategyID)

	strategyJob := &StrategyJob{}
	err := row.Scan(
		&strategyJob.StrategyID,
		&strategyJob.JobID,
		&strategyJob.UUID,
		&strategyJob.Name,
		&strategyJob.AlternateName,
		&strategyJob.StartDate,
		&strategyJob.EndDate,
		&strategyJob.Description,
		&strategyJob.StatusID,
		&strategyJob.ResponseStatus,
		&strategyJob.RequestUrl,
		&strategyJob.RequestBody,
		&strategyJob.RequestMethod,
		&strategyJob.ResponseData,
		&strategyJob.ResponseDataJson,
		&strategyJob.CreatedBy,
		&strategyJob.CreatedAt,
		&strategyJob.UpdatedBy,
		&strategyJob.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyJob, nil
}

func GetTopTenStrategyJobs() ([]StrategyJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategyJobs := make([]StrategyJob, 0)
	for results.Next() {
		var strategyJob StrategyJob
		results.Scan(
			&strategyJob.StrategyID,
			&strategyJob.JobID,
			&strategyJob.UUID,
			&strategyJob.Name,
			&strategyJob.AlternateName,
			&strategyJob.StartDate,
			&strategyJob.EndDate,
			&strategyJob.Description,
			&strategyJob.StatusID,
			&strategyJob.ResponseStatus,
			&strategyJob.RequestUrl,
			&strategyJob.RequestBody,
			&strategyJob.RequestMethod,
			&strategyJob.ResponseData,
			&strategyJob.ResponseDataJson,
			&strategyJob.CreatedBy,
			&strategyJob.CreatedAt,
			&strategyJob.UpdatedBy,
			&strategyJob.UpdatedAt,
		)

		strategyJobs = append(strategyJobs, strategyJob)
	}
	return strategyJobs, nil
}

func RemoveStrategyJob(strategyID int, jobID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM strategy_jobs WHERE 
	strategy_id = $1 AND job_id =$2`, strategyID, jobID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetStrategyJobList() ([]StrategyJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	strategyJobs := make([]StrategyJob, 0)
	for results.Next() {
		var strategyJob StrategyJob
		results.Scan(
			&strategyJob.StrategyID,
			&strategyJob.JobID,
			&strategyJob.UUID,
			&strategyJob.Name,
			&strategyJob.AlternateName,
			&strategyJob.StartDate,
			&strategyJob.EndDate,
			&strategyJob.Description,
			&strategyJob.StatusID,
			&strategyJob.ResponseStatus,
			&strategyJob.RequestUrl,
			&strategyJob.RequestBody,
			&strategyJob.RequestMethod,
			&strategyJob.ResponseData,
			&strategyJob.ResponseDataJson,
			&strategyJob.CreatedBy,
			&strategyJob.CreatedAt,
			&strategyJob.UpdatedBy,
			&strategyJob.UpdatedAt,
		)

		strategyJobs = append(strategyJobs, strategyJob)
	}
	return strategyJobs, nil
}

func UpdateStrategyJob(strategyJob StrategyJob) error {
	// if the strategyJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (strategyJob.StrategyID == nil || *strategyJob.StrategyID == 0) || (strategyJob.JobID == nil || *strategyJob.JobID == 0) {
		return errors.New("strategyJob has invalid ID")
	}
	layoutPostgres := utils.LayoutPostgres
	_, err := database.DbConnPgx.Query(ctx, `UPDATE strategy_jobs SET 
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
		WHERE strategy_id=$14 AND job_id=$15`,
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
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertStrategyJob(strategyJob StrategyJob) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var StrategyID int
	var JobID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO strategy_jobs  
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
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(StrategyID), int(JobID), nil
}

func InsertStrategyJobList(strategyJobList []StrategyJob) error {
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
			*strategyJob.StrategyID,    //1
			*strategyJob.JobID,         //2
			uuidString,                 //3
			strategyJob.Name,           //4
			strategyJob.AlternateName,  //5
			&strategyJob.StartDate,     //6
			&strategyJob.EndDate,       //7
			strategyJob.Description,    //8
			*strategyJob.StatusID,      //9
			strategyJob.ResponseStatus, //10
			strategyJob.RequestUrl,     //11
			strategyJob.RequestBody,    //12
			strategyJob.RequestMethod,  //13
			strategyJob.ResponseData,   //14
			// TODO: erroring out in json insert look into it later TAT-27
			// strategyJob.ResponseDataJson //15
			strategyJob.CreatedBy,  //16
			&strategyJob.CreatedAt, //17
			strategyJob.CreatedBy,  //18
			&now,                   //19
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"strategy_jobs"},
		[]string{
			"strategy_id",     //1
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
