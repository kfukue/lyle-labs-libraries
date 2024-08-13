package sourcejob

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

func GetSourceJob(dbConnPgx utils.PgxIface, sourceID, jobID *int) (*SourceJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		source_id,
		job_id,
		uuid, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM source_jobs 
	WHERE source_id = $1
	AND job_id = $2
	`, *sourceID, *jobID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	sourceJob, err := pgx.CollectOneRow(row, pgx.RowToStructByName[SourceJob])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &sourceJob, nil
}

func GetSourceJobBySourceID(dbConnPgx utils.PgxIface, sourceID *int) (*SourceJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
	source_id,
	job_id,
	uuid, 
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM source_jobs 
	WHERE source_id = $1
	`, *sourceID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	sourceJob, err := pgx.CollectOneRow(row, pgx.RowToStructByName[SourceJob])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &sourceJob, nil
}

func RemoveSourceJob(dbConnPgx utils.PgxIface, sourceID, jobID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemovePositionJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM source_jobs WHERE 
		source_id = $1 AND job_id =$2`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *sourceID, *jobID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetSourceJobList(dbConnPgx utils.PgxIface) ([]SourceJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	source_id,
	job_id,
	uuid, 
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM source_jobs`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	sourceJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[SourceJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return sourceJobs, nil
}

func UpdateSourceJob(dbConnPgx utils.PgxIface, sourceJob *SourceJob) error {
	// if the sourceJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (sourceJob.SourceID == nil || *sourceJob.SourceID == 0) || (sourceJob.JobID == nil || *sourceJob.JobID == 0) {
		return errors.New("sourceJob has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdatePositionJob DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE source_jobs SET 
		description=$1,
		updated_by=$2, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE source_id=$7 AND job_id=$8`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		sourceJob.Description, //1
		sourceJob.UpdatedBy,   //2
		sourceJob.SourceID,    //3
		sourceJob.JobID,       //4
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)

}

func InsertSourceJob(dbConnPgx utils.PgxIface, sourceJob *SourceJob) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertSourceJob DbConn.Begin   %s", err.Error())
		return -1, -1, err
	}
	var SourceID int
	var JobID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO source_jobs  
	(
		source_id,
		job_id,
		uuid, 
		description,
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
			current_timestamp at time zone 'UTC',
			$4,
			current_timestamp at time zone 'UTC'
		)
		RETURNING source_id, job_id`,
		sourceJob.SourceID,    //1
		sourceJob.JobID,       //2
		sourceJob.Description, //3
		sourceJob.CreatedBy,   //4
	).Scan(&SourceID, &JobID)
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
	return int(SourceID), int(JobID), nil
}

func InsertSourceJobs(dbConnPgx utils.PgxIface, sourceJobs []SourceJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range sourceJobs {
		sourceJob := sourceJobs[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(sourceJob.UUID)
		row := []interface{}{
			sourceJob.SourceID,    //1
			sourceJob.JobID,       //2
			uuidString,            //3
			sourceJob.Description, //4
			sourceJob.CreatedBy,   //5
			&sourceJob.CreatedAt,  //6
			sourceJob.UpdatedBy,   //7
			&now,                  //8
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"source_jobs"},
		[]string{
			"source_id",   //1
			"job_id",      //2
			"uuid",        //3
			"description", //4
			"created_by",  //5
			"created_at",  //6
			"updated_by",  //7
			"updated_at",  //8
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertSourceJobs: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetSourceJobListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]SourceJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		source_id,
		job_id,
		uuid, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM source_jobs 
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
	sourceJobs, err := pgx.CollectRows(results, pgx.RowToStructByName[SourceJob])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return sourceJobs, nil
}

func GetTotalSourceJobsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM source_jobs
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
