package structuredvalue

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/database"
)

func GetStructuredValue(structuredValueID int) (*StructuredValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	structured_value_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM structured_values 
	WHERE id = $1`, structuredValueID)

	structuredValue := &StructuredValue{}
	err := row.Scan(
		&structuredValue.ID,
		&structuredValue.UUID,
		&structuredValue.Name,
		&structuredValue.AlternateName,
		&structuredValue.StructuredValueTypeID,
		&structuredValue.CreatedBy,
		&structuredValue.CreatedAt,
		&structuredValue.UpdatedBy,
		&structuredValue.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return structuredValue, nil
}

func GetTopTenStructuredValues() ([]StructuredValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	structured_value_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM structured_values 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	structuredValues := make([]StructuredValue, 0)
	for results.Next() {
		var structuredValue StructuredValue
		results.Scan(
			&structuredValue.ID,
			&structuredValue.UUID,
			&structuredValue.Name,
			&structuredValue.AlternateName,
			&structuredValue.StructuredValueTypeID,
			&structuredValue.CreatedBy,
			&structuredValue.CreatedAt,
			&structuredValue.UpdatedBy,
			&structuredValue.UpdatedAt,
		)

		structuredValues = append(structuredValues, structuredValue)
	}
	return structuredValues, nil
}

func RemoveStructuredValue(structuredValueID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM structured_values WHERE id = $1`, structuredValueID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetStructuredValueByStructuredValueTypeIDList(structuredValueID int) ([]StructuredValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	structured_value_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM structured_values
	WHERE structured_value_type_id =$1`, structuredValueID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	structuredValues := make([]StructuredValue, 0)
	for results.Next() {
		var structuredValue StructuredValue
		results.Scan(
			&structuredValue.ID,
			&structuredValue.UUID,
			&structuredValue.Name,
			&structuredValue.AlternateName,
			&structuredValue.StructuredValueTypeID,
			&structuredValue.CreatedBy,
			&structuredValue.CreatedAt,
			&structuredValue.UpdatedBy,
			&structuredValue.UpdatedAt,
		)

		structuredValues = append(structuredValues, structuredValue)
	}
	return structuredValues, nil
}

func GetStructuredValueList() ([]StructuredValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	structured_value_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM structured_values`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	structuredValues := make([]StructuredValue, 0)
	for results.Next() {
		var structuredValue StructuredValue
		results.Scan(
			&structuredValue.ID,
			&structuredValue.UUID,
			&structuredValue.Name,
			&structuredValue.AlternateName,
			&structuredValue.StructuredValueTypeID,
			&structuredValue.CreatedBy,
			&structuredValue.CreatedAt,
			&structuredValue.UpdatedBy,
			&structuredValue.UpdatedAt,
		)

		structuredValues = append(structuredValues, structuredValue)
	}
	return structuredValues, nil
}

func UpdateStructuredValue(structuredValue StructuredValue) error {
	// if the structuredValue id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if structuredValue.ID == nil || *structuredValue.ID == 0 {
		return errors.New("structuredValue has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE structured_values SET 
		name=$1,  
		alternate_name=$2, 
		structured_value_type_id=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$5`,
		structuredValue.Name,
		structuredValue.AlternateName,
		structuredValue.StructuredValueTypeID,
		structuredValue.UpdatedBy,
		structuredValue.ID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertStructuredValue(structuredValue StructuredValue) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO structured_values  
	(
		name, 
		uuid,
		alternate_name, 
		structured_value_type_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES ($1,uuid_generate_v4(), $2, $3,$4, current_timestamp at time zone 'UTC', $5, current_timestamp at time zone 'UTC')
		RETURNING id`,
		structuredValue.Name,
		structuredValue.AlternateName,
		structuredValue.StructuredValueTypeID,
		structuredValue.CreatedBy,
		structuredValue.CreatedBy).Scan(&insertID)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
