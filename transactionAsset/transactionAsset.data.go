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
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/lib/pq"
)

func GetTransactionAsset(transactionID int, assetID int) (*TransactionAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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
	`, transactionID, assetID)

	transactionAsset := &TransactionAsset{}
	err := row.Scan(
		&transactionAsset.TransactionID,
		&transactionAsset.AssetID,
		&transactionAsset.UUID,
		&transactionAsset.Name,
		&transactionAsset.AlternateName,
		&transactionAsset.Description,
		&transactionAsset.Quantity,
		&transactionAsset.QuantityUSD,
		&transactionAsset.MarketDataID,
		&transactionAsset.ManualExchangeRateUSD,
		&transactionAsset.CreatedBy,
		&transactionAsset.CreatedAt,
		&transactionAsset.UpdatedBy,
		&transactionAsset.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionAsset, nil
}

func GetTransactionAssetByUUID(transactionAssetUUID string) (*TransactionAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
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

	transactionAsset := &TransactionAsset{}
	err := row.Scan(
		&transactionAsset.TransactionID,
		&transactionAsset.AssetID,
		&transactionAsset.UUID,
		&transactionAsset.Name,
		&transactionAsset.AlternateName,
		&transactionAsset.Description,
		&transactionAsset.Quantity,
		&transactionAsset.QuantityUSD,
		&transactionAsset.MarketDataID,
		&transactionAsset.ManualExchangeRateUSD,
		&transactionAsset.CreatedBy,
		&transactionAsset.CreatedAt,
		&transactionAsset.UpdatedBy,
		&transactionAsset.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionAsset, nil
}

func GetTransactionAssetsByUUIDs(UUIDList []string) ([]TransactionAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	transactionAssetList := make([]TransactionAsset, 0)
	for results.Next() {
		var transactionAsset TransactionAsset
		results.Scan(
			&transactionAsset.TransactionID,
			&transactionAsset.AssetID,
			&transactionAsset.UUID,
			&transactionAsset.Name,
			&transactionAsset.AlternateName,
			&transactionAsset.Description,
			&transactionAsset.Quantity,
			&transactionAsset.QuantityUSD,
			&transactionAsset.MarketDataID,
			&transactionAsset.ManualExchangeRateUSD,
			&transactionAsset.CreatedBy,
			&transactionAsset.CreatedAt,
			&transactionAsset.UpdatedBy,
			&transactionAsset.UpdatedAt,
		)

		transactionAssetList = append(transactionAssetList, transactionAsset)
	}
	return transactionAssetList, nil
}

func GetTopTenTransactionAssets() ([]TransactionAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	transactionAssets := make([]TransactionAsset, 0)
	for results.Next() {
		var transactionAsset TransactionAsset
		results.Scan(
			&transactionAsset.TransactionID,
			&transactionAsset.AssetID,
			&transactionAsset.UUID,
			&transactionAsset.Name,
			&transactionAsset.AlternateName,
			&transactionAsset.Description,
			&transactionAsset.Quantity,
			&transactionAsset.QuantityUSD,
			&transactionAsset.MarketDataID,
			&transactionAsset.ManualExchangeRateUSD,
			&transactionAsset.CreatedBy,
			&transactionAsset.CreatedAt,
			&transactionAsset.UpdatedBy,
			&transactionAsset.UpdatedAt,
		)

		transactionAssets = append(transactionAssets, transactionAsset)
	}
	return transactionAssets, nil
}

func RemoveTransactionAsset(transactionID int, assetID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM transaction_assets WHERE 
	transaction_id = $1 AND asset_id =$2`, transactionID, assetID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveTransactionAssetByUUID(transactionAssetUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM transaction_assets WHERE 
		WHERE text(uuid) = $1`,
		transactionAssetUUID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetTransactionAssetList() ([]TransactionAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	transactionAssets := make([]TransactionAsset, 0)
	for results.Next() {
		var transactionAsset TransactionAsset
		results.Scan(
			&transactionAsset.TransactionID,
			&transactionAsset.AssetID,
			&transactionAsset.UUID,
			&transactionAsset.Name,
			&transactionAsset.AlternateName,
			&transactionAsset.Description,
			&transactionAsset.Quantity,
			&transactionAsset.QuantityUSD,
			&transactionAsset.MarketDataID,
			&transactionAsset.ManualExchangeRateUSD,
			&transactionAsset.CreatedBy,
			&transactionAsset.CreatedAt,
			&transactionAsset.UpdatedBy,
			&transactionAsset.UpdatedAt,
		)

		transactionAssets = append(transactionAssets, transactionAsset)
	}
	return transactionAssets, nil
}

func UpdateTransactionAsset(transactionAsset TransactionAsset) error {
	// if the transactionAsset id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (transactionAsset.TransactionID == nil || *transactionAsset.TransactionID == 0) || (transactionAsset.AssetID == nil || *transactionAsset.AssetID == 0) {
		return errors.New("transactionAsset has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE transaction_assets SET 
		name=$1,
		alternate_name=$2,
		description=$3,
		quantity=$4,
		quantity_usd=$5,
		market_data_id=$6,
		manual_exchange_rate_usd=$7,
		updated_by=$8, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE transaction_id=$9 AND asset_id=$10`,

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
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateTransactionAssetByUUID(transactionAsset TransactionAsset) error {
	// if the transactionAsset id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (transactionAsset.TransactionID == nil || *transactionAsset.TransactionID == 0) || (transactionAsset.AssetID == nil || *transactionAsset.AssetID == 0) {
		return errors.New("transactionAsset has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE transaction_assets SET 
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
		`,
		transactionAsset.Name,                  //1
		transactionAsset.AlternateName,         //2
		transactionAsset.Description,           //3
		transactionAsset.Quantity,              //4
		transactionAsset.QuantityUSD,           //5
		transactionAsset.MarketDataID,          //6
		transactionAsset.ManualExchangeRateUSD, //7
		transactionAsset.UpdatedBy,             //8
		transactionAsset.UUID,                  //9
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertTransactionAsset(transactionAsset TransactionAsset) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var MarketDataID int
	var JobID int
	transactionAssetUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	if transactionAsset.UUID == "" {
		transactionAsset.UUID = transactionAssetUUID.String()
	}
	err = database.DbConnPgx.QueryRow(ctx, `INSERT INTO transaction_assets  
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
		RETURNING market_data_id, asset_id`,
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
	).Scan(&MarketDataID, &JobID)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(MarketDataID), int(JobID), nil
}

func InsertTransactionAssets(transactionAssets []TransactionAsset) error {
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

			*transactionAsset.TransactionID,         //1
			*transactionAsset.AssetID,               //2
			uuidString,                              //3
			transactionAsset.Name,                   //4
			transactionAsset.AlternateName,          //5
			transactionAsset.Description,            //6
			*transactionAsset.Quantity,              //7
			*transactionAsset.QuantityUSD,           //8
			*transactionAsset.MarketDataID,          //9
			*transactionAsset.ManualExchangeRateUSD, //10
			transactionAsset.CreatedBy,              //11
			&transactionAsset.CreatedAt,             //12
			transactionAsset.CreatedBy,              //13
			&now,                                    //14
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB

	copyCount, err := database.DbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}

	return nil
}
