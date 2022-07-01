package positionjob

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/lib/pq"
)

func GetPositionJob(positionID int, jobID int) (*PositionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	position_id,  
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
	FROM position_jobs 
	WHERE position_id = $1
	AND job_id = $2
	`, positionID, jobID)

	positionJob := &PositionJob{}
	err := row.Scan(
		&positionJob.PositionID,
		&positionJob.JobID,
		&positionJob.UUID,
		&positionJob.Name,
		&positionJob.AlternateName,
		&positionJob.StartDate,
		&positionJob.EndDate,
		&positionJob.Description,
		&positionJob.StatusID,
		&positionJob.ResponseStatus,
		&positionJob.RequestUrl,
		&positionJob.RequestBody,
		&positionJob.RequestMethod,
		&positionJob.ResponseData,
		&positionJob.ResponseDataJson,
		&positionJob.CreatedBy,
		&positionJob.CreatedAt,
		&positionJob.UpdatedBy,
		&positionJob.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return positionJob, nil
}

func GetPositionJobByUUID(positionJobUUID string) (*PositionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	position_id,  
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
	FROM position_jobs 
	WHERE text(uuid) = $1
	`, positionJobUUID)

	positionJob := &PositionJob{}
	err := row.Scan(
		&positionJob.PositionID,
		&positionJob.JobID,
		&positionJob.UUID,
		&positionJob.Name,
		&positionJob.AlternateName,
		&positionJob.StartDate,
		&positionJob.EndDate,
		&positionJob.Description,
		&positionJob.StatusID,
		&positionJob.ResponseStatus,
		&positionJob.RequestUrl,
		&positionJob.RequestBody,
		&positionJob.RequestMethod,
		&positionJob.ResponseData,
		&positionJob.ResponseDataJson,
		&positionJob.CreatedBy,
		&positionJob.CreatedAt,
		&positionJob.UpdatedBy,
		&positionJob.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return positionJob, nil
}

func GetPositionJobsByUUIDs(UUIDList []string) ([]PositionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	position_id,  
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
	FROM position_jobs
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	positionJobList := make([]PositionJob, 0)
	for results.Next() {
		var positionJob PositionJob
		results.Scan(
			&positionJob.PositionID,
			&positionJob.JobID,
			&positionJob.UUID,
			&positionJob.Name,
			&positionJob.AlternateName,
			&positionJob.StartDate,
			&positionJob.EndDate,
			&positionJob.Description,
			&positionJob.StatusID,
			&positionJob.ResponseStatus,
			&positionJob.RequestUrl,
			&positionJob.RequestBody,
			&positionJob.RequestMethod,
			&positionJob.ResponseData,
			&positionJob.ResponseDataJson,
			&positionJob.CreatedBy,
			&positionJob.CreatedAt,
			&positionJob.UpdatedBy,
			&positionJob.UpdatedAt,
		)

		positionJobList = append(positionJobList, positionJob)
	}
	return positionJobList, nil
}

func GetTopTenPositionJobs() ([]PositionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	position_id,  
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
	FROM position_jobs 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	positionJobs := make([]PositionJob, 0)
	for results.Next() {
		var positionJob PositionJob
		results.Scan(
			&positionJob.PositionID,
			&positionJob.JobID,
			&positionJob.UUID,
			&positionJob.Name,
			&positionJob.AlternateName,
			&positionJob.StartDate,
			&positionJob.EndDate,
			&positionJob.Description,
			&positionJob.StatusID,
			&positionJob.ResponseStatus,
			&positionJob.RequestUrl,
			&positionJob.RequestBody,
			&positionJob.RequestMethod,
			&positionJob.ResponseData,
			&positionJob.ResponseDataJson,
			&positionJob.CreatedBy,
			&positionJob.CreatedAt,
			&positionJob.UpdatedBy,
			&positionJob.UpdatedAt,
		)

		positionJobs = append(positionJobs, positionJob)
	}
	return positionJobs, nil
}

func RemovePositionJob(positionID int, jobID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM position_jobs WHERE 
	position_id = $1 AND job_id =$2`, positionID, jobID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemovePositionJobByUUID(positionJobUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM position_jobs WHERE 
		WHERE text(uuid) = $1`,
		positionJobUUID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetPositionJobList() ([]PositionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	position_id,  
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
	FROM position_jobs`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	positionJobs := make([]PositionJob, 0)
	for results.Next() {
		var positionJob PositionJob
		results.Scan(
			&positionJob.PositionID,
			&positionJob.JobID,
			&positionJob.UUID,
			&positionJob.Name,
			&positionJob.AlternateName,
			&positionJob.StartDate,
			&positionJob.EndDate,
			&positionJob.Description,
			&positionJob.StatusID,
			&positionJob.ResponseStatus,
			&positionJob.RequestUrl,
			&positionJob.RequestBody,
			&positionJob.RequestMethod,
			&positionJob.ResponseData,
			&positionJob.ResponseDataJson,
			&positionJob.CreatedBy,
			&positionJob.CreatedAt,
			&positionJob.UpdatedBy,
			&positionJob.UpdatedAt,
		)

		positionJobs = append(positionJobs, positionJob)
	}
	return positionJobs, nil
}

func UpdatePositionJob(positionJob PositionJob) error {
	// if the positionJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (positionJob.PositionID == nil || *positionJob.PositionID == 0) || (positionJob.JobID == nil || *positionJob.JobID == 0) {
		return errors.New("positionJob has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE position_jobs SET 
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
		WHERE position_id=$13 AND job_id=$14`,

		positionJob.Name,           //1
		positionJob.AlternateName,  //2
		positionJob.StartDate,      //3
		positionJob.EndDate,        //4
		positionJob.Description,    //5
		positionJob.StatusID,       //6
		positionJob.ResponseStatus, //7
		positionJob.RequestUrl,     //8
		positionJob.RequestBody,    //9
		positionJob.RequestMethod,  //10
		positionJob.ResponseData,   //11
		// positionJob.ResponseDataJson,//
		positionJob.UpdatedBy,  //12
		positionJob.PositionID, //13
		positionJob.JobID,      //14
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdatePositionJobByUUID(positionJob PositionJob) error {
	// if the positionJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (positionJob.PositionID == nil || *positionJob.PositionID == 0) || (positionJob.JobID == nil || *positionJob.JobID == 0) {
		return errors.New("positionJob has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE position_jobs SET 
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
		WHERE text(uuid) = $13
		`,
		positionJob.Name,           //1
		positionJob.AlternateName,  //2
		positionJob.StartDate,      //3
		positionJob.EndDate,        //4
		positionJob.Description,    //5
		positionJob.StatusID,       //6
		positionJob.ResponseStatus, //7
		positionJob.RequestUrl,     //8
		positionJob.RequestBody,    //9
		positionJob.RequestMethod,  //10
		positionJob.ResponseData,   //11
		// positionJob.ResponseDataJson,//
		positionJob.UpdatedBy, //12
		positionJob.UUID,      //13
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertPositionJob(positionJob PositionJob) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var PositionID int
	var JobID int
	positionJobUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	if positionJob.UUID == "" {
		positionJob.UUID = positionJobUUID.String()
	}
	err = database.DbConn.QueryRowContext(ctx, `INSERT INTO position_jobs  
	(
		position_id,  
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
		RETURNING position_id, job_id`,
		positionJob.PositionID,       //1
		positionJob.JobID,            //2
		positionJob.Name,             //3
		positionJob.AlternateName,    //4
		positionJob.StartDate,        //5
		positionJob.EndDate,          //6
		positionJob.Description,      //7
		positionJob.StatusID,         //8
		positionJob.ResponseStatus,   //9
		positionJob.RequestUrl,       //10
		positionJob.RequestBody,      //11
		positionJob.RequestMethod,    //12
		positionJob.ResponseData,     //13
		positionJob.ResponseDataJson, //14
		positionJob.CreatedBy,        //15
	).Scan(&PositionID, &JobID)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(PositionID), int(JobID), nil
}

func InsertPositionJobs(positionJobs []PositionJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range positionJobs {
		positionJob := positionJobs[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(positionJob.UUID)
		row := []interface{}{

			*positionJob.PositionID,    //1
			*positionJob.JobID,         //2
			uuidString,                 //3
			positionJob.Name,           //4
			positionJob.AlternateName,  //5
			&positionJob.StartDate,     //6
			&positionJob.EndDate,       //7
			positionJob.Description,    //8
			*positionJob.StatusID,      //9
			positionJob.ResponseStatus, //10
			positionJob.RequestUrl,     //11
			positionJob.RequestBody,    //12
			positionJob.RequestMethod,  //13
			positionJob.ResponseData,   //14
			// TODO: erroring out in json insert look into it later TAT-27
			// positionJob.ResponseDataJson //15
			positionJob.CreatedBy,  //16
			&positionJob.CreatedAt, //17
			positionJob.CreatedBy,  //18
			&now,                   //19
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"position_jobs"},
		[]string{
			"position_id",     //1
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
