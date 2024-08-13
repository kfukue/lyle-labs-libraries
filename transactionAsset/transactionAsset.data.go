package transactionasset

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

func GetTransactionAsset(dbConnPgx utils.PgxIface, transactionID, assetID *int) (*TransactionAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		transaction_id,  
		asset_id,
		uuid, 
		name,
		alternate_name,
		description,
		quantity,
		quantity_usd,
		market_data_id,
		manual_exchange_rate_usd,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM transaction_assets 
	WHERE transaction_id = $1
	AND asset_id = $2
	`, *transactionID, *assetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	transactionAsset, err := pgx.CollectOneRow(row, pgx.RowToStructByName[TransactionAsset])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &transactionAsset, nil
}

func GetTransactionAssetByUUID(dbConnPgx utils.PgxIface, transactionAssetUUID string) (*TransactionAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		transaction_id,  
		asset_id,
		uuid, 
		name,
		alternate_name,
		description,
		quantity,
		quantity_usd,
		market_data_id,
		manual_exchange_rate_usd,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM transaction_assets 
	WHERE text(uuid) = $1
	`, transactionAssetUUID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	transactionAsset, err := pgx.CollectOneRow(row, pgx.RowToStructByName[TransactionAsset])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &transactionAsset, nil
}

func GetTransactionAssetsByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]TransactionAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		transaction_id,  
		asset_id,
		uuid, 
		name,
		alternate_name,
		description,
		quantity,
		quantity_usd,
		market_data_id,
		manual_exchange_rate_usd,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM transaction_assets
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	transactionAssetList, err := pgx.CollectRows(results, pgx.RowToStructByName[TransactionAsset])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return transactionAssetList, nil
}

func RemoveTransactionAsset(dbConnPgx utils.PgxIface, transactionID, assetID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveTransactionAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM transaction_assets WHERE transaction_id = $1 AND asset_id = $2`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *transactionID, *assetID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemoveTransactionAssetByUUID(dbConnPgx utils.PgxIface, transactionAssetUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveTransactionAssetByUUID DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM transaction_assets WHERE text(uuid) = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, transactionAssetUUID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetTransactionAssets(dbConnPgx utils.PgxIface) ([]TransactionAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		transaction_id,  
		asset_id,
		uuid, 
		name,
		alternate_name,
		description,
		quantity,
		quantity_usd,
		market_data_id,
		manual_exchange_rate_usd,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM transaction_assets`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	transactionAssets, err := pgx.CollectRows(results, pgx.RowToStructByName[TransactionAsset])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return transactionAssets, nil
}

func UpdateTransactionAsset(dbConnPgx utils.PgxIface, transactionAsset *TransactionAsset) error {
	// if the transactionAsset id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (transactionAsset.TransactionID == nil || *transactionAsset.TransactionID <= 0) || (transactionAsset.AssetID == nil || *transactionAsset.AssetID <= 0) {
		return errors.New("transactionAsset has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateTransactionAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE transaction_assets SET 
		name=$1,
		alternate_name=$2,
		description=$3,
		quantity=$4,
		quantity_usd=$5,
		market_data_id=$6,
		manual_exchange_rate_usd=$7,
		updated_by=$8, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE transaction_id=$9 AND asset_id=$10`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		transactionAsset.Name,                  //1
		transactionAsset.AlternateName,         //2
		transactionAsset.Description,           //3
		transactionAsset.Quantity,              //4
		transactionAsset.QuantityUSD,           //5
		transactionAsset.MarketDataID,          //6
		transactionAsset.ManualExchangeRateUSD, //7
		transactionAsset.UpdatedBy,             //8
		transactionAsset.TransactionID,         //9
		transactionAsset.AssetID,               //10
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func UpdateTransactionAssetByUUID(dbConnPgx utils.PgxIface, transactionAsset *TransactionAsset) error {
	// if the transactionAsset id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (transactionAsset.TransactionID == nil || *transactionAsset.TransactionID == 0) || (transactionAsset.AssetID == nil || *transactionAsset.AssetID == 0) {
		return errors.New("transactionAsset has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateTransactionAssetByUUID DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE transaction_assets SET 
		name=$1,
		alternate_name=$2,
		description=$3,
		quantity=$4,
		quantity_usd=$5,
		market_data_id=$6,
		manual_exchange_rate_usd=$7,
		updated_by=$8,
		updated_at=current_timestamp at time zone 'UTC'
		WHERE text(uuid) = $9
		`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		transactionAsset.Name,                  //1
		transactionAsset.AlternateName,         //2
		transactionAsset.Description,           //3
		transactionAsset.Quantity,              //4
		transactionAsset.QuantityUSD,           //5
		transactionAsset.MarketDataID,          //6
		transactionAsset.ManualExchangeRateUSD, //7
		transactionAsset.UpdatedBy,             //8
		transactionAsset.UUID,                  //9
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertTransactionAsset(dbConnPgx utils.PgxIface, transactionAsset *TransactionAsset) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertTransactionAsset DbConn.Begin   %s", err.Error())
		return -1, -1, err
	}
	var transactionID int
	var assetID int
	transactionAssetUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	if transactionAsset.UUID == "" {
		transactionAsset.UUID = transactionAssetUUID.String()
	}
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO transaction_assets  
	(
		transaction_id,  
		asset_id,
		uuid, 
		name,
		alternate_name,
		description,
		quantity,
		quantity_usd,
		market_data_id,
		manual_exchange_rate_usd,
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
			current_timestamp at time zone 'UTC',
			$11,
			current_timestamp at time zone 'UTC'
		)
		RETURNING transaction_id, asset_id`,
		transactionAsset.TransactionID,         //1
		transactionAsset.AssetID,               //2
		transactionAsset.UUID,                  //3
		transactionAsset.Name,                  //4
		transactionAsset.AlternateName,         //5
		transactionAsset.Description,           //6
		transactionAsset.Quantity,              //7
		transactionAsset.QuantityUSD,           //8
		transactionAsset.MarketDataID,          //9
		transactionAsset.ManualExchangeRateUSD, //10
		transactionAsset.CreatedBy,             //11
	).Scan(&transactionID, &assetID)
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
	return int(transactionID), int(assetID), nil
}

func InsertTransactionAssets(dbConnPgx utils.PgxIface, transactionAssets []TransactionAsset) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range transactionAssets {
		transactionAsset := transactionAssets[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(transactionAsset.UUID)
		row := []interface{}{
			transactionAsset.TransactionID,         //1
			transactionAsset.AssetID,               //2
			uuidString,                             //3
			transactionAsset.Name,                  //4
			transactionAsset.AlternateName,         //5
			transactionAsset.Description,           //6
			transactionAsset.Quantity,              //7
			transactionAsset.QuantityUSD,           //8
			transactionAsset.MarketDataID,          //9
			transactionAsset.ManualExchangeRateUSD, //10
			transactionAsset.CreatedBy,             //11
			&transactionAsset.CreatedAt,            //12
			transactionAsset.CreatedBy,             //13
			&now,                                   //14
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"transaction_assets"},
		[]string{
			"transaction_id",           //1
			"asset_id",                 //2
			"uuid",                     //3
			"name",                     //4
			"alternate_name",           //5
			"description",              //6
			"quantity",                 //7
			"quantity_usd",             //8
			"market_data_id",           //9
			"manual_exchange_rate_usd", //10
			"created_by",               //11
			"created_at",               //12
			"updated_by",               //13
			"updated_at",               //14
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertTransactionAssets: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetTransactionAssetListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]TransactionAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		transaction_id,  
		asset_id,
		uuid, 
		name,
		alternate_name,
		description,
		quantity,
		quantity_usd,
		market_data_id,
		manual_exchange_rate_usd,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM transaction_assets 
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
	transactionAssets, err := pgx.CollectRows(results, pgx.RowToStructByName[TransactionAsset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionAssets, nil
}

func GetTotalTransactionAssetsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM transaction_assets
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
