package structuredvaluetype

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetStructuredValueType(structuredValueTypeID int) (*StructuredValueType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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
	if errors.Is(err, pgx.ErrNoRows) {
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
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM structured_value_types WHERE id = $1`, structuredValueTypeID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetStructuredValueTypeList(ids []int) ([]StructuredValueType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM structured_value_types`
	if len(ids) > 0 {
		strIds := utils.SplitToString(ids, ",")
		additionalQuery := fmt.Sprintf(` WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := database.DbConnPgx.Query(ctx, sql)
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
	_, err := database.DbConnPgx.Query(ctx, `UPDATE structured_value_types SET 
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
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO structured_value_types  
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

func GetStructuredValueListByPagination(_start, _end *int, _order, _sort string, _filters []string) ([]StructuStructuredValueTyperedValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
	iid,
	uuid, 
	name, 
	alternate_name, 
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM structured_value_types
	`
	if len(_filters) > 0 {
		sql += "WHERE "
		for i, filter := range _filters {
			sql += filter
			if i < len(_filters)-1 {
				sql += " AND "
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

	results, err := database.DbConnPgx.Query(ctx, sql)
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

func GetTotalStructuredValueCount() (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM structured_value_types
	`)
	totalCount := 0
	err := row.Scan(
		&totalCount,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &totalCount, nil
}
