package gethlyletransfers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	gethlyleaddresses "github.com/kfukue/lyle-labs-libraries/v2/gethlyle/address"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
)

func GetGethTransfer(dbConnPgx utils.PgxIface, gethTransferID *int) (*GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
	FROM geth_transfers
	WHERE id = $1
	`, *gethTransferID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethTransfer, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethTransfer])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethTransfer, nil
}

func GetGethTransferByBlockChain(dbConnPgx utils.PgxIface, txnHash string, blockNumber *uint64, indexNumber *uint) (*GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
	FROM geth_transfers
	WHERE txn_hash = $1
	AND block_number = $2
	AND index_number = $3
	`, txnHash, *blockNumber, *indexNumber)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethTransfer, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethTransfer])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethTransfer, nil
}

func GetTransfersTransactionHashByUserAddress(dbConnPgx utils.PgxIface, userAddressID, assetID *int, blockNumber *uint64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT DISTINCT txn_hash FROM geth_transfers
	WHERE
	(to_address_id =$1 OR sender_address_id = $1)
	AND asset_id = $2
	AND block_number > $3
	`,
		*userAddressID, *assetID, *blockNumber)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	txnAddresses := make([]string, 0)
	for results.Next() {
		var txnAddress string
		results.Scan(
			&txnAddress,
		)

		txnAddresses = append(txnAddresses, txnAddress)
	}
	return txnAddresses, nil
}

func GetDistinctAddressesFromAssetId(dbConnPgx utils.PgxIface, assetID *int) ([]gethlyleaddresses.GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	WITH sender_table as(
		SELECT DISTINCT sender_address_id as address  
		FROM geth_transfers
		WHERE asset_id = $1
	),
	to_table as(
		SELECT DISTINCT to_address_id as address FROM geth_transfers
		WHERE asset_id = $1
	)
	SELECT 
	id,  
		uuid, 
		name,
		alternate_name,
		description,
		address_str,
		address_type_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		FROM geth_addresses
		where id IN (
	SELECT * FROM sender_table
	UNION
	SELECT * FROM to_table
		)
	`,
		*assetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethAddresses, err := pgx.CollectRows(results, pgx.RowToStructByName[gethlyleaddresses.GethAddress])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethAddresses, nil

}

func GetDistinctTransactionHashesFromAssetIdAndStartingBlock(dbConnPgx utils.PgxIface, assetID *int, startingBlock *uint64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		DISTINCT txn_hash 
	FROM geth_transfers
	WHERE
		block_number >= $1
		AND asset_id = $2
	`,
		*startingBlock, *assetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	txnHashes := make([]string, 0)
	for results.Next() {
		var txnHash string
		results.Scan(
			&txnHash,
		)
		txnHashes = append(txnHashes, txnHash)
	}
	return txnHashes, nil

}

func GetHighestBlockFromBaseAssetId(dbConnPgx utils.PgxIface, baseAssetID *int) (*uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := dbConnPgx.QueryRow(ctx, `SELECT COALESCE (MAX(block_number), 0) FROM geth_transfers
	WHERE base_asset_id=$1
		`,
		*baseAssetID)
	var maxBlockNumber uint64
	err := row.Scan(
		&maxBlockNumber)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &maxBlockNumber, nil
}

func GetGethTransferByFromTokenAddress(dbConnPgx utils.PgxIface, tokenAddressID *int) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
		FROM geth_transfers
		WHERE
		token_address_id = $1
		`,
		*tokenAddressID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransfer])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethTransfers, nil
}

func GetGethTransferByFromMakerAddressAndTokenAddressID(dbConnPgx utils.PgxIface, makerAddressID, tokenAddressID *int) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
		FROM geth_transfers
		WHERE
		token_address_id = $1
		AND (sender_address_id =$2 OR 
			to_address_id = $2)
		`,
		*tokenAddressID, *makerAddressID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransfer])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethTransfers, nil
}

func GetGethTransferByFromMakerAddressAndTokenAddressIDAndBeforeBlockNumber(dbConnPgx utils.PgxIface, makerAddressID, baseAssetID *int, blockNumber *uint64) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
		FROM geth_transfers
		WHERE
		asset_id = $1
		AND
		base_asset_id =$1
		AND (sender_address_id =$2 OR 
			to_address_id = $2)
		AND block_number <= $3
		`,
		*baseAssetID, *makerAddressID, *blockNumber,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransfer])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethTransfers, nil
}

func GetGethTransferByFromBaseAssetIDAndBeforeBlockNumber(dbConnPgx utils.PgxIface, baseAssetID *int, blockNumber *uint64) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
		FROM geth_transfers
		WHERE
		asset_id = $1
		AND
		base_asset_id =$1
		AND block_number <= $2
		AND sender_address_id IS NOT NULL
		AND to_address_id IS NOT NULL
		`,
		*baseAssetID, *blockNumber,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransfer])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethTransfers, nil
}

func GetGethTransfersByTxnHash(dbConnPgx utils.PgxIface, txnHash string, baseAssetID *int) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
		FROM geth_transfers
		WHERE
		txn_hash = $1
		AND base_asset_id = $2
		`,
		txnHash, *baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransfer])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethTransfers, nil
}

func GetGethTransfersByTxnHashes(dbConnPgx utils.PgxIface, txnHashes []string, baseAssetID *int) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
		FROM geth_transfers
		WHERE
		txn_hash = ANY($1)
		AND base_asset_id = $2
		`,
		pq.Array(txnHashes), *baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransfer])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethTransfers, nil
}

func RemoveGethTransfer(dbConnPgx utils.PgxIface, gethTransferID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethTransfer DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_transfers WHERE id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *gethTransferID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemoveGethTransfersFromBaseAssetIDAndStartBlockNumber(dbConnPgx utils.PgxIface, baseAssetID *int, startBlockNumber *uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethTransfersFromBaseAssetIDAndStartBlockNumber DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_transfers WHERE base_asset_id = $1 AND block_number >=  $2`

	if _, err := dbConnPgx.Exec(ctx, sql, *baseAssetID, *startBlockNumber); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemoveGethTransfersFromBaseAssetID(dbConnPgx utils.PgxIface, baseAssetID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethTransfersFromBaseAssetID DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_transfers WHERE base_asset_id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *baseAssetID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetGethTransferList(dbConnPgx utils.PgxIface) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
	FROM geth_transfers `)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransfer])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return gethTransfers, nil
}

func UpdateGethTransfer(dbConnPgx utils.PgxIface, gethTransfer *GethTransfer) error {
	// if the gethTransfer id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethTransfer.ID == nil {
		return errors.New("gethTransfer has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethTransfer DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_transfers SET 
		chain_id=$1,
		token_address=$2,
		token_address_id = $3,
		asset_id=$4,
		block_number=$5,
		index_number=$6,
		transfer_date=$7,
		txn_hash=$8,
		sender_address=$9,
		sender_address_id=$10,
		to_address=$11,
		to_address_id=$12,
		amount=$13,
		description=$14,
		updated_by=$15,
		updated_at=current_timestamp at time zone 'UTC',
		geth_process_job_id =$16,
		topics_str=$17,
		status_id=$18,
		base_asset_id=$19,
		transfer_type_id=$20
		WHERE id=$21`

	if _, err := dbConnPgx.Exec(ctx, sql,
		gethTransfer.ChainID,             //1
		gethTransfer.TokenAddress,        //2
		gethTransfer.TokenAddressID,      //3
		gethTransfer.AssetID,             //4
		gethTransfer.BlockNumber,         //5
		gethTransfer.IndexNumber,         //6
		gethTransfer.TransferDate,        //7
		gethTransfer.TxnHash,             //8
		gethTransfer.SenderAddress,       //9
		gethTransfer.SenderAddressID,     //10
		gethTransfer.ToAddress,           //11
		gethTransfer.ToAddressID,         //12
		gethTransfer.Amount,              //13
		gethTransfer.Description,         //14
		gethTransfer.UpdatedBy,           //15
		gethTransfer.GethProcessJobID,    //16
		pq.Array(gethTransfer.TopicsStr), //17
		gethTransfer.StatusID,            //18
		gethTransfer.BaseAssetID,         //19
		gethTransfer.TransferTypeID,      //20
		gethTransfer.ID,                  //21
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethTransfer(dbConnPgx utils.PgxIface, gethTransfer *GethTransfer) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethTransfer DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var gethTransferID int
	var gethTransferUUID string
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_transfers
	(
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
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
		current_timestamp at time zone 'UTC',
		$15,
		current_timestamp at time zone 'UTC', 
		$16,
		$17,
		$18,
		$19,
		$20
		)
		RETURNING id, uuid`,
		gethTransfer.ChainID,             //1
		gethTransfer.TokenAddress,        //2
		gethTransfer.TokenAddressID,      //3
		gethTransfer.AssetID,             //4
		gethTransfer.BlockNumber,         //5
		gethTransfer.IndexNumber,         //6
		gethTransfer.TransferDate,        //7
		gethTransfer.TxnHash,             //8
		gethTransfer.SenderAddress,       //9
		gethTransfer.SenderAddressID,     //10
		gethTransfer.ToAddress,           //11
		gethTransfer.ToAddressID,         //12
		gethTransfer.Amount,              //13
		gethTransfer.Description,         //14
		gethTransfer.CreatedBy,           //15
		gethTransfer.GethProcessJobID,    //16
		pq.Array(gethTransfer.TopicsStr), //17
		gethTransfer.StatusID,            //18
		gethTransfer.BaseAssetID,         //19
		gethTransfer.TransferTypeID,      //20
	).Scan(&gethTransferID, &gethTransferUUID)
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
	return int(gethTransferID), gethTransferUUID, nil
}
func InsertGethTransfers(dbConnPgx utils.PgxIface, gethTransfers []GethTransfer) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethTransfers {
		gethTransfer := gethTransfers[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(gethTransfer.UUID)
		row := []interface{}{
			uuidString,                       //1
			gethTransfer.ChainID,             //2
			gethTransfer.TokenAddress,        //3
			gethTransfer.TokenAddressID,      //4
			gethTransfer.AssetID,             //5
			gethTransfer.BlockNumber,         //6
			gethTransfer.IndexNumber,         //7
			gethTransfer.TransferDate,        //8
			gethTransfer.TxnHash,             //9
			gethTransfer.SenderAddress,       //10
			gethTransfer.SenderAddressID,     //11
			gethTransfer.ToAddress,           //12
			gethTransfer.ToAddressID,         //13
			gethTransfer.Amount,              //14
			gethTransfer.Description,         //15
			gethTransfer.CreatedBy,           //16
			&now,                             //17
			gethTransfer.CreatedBy,           //18
			&now,                             //19
			gethTransfer.GethProcessJobID,    //20
			pq.Array(gethTransfer.TopicsStr), //21
			gethTransfer.StatusID,            //22
			gethTransfer.BaseAssetID,         //23
			gethTransfer.TransferTypeID,      //24

		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_transfers"},
		[]string{
			"uuid",                //1
			"chain_id",            //2
			"token_address",       //3
			"token_address_id",    //4
			"asset_id",            //5
			"block_number",        //6
			"index_number",        //7
			"transfer_date",       //8
			"txn_hash",            //9
			"sender_address",      //10
			"sender_address_id",   //11
			"to_address",          //12
			"to_address_id",       //13
			"amount",              //14
			"description",         //15
			"created_by",          //16
			"created_at",          //17
			"updated_by",          //18
			"updated_at",          //19
			"geth_process_job_id", //20
			"topics_str",          //21
			"status_id",           //22
			"base_asset_id",       //23
			"transfer_type_id",    //24
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethTransfers: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateGethTransferAddresses(dbConnPgx utils.PgxIface, baseAssetID *int) error {
	// update address ids from existing addresses in geth_addresses
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethTransferAddresses DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_transfers as gt SET
			sender_address_id = ga.id from geth_addresses as ga
			WHERE 
				gt.sender_address_id IS NULL
				AND	gt.base_asset_id = $1
				AND LOWER(gt.sender_address) = LOWER(ga.address_str)
			`

	if _, err := dbConnPgx.Exec(ctx, sql, *baseAssetID); err != nil {
		log.Println(fmt.Printf("UpdateGethTransferAddresses: Error at sql1 %v", err))
		tx.Rollback(ctx)
		return err
	}
	sql2 := `UPDATE geth_transfers as gt SET
			to_address_id = ga.id from geth_addresses as ga
			WHERE
				gt.to_address_id IS NULL
				AND	gt.base_asset_id = $1
				AND LOWER(gt.to_address) = LOWER(ga.address_str)
			`
	if _, err := dbConnPgx.Exec(ctx, sql2, *baseAssetID); err != nil {
		log.Println(fmt.Printf("UpdateGethTransferAddresses: Error at sql2 %v", err))
		tx.Rollback(ctx)
		return err
	}
	sql3 := `UPDATE geth_transfers as gt SET
			asset_id = assets.id
			from assets as assets
			WHERE
				gt.asset_id IS NULL
				AND	gt.base_asset_id = $1
				AND LOWER(gt.token_address) = LOWER(assets.contract_address)
	`
	if _, err := dbConnPgx.Exec(ctx, sql3, *baseAssetID); err != nil {
		log.Println(fmt.Printf("UpdateGethTransferAddresses: Error at sql3 %v", err))
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetNullAddressStrsFromTransfers(dbConnPgx utils.PgxIface, baseAssetID *int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
		WITH sender_table as(
			SELECT DISTINCT LOWER(gt.sender_address) as address  
			FROM geth_transfers gt
				LEFT JOIN geth_addresses as ga
				ON LOWER(gt.sender_address) = LOWER(ga.address_str)
			WHERE gt.sender_address_id IS NULL
			AND gt.base_asset_id = $1
			AND ga.id IS NULL
		),
		to_table as(
			SELECT DISTINCT LOWER(gt.to_address) as address 
			FROM geth_transfers gt
				LEFT JOIN geth_addresses as ga
				ON LOWER(gt.to_address) = LOWER(ga.address_str)
			WHERE gt.to_address_id IS NULL
			AND gt.base_asset_id = $1
			AND ga.id IS NULL
		)
		SELECT * FROM sender_table
		UNION
		SELECT * FROM to_table
		`, *baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

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

func UpdateGethTransfersAssetIDs(dbConnPgx utils.PgxIface) error {
	// update address ids from existing addresses in geth_addresses
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethSwap DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_transfers as gt SET
		asset_id = assets.id from assets as assets
			WHERE LOWER(gt.token_address) = LOWER(assets.contract_address) AND
			gt.asset_id is NULL
	`

	if _, err := dbConnPgx.Exec(ctx, sql); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

// for refinedev
func GetGethTransferListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
	FROM geth_transfers 
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

	gethTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransfer])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransfers, nil
}

func GetTotalTransfersCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_transfers
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

func GetClosestBlockNumberFromGethTransferFromChainAndDate(dbConnPgx utils.PgxIface, chainID *int, asOfDate time.Time, isBefore *bool) (*GethTrade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	selectSQL := `SELECT
		id,
		uuid,
		chain_id,
		token_address,
		token_address_id,
		asset_id,
		block_number,
		index_number,
		transfer_date,
		txn_hash,
		sender_address,
		sender_address_id,
		to_address,
		to_address_id,
		amount,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		geth_process_job_id,
		topics_str,
		status_id,
		base_asset_id,
		transfer_type_id
	FROM geth_transfers 
	WHERE 
	`
	tradeDateSQL := ``
	if *isBefore {
		tradeDateSQL = ` transfer_date <= $2 
		AND chain_id = $1
		AND block_number IS NOT NULL 
		ORDER BY transfer_date desc`
	} else {
		tradeDateSQL = `  transfer_date >= $2
		AND chain_id = $1
		AND block_number IS NOT NULL 
		 ORDER BY transfer_date asc`
	}
	sql := selectSQL + tradeDateSQL + ` LIMIT 1`
	row, err := dbConnPgx.Query(ctx, sql, *chainID, asOfDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethTrade, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethTransfer])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethTrade, nil
}
