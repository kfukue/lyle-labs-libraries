package gethlyleminers

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

func GetAllGethMinerTransactionsByMinerID(dbConnPgx utils.PgxIface, minerID *int) ([]GethMinerTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT 
		miner_id,
		transaction_id,
		uuid,
		name,
		alternate_name,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
	FROM geth_miners_transactions 
	WHERE 
	miner_id = $1
	`, *minerID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethMinerTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethMinerTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethMinerTransactions, nil
}

func GetMinAndMaxDatesFromTransactionsByMinerID(dbConnPgx utils.PgxIface, minerID *int) (*time.Time, *time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := dbConnPgx.QueryRow(ctx, `
	SELECT 
		MIN(gt.txn_date) as min_date,
		MAX(gt.txn_date) as max_date
	FROM geth_miners_transactions gmt
	LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id
	WHERE 
	gmt.miner_id = $1
	`, *minerID)
	var minDate, maxDate *time.Time
	err := row.Scan(
		&minDate,
		&maxDate,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	return minDate, maxDate, nil
}
func GetDistinctAddressesFromGethTransactionsByMinerIDAndBeforeDate(dbConnPgx utils.PgxIface, minerID *int, beforeDate *time.Time) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT 
		DISTINCT
		gt.from_address
	FROM geth_miners_transactions gmt
	LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id
	WHERE 
	gmt.miner_id = $1
	AND 
	gt.txn_date <= $2
	`, *minerID, beforeDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	addressesStr := make([]string, 0)
	for results.Next() {
		var addressStr string
		results.Scan(
			&addressStr,
		)

		addressesStr = append(addressesStr, addressStr)
	}
	return addressesStr, nil
}

func GetAllGethMinerTransactionsByTransactionID(dbConnPgx utils.PgxIface, transactionID *int) ([]GethMinerTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT 
		miner_id,
		transaction_id,
		uuid,
		name,
		alternate_name,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
	FROM geth_miners_transactions 
	WHERE 
	transaction_id = $1
	`, *transactionID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethMinerTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethMinerTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethMinerTransactions, nil
}

func GetGethMinerTransaction(dbConnPgx utils.PgxIface, minerID, transactionID *int) (*GethMinerTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `
	SELECT 
		miner_id,
		transaction_id,
		uuid,
		name,
		alternate_name,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
	FROM geth_miners_transactions 
	WHERE miner_id = $1
	AND transaction_id = $2
	`, *minerID, *transactionID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethMinerTransaction, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethMinerTransaction])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethMinerTransaction, nil
}

func RemoveGethMinerTransaction(dbConnPgx utils.PgxIface, minerID, transactionID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethMinerTransaction DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_miners_transactions WHERE miner_id = $1 AND transaction_id =$2`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *minerID, *transactionID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetGethMinerTransactionList(dbConnPgx utils.PgxIface, minerIDs, transactionIDs []int) ([]GethMinerTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT 
		miner_id,
		transaction_id,
		uuid,
		name,
		alternate_name,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
	FROM geth_miners_transactions `
	if len(minerIDs) > 0 || len(transactionIDs) > 0 {
		additionalQuery := ` WHERE`
		if len(minerIDs) > 0 {
			strIds := utils.SplitToString(minerIDs, ",")
			additionalQuery += fmt.Sprintf(`miner_id IN (%s)`, strIds)
		}
		if len(transactionIDs) > 0 {
			if len(minerIDs) > 0 {
				additionalQuery += `AND `
			}
			strIds := utils.SplitToString(transactionIDs, ",")
			additionalQuery += fmt.Sprintf(`transaction_id IN (%s)`, strIds)
		}
		sql += additionalQuery
	}
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethMinerTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethMinerTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethMinerTransactions, nil
}

func UpdateGethMinerTransaction(dbConnPgx utils.PgxIface, minerTransactionInput *GethMinerTransaction) error {
	// if the minerTransactionInput id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (minerTransactionInput.MinerID == nil || *minerTransactionInput.MinerID == 0) || (minerTransactionInput.TransactionID == nil || *minerTransactionInput.TransactionID == 0) {
		return errors.New("minerTransactionInput has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethMinerTransaction DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_miners_transactions SET 
		name=$1,  
		alternate_name=$2, 
		description=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE miner_id=$5 AND transaction_id=$6`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		minerTransactionInput.Name,          //1
		minerTransactionInput.AlternateName, //2
		minerTransactionInput.Description,   //3
		minerTransactionInput.UpdatedBy,     //4
		minerTransactionInput.MinerID,       //5
		minerTransactionInput.TransactionID, //6
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethMinerTransaction(dbConnPgx utils.PgxIface, minerTransactionInput *GethMinerTransaction) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethMinerTransaction DbConn.Begin   %s", err.Error())
		return -1, -1, err
	}
	var MinerID int
	var TransactionID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_miners_transactions  
	(
		miner_id,
		transaction_id,
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
			uuid_generate_v4(),
			$3,
			$4,
			$5,
			$6,
			current_timestamp at time zone 'UTC',
			$6,
			current_timestamp at time zone 'UTC'
		)
		RETURNING miner_id, transaction_id`,
		minerTransactionInput.MinerID,       //1
		minerTransactionInput.TransactionID, //2
		minerTransactionInput.Name,          //3
		minerTransactionInput.AlternateName, //4
		minerTransactionInput.Description,   //5
		minerTransactionInput.CreatedBy,     //6
	).Scan(&MinerID, &TransactionID)
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
	return int(MinerID), int(TransactionID), nil
}

func InsertGethMinersTransactions(dbConnPgx utils.PgxIface, gethMinersTransaction []GethMinerTransaction) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethMinersTransaction {
		gethTransactionInput := gethMinersTransaction[i]
		uuidString := pgtype.UUID{}
		uuidString.Set(gethTransactionInput.UUID)
		row := []interface{}{
			gethTransactionInput.MinerID,       //1
			gethTransactionInput.TransactionID, //2
			uuidString,                         //3
			gethTransactionInput.Name,          //4
			gethTransactionInput.AlternateName, //5
			gethTransactionInput.Description,   //6
			gethTransactionInput.CreatedBy,     //7
			now,                                //8
			gethTransactionInput.CreatedBy,     //9
			now,                                //10

		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_miners_transactions"},
		[]string{
			"miner_id",       //1
			"transaction_id", //2
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
	log.Println(fmt.Printf("InsertGethMinersTransactions: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveAllTransactionsAndTransactionInputs(dbConnPgx utils.PgxIface) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveAllTransactionsAndTransactionInputs DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_miners_transactions;
		DELETE FROM geth_transactions;`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

// for refinedev
func GetMinerTransactionListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethMinerTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
	miner_id,
	transaction_id,
	uuid, 
	name, 
	alternate_name, 
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_miners_transactions 
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
	gethMinerTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethMinerTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethMinerTransactions, nil
}

func GetTotalMinerTransactionCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_miners_transactions
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
