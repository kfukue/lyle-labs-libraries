package transactionstep

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
)

func GetTransactionStep(dbConnPgx utils.PgxIface, transactionID, stepID *int) (*TransactionStep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		transaction_id,  
		step_id,
		uuid, 
		name,
		alternate_name,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM transaction_steps 
	WHERE transaction_id = $1
	AND step_id = $2
	`, *transactionID, *stepID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	transactionStep, err := pgx.CollectOneRow(row, pgx.RowToStructByName[TransactionStep])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &transactionStep, nil
}

func GetTransactionStepByUUID(dbConnPgx utils.PgxIface, transactionStepUUID string) (*TransactionStep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
	transaction_id,  
	step_id,
	uuid, 
	name,
	alternate_name,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM transaction_steps 
	WHERE text(uuid) = $1
	`, transactionStepUUID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	transactionStep, err := pgx.CollectOneRow(row, pgx.RowToStructByName[TransactionStep])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &transactionStep, nil
}

func GetTransactionStepsByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]TransactionStep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	transaction_id,  
	step_id,
	uuid, 
	name,
	alternate_name,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM transaction_steps
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	transactionStepList, err := pgx.CollectRows(results, pgx.RowToStructByName[TransactionStep])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionStepList, nil
}

func RemoveTransactionStep(dbConnPgx utils.PgxIface, transactionID, stepID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveTransactionStep DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM transaction_steps WHERE 
		transaction_id = $1 AND step_id =$2`

	if _, err := dbConnPgx.Exec(ctx, sql, *transactionID, *stepID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemoveTransactionStepByUUID(dbConnPgx utils.PgxIface, transactionStepUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveTransactionStepByUUID DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM transaction_steps WHERE text(uuid) = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, transactionStepUUID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetTransactionStepList(dbConnPgx utils.PgxIface) ([]TransactionStep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	transaction_id,  
	step_id,
	uuid, 
	name,
	alternate_name,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM transaction_steps`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	transactionSteps, err := pgx.CollectRows(results, pgx.RowToStructByName[TransactionStep])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionSteps, nil
}

func UpdateTransactionStep(dbConnPgx utils.PgxIface, transactionStep *TransactionStep) error {
	// if the transactionStep id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (transactionStep.TransactionID == nil || *transactionStep.TransactionID == 0) || (transactionStep.StepID == nil || *transactionStep.StepID == 0) {
		return errors.New("transactionStep has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateTransactionStep DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE transaction_steps SET 
		name=$1,
		alternate_name=$2,
		description=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE transaction_id=$5 AND step_id=$6`

	if _, err := dbConnPgx.Exec(ctx, sql,
		transactionStep.Name,          //1
		transactionStep.AlternateName, //2
		transactionStep.Description,   //3
		transactionStep.UpdatedBy,     //4
		transactionStep.TransactionID, //5
		transactionStep.StepID,        //6
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func UpdateTransactionStepByUUID(dbConnPgx utils.PgxIface, transactionStep *TransactionStep) error {
	// if the transactionStep id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (transactionStep.TransactionID == nil || *transactionStep.TransactionID == 0) || (transactionStep.StepID == nil || *transactionStep.StepID == 0) {
		return errors.New("transactionStep has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateTransactionStepByUUID DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE transaction_steps SET 
		name=$1,
		alternate_name=$2,
		description=$3,
		updated_by=$4,
		updated_at=current_timestamp at time zone 'UTC'
		WHERE text(uuid) = $5
		`

	if _, err := dbConnPgx.Exec(ctx, sql,
		transactionStep.Name,          //1
		transactionStep.AlternateName, //2
		transactionStep.Description,   //3
		transactionStep.UpdatedBy,     //4
		transactionStep.UUID,          //5
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertTransactionStep(dbConnPgx utils.PgxIface, transactionStep *TransactionStep) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertTransactionStep DbConn.Begin   %s", err.Error())
		return -1, -1, err
	}
	var MarketDataID int
	var StepID int
	transactionStepUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	if transactionStep.UUID == "" {
		transactionStep.UUID = transactionStepUUID.String()
	}
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO transaction_steps  
	(
		transaction_id,  
		step_id,
		uuid, 
		name,
		alternate_name,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			current_timestamp at time zone 'UTC',
			$7,
			current_timestamp at time zone 'UTC'
		)
		RETURNING transaction_id, step_id`,
		transactionStep.TransactionID, //1
		transactionStep.StepID,        //2
		transactionStep.UUID,          //3
		transactionStep.Name,          //4
		transactionStep.AlternateName, //5
		transactionStep.Description,   //6
		transactionStep.CreatedBy,     //7
	).Scan(&MarketDataID, &StepID)
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
	return int(MarketDataID), int(StepID), nil
}

func InsertTransactionSteps(dbConnPgx utils.PgxIface, transactionSteps []TransactionStep) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range transactionSteps {
		transactionStep := transactionSteps[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(transactionStep.UUID)
		row := []interface{}{
			transactionStep.TransactionID, //1
			transactionStep.StepID,        //2
			uuidString,                    //3
			transactionStep.Name,          //4
			transactionStep.AlternateName, //5
			transactionStep.Description,   //6
			transactionStep.CreatedBy,     //7
			&transactionStep.CreatedAt,    //8
			transactionStep.CreatedBy,     //9
			&now,                          //10
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"transaction_steps"},
		[]string{
			"transaction_id", //1
			"step_id",        //2
			"uuid",           //3
			"name",           //4
			"alternate_name", //5
			"description",    //6
			"created_by",     //7
			"created_at",     //8
			"updated_by",     //9
			"updated_at",     //10
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertTransactionSteps: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetTransactionStepListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]TransactionStep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT 
	transaction_id,  
	step_id,
	uuid, 
	name,
	alternate_name,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM transaction_steps 
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

	transactionSteps, err := pgx.CollectRows(results, pgx.RowToStructByName[TransactionStep])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionSteps, nil
}

func GetTotalTransactionStepsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM transaction_steps
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
