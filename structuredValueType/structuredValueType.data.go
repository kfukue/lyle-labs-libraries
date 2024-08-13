package structuredvaluetype

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
)

func GetStructuredValueType(dbConnPgx utils.PgxIface, structuredValueTypeID *int) (*StructuredValueType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid, 
		name, 
		alternate_name, 
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM structured_value_types 
	WHERE id = $1`, *structuredValueTypeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	structuredValueType, err := pgx.CollectOneRow(row, pgx.RowToStructByName[StructuredValueType])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &structuredValueType, nil
}

func RemoveStructuredValueType(dbConnPgx utils.PgxIface, structuredValueTypeID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveStructuredValueType DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM structured_value_types WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *structuredValueTypeID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetStructuredValueTypeList(dbConnPgx utils.PgxIface, ids []int) ([]StructuredValueType, error) {
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
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	structuredValueTypes, err := pgx.CollectRows(results, pgx.RowToStructByName[StructuredValueType])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return structuredValueTypes, nil
}

func UpdateStructuredValueType(dbConnPgx utils.PgxIface, structuredValueType *StructuredValueType) error {
	// if the structuredValueType id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if structuredValueType.ID == nil || *structuredValueType.ID == 0 {
		return errors.New("structuredValueType has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateStructuredValueType DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE structured_value_types SET 
		name=$1,  
		alternate_name=$2, 
		updated_by=$3, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$4`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		structuredValueType.Name,          //1
		structuredValueType.AlternateName, //2
		structuredValueType.UpdatedBy,     //3
		structuredValueType.ID,            //4
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertStructuredValueType(dbConnPgx utils.PgxIface, structuredValueType *StructuredValueType) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertStructuredValueType DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO structured_value_types  
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
		structuredValueType.Name,          //1
		structuredValueType.AlternateName, //2
		structuredValueType.CreatedBy,     //3
		structuredValueType.CreatedBy,     //4
	).Scan(&insertID)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, err
	}
	return int(insertID), nil
}

func InsertStructuredValueTypes(dbConnPgx utils.PgxIface, structuredValueTypes []StructuredValueType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range structuredValueTypes {
		structuredValueType := structuredValueTypes[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(structuredValueType.UUID)
		row := []interface{}{
			uuidString,                        //1
			structuredValueType.Name,          //2
			structuredValueType.AlternateName, //3
			structuredValueType.CreatedBy,     //4
			&now,                              //5
			structuredValueType.CreatedBy,     //6
			&now,                              //7
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"structured_value_types"},
		[]string{
			"uuid",           //1
			"name",           //2
			"alternate_name", //3
			"created_by",     //4
			"created_at",     //5
			"updated_by",     //6
			"updated_at",     //7
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertStructuredValues: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func GetStructuredValueTypeListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]StructuredValueType, error) {
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

	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	structuredValueTypes, err := pgx.CollectRows(results, pgx.RowToStructByName[StructuredValueType])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return structuredValueTypes, nil
}

func GetTotalStructuredValueTypeCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM structured_value_types
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
