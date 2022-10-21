package step

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetStep(stepID int) (*Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	WHERE id = $1`, stepID)

	step := &Step{}
	err := row.Scan(
		&step.ID,
		&step.PoolID,
		&step.ParentStepId,
		&step.Name,
		&step.UUID,
		&step.AlternateName,
		&step.StartDate,
		&step.EndDate,
		&step.Description,
		&step.ActionTypeID,
		&step.FunctionName,
		&step.StepOrder,
		&step.CreatedBy,
		&step.CreatedAt,
		&step.UpdatedBy,
		&step.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return step, nil
}

func GetTopTenStrategies() ([]Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	steps := make([]Step, 0)
	for results.Next() {
		var step Step
		results.Scan(
			&step.ID,
			&step.PoolID,
			&step.ParentStepId,
			&step.Name,
			&step.UUID,
			&step.AlternateName,
			&step.StartDate,
			&step.EndDate,
			&step.Description,
			&step.ActionTypeID,
			&step.FunctionName,
			&step.StepOrder,
			&step.CreatedBy,
			&step.CreatedAt,
			&step.UpdatedBy,
			&step.UpdatedAt,
		)

		steps = append(steps, step)
	}
	return steps, nil
}

func RemoveStep(stepID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM steps WHERE id = $1`, stepID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetSteps(ids []int) ([]Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
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
	FROM steps`
	if len(ids) > 0 {
		strIds := utils.SplitToString(ids, ",")
		additionalQuery := fmt.Sprintf(`WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := database.DbConn.QueryContext(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	steps := make([]Step, 0)
	for results.Next() {
		var step Step
		results.Scan(
			&step.ID,
			&step.PoolID,
			&step.ParentStepId,
			&step.Name,
			&step.UUID,
			&step.AlternateName,
			&step.StartDate,
			&step.EndDate,
			&step.Description,
			&step.ActionTypeID,
			&step.FunctionName,
			&step.StepOrder,
			&step.CreatedBy,
			&step.CreatedAt,
			&step.UpdatedBy,
			&step.UpdatedAt,
		)

		steps = append(steps, step)
	}
	return steps, nil
}

func GetStepsByUUIDs(UUIDList []string) ([]Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	steps := make([]Step, 0)
	for results.Next() {
		var step Step
		results.Scan(
			&step.ID,
			&step.PoolID,
			&step.ParentStepId,
			&step.Name,
			&step.UUID,
			&step.AlternateName,
			&step.StartDate,
			&step.EndDate,
			&step.Description,
			&step.ActionTypeID,
			&step.FunctionName,
			&step.StepOrder,
			&step.CreatedBy,
			&step.CreatedAt,
			&step.UpdatedBy,
			&step.UpdatedAt,
		)

		steps = append(steps, step)
	}
	return steps, nil
}

func GetStepsFromPoolID(poolID *int) ([]Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	WHERE pool_id =$1
	`, *poolID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	steps := make([]Step, 0)
	for results.Next() {
		var step Step
		results.Scan(
			&step.ID,
			&step.PoolID,
			&step.ParentStepId,
			&step.Name,
			&step.UUID,
			&step.AlternateName,
			&step.StartDate,
			&step.EndDate,
			&step.Description,
			&step.ActionTypeID,
			&step.FunctionName,
			&step.StepOrder,
			&step.CreatedBy,
			&step.CreatedAt,
			&step.UpdatedBy,
			&step.UpdatedAt,
		)

		steps = append(steps, step)
	}
	return steps, nil
}

func GetStartAndEndDateDiffSteps(diffInDate int) ([]Step, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	steps := make([]Step, 0)
	for results.Next() {
		var step Step
		results.Scan(
			&step.ID,
			&step.PoolID,
			&step.ParentStepId,
			&step.Name,
			&step.UUID,
			&step.AlternateName,
			&step.StartDate,
			&step.EndDate,
			&step.Description,
			&step.ActionTypeID,
			&step.FunctionName,
			&step.StepOrder,
			&step.CreatedBy,
			&step.CreatedAt,
			&step.UpdatedBy,
			&step.UpdatedAt,
		)

		steps = append(steps, step)
	}
	return steps, nil
}

func UpdateStep(step Step) error {
	// if the step id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if step.ID == nil || *step.ID == 0 {
		return errors.New("step has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE steps SET 
		pool_id=$1, 
		parent_step_id=$2,
		name=$3, 
		uuid=$4, 
		alternate_name=$5, 
		start_date=$6, 
		end_date=$7, 
		description=$8, 
		action_type_id=$9,
		function_name=$10,
		step_order=$11,
		updated_by=$12, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$13`,
		step.PoolID,        //1
		step.ParentStepId,  //2
		step.Name,          //3
		step.UUID,          //4
		step.AlternateName, //5
		step.StartDate,     //6
		step.EndDate,       //7
		step.Description,   //8
		step.ActionTypeID,  //9
		step.FunctionName,  //10
		step.StepOrder,     //11
		step.UpdatedBy,     //12
		step.ID)            //13
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertStep(step Step) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	// layoutPostgres := utils.LayoutPostgres
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO steps 
	(
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
		) VALUES (
			$1,
			$2, 
			$3, 
			uuid_generate_v4(), 
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
		&step.PoolID,       //1
		&step.ParentStepId, //2
		&step.Name,         //3
		// &step.UUID, //
		&step.AlternateName, //4
		&step.StartDate,     //5
		&step.EndDate,       //6
		&step.Description,   //7
		&step.ActionTypeID,  //8
		&step.FunctionName,  //9
		&step.StepOrder,     //10
		&step.CreatedBy,     //11
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertSteps(steps []Step) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range steps {
		step := steps[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(step.UUID)
		row := []interface{}{
			*step.PoolID,       //1
			*step.ParentStepId, //2
			step.Name,          //3
			uuidString,         //4
			step.AlternateName, //5
			&step.StartDate,    //6
			&step.EndDate,      //7
			step.Description,   //8
			*step.ActionTypeID, //9
			step.FunctionName,  //10
			*step.StepOrder,    //11
			step.CreatedBy,     //12
			&now,               //13
			step.CreatedBy,     //14
			&now,               //15
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"steps"},
		[]string{
			"pool_id",        //1
			"parent_step_id", //2
			"name",           //3
			"uuid",           //4
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
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}

	return nil
}
