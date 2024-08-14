package step

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

func GetStep(dbConnPgx utils.PgxIface, stepID *int) (*Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		pool_id,
		parent_step_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_type_id,
		function_name,
		step_order,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM steps 
	WHERE id = $1`, *stepID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	step, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Step])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &step, nil
}

func RemoveStep(dbConnPgx utils.PgxIface, stepID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveStep DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM steps WHERE id = $1`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *stepID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetSteps(dbConnPgx utils.PgxIface, ids []int) ([]Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		pool_id,
		parent_step_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_type_id,
		function_name,
		step_order,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM steps`
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
	steps, err := pgx.CollectRows(results, pgx.RowToStructByName[Step])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return steps, nil
}

func GetStepsByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		pool_id,
		parent_step_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_type_id,
		function_name,
		step_order,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM steps
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	steps, err := pgx.CollectRows(results, pgx.RowToStructByName[Step])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return steps, nil
}

func GetStepsFromPoolID(dbConnPgx utils.PgxIface, poolID *int) ([]Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		pool_id,
		parent_step_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_type_id,
		function_name,
		step_order,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM steps
	WHERE pool_id =$1
	`, *poolID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	steps, err := pgx.CollectRows(results, pgx.RowToStructByName[Step])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return steps, nil
}

func GetStartAndEndDateDiffSteps(dbConnPgx utils.PgxIface, diffInDate *int) ([]Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		pool_id,
		parent_step_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_type_id,
		function_name,
		step_order,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM steps
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, *diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	steps, err := pgx.CollectRows(results, pgx.RowToStructByName[Step])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return steps, nil
}

func UpdateStep(dbConnPgx utils.PgxIface, step *Step) error {
	// if the step id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if step.ID == nil || *step.ID == 0 {
		return errors.New("step has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateStep DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE steps SET 
		pool_id=$1, 
		parent_step_id=$2,
		uuid=$3, 
		name=$4, 
		alternate_name=$5, 
		start_date=$6, 
		end_date=$7, 
		description=$8, 
		action_type_id=$9,
		function_name=$10,
		step_order=$11,
		updated_by=$12, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$13`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		step.PoolID,        //1
		step.ParentStepId,  //2
		step.UUID,          //3
		step.Name,          //4
		step.AlternateName, //5
		step.StartDate,     //6
		step.EndDate,       //7
		step.Description,   //8
		step.ActionTypeID,  //9
		step.FunctionName,  //10
		step.StepOrder,     //11
		step.UpdatedBy,     //12
		step.ID,            //13
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertStep(dbConnPgx utils.PgxIface, step *Step) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertStep DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO steps 
	(
		pool_id,
		parent_step_id,
		uuid,
		name,
		alternate_name,
		start_date,
		end_date,
		description,
		action_type_id,
		function_name,
		step_order,
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
			$5, 
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			current_timestamp at time zone 'UTC',
			$11,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		step.PoolID,       //1
		step.ParentStepId, //2
		// &step.UUID, //
		step.Name,          //3
		step.AlternateName, //4
		step.StartDate,     //5
		step.EndDate,       //6
		step.Description,   //7
		step.ActionTypeID,  //8
		step.FunctionName,  //9
		step.StepOrder,     //10
		step.CreatedBy,     //11
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
func InsertSteps(dbConnPgx utils.PgxIface, steps []Step) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range steps {
		step := steps[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(step.UUID)
		row := []interface{}{
			step.PoolID,        //1
			step.ParentStepId,  //2
			uuidString,         //3
			step.Name,          //4
			step.AlternateName, //5
			&step.StartDate,    //6
			&step.EndDate,      //7
			step.Description,   //8
			step.ActionTypeID,  //9
			step.FunctionName,  //10
			step.StepOrder,     //11
			step.CreatedBy,     //12
			&now,               //13
			step.CreatedBy,     //14
			&now,               //15
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"steps"},
		[]string{
			"pool_id",        //1
			"parent_step_id", //2
			"uuid",           //3
			"name",           //4
			"alternate_name", //5
			"start_date",     //6
			"end_date",       //7
			"description",    //8
			"action_type_id", //9
			"function_name",  //10
			"step_order",     //11
			"created_by",     //12
			"created_at",     //13
			"updated_by",     //14
			"updated_at",     //15
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertSteps: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetStepListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		pool_id,
		parent_step_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		description,
		action_type_id,
		function_name,
		step_order,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM steps 
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
	steps, err := pgx.CollectRows(results, pgx.RowToStructByName[Step])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return steps, nil
}

func GetTotalStepsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM steps
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
