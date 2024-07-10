package gethlyletransactions

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetGethTransaction(dbConnPgx utils.PgxIface, gethTransactionID *int) (*GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		txn_date,
		txn_hash,
		from_address,
		from_address_id,
		to_address,
		to_address_id,
		interacted_contract_address,
		interacted_contract_address_id,
		native_asset_id,
		geth_process_job_id,
		value,
		geth_transction_input_id,
		status_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_transactions
	WHERE id = $1
	`, *gethTransactionID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethTransaction, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethTransaction])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethTransaction, nil
}

func GetGethTransactionByFromToAddress(dbConnPgx utils.PgxIface, fromToAddressID *int) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		txn_date,
		txn_hash,
		from_address,
		from_address_id,
		to_address,
		to_address_id,
		interacted_contract_address,
		interacted_contract_address_id,
		native_asset_id,
		geth_process_job_id,
		value,
		geth_transction_input_id,
		status_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_transactions
	WHERE
		AND (from_address_id =$1 OR 
			to_address_id = $1)
		`,
		*fromToAddressID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactions, nil
}

func GetGethTransactionByFromAddressAndBeforeBlockNumber(dbConnPgx utils.PgxIface, fromAddressID *int, blockNumber *uint64) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		txn_date,
		txn_hash,
		from_address,
		from_address_id,
		to_address,
		to_address_id,
		interacted_contract_address,
		interacted_contract_address_id,
		native_asset_id,
		geth_process_job_id,
		value,
		geth_transction_input_id,
		status_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_transactions
	WHERE
		(from_address_id =$1 OR 
			to_address_id = $1)
		AND block_number <= $2
		`,
		*fromAddressID, *blockNumber,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactions, nil
}

func GetGethTransactionsByTxnHash(dbConnPgx utils.PgxIface, txnHash string) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		txn_date,
		txn_hash,
		from_address,
		from_address_id,
		to_address,
		to_address_id,
		interacted_contract_address,
		interacted_contract_address_id,
		native_asset_id,
		geth_process_job_id,
		value,
		geth_transction_input_id,
		status_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_transactions
		WHERE
		txn_hash = $1
		`,
		txnHash,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactions, nil
}

func GetGethTransactionsByTxnHashes(dbConnPgx utils.PgxIface, txnHashes []string) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		txn_date,
		txn_hash,
		from_address,
		from_address_id,
		to_address,
		to_address_id,
		interacted_contract_address,
		interacted_contract_address_id,
		native_asset_id,
		geth_process_job_id,
		value,
		geth_transction_input_id,
		status_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
		FROM geth_transactions
		WHERE
		txn_hash = ANY($1)
		`,
		pq.Array(txnHashes),
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactions, nil
}

func GetGethTransactionsByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		txn_date,
		txn_hash,
		from_address,
		from_address_id,
		to_address,
		to_address_id,
		interacted_contract_address,
		interacted_contract_address_id,
		native_asset_id,
		geth_process_job_id,
		value,
		geth_transction_input_id,
		status_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
		FROM geth_transactions
		WHERE text(uuid) = ANY($1)
		`,
		pq.Array(UUIDList),
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactions, nil
}

func RemoveGethTransaction(dbConnPgx utils.PgxIface, gethTransactionID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethMinerTransaction DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_transactions WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *gethTransactionID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemoveGethTransactionsFromChainIDAndStartBlockNumber(dbConnPgx utils.PgxIface, chainID *int, startBlockNumber *uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethMinerTransaction DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_transactions WHERE chain_id = $1 AND block_number >=  $2`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *chainID, *startBlockNumber); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemoveGethTransactionsFromChainID(dbConnPgx utils.PgxIface, chainID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethMinerTransaction DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_transactions WHERE chain_id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *chainID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetGethTransactionList(dbConnPgx utils.PgxIface) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		txn_date,
		txn_hash,
		from_address,
		from_address_id,
		to_address,
		to_address_id,
		interacted_contract_address,
		interacted_contract_address_id,
		native_asset_id,
		geth_process_job_id,
		value,
		geth_transction_input_id,
		status_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_transactions `)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactions, nil
}

func UpdateGethTransaction(dbConnPgx utils.PgxIface, gethTransaction *GethTransaction) error {
	// if the gethTransaction id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethTransaction.ID == nil {
		return errors.New("gethTransaction has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethTransaction DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_transactions SET 
		chain_id=$1,
		exchange_id=$2,
		block_number=$3,
		index_number=$4,
		txn_date=$5,
		txn_hash=$6,
		from_address=$7,
		from_address_id=$8,
		to_address=$9,
		to_address_id=$10,
		interacted_contract_address=$11,
		interacted_contract_address_id=$12,
		native_asset_id=$13,
		geth_process_job_id=$14,
		value=$15,
		geth_transction_input_id=$16,
		status_id=$17,
		description=$18,
		updated_by=$19,
		updated_at=current_timestamp at time zone 'UTC',
		WHERE id=$20`

	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		gethTransaction.ChainID,                     //1
		gethTransaction.ExchangeID,                  //2
		gethTransaction.BlockNumber,                 //3
		gethTransaction.IndexNumber,                 //4
		gethTransaction.TxnDate,                     //5
		gethTransaction.TxnHash,                     //6
		gethTransaction.FromAddress,                 //7
		gethTransaction.FromAddressID,               //8
		gethTransaction.ToAddress,                   //9
		gethTransaction.ToAddressID,                 //10
		gethTransaction.InteractedContractAddress,   //11
		gethTransaction.InteractedContractAddressID, //12
		gethTransaction.NativeAssetID,               //13
		gethTransaction.GethProcessJobID,            //14
		gethTransaction.Value,                       //15
		gethTransaction.GethTransctionInputId,       //16
		gethTransaction.StatusID,                    //17
		gethTransaction.Description,                 //18
		gethTransaction.UpdatedBy,                   //19
		gethTransaction.ID,                          //20
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethTransaction(dbConnPgx utils.PgxIface, gethTransaction *GethTransaction) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethTransaction DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var gethTransactionID int
	var gethTransactionUUID string
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_transactions
	(
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		txn_date,
		txn_hash,
		from_address,
		from_address_id,
		to_address,
		to_address_id,
		interacted_contract_address,
		interacted_contract_address_id,
		native_asset_id,
		geth_process_job_id,
		value,
		geth_transction_input_id,
		status_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
		) VALUES (
		uuid_generate_v4(),
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
		$13,
		$14,
		$15,
		$16,
		$17,
		$18,
		$19,
		current_timestamp at time zone 'UTC', 
		$19,
		current_timestamp at time zone 'UTC'
		)
		RETURNING id, uuid`,
		gethTransaction.ChainID,                     //1
		gethTransaction.ExchangeID,                  //2
		gethTransaction.BlockNumber,                 //3
		gethTransaction.IndexNumber,                 //4
		gethTransaction.TxnDate,                     //5
		gethTransaction.TxnHash,                     //6
		gethTransaction.FromAddress,                 //7
		gethTransaction.FromAddressID,               //8
		gethTransaction.ToAddress,                   //9
		gethTransaction.ToAddressID,                 //10
		gethTransaction.InteractedContractAddress,   //11
		gethTransaction.InteractedContractAddressID, //12
		gethTransaction.NativeAssetID,               //13
		gethTransaction.GethProcessJobID,            //14
		gethTransaction.Value,                       //15
		gethTransaction.GethTransctionInputId,       //16
		gethTransaction.StatusID,                    //17
		gethTransaction.Description,                 //18
		gethTransaction.CreatedBy,                   //19
	).Scan(&gethTransactionID, &gethTransactionUUID)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, "", err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, "", err
	}
	return int(gethTransactionID), gethTransactionUUID, nil
}
func InsertGethTransactions(dbConnPgx utils.PgxIface, gethTransactions []GethTransaction) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethTransactions {
		gethTransaction := gethTransactions[i]
		uuidString := pgtype.UUID{}
		uuidString.Set(gethTransaction.UUID)
		row := []interface{}{
			uuidString,                                  //1
			gethTransaction.ChainID,                     //2
			gethTransaction.ExchangeID,                  //3
			gethTransaction.BlockNumber,                 //4
			gethTransaction.IndexNumber,                 //5
			gethTransaction.TxnDate,                     //6
			gethTransaction.TxnHash,                     //7
			gethTransaction.FromAddress,                 //8
			gethTransaction.FromAddressID,               //9
			gethTransaction.ToAddress,                   //10
			gethTransaction.ToAddressID,                 //11
			gethTransaction.InteractedContractAddress,   //12
			gethTransaction.InteractedContractAddressID, //13
			gethTransaction.NativeAssetID,               //14
			gethTransaction.GethProcessJobID,            //15
			gethTransaction.Value,                       //16
			gethTransaction.GethTransctionInputId,       //17
			gethTransaction.StatusID,                    //18
			gethTransaction.Description,                 //19
			gethTransaction.CreatedBy,                   //20
			now,                                         //21
			gethTransaction.CreatedBy,                   //22
			now,                                         //23

		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_transactions"},
		[]string{
			"uuid",                           //1
			"chain_id",                       //2
			"exchange_id",                    //3
			"block_number",                   //4
			"index_number",                   //5
			"txn_date",                       //6
			"txn_hash",                       //7
			"from_address",                   //8
			"from_address_id",                //9
			"to_address",                     //10
			"to_address_id",                  //11
			"interacted_contract_address",    //12
			"interacted_contract_address_id", //13
			"native_asset_id",                //14
			"geth_process_job_id",            //15
			"value",                          //16
			"geth_transction_input_id",       //17
			"status_id",                      //18
			"description",                    //19
			"created_by",                     //20
			"created_at",                     //21
			"updated_by",                     //22
			"updated_at",                     //23
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethTransactions: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateGethTransactionAddresses(dbConnPgx utils.PgxIface) error {
	// update address ids from existing addresses in geth_addresses
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethTransactionAddresses DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_transactions as gt SET
			from_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gt.from_address) = LOWER(ga.address_str)
			AND gt.from_address_id IS NULL
			`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql); err != nil {
		tx.Rollback(ctx)
		return err
	}

	sql2 := `UPDATE geth_transactions as gt SET
			to_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gt.to_address) = LOWER(ga.address_str)
			AND gt.to_address_id IS NULL
			`
	if _, err := dbConnPgx.Exec(ctx, sql2); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetNullAddressStrsFromTransactions(dbConnPgx utils.PgxIface) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
		WITH sender_table as(
			SELECT DISTINCT LOWER(gt.from_address) as address  
			FROM geth_transactions gt
				LEFT JOIN geth_addresses as ga
				ON LOWER(gt.from_address) = LOWER(ga.address_str)
			WHERE gt.from_address_id IS NULL
			AND ga.id IS NULL
		),
		to_table as(
			SELECT DISTINCT LOWER(gt.to_address) as address 
			FROM geth_transactions gt
				LEFT JOIN geth_addresses as ga
				ON LOWER(gt.to_address) = LOWER(ga.address_str)
			WHERE gt.to_address_id IS NULL
			AND ga.id IS NULL
		)
		SELECT * FROM sender_table
		UNION
		SELECT * FROM to_table
		`,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethNullAddressStrs := make([]string, 0)
	for results.Next() {
		var gethNullAddressStr string
		results.Scan(
			&gethNullAddressStr,
		)
		gethNullAddressStrs = append(gethNullAddressStrs, gethNullAddressStr)
	}
	return gethNullAddressStrs, nil
}

// for refinedev
func GetTransactionListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		uuid,
		chain_id,
		exchange_id,
		block_number,
		index_number,
		txn_date,
		txn_hash,
		from_address,
		from_address_id,
		to_address,
		to_address_id,
		interacted_contract_address,
		interacted_contract_address_id,
		native_asset_id,
		geth_process_job_id,
		value,
		geth_transction_input_id,
		status_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_transactions 
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
	gethTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactions, nil
}

func GetTotalTransactionsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_transactions
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

func GetAllGethTransactionsByMinerIDAndFromAddress(dbConnPgx utils.PgxIface, minerID *int, fromAddress string) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT 
		gt.id,
		gt.uuid,
		gt.chain_id,
		gt.exchange_id,
		gt.block_number,
		gt.index_number,
		gt.txn_date,
		gt.txn_hash,
		gt.from_address,
		gt.from_address_id,
		gt.to_address,
		gt.to_address_id,
		gt.interacted_contract_address,
		gt.interacted_contract_address_id,
		gt.native_asset_id,
		gt.geth_process_job_id,
		gt.value,
		gt.geth_transction_input_id,
		gt.status_id,
		gt.description,
		gt.created_by,
		gt.created_at,
		gt.updated_by,
		gt.updated_at
	FROM geth_miners_transactions gmt
	LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id
	WHERE 
	gmt.miner_id = $1
	AND 
	gt.from_address = $2
	ORDER BY gt.txn_date asc
	`, *minerID, fromAddress)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactions, nil
}

func GetAllGethTransactionsByMinerIDAndFromAddressToDate(dbConnPgx utils.PgxIface, minerID *int, fromAddress string, toDate *time.Time) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT 
		gt.id,
		gt.uuid,
		gt.chain_id,
		gt.exchange_id,
		gt.block_number,
		gt.index_number,
		gt.txn_date,
		gt.txn_hash,
		gt.from_address,
		gt.from_address_id,
		gt.to_address,
		gt.to_address_id,
		gt.interacted_contract_address,
		gt.interacted_contract_address_id,
		gt.native_asset_id,
		gt.geth_process_job_id,
		gt.value,
		gt.geth_transction_input_id,
		gt.status_id,
		gt.description,
		gt.created_by,
		gt.created_at,
		gt.updated_by,
		gt.updated_at
	FROM geth_miners_transactions gmt
	LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id
	WHERE 
		gmt.miner_id = $1
			AND 
		gt.from_address = $2
			AND 
		gt.txn_date <= $3
	ORDER BY gt.txn_date asc
	`, *minerID, fromAddress, toDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactions, nil
}

// fromDate inclusive and toDate exclusive
func GetAllGethTransactionsByMinerIDAndFromAddressFromToDate(dbConnPgx utils.PgxIface, minerID *int, fromAddress string, fromDate, toDate *time.Time) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT 
		gt.id,
		gt.uuid,
		gt.chain_id,
		gt.exchange_id,
		gt.block_number,
		gt.index_number,
		gt.txn_date,
		gt.txn_hash,
		gt.from_address,
		gt.from_address_id,
		gt.to_address,
		gt.to_address_id,
		gt.interacted_contract_address,
		gt.interacted_contract_address_id,
		gt.native_asset_id,
		gt.geth_process_job_id,
		gt.value,
		gt.geth_transction_input_id,
		gt.status_id,
		gt.description,
		gt.created_by,
		gt.created_at,
		gt.updated_by,
		gt.updated_at
	FROM geth_miners_transactions gmt
	LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id
	WHERE 
		gmt.miner_id = $1
			AND 
		gt.from_address = $2
			AND 
		gt.txn_date >= $3
			AND 
		gt.txn_date < $4
	ORDER BY gt.txn_date asc
	`, *minerID, fromAddress, fromDate.Format(utils.LayoutPostgres), toDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransaction])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactions, nil
}
