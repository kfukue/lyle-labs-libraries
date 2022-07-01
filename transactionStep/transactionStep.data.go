package transactionstep

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/lib/pq"
)

func GetTransactionStep(transactionID int, stepID int) (*TransactionStep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	`, transactionID, stepID)

	transactionStep := &TransactionStep{}
	err := row.Scan(
		&transactionStep.TransactionID,
		&transactionStep.StepID,
		&transactionStep.UUID,
		&transactionStep.Name,
		&transactionStep.AlternateName,
		&transactionStep.Description,
		&transactionStep.CreatedBy,
		&transactionStep.CreatedAt,
		&transactionStep.UpdatedBy,
		&transactionStep.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionStep, nil
}

func GetTransactionStepByUUID(transactionStepUUID string) (*TransactionStep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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

	transactionStep := &TransactionStep{}
	err := row.Scan(
		&transactionStep.TransactionID,
		&transactionStep.StepID,
		&transactionStep.UUID,
		&transactionStep.Name,
		&transactionStep.AlternateName,
		&transactionStep.Description,
		&transactionStep.CreatedBy,
		&transactionStep.CreatedAt,
		&transactionStep.UpdatedBy,
		&transactionStep.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionStep, nil
}

func GetTransactionStepsByUUIDs(UUIDList []string) ([]TransactionStep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	defer results.Close()
	transactionStepList := make([]TransactionStep, 0)
	for results.Next() {
		var transactionStep TransactionStep
		results.Scan(
			&transactionStep.TransactionID,
			&transactionStep.StepID,
			&transactionStep.UUID,
			&transactionStep.Name,
			&transactionStep.AlternateName,
			&transactionStep.Description,
			&transactionStep.CreatedBy,
			&transactionStep.CreatedAt,
			&transactionStep.UpdatedBy,
			&transactionStep.UpdatedAt,
		)

		transactionStepList = append(transactionStepList, transactionStep)
	}
	return transactionStepList, nil
}

func GetTopTenTransactionSteps() ([]TransactionStep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	transactionSteps := make([]TransactionStep, 0)
	for results.Next() {
		var transactionStep TransactionStep
		results.Scan(
			&transactionStep.TransactionID,
			&transactionStep.StepID,
			&transactionStep.UUID,
			&transactionStep.Name,
			&transactionStep.AlternateName,
			&transactionStep.Description,
			&transactionStep.CreatedBy,
			&transactionStep.CreatedAt,
			&transactionStep.UpdatedBy,
			&transactionStep.UpdatedAt,
		)

		transactionSteps = append(transactionSteps, transactionStep)
	}
	return transactionSteps, nil
}

func RemoveTransactionStep(transactionID int, stepID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM transaction_steps WHERE 
	transaction_id = $1 AND step_id =$2`, transactionID, stepID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveTransactionStepByUUID(transactionStepUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM transaction_steps WHERE 
		WHERE text(uuid) = $1`,
		transactionStepUUID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetTransactionStepList() ([]TransactionStep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	defer results.Close()
	transactionSteps := make([]TransactionStep, 0)
	for results.Next() {
		var transactionStep TransactionStep
		results.Scan(
			&transactionStep.TransactionID,
			&transactionStep.StepID,
			&transactionStep.UUID,
			&transactionStep.Name,
			&transactionStep.AlternateName,
			&transactionStep.Description,
			&transactionStep.CreatedBy,
			&transactionStep.CreatedAt,
			&transactionStep.UpdatedBy,
			&transactionStep.UpdatedAt,
		)

		transactionSteps = append(transactionSteps, transactionStep)
	}
	return transactionSteps, nil
}

func UpdateTransactionStep(transactionStep TransactionStep) error {
	// if the transactionStep id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (transactionStep.TransactionID == nil || *transactionStep.TransactionID == 0) || (transactionStep.StepID == nil || *transactionStep.StepID == 0) {
		return errors.New("transactionStep has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE transaction_steps SET 
		name=$1,
		alternate_name=$2,
		description=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE transaction_id=$5 AND step_id=$6`,

		transactionStep.Name,          //1
		transactionStep.AlternateName, //2
		transactionStep.Description,   //3
		transactionStep.UpdatedBy,     //4
		transactionStep.TransactionID, //5
		transactionStep.StepID,        //6
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateTransactionStepByUUID(transactionStep TransactionStep) error {
	// if the transactionStep id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (transactionStep.TransactionID == nil || *transactionStep.TransactionID == 0) || (transactionStep.StepID == nil || *transactionStep.StepID == 0) {
		return errors.New("transactionStep has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE transaction_steps SET 
		name=$1,
		alternate_name=$2,
		description=$3,
		updated_by=$4,
		updated_at=current_timestamp at time zone 'UTC'
		WHERE text(uuid) = $5
		`,
		transactionStep.Name,          //1
		transactionStep.AlternateName, //2
		transactionStep.Description,   //3
		transactionStep.UpdatedBy,     //4
		transactionStep.UUID,          //5
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertTransactionStep(transactionStep TransactionStep) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var MarketDataID int
	var JobID int
	transactionStepUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	if transactionStep.UUID == "" {
		transactionStep.UUID = transactionStepUUID.String()
	}
	err = database.DbConn.QueryRowContext(ctx, `INSERT INTO transaction_steps  
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
	).Scan(&MarketDataID, &JobID)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(MarketDataID), int(JobID), nil
}

func InsertTransactionSteps(transactionSteps []TransactionStep) error {
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

			*transactionStep.TransactionID, //1
			*transactionStep.StepID,        //2
			uuidString,                     //3
			transactionStep.Name,           //4
			transactionStep.AlternateName,  //5
			transactionStep.Description,    //6
			transactionStep.CreatedBy,      //7
			&transactionStep.CreatedAt,     //8
			transactionStep.CreatedBy,      //9
			&now,                           //10
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}

	return nil
}
