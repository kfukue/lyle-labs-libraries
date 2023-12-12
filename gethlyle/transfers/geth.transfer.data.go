package gethlyletransfers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	gethlyleaddresses "github.com/kfukue/lyle-labs-libraries/gethlyle/address"
	"github.com/lib/pq"
)

func GetGethTransfer(gethTransferID int) (*GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT
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
	`, gethTransferID)

	gethTransfer := &GethTransfer{}
	err := row.Scan(
		&gethTransfer.ID,
		&gethTransfer.UUID,
		&gethTransfer.ChainID,
		&gethTransfer.TokenAddress,
		&gethTransfer.TokenAddressID,
		&gethTransfer.AssetID,
		&gethTransfer.BlockNumber,
		&gethTransfer.IndexNumber,
		&gethTransfer.TransferDate,
		&gethTransfer.TxnHash,
		&gethTransfer.SenderAddress,
		&gethTransfer.SenderAddressID,
		&gethTransfer.ToAddress,
		&gethTransfer.ToAddressID,
		&gethTransfer.Amount,
		&gethTransfer.Description,
		&gethTransfer.CreatedBy,
		&gethTransfer.CreatedAt,
		&gethTransfer.UpdatedBy,
		&gethTransfer.UpdatedAt,
		&gethTransfer.GethProcessJobID,
		&gethTransfer.TopicsStr,
		&gethTransfer.StatusID,
		&gethTransfer.BaseAssetID,
		&gethTransfer.TransferTypeID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransfer, nil
}

func GetGethTransferByBlockChain(txnHash string, blockNumber *uint64, indexNumber *uint) (*GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT
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
	`, txnHash, blockNumber, indexNumber)

	gethTransfer := &GethTransfer{}
	err := row.Scan(
		&gethTransfer.ID,
		&gethTransfer.UUID,
		&gethTransfer.ChainID,
		&gethTransfer.TokenAddress,
		&gethTransfer.TokenAddressID,
		&gethTransfer.AssetID,
		&gethTransfer.BlockNumber,
		&gethTransfer.IndexNumber,
		&gethTransfer.TransferDate,
		&gethTransfer.TxnHash,
		&gethTransfer.SenderAddress,
		&gethTransfer.SenderAddressID,
		&gethTransfer.ToAddress,
		&gethTransfer.ToAddressID,
		&gethTransfer.Amount,
		&gethTransfer.Description,
		&gethTransfer.CreatedBy,
		&gethTransfer.CreatedAt,
		&gethTransfer.UpdatedBy,
		&gethTransfer.UpdatedAt,
		&gethTransfer.GethProcessJobID,
		&gethTransfer.TopicsStr,
		&gethTransfer.StatusID,
		&gethTransfer.StatusID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransfer, nil
}

func GetTransfersTransactionHashByUserAddress(userAddressID *int, assetID *int, blockNumber *uint64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT DISTINCT txn_hash FROM geth_transfers
	WHERE
	(to_address_id =$1 OR sender_address_id = $1)
	AND asset_id = $2
	AND block_number > $3
	`,
		userAddressID, assetID, blockNumber)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
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

func GetDistinctAddressesFromAssetId(assetID *int) ([]gethlyleaddresses.GethAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
		
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
		assetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethAddresses := make([]gethlyleaddresses.GethAddress, 0)
	for results.Next() {
		var gethAddress gethlyleaddresses.GethAddress
		results.Scan(
			&gethAddress.ID,
			&gethAddress.UUID,
			&gethAddress.Name,
			&gethAddress.AlternateName,
			&gethAddress.Description,
			&gethAddress.AddressStr,
			&gethAddress.AddressTypeID,
			&gethAddress.CreatedBy,
			&gethAddress.CreatedAt,
			&gethAddress.UpdatedBy,
			&gethAddress.UpdatedAt,
		)

		gethAddresses = append(gethAddresses, gethAddress)
	}
	return gethAddresses, nil

}

func GetDistinctTransactionHashesFromAssetIdAndStartingBlock(assetID *int, startingBlock *uint64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
		
	SELECT DISTINCT txn_hash FROM geth_transfers
	WHERE block_number >= $1
	AND asset_id = $2
	`,
		*startingBlock, *assetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
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

func GetHighestBlockFromBaseAssetId(baseAssetID *int) (*uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := database.DbConnPgx.Query(ctx, `SELECT COALESCE (MAX(block_number), 0) FROM geth_transfers
	WHERE base_asset_id=$1
		`,
		baseAssetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()

	var maxBlockNumber uint64
	for row.Next() {
		err = row.Scan(
			&maxBlockNumber)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		// no transfers
		zeroHeight := uint64(0)
		return &zeroHeight, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &maxBlockNumber, nil
}

func GetGethTransferByFromTokenAddress(tokenAddressID *int) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT
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
		tokenAddressID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransfers := make([]GethTransfer, 0)
	for results.Next() {
		var gethTransfer GethTransfer
		results.Scan(
			&gethTransfer.ID,
			&gethTransfer.UUID,
			&gethTransfer.ChainID,
			&gethTransfer.TokenAddress,
			&gethTransfer.TokenAddressID,
			&gethTransfer.AssetID,
			&gethTransfer.BlockNumber,
			&gethTransfer.IndexNumber,
			&gethTransfer.TransferDate,
			&gethTransfer.TxnHash,
			&gethTransfer.SenderAddress,
			&gethTransfer.SenderAddressID,
			&gethTransfer.ToAddress,
			&gethTransfer.ToAddressID,
			&gethTransfer.Amount,
			&gethTransfer.Description,
			&gethTransfer.CreatedBy,
			&gethTransfer.CreatedAt,
			&gethTransfer.UpdatedBy,
			&gethTransfer.UpdatedAt,
			&gethTransfer.GethProcessJobID,
			&gethTransfer.TopicsStr,
			&gethTransfer.StatusID,
		)
		gethTransfers = append(gethTransfers, gethTransfer)
	}
	return gethTransfers, nil
}

func GetGethTransferByFromMakerAddressAndTokenAddressID(makerAddressID *int, tokenAddressID *int) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT
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
		AND (sender_address_id =$2 OR 
			to_address_id = $2)
		`,
		tokenAddressID, makerAddressID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransfers := make([]GethTransfer, 0)
	for results.Next() {
		var gethTransfer GethTransfer
		results.Scan(
			&gethTransfer.ID,
			&gethTransfer.UUID,
			&gethTransfer.ChainID,
			&gethTransfer.TokenAddress,
			&gethTransfer.TokenAddressID,
			&gethTransfer.AssetID,
			&gethTransfer.BlockNumber,
			&gethTransfer.IndexNumber,
			&gethTransfer.TransferDate,
			&gethTransfer.TxnHash,
			&gethTransfer.SenderAddress,
			&gethTransfer.SenderAddressID,
			&gethTransfer.ToAddress,
			&gethTransfer.ToAddressID,
			&gethTransfer.Amount,
			&gethTransfer.Description,
			&gethTransfer.CreatedBy,
			&gethTransfer.CreatedAt,
			&gethTransfer.UpdatedBy,
			&gethTransfer.UpdatedAt,
			&gethTransfer.GethProcessJobID,
			&gethTransfer.TopicsStr,
			&gethTransfer.StatusID,
			&gethTransfer.BaseAssetID,
			&gethTransfer.TransferTypeID,
		)
		gethTransfers = append(gethTransfers, gethTransfer)
	}
	return gethTransfers, nil
}

func GetGethTransferByFromMakerAddressAndTokenAddressIDAndBeforeBlockNumber(makerAddressID, baseAssetID, blockNumber *int) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT
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
	defer results.Close()
	gethTransfers := make([]GethTransfer, 0)
	for results.Next() {
		var gethTransfer GethTransfer
		results.Scan(
			&gethTransfer.ID,
			&gethTransfer.UUID,
			&gethTransfer.ChainID,
			&gethTransfer.TokenAddress,
			&gethTransfer.TokenAddressID,
			&gethTransfer.AssetID,
			&gethTransfer.BlockNumber,
			&gethTransfer.IndexNumber,
			&gethTransfer.TransferDate,
			&gethTransfer.TxnHash,
			&gethTransfer.SenderAddress,
			&gethTransfer.SenderAddressID,
			&gethTransfer.ToAddress,
			&gethTransfer.ToAddressID,
			&gethTransfer.Amount,
			&gethTransfer.Description,
			&gethTransfer.CreatedBy,
			&gethTransfer.CreatedAt,
			&gethTransfer.UpdatedBy,
			&gethTransfer.UpdatedAt,
			&gethTransfer.GethProcessJobID,
			&gethTransfer.TopicsStr,
			&gethTransfer.StatusID,
			&gethTransfer.BaseAssetID,
			&gethTransfer.TransferTypeID,
		)
		gethTransfers = append(gethTransfers, gethTransfer)
	}
	return gethTransfers, nil
}

func GetGethTransferByFromBaseAssetIDAndBeforeBlockNumber(baseAssetID, blockNumber *int) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT
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
		`,
		*baseAssetID, *blockNumber,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransfers := make([]GethTransfer, 0)
	for results.Next() {
		var gethTransfer GethTransfer
		results.Scan(
			&gethTransfer.ID,
			&gethTransfer.UUID,
			&gethTransfer.ChainID,
			&gethTransfer.TokenAddress,
			&gethTransfer.TokenAddressID,
			&gethTransfer.AssetID,
			&gethTransfer.BlockNumber,
			&gethTransfer.IndexNumber,
			&gethTransfer.TransferDate,
			&gethTransfer.TxnHash,
			&gethTransfer.SenderAddress,
			&gethTransfer.SenderAddressID,
			&gethTransfer.ToAddress,
			&gethTransfer.ToAddressID,
			&gethTransfer.Amount,
			&gethTransfer.Description,
			&gethTransfer.CreatedBy,
			&gethTransfer.CreatedAt,
			&gethTransfer.UpdatedBy,
			&gethTransfer.UpdatedAt,
			&gethTransfer.GethProcessJobID,
			&gethTransfer.TopicsStr,
			&gethTransfer.StatusID,
			&gethTransfer.BaseAssetID,
			&gethTransfer.TransferTypeID,
		)
		gethTransfers = append(gethTransfers, gethTransfer)
	}
	return gethTransfers, nil
}

func GetGethTransfersByTxnHash(txnHash string, baseAssetID *int) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT
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
	defer results.Close()
	gethTransfers := make([]GethTransfer, 0)
	for results.Next() {
		var gethTransfer GethTransfer
		results.Scan(
			&gethTransfer.ID,
			&gethTransfer.UUID,
			&gethTransfer.ChainID,
			&gethTransfer.TokenAddress,
			&gethTransfer.TokenAddressID,
			&gethTransfer.AssetID,
			&gethTransfer.BlockNumber,
			&gethTransfer.IndexNumber,
			&gethTransfer.TransferDate,
			&gethTransfer.TxnHash,
			&gethTransfer.SenderAddress,
			&gethTransfer.SenderAddressID,
			&gethTransfer.ToAddress,
			&gethTransfer.ToAddressID,
			&gethTransfer.Amount,
			&gethTransfer.Description,
			&gethTransfer.CreatedBy,
			&gethTransfer.CreatedAt,
			&gethTransfer.UpdatedBy,
			&gethTransfer.UpdatedAt,
			&gethTransfer.GethProcessJobID,
			&gethTransfer.TopicsStr,
			&gethTransfer.StatusID,
			&gethTransfer.BaseAssetID,
			&gethTransfer.TransferTypeID,
		)
		gethTransfers = append(gethTransfers, gethTransfer)
	}
	return gethTransfers, nil
}

func GetGethTransfersByTxnHashes(txnHashes []string, baseAssetID *int) ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT
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
	defer results.Close()
	gethTransfers := make([]GethTransfer, 0)
	for results.Next() {
		var gethTransfer GethTransfer
		results.Scan(
			&gethTransfer.ID,
			&gethTransfer.UUID,
			&gethTransfer.ChainID,
			&gethTransfer.TokenAddress,
			&gethTransfer.TokenAddressID,
			&gethTransfer.AssetID,
			&gethTransfer.BlockNumber,
			&gethTransfer.IndexNumber,
			&gethTransfer.TransferDate,
			&gethTransfer.TxnHash,
			&gethTransfer.SenderAddress,
			&gethTransfer.SenderAddressID,
			&gethTransfer.ToAddress,
			&gethTransfer.ToAddressID,
			&gethTransfer.Amount,
			&gethTransfer.Description,
			&gethTransfer.CreatedBy,
			&gethTransfer.CreatedAt,
			&gethTransfer.UpdatedBy,
			&gethTransfer.UpdatedAt,
			&gethTransfer.GethProcessJobID,
			&gethTransfer.TopicsStr,
			&gethTransfer.StatusID,
			&gethTransfer.BaseAssetID,
			&gethTransfer.TransferTypeID,
		)
		gethTransfers = append(gethTransfers, gethTransfer)
	}
	return gethTransfers, nil
}

func RemoveGethTransfer(gethTransferID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_transfers WHERE id = $1`, gethTransferID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveGethTransfersFromBaseAssetIDAndStartBlockNumber(baseAssetID *int, startBlockNumber *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_transfers WHERE base_asset_id = $1 AND block_number >=  $2`, *baseAssetID, *startBlockNumber)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetGethTransferList() ([]GethTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT
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
	defer results.Close()
	gethTransfers := make([]GethTransfer, 0)
	for results.Next() {
		var gethTransfer GethTransfer
		results.Scan(
			&gethTransfer.ID,
			&gethTransfer.UUID,
			&gethTransfer.ChainID,
			&gethTransfer.TokenAddress,
			&gethTransfer.TokenAddressID,
			&gethTransfer.AssetID,
			&gethTransfer.BlockNumber,
			&gethTransfer.IndexNumber,
			&gethTransfer.TransferDate,
			&gethTransfer.TxnHash,
			&gethTransfer.SenderAddress,
			&gethTransfer.SenderAddressID,
			&gethTransfer.ToAddress,
			&gethTransfer.ToAddressID,
			&gethTransfer.Amount,
			&gethTransfer.Description,
			&gethTransfer.CreatedBy,
			&gethTransfer.CreatedAt,
			&gethTransfer.UpdatedBy,
			&gethTransfer.GethProcessJobID,
			&gethTransfer.TopicsStr,
			&gethTransfer.StatusID,
			&gethTransfer.BaseAssetID,
			&gethTransfer.TransferTypeID,
		)

		gethTransfers = append(gethTransfers, gethTransfer)
	}
	return gethTransfers, nil
}

func UpdateGethTransfer(gethTransfer GethTransfer) error {
	// if the gethTransfer id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethTransfer.ID == nil {
		return errors.New("gethTransfer has invalid ID")
	}
	_, err := database.DbConnPgx.Exec(ctx, `UPDATE geth_transfers SET
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
		WHERE id=$21`,
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
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertGethTransfer(gethTransfer GethTransfer) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var gethTransferID int
	var gethTransferUUID string
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_transfers
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
		log.Println(err.Error())
		return 0, "", err
	}
	return int(gethTransferID), gethTransferUUID, nil
}
func InsertGethTransfers(gethTransfers []*GethTransfer) error {
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
	copyCount, err := database.DbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}
	return nil
}

func UpdateGethTransferAddresses(baseAssetID *int) error {
	// update address ids from existing addresses in geth_addresses
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `
		UPDATE geth_transfers as gt SET
			sender_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gt.sender_address) = LOWER(ga.address_str)
			AND gt.sender_address_id IS NULL
			AND gt.base_asset_id = $1
			`, *baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	_, err = database.DbConnPgx.Exec(ctx, `
			UPDATE geth_transfers as gt SET
			to_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gt.to_address) = LOWER(ga.address_str)
			AND gt.to_address_id IS NULL
			AND gt.base_asset_id = $1
			`, *baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	_, err = database.DbConnPgx.Exec(ctx, `UPDATE geth_transfers as gt SET
			asset_id = assets.id
			from assets as assets
			WHERE LOWER(gt.token_address) = LOWER(assets.contract_address)
			AND gt.asset_id IS NULL
			AND gt.base_asset_id = $1
	`, *baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetNullAddressStrsFromTransfers(baseAssetID *int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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

func UpdateGethTransfersAssetIDs() error {
	// update address ids from existing addresses in geth_addresses
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `
		UPDATE geth_transfers as gt SET
		asset_id = assets.id from assets as assets
			WHERE LOWER(gt.token_address) = LOWER(assets.contract_address) AND
			gt.asset_id is NULL
	`,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
