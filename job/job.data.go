package job

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetJob(jobID int) (*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	FROM jobs 
	WHERE id = $1`, jobID)

	job := &Job{}
	err := row.Scan(
		&job.ID,
		&job.UUID,
		&job.Name,
		&job.AlternateName,
		&job.StartDate,
		&job.EndDate,
		&job.Description,
		&job.StatusID,
		&job.ResponseStatus,
		&job.RequestUrl,
		&job.RequestBody,
		&job.RequestMethod,
		&job.ResponseData,
		&job.ResponseDataJson,
		&job.JobCategoryID,
		&job.CreatedBy,
		&job.CreatedAt,
		&job.UpdatedBy,
		&job.UpdatedAt,
		&job.ResponseDataJson,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return job, nil
}

func GetTopTenJobs() ([]Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	FROM jobs 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	jobList := make([]Job, 0)
	for results.Next() {
		var job Job
		results.Scan(
			&job.ID,
			&job.UUID,
			&job.Name,
			&job.AlternateName,
			&job.StartDate,
			&job.EndDate,
			&job.Description,
			&job.StatusID,
			&job.ResponseStatus,
			&job.RequestUrl,
			&job.RequestBody,
			&job.RequestMethod,
			&job.ResponseData,
			&job.ResponseDataJson,
			&job.JobCategoryID,
			&job.CreatedBy,
			&job.CreatedAt,
			&job.UpdatedBy,
			&job.UpdatedAt,
			&job.ResponseDataJson,
		)

		jobList = append(jobList, job)
	}
	return jobList, nil
}

func RemoveJob(jobID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM jobs WHERE id = $1`, jobID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetJobList(ids []int) ([]Job, error) {
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
	FROM jobs`
	if len(ids) > 0 {
		strIds := utils.SplitToString(ids, ",")
		additionalQuery := fmt.Sprintf(`WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := database.DbConn.QueryContext(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	jobList := make([]Job, 0)
	for results.Next() {
		var job Job
		results.Scan(
			&job.ID,
			&job.UUID,
			&job.Name,
			&job.AlternateName,
			&job.StartDate,
			&job.EndDate,
			&job.Description,
			&job.StatusID,
			&job.ResponseStatus,
			&job.RequestUrl,
			&job.RequestBody,
			&job.RequestMethod,
			&job.ResponseData,
			&job.ResponseDataJson,
			&job.JobCategoryID,
			&job.CreatedBy,
			&job.CreatedAt,
			&job.UpdatedBy,
			&job.UpdatedAt,
			&job.ResponseDataJson,
		)

		jobList = append(jobList, job)
	}
	return jobList, nil
}

func GetJobListByUUIDs(UUIDList []string) ([]Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	FROM jobs
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	jobList := make([]Job, 0)
	for results.Next() {
		var job Job
		results.Scan(
			&job.ID,
			&job.UUID,
			&job.Name,
			&job.AlternateName,
			&job.StartDate,
			&job.EndDate,
			&job.Description,
			&job.StatusID,
			&job.ResponseStatus,
			&job.RequestUrl,
			&job.RequestBody,
			&job.RequestMethod,
			&job.ResponseData,
			&job.ResponseDataJson,
			&job.JobCategoryID,
			&job.CreatedBy,
			&job.CreatedAt,
			&job.UpdatedBy,
			&job.UpdatedAt,
			&job.ResponseDataJson,
		)

		jobList = append(jobList, job)
	}
	return jobList, nil
}

func GetStartAndEndDateDiffJobList(diffInDate int) ([]Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	updated_at
	FROM jobs
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	jobList := make([]Job, 0)
	for results.Next() {
		var job Job
		results.Scan(
			&job.ID,
			&job.UUID,
			&job.Name,
			&job.AlternateName,
			&job.StartDate,
			&job.EndDate,
			&job.Description,
			&job.StatusID,
			&job.ResponseStatus,
			&job.RequestUrl,
			&job.RequestBody,
			&job.RequestMethod,
			&job.ResponseData,
			&job.ResponseDataJson,
			&job.JobCategoryID,
			&job.CreatedBy,
			&job.CreatedAt,
			&job.UpdatedBy,
			&job.UpdatedAt,
			&job.ResponseDataJson,
		)

		jobList = append(jobList, job)
	}
	return jobList, nil
}

func UpdateJob(job Job) error {
	// if the job id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if job.ID == nil || *job.ID == 0 {
		return errors.New("job has invalid ID")
	}
	layoutPostgres := utils.LayoutPostgres
	startDate := job.StartDate
	endDate := job.EndDate
	log.Println(fmt.Sprintf("Updating start: %s, end : %s", startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)))
	_, err := database.DbConn.ExecContext(ctx, `UPDATE jobs SET 
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
		job_category_id=$13,
		updated_by=$14, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$15`,
		job.Name,                             //1
		job.AlternateName,                    //2
		job.StartDate.Format(layoutPostgres), //3
		job.EndDate.Format(layoutPostgres),   //4
		job.Description,                      //5
		job.StatusID,                         //6
		job.ResponseStatus,                   //7
		job.RequestUrl,                       //8
		job.RequestBody,                      //9
		job.RequestMethod,                    //10
		job.ResponseData,                     //11
		job.ResponseDataJson,                 //12
		job.JobCategoryID,                    //13
		job.UpdatedBy,                        //14
		job.ID)                               //15
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertJob(job Job) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	layoutPostgres := utils.LayoutPostgres
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO jobs 
	(
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
		updated_at
		) VALUES (
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
			$13,
			$14,
			$15,
			current_timestamp at time zone 'UTC',
			$15,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		job.UUID,                             //1
		job.Name,                             //2
		job.AlternateName,                    //3
		job.StartDate.Format(layoutPostgres), //4
		job.EndDate.Format(layoutPostgres),   //5
		job.Description,                      //6
		job.StatusID,                         //7
		job.ResponseStatus,                   //8
		job.RequestUrl,                       //9
		job.RequestBody,                      //10
		job.RequestMethod,                    //11
		job.ResponseData,                     //12
		job.ResponseDataJson,                 //13
		job.JobCategoryID,                    //14
		job.CreatedBy,                        //15
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

func InsertJobListManual(jobList []Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range jobList {
		job := jobList[i]

		uuidString := &pgtype.UUID{}
		uuidString.Set(job.UUID)
		row := []interface{}{
			job.Name,           //1
			uuidString,         //2
			job.AlternateName,  //3
			&job.StartDate,     //4
			&job.EndDate,       //5
			job.Description,    //6
			*job.StatusID,      //7
			job.ResponseStatus, //8
			job.RequestUrl,     //9
			job.RequestBody,    //10
			job.RequestMethod,  //11
			job.ResponseData,   //12
			// TODO: erroring out in json insert look into it later TAT-27
			// job.ResponseDataJson                 //13
			*job.JobCategoryID, //14
			job.CreatedBy,      //15
			&now,               //16
			job.CreatedBy,      //17
			&now,               //18
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"jobs"},
		[]string{
			"name",            //1
			"uuid",            //2
			"alternate_name",  //3
			"start_date",      //4
			"end_date",        //5
			"description",     //6
			"status_id",       //7
			"response_status", //8
			"request_url",     //9
			"request_body",    //10
			"request_method",  //11
			"response_data",   //12
			// "response_data_json", //13
			"job_category_id", //14
			"created_by",      //15
			"created_at",      //16
			"updated_by",      //17
			"updated_at",      //18
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

func InsertJobList(jobList []Job) error {
	txn, err := database.DbConn.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	layoutPostgres := utils.LayoutPostgres
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn(
		"jobs",
		"name",            //1
		"uuid",            //2
		"alternate_name",  //3
		"start_date",      //4
		"end_date",        //5
		"description",     //6
		"status_id",       //7
		"response_status", //8
		"request_url",     //9
		"request_body",    //10
		"request_method",  //11
		"response_data",   //12
		// "response_data_json", //13
		"job_category_id", //14
		"created_by",      //15
		"created_at",      //16
		"updated_by",      //17
		"updated_at",      //18
	))
	if err != nil {
		log.Fatal(err)
	}
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	for _, job := range jobList {
		_, err = stmt.Exec(
			job.Name,                             //1
			job.UUID,                             //2
			job.AlternateName,                    //3
			job.StartDate.Format(layoutPostgres), //4
			job.EndDate.Format(layoutPostgres),   //5
			job.Description,                      //6
			job.StatusID,                         //7
			job.ResponseStatus,                   //8
			job.RequestUrl,                       //9
			job.RequestBody,                      //10
			job.RequestMethod,                    //11
			job.ResponseData,                     //12
			// TODO: erroring out in json insert look into it later TAT-27
			// job.ResponseDataJson,                 //13
			job.JobCategoryID,          //14
			job.CreatedBy,              //15
			now.Format(layoutPostgres), //16
			job.CreatedBy,              //17
			now.Format(layoutPostgres), //18
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
