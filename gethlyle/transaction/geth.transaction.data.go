package gethlyletransactions

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/lib/pq"
)

func GetGethTransaction(gethTransactionID int) (*GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT
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
	`, gethTransactionID)

	gethTransaction := GethTransaction{}
	err := row.Scan(
		&gethTransaction.ID,
		&gethTransaction.UUID,
		&gethTransaction.ChainID,
		&gethTransaction.ExchangeID,
		&gethTransaction.BlockNumber,
		&gethTransaction.IndexNumber,
		&gethTransaction.TxnDate,
		&gethTransaction.TxnHash,
		&gethTransaction.FromAddress,
		&gethTransaction.FromAddressID,
		&gethTransaction.ToAddress,
		&gethTransaction.ToAddressID,
		&gethTransaction.InteractedContractAddress,
		&gethTransaction.InteractedContractAddressID,
		&gethTransaction.NativeAssetID,
		&gethTransaction.GethProcessJobID,
		&gethTransaction.Value,
		&gethTransaction.GethTransctionInputId,
		&gethTransaction.StatusID,
		&gethTransaction.Description,
		&gethTransaction.CreatedBy,
		&gethTransaction.CreatedAt,
		&gethTransaction.UpdatedBy,
		&gethTransaction.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransaction, nil
}

func GetGethTransactionByFromToAddress(fromToAddressID *int) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
		fromToAddressID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions := make([]GethTransaction, 0)
	for results.Next() {
		var gethTransaction GethTransaction
		results.Scan(
			&gethTransaction.ID,
			&gethTransaction.UUID,
			&gethTransaction.ChainID,
			&gethTransaction.ExchangeID,
			&gethTransaction.BlockNumber,
			&gethTransaction.IndexNumber,
			&gethTransaction.TxnDate,
			&gethTransaction.TxnHash,
			&gethTransaction.FromAddress,
			&gethTransaction.FromAddressID,
			&gethTransaction.ToAddress,
			&gethTransaction.ToAddressID,
			&gethTransaction.InteractedContractAddress,
			&gethTransaction.InteractedContractAddressID,
			&gethTransaction.NativeAssetID,
			&gethTransaction.GethProcessJobID,
			&gethTransaction.Value,
			&gethTransaction.GethTransctionInputId,
			&gethTransaction.StatusID,
			&gethTransaction.Description,
			&gethTransaction.CreatedBy,
			&gethTransaction.CreatedAt,
			&gethTransaction.UpdatedBy,
			&gethTransaction.UpdatedAt,
		)
		gethTransactions = append(gethTransactions, gethTransaction)
	}
	return gethTransactions, nil
}

func GetGethTransactionByFromAddressAndBeforeBlockNumber(fromAddressID, blockNumber *int) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
		fromAddressID, blockNumber,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions := make([]GethTransaction, 0)
	for results.Next() {
		var gethTransaction GethTransaction
		results.Scan(
			&gethTransaction.ID,
			&gethTransaction.UUID,
			&gethTransaction.ChainID,
			&gethTransaction.ExchangeID,
			&gethTransaction.BlockNumber,
			&gethTransaction.IndexNumber,
			&gethTransaction.TxnDate,
			&gethTransaction.TxnHash,
			&gethTransaction.FromAddress,
			&gethTransaction.FromAddressID,
			&gethTransaction.ToAddress,
			&gethTransaction.ToAddressID,
			&gethTransaction.InteractedContractAddress,
			&gethTransaction.InteractedContractAddressID,
			&gethTransaction.NativeAssetID,
			&gethTransaction.GethProcessJobID,
			&gethTransaction.Value,
			&gethTransaction.GethTransctionInputId,
			&gethTransaction.StatusID,
			&gethTransaction.Description,
			&gethTransaction.CreatedBy,
			&gethTransaction.CreatedAt,
			&gethTransaction.UpdatedBy,
			&gethTransaction.UpdatedAt,
		)
		gethTransactions = append(gethTransactions, gethTransaction)
	}
	return gethTransactions, nil
}

func GetGethTransactionsByTxnHash(txnHash string) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	gethTransactions := make([]GethTransaction, 0)
	for results.Next() {
		var gethTransaction GethTransaction
		results.Scan(
			&gethTransaction.ID,
			&gethTransaction.UUID,
			&gethTransaction.ChainID,
			&gethTransaction.ExchangeID,
			&gethTransaction.BlockNumber,
			&gethTransaction.IndexNumber,
			&gethTransaction.TxnDate,
			&gethTransaction.TxnHash,
			&gethTransaction.FromAddress,
			&gethTransaction.FromAddressID,
			&gethTransaction.ToAddress,
			&gethTransaction.ToAddressID,
			&gethTransaction.InteractedContractAddress,
			&gethTransaction.InteractedContractAddressID,
			&gethTransaction.NativeAssetID,
			&gethTransaction.GethProcessJobID,
			&gethTransaction.Value,
			&gethTransaction.GethTransctionInputId,
			&gethTransaction.StatusID,
			&gethTransaction.Description,
			&gethTransaction.CreatedBy,
			&gethTransaction.CreatedAt,
			&gethTransaction.UpdatedBy,
			&gethTransaction.UpdatedAt,
		)
		gethTransactions = append(gethTransactions, gethTransaction)
	}
	return gethTransactions, nil
}

func GetGethTransactionsByTxnHashes(txnHashes []string) ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	gethTransactions := make([]GethTransaction, 0)
	for results.Next() {
		var gethTransaction GethTransaction
		results.Scan(
			&gethTransaction.ID,
			&gethTransaction.UUID,
			&gethTransaction.ChainID,
			&gethTransaction.ExchangeID,
			&gethTransaction.BlockNumber,
			&gethTransaction.IndexNumber,
			&gethTransaction.TxnDate,
			&gethTransaction.TxnHash,
			&gethTransaction.FromAddress,
			&gethTransaction.FromAddressID,
			&gethTransaction.ToAddress,
			&gethTransaction.ToAddressID,
			&gethTransaction.InteractedContractAddress,
			&gethTransaction.InteractedContractAddressID,
			&gethTransaction.NativeAssetID,
			&gethTransaction.GethProcessJobID,
			&gethTransaction.Value,
			&gethTransaction.GethTransctionInputId,
			&gethTransaction.StatusID,
			&gethTransaction.Description,
			&gethTransaction.CreatedBy,
			&gethTransaction.CreatedAt,
			&gethTransaction.UpdatedBy,
			&gethTransaction.UpdatedAt,
		)
		gethTransactions = append(gethTransactions, gethTransaction)
	}
	return gethTransactions, nil
}

func RemoveGethTransaction(gethTransactionID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_transactions WHERE id = $1`, gethTransactionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveGethTransactionsFromChainIDAndStartBlockNumber(chainID, startBlockNumber *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_transactions WHERE chain_id = $1 AND block_number >=  $2`, chainID, startBlockNumber)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveGethTransactionsFromChainID(chainID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_transactions WHERE chain_id = $1`, *chainID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetGethTransactionList() ([]GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	gethTransactions := make([]GethTransaction, 0)
	for results.Next() {
		var gethTransaction GethTransaction
		results.Scan(
			&gethTransaction.ID,
			&gethTransaction.UUID,
			&gethTransaction.ChainID,
			&gethTransaction.ExchangeID,
			&gethTransaction.BlockNumber,
			&gethTransaction.IndexNumber,
			&gethTransaction.TxnDate,
			&gethTransaction.TxnHash,
			&gethTransaction.FromAddress,
			&gethTransaction.FromAddressID,
			&gethTransaction.ToAddress,
			&gethTransaction.ToAddressID,
			&gethTransaction.InteractedContractAddress,
			&gethTransaction.InteractedContractAddressID,
			&gethTransaction.NativeAssetID,
			&gethTransaction.GethProcessJobID,
			&gethTransaction.Value,
			&gethTransaction.GethTransctionInputId,
			&gethTransaction.StatusID,
			&gethTransaction.Description,
			&gethTransaction.CreatedBy,
			&gethTransaction.CreatedAt,
			&gethTransaction.UpdatedBy,
			&gethTransaction.UpdatedAt,
		)

		gethTransactions = append(gethTransactions, gethTransaction)
	}
	return gethTransactions, nil
}

func UpdateGethTransaction(gethTransaction GethTransaction) error {
	// if the gethTransaction id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethTransaction.ID == nil {
		return errors.New("gethTransaction has invalid ID")
	}
	_, err := database.DbConnPgx.Exec(ctx, `UPDATE geth_transactions SET
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
		WHERE id=$20`,
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
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertGethTransaction(gethTransaction GethTransaction) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var gethTransactionID int
	var gethTransactionUUID string
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_transactions
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
		log.Println(err.Error())
		return 0, "", err
	}
	return int(gethTransactionID), gethTransactionUUID, nil
}
func InsertGethTransactions(gethTransactions []*GethTransaction) error {
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
	copyCount, err := database.DbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("geth_transactions copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}
	return nil
}

func UpdateGethTransactionAddresses() error {
	// update address ids from existing addresses in geth_addresses
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `
		UPDATE geth_transactions as gt SET
			from_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gt.from_address) = LOWER(ga.address_str)
			AND gt.from_address_id IS NULL
			`,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	_, err = database.DbConnPgx.Exec(ctx, `
			UPDATE geth_transactions as gt SET
			to_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gt.to_address) = LOWER(ga.address_str)
			AND gt.to_address_id IS NULL
			`,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetNullAddressStrsFromTransactions() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
