package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetTransaction(transactionID int) (*Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	description,
	tx_hash,
	status_id,
	from_account_id,
	to_account_id,
	chain_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at
	FROM transactions 
	WHERE id = $1`, transactionID)

	transaction := &Transaction{}
	err := row.Scan(
		&transaction.ID,
		&transaction.UUID,
		&transaction.Name,
		&transaction.AlternateName,
		&transaction.StartDate,
		&transaction.EndDate,
		&transaction.Description,
		&transaction.TxHash,
		&transaction.StatusID,
		&transaction.FromAccountID,
		&transaction.ToAccountID,
		&transaction.ChainID,
		&transaction.CreatedBy,
		&transaction.CreatedAt,
		&transaction.UpdatedBy,
		&transaction.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return transaction, nil
}

func GetTopTenTransactions() ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	description,
	tx_hash,
	status_id,
	from_account_id,
	to_account_id,
	chain_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at
	FROM transactions 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	transactionList := make([]Transaction, 0)
	for results.Next() {
		var transaction Transaction
		results.Scan(
			&transaction.ID,
			&transaction.UUID,
			&transaction.Name,
			&transaction.AlternateName,
			&transaction.StartDate,
			&transaction.EndDate,
			&transaction.Description,
			&transaction.TxHash,
			&transaction.StatusID,
			&transaction.FromAccountID,
			&transaction.ToAccountID,
			&transaction.ChainID,
			&transaction.CreatedBy,
			&transaction.CreatedAt,
			&transaction.UpdatedBy,
			&transaction.UpdatedAt,
		)

		transactionList = append(transactionList, transaction)
	}
	return transactionList, nil
}

func RemoveTransaction(transactionID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM transactions WHERE id = $1`, transactionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetTransactions(ids []int) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	description,
	tx_hash,
	status_id,
	from_account_id,
	to_account_id,
	chain_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at
	FROM transactions`
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
	transactionList := make([]Transaction, 0)
	for results.Next() {
		var transaction Transaction
		results.Scan(
			&transaction.ID,
			&transaction.UUID,
			&transaction.Name,
			&transaction.AlternateName,
			&transaction.StartDate,
			&transaction.EndDate,
			&transaction.Description,
			&transaction.TxHash,
			&transaction.StatusID,
			&transaction.FromAccountID,
			&transaction.ToAccountID,
			&transaction.ChainID,
			&transaction.CreatedBy,
			&transaction.CreatedAt,
			&transaction.UpdatedBy,
			&transaction.UpdatedAt,
		)

		transactionList = append(transactionList, transaction)
	}
	return transactionList, nil
}

func GetTransactionsByUUIDs(UUIDList []string) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	description,
	tx_hash,
	status_id,
	from_account_id,
	to_account_id,
	chain_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at
	FROM transactions
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	transactionList := make([]Transaction, 0)
	for results.Next() {
		var transaction Transaction
		results.Scan(
			&transaction.ID,
			&transaction.UUID,
			&transaction.Name,
			&transaction.AlternateName,
			&transaction.StartDate,
			&transaction.EndDate,
			&transaction.Description,
			&transaction.TxHash,
			&transaction.StatusID,
			&transaction.FromAccountID,
			&transaction.ToAccountID,
			&transaction.ChainID,
			&transaction.CreatedBy,
			&transaction.CreatedAt,
			&transaction.UpdatedBy,
			&transaction.UpdatedAt,
		)

		transactionList = append(transactionList, transaction)
	}
	return transactionList, nil
}

func GetStartAndEndDateDiffTransactions(diffInDate int) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	description,
	tx_hash,
	status_id,
	from_account_id,
	to_account_id,
	chain_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at
	FROM transactions
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	transactionList := make([]Transaction, 0)
	for results.Next() {
		var transaction Transaction
		results.Scan(
			&transaction.ID,
			&transaction.UUID,
			&transaction.Name,
			&transaction.AlternateName,
			&transaction.StartDate,
			&transaction.EndDate,
			&transaction.Description,
			&transaction.TxHash,
			&transaction.StatusID,
			&transaction.FromAccountID,
			&transaction.ToAccountID,
			&transaction.ChainID,
			&transaction.CreatedBy,
			&transaction.CreatedAt,
			&transaction.UpdatedBy,
			&transaction.UpdatedAt,
		)

		transactionList = append(transactionList, transaction)
	}
	return transactionList, nil
}

func UpdateTransaction(transaction Transaction) error {
	// if the transaction id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if transaction.ID == nil || *transaction.ID == 0 {
		return errors.New("transaction has invalid ID")
	}
	startDate := transaction.StartDate
	endDate := transaction.EndDate
	log.Println(fmt.Sprintf("Updating start: %s, end : %s", startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)))
	_, err := database.DbConnPgx.Query(ctx, `UPDATE transactions SET 
		name=$1,  
		alternate_name=$2, 
		start_date =$3,
		end_date =$4,
		description=$5, 
		tx_hash=$6,
		status_id=$7, 
		from_account_id=$8,
		to_account_id=$9,
		chain_id=$10,
		updated_by=$11, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$12`,
		transaction.Name,          //1
		transaction.AlternateName, //2
		transaction.StartDate,     //3
		transaction.EndDate,       //4
		transaction.Description,   //5
		transaction.TxHash,        //6
		transaction.StatusID,      //7
		transaction.FromAccountID, //8
		transaction.ToAccountID,   //9
		transaction.ChainID,       //10
		transaction.UpdatedBy,     //11
		transaction.ID)            //12
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertTransaction(transaction Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	transactionUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	if transaction.UUID == "" {
		transaction.UUID = transactionUUID.String()
	}
	err = database.DbConn.QueryRowContext(ctx, `INSERT INTO transactions 
	(
		uuid, 
		name, 
		alternate_name, 
		start_date,
		end_date,
		description,
		tx_hash,
		status_id,
		from_account_id,
		to_account_id,
		chain_id,
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
			$8,
			$9,
			$10,
			$11,
			$12,
			current_timestamp at time zone 'UTC',
			$12,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		&transaction.UUID,          //1
		&transaction.Name,          //2
		&transaction.AlternateName, //3
		&transaction.StartDate,     //4
		&transaction.EndDate,       //5
		&transaction.Description,   //6
		&transaction.TxHash,        //7
		&transaction.StatusID,      //8
		&transaction.FromAccountID, //9
		&transaction.ToAccountID,   //10
		&transaction.ChainID,       //11
		&transaction.CreatedBy,     //12
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

func insertTransactions(transactions []Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range transactions {
		transaction := transactions[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(transaction.UUID)
		row := []interface{}{
			uuidString,                 //1
			transaction.Name,           //2
			transaction.AlternateName,  //3
			&transaction.StartDate,     //4
			&transaction.EndDate,       //5
			transaction.Description,    //6
			transaction.TxHash,         //7
			*transaction.StatusID,      //8
			*transaction.FromAccountID, //9
			*transaction.ToAccountID,   //10
			*transaction.ChainID,       //11
			transaction.CreatedBy,      //12
			&now,                       //13
			transaction.CreatedBy,      //14
			&now,                       //15
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"transactions"},
		[]string{
			"uuid",            //1
			"name",            //2
			"alternate_name",  //3
			"start_date",      //4
			"end_date",        //5
			"description",     //6
			"tx_hash",         //7
			"status_id",       //8
			"from_account_id", //9
			"to_account_id",   //10
			"chain_id",        //11
			"created_by",      //12
			"created_at",      //13
			"updated_by",      //14
			"updated_at",      //15
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
