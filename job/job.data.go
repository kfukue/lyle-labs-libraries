package job

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
)

func GetJob(dbConnPgx utils.PgxIface, jobID *int) (*Job, error) {
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
	WHERE id = $1`, *jobID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	job, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Job])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &job, nil
}

func RemoveJob(dbConnPgx utils.PgxIface, jobID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM jobs WHERE id = $1`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *jobID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetJobList(dbConnPgx utils.PgxIface, ids []int) ([]Job, error) {
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
		additionalQuery := fmt.Sprintf(` WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	jobList, err := pgx.CollectRows(results, pgx.RowToStructByName[Job])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return jobList, nil
}

func GetJobListByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]Job, error) {
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
	jobList, err := pgx.CollectRows(results, pgx.RowToStructByName[Job])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return jobList, nil
}

func GetStartAndEndDateDiffJobList(dbConnPgx utils.PgxIface, diffInDate *int) ([]Job, error) {
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
	`, *diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	jobList, err := pgx.CollectRows(results, pgx.RowToStructByName[Job])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return jobList, nil
}

func UpdateJob(dbConnPgx utils.PgxIface, job *Job) error {
	// if the job id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if job.ID == nil || *job.ID == 0 {
		return errors.New("job has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateJob DbConn.Begin   %s", err.Error())
		return err
	}
	layoutPostgres := utils.LayoutPostgres
	startDate := job.StartDate
	endDate := job.EndDate
	log.Println(fmt.Sprintf("Updating start: %s, end : %s", startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)))
	sql := `UPDATE jobs SET 
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
		WHERE id=$15`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
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
		job.ID,                               //15
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertJob(dbConnPgx utils.PgxIface, job *Job) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertJob DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var insertID int
	var jobUUID string
	layoutPostgres := utils.LayoutPostgres
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO jobs 
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
		RETURNING id, uuid`,
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
	).Scan(&insertID, &jobUUID)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, "", err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, "", err
	}
	return int(insertID), jobUUID, nil
}

func InsertJobList(dbConnPgx utils.PgxIface, jobList []Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range jobList {
		job := jobList[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(job.UUID)
		row := []interface{}{
			uuidString,         //1
			job.Name,           //2
			job.AlternateName,  //3
			job.StartDate,      //4
			job.EndDate,        //5
			job.Description,    //6
			job.StatusID,       //7
			job.ResponseStatus, //8
			job.RequestUrl,     //9
			job.RequestBody,    //10
			job.RequestMethod,  //11
			job.ResponseData,   //12
			// TODO: erroring out in json insert look into it later TAT-27
			// job.ResponseDataJson                 //13
			nil,               //13
			job.JobCategoryID, //14
			job.CreatedBy,     //15
			&now,              //16
			job.CreatedBy,     //17
			&now,              //18
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"jobs"},
		[]string{
			"uuid",               //1
			"name",               //2
			"alternate_name",     //3
			"start_date",         //4
			"end_date",           //5
			"description",        //6
			"status_id",          //7
			"response_status",    //8
			"request_url",        //9
			"request_body",       //10
			"request_method",     //11
			"response_data",      //12
			"response_data_json", //13
			"job_category_id",    //14
			"created_by",         //15
			"created_at",         //16
			"updated_by",         //17
			"updated_at",         //18
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertJobList: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetJobListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Job, error) {
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
	FROM jobs 
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
	jobs, err := pgx.CollectRows(results, pgx.RowToStructByName[Job])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return jobs, nil
}

func GetTotalJobsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM jobs
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
