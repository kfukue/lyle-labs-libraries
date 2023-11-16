package sourcejob

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/kfukue/lyle-labs-libraries/database"
)

func GetSourceJob(sourceID int, jobID int) (*SourceJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	market_data_id,
	asset_id,
	uuid, 
	description,
	source_data,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM market_data_jobs 
	WHERE market_data_id = $1
	AND job_id = $2
	`, sourceID, jobID)

	sourceJob := &SourceJob{}
	err := row.Scan(
		&sourceJob.SourceID,
		&sourceJob.JobID,
		&sourceJob.UUID,
		&sourceJob.Description,
		&sourceJob.CreatedBy,
		&sourceJob.CreatedAt,
		&sourceJob.UpdatedBy,
		&sourceJob.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return sourceJob, nil
}

func GetSourceJobBySourceID(sourceID int) (*SourceJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	market_data_id,
	job_id,
	uuid, 
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM market_data_jobs 
	WHERE market_data_id = $1
	`, sourceID)

	sourceJob := &SourceJob{}
	err := row.Scan(
		&sourceJob.SourceID,
		&sourceJob.JobID,
		&sourceJob.UUID,
		&sourceJob.Description,
		&sourceJob.CreatedBy,
		&sourceJob.CreatedAt,
		&sourceJob.UpdatedBy,
		&sourceJob.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return sourceJob, nil
}

func GetTopTenSourceJobs() ([]SourceJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	market_data_id,
	job_id,
	uuid, 
	description,
	source_data,
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
	sourceJobs := make([]SourceJob, 0)
	for results.Next() {
		var sourceJob SourceJob
		results.Scan(
			&sourceJob.SourceID,
			&sourceJob.JobID,
			&sourceJob.UUID,
			&sourceJob.Description,
			&sourceJob.CreatedBy,
			&sourceJob.CreatedAt,
			&sourceJob.UpdatedBy,
			&sourceJob.UpdatedAt,
		)

		sourceJobs = append(sourceJobs, sourceJob)
	}
	return sourceJobs, nil
}

func RemoveSourceJob(sourceID int, jobID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM market_data_jobs WHERE 
	market_data_id = $1 AND job_id =$2`, sourceID, jobID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetSourceJobList() ([]SourceJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	market_data_id,
	job_id,
	uuid, 
	description,
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
	sourceJobs := make([]SourceJob, 0)
	for results.Next() {
		var sourceJob SourceJob
		results.Scan(
			&sourceJob.SourceID,
			&sourceJob.JobID,
			&sourceJob.UUID,
			&sourceJob.Description,
			&sourceJob.CreatedBy,
			&sourceJob.CreatedAt,
			&sourceJob.UpdatedBy,
			&sourceJob.UpdatedAt,
		)

		sourceJobs = append(sourceJobs, sourceJob)
	}
	return sourceJobs, nil
}

func UpdateSourceJob(sourceJob SourceJob) error {
	// if the sourceJob id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (sourceJob.SourceID == nil || *sourceJob.SourceID == 0) || (sourceJob.JobID == nil || *sourceJob.JobID == 0) {
		return errors.New("sourceJob has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE market_data_jobs SET 
		description=$1,
		updated_by=$2, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE market_data_id=$7 AND job_id=$8`,
		sourceJob.Description,
		sourceJob.UpdatedBy,
		sourceJob.SourceID,
		sourceJob.JobID,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertSourceJob(sourceJob SourceJob) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var SourceID int
	var JobID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO market_data_jobs  
	(
		market_data_id,
		job_id,
		uuid, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES ($1,
			$2,
			uuid_generate_v4(),
			$3,
			$4,
			current_timestamp at time zone 'UTC',
			$4,
			current_timestamp at time zone 'UTC'
		)
		RETURNING market_data_id, job_id`,
		sourceJob.SourceID,    // 1
		sourceJob.JobID,       // 2
		sourceJob.Description, //3
		sourceJob.CreatedBy,   //4
	).Scan(&SourceID, &JobID)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(SourceID), int(JobID), nil
}
