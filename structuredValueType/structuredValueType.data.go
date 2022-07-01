package structuredvaluetype

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/database"
)

func GetStructuredValueType(structuredValueTypeID int) (*StructuredValueType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM structured_value_types 
	WHERE id = $1`, structuredValueTypeID)

	structuredValueType := &StructuredValueType{}
	err := row.Scan(
		&structuredValueType.ID,
		&structuredValueType.UUID,
		&structuredValueType.Name,
		&structuredValueType.AlternateName,
		&structuredValueType.CreatedBy,
		&structuredValueType.CreatedAt,
		&structuredValueType.UpdatedBy,
		&structuredValueType.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return structuredValueType, nil
}

func GetTopTenStructuredValueTypes() ([]StructuredValueType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM structured_value_types
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	structuredValueTypes := make([]StructuredValueType, 0)
	for results.Next() {
		var structuredValueType StructuredValueType
		results.Scan(
			&structuredValueType.ID,
			&structuredValueType.UUID,
			&structuredValueType.Name,
			&structuredValueType.AlternateName,
			&structuredValueType.CreatedBy,
			&structuredValueType.CreatedAt,
			&structuredValueType.UpdatedBy,
			&structuredValueType.UpdatedAt,
		)

		structuredValueTypes = append(structuredValueTypes, structuredValueType)
	}
	return structuredValueTypes, nil
}

func RemoveStructuredValueType(structuredValueTypeID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM structured_value_types WHERE id = $1`, structuredValueTypeID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetStructuredValueTypeList() ([]StructuredValueType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM structured_value_types`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	structuredValueTypes := make([]StructuredValueType, 0)
	for results.Next() {
		var structuredValueType StructuredValueType
		results.Scan(
			&structuredValueType.ID,
			&structuredValueType.UUID,
			&structuredValueType.Name,
			&structuredValueType.AlternateName,
			&structuredValueType.CreatedBy,
			&structuredValueType.CreatedAt,
			&structuredValueType.UpdatedBy,
			&structuredValueType.UpdatedAt,
		)

		structuredValueTypes = append(structuredValueTypes, structuredValueType)
	}
	return structuredValueTypes, nil
}

func UpdateStructuredValueType(structuredValueType StructuredValueType) error {
	// if the structuredValueType id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if structuredValueType.ID == nil || *structuredValueType.ID == 0 {
		return errors.New("structuredValueType has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE structured_value_types SET 
		name=$1,  
		alternate_name=$2, 
		updated_by=$3, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$4`,
		structuredValueType.Name,
		structuredValueType.AlternateName,
		structuredValueType.UpdatedBy,
		structuredValueType.ID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertStructuredValueType(structuredValueType StructuredValueType) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO structured_value_types  
	(
		name, 
		uuid,
		alternate_name, 
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES ($1,uuid_generate_v4(), $2, $3, current_timestamp at time zone 'UTC', $4, current_timestamp at time zone 'UTC')
		RETURNING id`,
		structuredValueType.Name,
		structuredValueType.AlternateName,
		structuredValueType.CreatedBy,
		structuredValueType.CreatedBy).Scan(&insertID)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
