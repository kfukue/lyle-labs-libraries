package transaction

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

func GetTransaction(dbConnPgx utils.PgxIface, transactionID *int) (*Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
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
	WHERE id = $1`, *transactionID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	transaction, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Transaction])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &transaction, nil
}

func RemoveTransaction(dbConnPgx utils.PgxIface, transactionID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveTransaction DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM transactions WHERE id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *transactionID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetTransactions(dbConnPgx utils.PgxIface, ids []int) ([]Transaction, error) {
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
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	transactionList, err := pgx.CollectRows(results, pgx.RowToStructByName[Transaction])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return transactionList, nil
}

func GetTransactionsByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
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

	transactionList, err := pgx.CollectRows(results, pgx.RowToStructByName[Transaction])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return transactionList, nil
}

func GetStartAndEndDateDiffTransactions(dbConnPgx utils.PgxIface, diffInDate *int) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
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
	`, *diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	transactionList, err := pgx.CollectRows(results, pgx.RowToStructByName[Transaction])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return transactionList, nil
}

func UpdateTransaction(dbConnPgx utils.PgxIface, transaction *Transaction) error {
	// if the transaction id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if transaction.ID == nil || *transaction.ID == 0 {
		return errors.New("transaction has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateTransaction DbConn.Begin   %s", err.Error())
		return err
	}
	startDate := transaction.StartDate
	endDate := transaction.EndDate
	log.Println(fmt.Sprintf("Updating start: %s, end : %s", startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)))
	sql := `UPDATE transactions SET 
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
		WHERE id=$12`

	if _, err := dbConnPgx.Exec(ctx, sql,
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
		transaction.ID,            //12
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertTransaction(dbConnPgx utils.PgxIface, transaction *Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertTransaction DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	transactionUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	if transaction.UUID == "" {
		transaction.UUID = transactionUUID.String()
	}
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO transactions 
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
		transaction.UUID,          //1
		transaction.Name,          //2
		transaction.AlternateName, //3
		transaction.StartDate,     //4
		transaction.EndDate,       //5
		transaction.Description,   //6
		transaction.TxHash,        //7
		transaction.StatusID,      //8
		transaction.FromAccountID, //9
		transaction.ToAccountID,   //10
		transaction.ChainID,       //11
		transaction.CreatedBy,     //12
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

func InsertTransactions(dbConnPgx utils.PgxIface, transactions []Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range transactions {
		transaction := transactions[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(transaction.UUID)
		row := []interface{}{
			uuidString,                //1
			transaction.Name,          //2
			transaction.AlternateName, //3
			&transaction.StartDate,    //4
			&transaction.EndDate,      //5
			transaction.Description,   //6
			transaction.TxHash,        //7
			transaction.StatusID,      //8
			transaction.FromAccountID, //9
			transaction.ToAccountID,   //10
			transaction.ChainID,       //11
			transaction.CreatedBy,     //12
			&now,                      //13
			transaction.CreatedBy,     //14
			&now,                      //15
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("InsertTransactions: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetTransactionsByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
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

	transactions, err := pgx.CollectRows(results, pgx.RowToStructByName[Transaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactions, nil
}

func GetTotalTransactionsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM transactions
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
