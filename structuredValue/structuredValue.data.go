package structuredvalue

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

func GetStructuredValue(dbConnPgx utils.PgxIface, structuredValueID *int) (*StructuredValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
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
	WHERE id = $1`, *structuredValueID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	structuredValue, err := pgx.CollectOneRow(row, pgx.RowToStructByName[StructuredValue])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &structuredValue, nil
}

func RemoveStructuredValue(dbConnPgx utils.PgxIface, structuredValueID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveStructuredValue DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM structured_values WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *structuredValueID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetStructuredValueByStructuredValueTypeIDList(dbConnPgx utils.PgxIface, structuredValueTypeID *int) ([]StructuredValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT 
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
	WHERE structured_value_type_id =$1`, *structuredValueTypeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	structuredValues, err := pgx.CollectRows(results, pgx.RowToStructByName[StructuredValue])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return structuredValues, nil
}

func GetStructuredValueList(dbConnPgx utils.PgxIface, ids []int) ([]StructuredValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT  
	id,
	uuid, 
	name, 
	alternate_name, 
	structured_value_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM structured_values`
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
	structuredValues, err := pgx.CollectRows(results, pgx.RowToStructByName[StructuredValue])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return structuredValues, nil
}

func UpdateStructuredValue(dbConnPgx utils.PgxIface, structuredValue *StructuredValue) error {
	// if the structuredValue id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if structuredValue.ID == nil || *structuredValue.ID == 0 {
		return errors.New("structuredValue has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateStructuredValue DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE structured_values SET 
		name=$1,  
		alternate_name=$2, 
		structured_value_type_id=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$5`

	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		structuredValue.Name,                  //1
		structuredValue.AlternateName,         //2
		structuredValue.StructuredValueTypeID, //3
		structuredValue.UpdatedBy,             //4
		structuredValue.ID,                    //5
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)

}

func InsertStructuredValue(dbConnPgx utils.PgxIface, structuredValue *StructuredValue) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertStructuredValue DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO structured_values  
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
		structuredValue.Name,                  //1
		structuredValue.AlternateName,         //2
		structuredValue.StructuredValueTypeID, //3
		structuredValue.CreatedBy,             //4
		structuredValue.CreatedBy,             //5
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
func InsertStructuredValues(dbConnPgx utils.PgxIface, structuredValues []StructuredValue) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range structuredValues {
		structuredValue := structuredValues[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(structuredValue.UUID)
		row := []interface{}{
			uuidString,                            //1
			structuredValue.Name,                  //2
			structuredValue.AlternateName,         //3
			structuredValue.StructuredValueTypeID, //4
			structuredValue.CreatedBy,             //5
			&now,                                  //6
			structuredValue.CreatedBy,             //7
			&now,                                  //8
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"structured_values"},
		[]string{
			"uuid",                     //1
			"name",                     //2
			"alternate_name",           //3
			"structured_value_type_id", //4
			"created_by",               //5
			"created_at",               //6
			"updated_by",               //7
			"updated_at",               //8
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

func GetStructuredValueListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]StructuredValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
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
	structuredValues, err := pgx.CollectRows(results, pgx.RowToStructByName[StructuredValue])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return structuredValues, nil
}

func GetTotalStructuredValuesCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM structured_values
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
