package transactionjob

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/lib/pq"
)

func GetTransactionJob(transactionID int, jobID int) (*TransactionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	transaction_id,  
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
	FROM transaction_jobs 
	WHERE transaction_id = $1
	AND job_id = $2
	`, transactionID, jobID)

	transactionJob := &TransactionJob{}
	err := row.Scan(
		&transactionJob.TransactionID,
		&transactionJob.JobID,
		&transactionJob.UUID,
		&transactionJob.Name,
		&transactionJob.AlternateName,
		&transactionJob.StartDate,
		&transactionJob.EndDate,
		&transactionJob.Description,
		&transactionJob.StatusID,
		&transactionJob.ResponseStatus,
		&transactionJob.RequestUrl,
		&transactionJob.RequestBody,
		&transactionJob.RequestMethod,
		&transactionJob.ResponseData,
		&transactionJob.ResponseDataJson,
		&transactionJob.CreatedBy,
		&transactionJob.CreatedAt,
		&transactionJob.UpdatedBy,
		&transactionJob.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionJob, nil
}

func GetTransactionJobByUUID(transactionJobUUID string) (*TransactionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	transaction_id,  
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
	FROM transaction_jobs 
	WHERE text(uuid) = $1
	`, transactionJobUUID)

	transactionJob := &TransactionJob{}
	err := row.Scan(
		&transactionJob.TransactionID,
		&transactionJob.JobID,
		&transactionJob.UUID,
		&transactionJob.Name,
		&transactionJob.AlternateName,
		&transactionJob.StartDate,
		&transactionJob.EndDate,
		&transactionJob.Description,
		&transactionJob.StatusID,
		&transactionJob.ResponseStatus,
		&transactionJob.RequestUrl,
		&transactionJob.RequestBody,
		&transactionJob.RequestMethod,
		&transactionJob.ResponseData,
		&transactionJob.ResponseDataJson,
		&transactionJob.CreatedBy,
		&transactionJob.CreatedAt,
		&transactionJob.UpdatedBy,
		&transactionJob.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionJob, nil
}

func GetTransactionJobsByUUIDs(UUIDList []string) ([]TransactionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	transaction_id,  
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
	FROM transaction_jobs
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	transactionJobList := make([]TransactionJob, 0)
	for results.Next() {
		var transactionJob TransactionJob
		results.Scan(
			&transactionJob.TransactionID,
			&transactionJob.JobID,
			&transactionJob.UUID,
			&transactionJob.Name,
			&transactionJob.AlternateName,
			&transactionJob.StartDate,
			&transactionJob.EndDate,
			&transactionJob.Description,
			&transactionJob.StatusID,
			&transactionJob.ResponseStatus,
			&transactionJob.RequestUrl,
			&transactionJob.RequestBody,
			&transactionJob.RequestMethod,
			&transactionJob.ResponseData,
			&transactionJob.ResponseDataJson,
			&transactionJob.CreatedBy,
			&transactionJob.CreatedAt,
			&transactionJob.UpdatedBy,
			&transactionJob.UpdatedAt,
		)

		transactionJobList = append(transactionJobList, transactionJob)
	}
	return transactionJobList, nil
}

func GetTopTenTransactionJobs() ([]TransactionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	transaction_id,  
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
	FROM transaction_jobs 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	transactionJobs := make([]TransactionJob, 0)
	for results.Next() {
		var transactionJob TransactionJob
		results.Scan(
			&transactionJob.TransactionID,
			&transactionJob.JobID,
			&transactionJob.UUID,
			&transactionJob.Name,
			&transactionJob.AlternateName,
			&transactionJob.StartDate,
			&transactionJob.EndDate,
			&transactionJob.Description,
			&transactionJob.StatusID,
			&transactionJob.ResponseStatus,
			&transactionJob.RequestUrl,
			&transactionJob.RequestBody,
			&transactionJob.RequestMethod,
			&transactionJob.ResponseData,
			&transactionJob.ResponseDataJson,
			&transactionJob.CreatedBy,
			&transactionJob.CreatedAt,
			&transactionJob.UpdatedBy,
			&transactionJob.UpdatedAt,
		)

		transactionJobs = append(transactionJobs, transactionJob)
	}
	return transactionJobs, nil
}

func RemoveTransactionJob(transactionID int, jobID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM transaction_jobs WHERE 
	transaction_id = $1 AND job_id =$2`, transactionID, jobID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveTransactionJobByUUID(transactionJobUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM transaction_jobs WHERE 
		WHERE text(uuid) = $1`,
		transactionJobUUID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetTransactionJobList() ([]TransactionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	transaction_id,  
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
	FROM transaction_jobs`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	transactionJobs := make([]TransactionJob, 0)
	for results.Next() {
		var transactionJob TransactionJob
		results.Scan(
			&transactionJob.TransactionID,
			&transactionJob.JobID,
			&transactionJob.UUID,
			&transactionJob.Name,
			&transactionJob.AlternateName,
			&transactionJob.StartDate,
			&transactionJob.EndDate,
			&transactionJob.Description,
			&transactionJob.StatusID,
			&transactionJob.ResponseStatus,
			&transactionJob.RequestUrl,
			&transactionJob.RequestBody,
			&transactionJob.RequestMethod,
			&transactionJob.ResponseData,
			&transactionJob.ResponseDataJson,
			&transactionJob.CreatedBy,
			&transactionJob.CreatedAt,
			&transactionJob.UpdatedBy,
			&transactionJob.UpdatedAt,
		)

		transactionJobs = append(transactionJobs, transactionJob)
	}
	return transactionJobs, nil
}

func UpdateTransactionJob(transactionJob TransactionJob) error {
	// if the transactionJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (transactionJob.TransactionID == nil || *transactionJob.TransactionID == 0) || (transactionJob.JobID == nil || *transactionJob.JobID == 0) {
		return errors.New("transactionJob has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE transaction_jobs SET 
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
		WHERE transaction_id=$13 AND job_id=$14`,

		transactionJob.Name,           //1
		transactionJob.AlternateName,  //2
		transactionJob.StartDate,      //3
		transactionJob.EndDate,        //4
		transactionJob.Description,    //5
		transactionJob.StatusID,       //6
		transactionJob.ResponseStatus, //7
		transactionJob.RequestUrl,     //8
		transactionJob.RequestBody,    //9
		transactionJob.RequestMethod,  //10
		transactionJob.ResponseData,   //11
		// transactionJob.ResponseDataJson,//
		transactionJob.UpdatedBy,     //12
		transactionJob.TransactionID, //13
		transactionJob.JobID,         //14
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateTransactionJobByUUID(transactionJob TransactionJob) error {
	// if the transactionJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (transactionJob.TransactionID == nil || *transactionJob.TransactionID == 0) || (transactionJob.JobID == nil || *transactionJob.JobID == 0) {
		return errors.New("transactionJob has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE transaction_jobs SET 
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
		transactionJob.Name,           //1
		transactionJob.AlternateName,  //2
		transactionJob.StartDate,      //3
		transactionJob.EndDate,        //4
		transactionJob.Description,    //5
		transactionJob.StatusID,       //6
		transactionJob.ResponseStatus, //7
		transactionJob.RequestUrl,     //8
		transactionJob.RequestBody,    //9
		transactionJob.RequestMethod,  //10
		transactionJob.ResponseData,   //11
		// transactionJob.ResponseDataJson,//
		transactionJob.UpdatedBy, //12
		transactionJob.UUID,      //13
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertTransactionJob(transactionJob TransactionJob) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var TransactionID int
	var JobID int
	transactionJobUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	if transactionJob.UUID == "" {
		transactionJob.UUID = transactionJobUUID.String()
	}
	err = database.DbConn.QueryRowContext(ctx, `INSERT INTO transaction_jobs  
	(
		transaction_id,  
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
		RETURNING transaction_id, job_id`,
		transactionJob.TransactionID,    //1
		transactionJob.JobID,            //2
		transactionJob.Name,             //3
		transactionJob.AlternateName,    //4
		transactionJob.StartDate,        //5
		transactionJob.EndDate,          //6
		transactionJob.Description,      //7
		transactionJob.StatusID,         //8
		transactionJob.ResponseStatus,   //9
		transactionJob.RequestUrl,       //10
		transactionJob.RequestBody,      //11
		transactionJob.RequestMethod,    //12
		transactionJob.ResponseData,     //13
		transactionJob.ResponseDataJson, //14
		transactionJob.CreatedBy,        //15
	).Scan(&TransactionID, &JobID)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(TransactionID), int(JobID), nil
}

func InsertTransactionJobs(transactionJobs []TransactionJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range transactionJobs {
		transactionJob := transactionJobs[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(transactionJob.UUID)
		row := []interface{}{

			*transactionJob.TransactionID, //1
			*transactionJob.JobID,         //2
			uuidString,                    //3
			transactionJob.Name,           //4
			transactionJob.AlternateName,  //5
			&transactionJob.StartDate,     //6
			&transactionJob.EndDate,       //7
			transactionJob.Description,    //8
			*transactionJob.StatusID,      //9
			transactionJob.ResponseStatus, //10
			transactionJob.RequestUrl,     //11
			transactionJob.RequestBody,    //12
			transactionJob.RequestMethod,  //13
			transactionJob.ResponseData,   //14
			// TODO: erroring out in json insert look into it later TAT-27
			// transactionJob.ResponseDataJson //15
			transactionJob.CreatedBy,  //16
			&transactionJob.CreatedAt, //17
			transactionJob.CreatedBy,  //18
			&now,                      //19
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"transaction_jobs"},
		[]string{
			"transaction_id",  //1
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
