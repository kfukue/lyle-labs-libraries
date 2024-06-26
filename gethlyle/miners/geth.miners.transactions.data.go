package gethlyleminers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	gethlyletransactions "github.com/kfukue/lyle-labs-libraries/gethlyle/transactions"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetAllGethMinerTransactionsByMinerID(minerID *int) ([]*GethMinerTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	`, minerID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethMinerTransactions := make([]*GethMinerTransaction, 0)
	for results.Next() {
		var gethMinerTransaction GethMinerTransaction
		results.Scan(
			&gethMinerTransaction.MinerID,
			&gethMinerTransaction.TransactionID,
			&gethMinerTransaction.UUID,
			&gethMinerTransaction.Name,
			&gethMinerTransaction.AlternateName,
			&gethMinerTransaction.Description,
			&gethMinerTransaction.CreatedBy,
			&gethMinerTransaction.CreatedAt,
			&gethMinerTransaction.UpdatedBy,
			&gethMinerTransaction.UpdatedAt,
		)

		gethMinerTransactions = append(gethMinerTransactions, &gethMinerTransaction)
	}
	return gethMinerTransactions, nil
}

func GetMinAndMaxDatesFromTransactionsByMinerID(minerID *int) (*time.Time, *time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `
	SELECT 
		MIN(gt.txn_date) as min_date,
		MAX(gt.txn_date) as max_date
	FROM geth_miners_transactions gmt
	LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id
	WHERE 
	gmt.miner_id = $1
	`, minerID)
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
func GetDistinctAddressesFromGethTransactionsByMinerIDAndBeforeDate(minerID *int, beforeDate *time.Time) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
	SELECT 
		DISTINCT
		gt.from_address
	FROM geth_miners_transactions gmt
	LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id
	WHERE 
	gmt.miner_id = $1
	AND 
	gt.txn_date <= $2
	`, minerID, beforeDate.Format(utils.LayoutPostgres))
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

func GetAllGethTransactionsByMinerIDAndFromAddress(minerID *int, fromAddress string) ([]*gethlyletransactions.GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	`, minerID, fromAddress)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions := make([]*gethlyletransactions.GethTransaction, 0)
	for results.Next() {
		var gethTransaction gethlyletransactions.GethTransaction
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

		gethTransactions = append(gethTransactions, &gethTransaction)
	}
	return gethTransactions, nil
}

func GetAllGethTransactionsByMinerIDAndFromAddressToDate(minerID *int, fromAddress string, toDate *time.Time) ([]*gethlyletransactions.GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	`, minerID, fromAddress, toDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions := make([]*gethlyletransactions.GethTransaction, 0)
	for results.Next() {
		var gethTransaction gethlyletransactions.GethTransaction
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

		gethTransactions = append(gethTransactions, &gethTransaction)
	}
	return gethTransactions, nil
}

// fromDate inclusive and toDate exclusive
func GetAllGethTransactionsByMinerIDAndFromAddressFromToDate(minerID *int, fromAddress string, fromDate, toDate *time.Time) ([]*gethlyletransactions.GethTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	`, minerID, fromAddress, fromDate.Format(utils.LayoutPostgres), toDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactions := make([]*gethlyletransactions.GethTransaction, 0)
	for results.Next() {
		var gethTransaction gethlyletransactions.GethTransaction
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

		gethTransactions = append(gethTransactions, &gethTransaction)
	}
	return gethTransactions, nil
}

func GetAllGethMinerTransactionsByTransactionID(transactionID *int) ([]*GethMinerTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	`, transactionID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethMinerTransactions := make([]*GethMinerTransaction, 0)
	for results.Next() {
		var gethMinerTransaction GethMinerTransaction
		results.Scan(
			&gethMinerTransaction.MinerID,
			&gethMinerTransaction.TransactionID,
			&gethMinerTransaction.UUID,
			&gethMinerTransaction.Name,
			&gethMinerTransaction.AlternateName,
			&gethMinerTransaction.Description,
			&gethMinerTransaction.CreatedBy,
			&gethMinerTransaction.CreatedAt,
			&gethMinerTransaction.UpdatedBy,
			&gethMinerTransaction.UpdatedAt,
		)

		gethMinerTransactions = append(gethMinerTransactions, &gethMinerTransaction)
	}
	return gethMinerTransactions, nil
}

func GetMinerTransaction(minerID, transactionID *int) (*GethMinerTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `
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
	`, minerID, transactionID)

	gethMinerTransaction := &GethMinerTransaction{}
	err := row.Scan(
		&gethMinerTransaction.MinerID,
		&gethMinerTransaction.TransactionID,
		&gethMinerTransaction.UUID,
		&gethMinerTransaction.Name,
		&gethMinerTransaction.AlternateName,
		&gethMinerTransaction.Description,
		&gethMinerTransaction.CreatedBy,
		&gethMinerTransaction.CreatedAt,
		&gethMinerTransaction.UpdatedBy,
		&gethMinerTransaction.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethMinerTransaction, nil
}

func RemoveMinerTransaction(minerID, transactionID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM geth_miners_transactions WHERE 
	miner_id = $1 AND transaction_id =$2`, minerID, transactionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetMinerTransactionList(minerIDs, transactionIDs []int) ([]*GethMinerTransaction, error) {
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
	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	minerTransactionInputs := make([]*GethMinerTransaction, 0)
	for results.Next() {
		var minerTransactionInput GethMinerTransaction
		results.Scan(
			&minerTransactionInput.MinerID,
			&minerTransactionInput.TransactionID,
			&minerTransactionInput.UUID,
			&minerTransactionInput.Name,
			&minerTransactionInput.AlternateName,
			&minerTransactionInput.Description,
			&minerTransactionInput.CreatedBy,
			&minerTransactionInput.CreatedAt,
			&minerTransactionInput.UpdatedBy,
			&minerTransactionInput.UpdatedAt,
		)

		minerTransactionInputs = append(minerTransactionInputs, &minerTransactionInput)
	}
	return minerTransactionInputs, nil
}

func UpdateMinerTransaction(minerTransactionInput GethMinerTransaction) error {
	// if the minerTransactionInput id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (minerTransactionInput.MinerID == nil || *minerTransactionInput.MinerID == 0) || (minerTransactionInput.TransactionID == nil || *minerTransactionInput.TransactionID == 0) {
		return errors.New("minerTransactionInput has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE geth_miners_transactions SET 
		name=$1,  
		alternate_name=$2, 
		description=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE miner_id=$5 AND transaction_id=$6`,
		minerTransactionInput.Name,          //1
		minerTransactionInput.AlternateName, //2
		minerTransactionInput.Description,   //3
		minerTransactionInput.UpdatedBy,     //4
		minerTransactionInput.MinerID,       //5
		minerTransactionInput.TransactionID, //6
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertMinerTransaction(minerTransactionInput GethMinerTransaction) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var MinerID int
	var TransactionID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_miners_transactions  
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
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(MinerID), int(TransactionID), nil
}

func InsertGethMinersTransactions(gethMinersTransaction []*GethMinerTransaction) error {
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
	copyCount, err := database.DbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("geth_miners_transactions copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}
	return nil
}

func RemoveAllTransactionsAndTransactionInputs() error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `
	DELETE FROM geth_miners_transactions;
	DELETE FROM geth_transactions;
			`,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetMinerTransactionListByPagination(_start, _end *int, _order, _sort string, _filters []string) ([]*GethMinerTransaction, error) {
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

	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	minerTransactionInputs := make([]*GethMinerTransaction, 0)
	for results.Next() {
		var minerTransactionInput GethMinerTransaction
		results.Scan(
			&minerTransactionInput.MinerID,
			&minerTransactionInput.TransactionID,
			&minerTransactionInput.UUID,
			&minerTransactionInput.Name,
			&minerTransactionInput.AlternateName,
			&minerTransactionInput.Description,
			&minerTransactionInput.CreatedBy,
			&minerTransactionInput.CreatedAt,
			&minerTransactionInput.UpdatedBy,
			&minerTransactionInput.UpdatedAt,
		)

		minerTransactionInputs = append(minerTransactionInputs, &minerTransactionInput)
	}
	return minerTransactionInputs, nil
}

func GetTotalMinerTransactionCount() (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_miners_transactions
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
