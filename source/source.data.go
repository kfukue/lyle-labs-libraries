package source

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/database"
)

func GetSource(sourceID int) (*Source, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	url,
	ticker,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM sources 
	WHERE id = $1`, sourceID)

	source := &Source{}
	err := row.Scan(
		&source.ID,
		&source.UUID,
		&source.Name,
		&source.AlternateName,
		&source.URL,
		&source.Ticker,
		&source.Description,
		&source.CreatedBy,
		&source.CreatedAt,
		&source.UpdatedBy,
		&source.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return source, nil
}

func GetTopTenSources() ([]Source, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	url,
	ticker,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM sources 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	sources := make([]Source, 0)
	for results.Next() {
		var source Source
		results.Scan(
			&source.ID,
			&source.UUID,
			&source.Name,
			&source.AlternateName,
			&source.URL,
			&source.Ticker,
			&source.Description,
			&source.CreatedBy,
			&source.CreatedAt,
			&source.UpdatedBy,
			&source.UpdatedAt,
		)

		sources = append(sources, source)
	}
	return sources, nil
}

func RemoveSource(sourceID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM sources WHERE id = $1`, sourceID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetSourceList() ([]Source, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	url,
	ticker,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM sources`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	sources := make([]Source, 0)
	for results.Next() {
		var source Source
		results.Scan(
			&source.ID,
			&source.UUID,
			&source.Name,
			&source.AlternateName,
			&source.URL,
			&source.Ticker,
			&source.Description,
			&source.CreatedBy,
			&source.CreatedAt,
			&source.UpdatedBy,
			&source.UpdatedAt,
		)

		sources = append(sources, source)
	}
	return sources, nil
}

func UpdateSource(source Source) error {
	// if the source id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if source.ID == nil || *source.ID == 0 {
		return errors.New("source has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE sources SET 
		name=$1,  
		alternate_name=$2, 
		url=$3,
		ticker=$4,
		description=$5,
		updated_by=$6, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$7`,
		source.Name,
		source.AlternateName,
		source.URL,
		source.Ticker,
		source.Description,
		source.UpdatedBy,
		source.ID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertSource(source Source) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO sources  
	(
		name, 
		uuid,
		alternate_name, 
		url,
		ticker,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES ($1,uuid_generate_v4(),
			$2, $3, 
			$4, 
			$5, $6, 
			current_timestamp at time zone 'UTC', $7, 
			current_timestamp at time zone 'UTC')
		RETURNING id`,
		source.Name,
		source.AlternateName,
		source.URL,
		source.Ticker,
		source.Description,
		source.CreatedBy,
		source.CreatedBy).Scan(&insertID)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
