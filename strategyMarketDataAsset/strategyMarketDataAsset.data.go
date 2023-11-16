package strategymarketdataasset

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetStrategyMarketDataAsset(strategyMarketDataAssetID int) (*StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM strategy_market_data_assets 
	WHERE id = $1`, strategyMarketDataAssetID)

	strategyMarketDataAsset := &StrategyMarketDataAsset{}
	err := row.Scan(
		&strategyMarketDataAsset.ID,
		&strategyMarketDataAsset.StrategyID,
		&strategyMarketDataAsset.BaseAssetID,
		&strategyMarketDataAsset.QuoteAssetID,
		&strategyMarketDataAsset.Name,
		&strategyMarketDataAsset.UUID,
		&strategyMarketDataAsset.AlternateName,
		&strategyMarketDataAsset.StartDate,
		&strategyMarketDataAsset.EndDate,
		&strategyMarketDataAsset.Ticker,
		&strategyMarketDataAsset.Description,
		&strategyMarketDataAsset.SourceID,
		&strategyMarketDataAsset.FrequencyID,
		&strategyMarketDataAsset.CreatedBy,
		&strategyMarketDataAsset.CreatedAt,
		&strategyMarketDataAsset.UpdatedBy,
		&strategyMarketDataAsset.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return strategyMarketDataAsset, nil
}

func GetTopTenStrategies() ([]StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,
	strategy_id,
	base_asset_id,
	quote_asset_id,
	name,
	uuid,
	alternate_name,
	start_date,
	end_date,
	ticker,
	description,
	source_id ,
	frequency_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM strategy_market_data_assets 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategyMarketDataAssets := make([]StrategyMarketDataAsset, 0)
	for results.Next() {
		var strategyMarketDataAsset StrategyMarketDataAsset
		results.Scan(
			&strategyMarketDataAsset.ID,
			&strategyMarketDataAsset.StrategyID,
			&strategyMarketDataAsset.BaseAssetID,
			&strategyMarketDataAsset.QuoteAssetID,
			&strategyMarketDataAsset.Name,
			&strategyMarketDataAsset.UUID,
			&strategyMarketDataAsset.AlternateName,
			&strategyMarketDataAsset.StartDate,
			&strategyMarketDataAsset.EndDate,
			&strategyMarketDataAsset.Ticker,
			&strategyMarketDataAsset.Description,
			&strategyMarketDataAsset.SourceID,
			&strategyMarketDataAsset.FrequencyID,
			&strategyMarketDataAsset.CreatedBy,
			&strategyMarketDataAsset.CreatedAt,
			&strategyMarketDataAsset.UpdatedBy,
			&strategyMarketDataAsset.UpdatedAt,
		)

		strategyMarketDataAssets = append(strategyMarketDataAssets, strategyMarketDataAsset)
	}
	return strategyMarketDataAssets, nil
}

func RemoveStrategyMarketDataAsset(strategyMarketDataAssetID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM strategy_market_data_assets WHERE id = $1`, strategyMarketDataAssetID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetStrategyMarketDataAssets(ids []int) ([]StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM strategy_market_data_assets`
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
	strategyMarketDataAssets := make([]StrategyMarketDataAsset, 0)
	for results.Next() {
		var strategyMarketDataAsset StrategyMarketDataAsset
		results.Scan(
			&strategyMarketDataAsset.ID,
			&strategyMarketDataAsset.StrategyID,
			&strategyMarketDataAsset.BaseAssetID,
			&strategyMarketDataAsset.QuoteAssetID,
			&strategyMarketDataAsset.Name,
			&strategyMarketDataAsset.UUID,
			&strategyMarketDataAsset.AlternateName,
			&strategyMarketDataAsset.StartDate,
			&strategyMarketDataAsset.EndDate,
			&strategyMarketDataAsset.Ticker,
			&strategyMarketDataAsset.Description,
			&strategyMarketDataAsset.SourceID,
			&strategyMarketDataAsset.FrequencyID,
			&strategyMarketDataAsset.CreatedBy,
			&strategyMarketDataAsset.CreatedAt,
			&strategyMarketDataAsset.UpdatedBy,
			&strategyMarketDataAsset.UpdatedAt,
		)

		strategyMarketDataAssets = append(strategyMarketDataAssets, strategyMarketDataAsset)
	}
	return strategyMarketDataAssets, nil
}

func GetStrategyMarketDataAssetsByUUIDs(UUIDList []string) ([]StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
	FROM strategy_market_data_assets
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategyMarketDataAssets := make([]StrategyMarketDataAsset, 0)
	for results.Next() {
		var strategyMarketDataAsset StrategyMarketDataAsset
		results.Scan(
			&strategyMarketDataAsset.ID,
			&strategyMarketDataAsset.StrategyID,
			&strategyMarketDataAsset.BaseAssetID,
			&strategyMarketDataAsset.QuoteAssetID,
			&strategyMarketDataAsset.Name,
			&strategyMarketDataAsset.UUID,
			&strategyMarketDataAsset.AlternateName,
			&strategyMarketDataAsset.StartDate,
			&strategyMarketDataAsset.EndDate,
			&strategyMarketDataAsset.Ticker,
			&strategyMarketDataAsset.Description,
			&strategyMarketDataAsset.SourceID,
			&strategyMarketDataAsset.FrequencyID,
			&strategyMarketDataAsset.CreatedBy,
			&strategyMarketDataAsset.CreatedAt,
			&strategyMarketDataAsset.UpdatedBy,
			&strategyMarketDataAsset.UpdatedAt,
		)

		strategyMarketDataAssets = append(strategyMarketDataAssets, strategyMarketDataAsset)
	}
	return strategyMarketDataAssets, nil
}

func GetStrategyMarketDataAssetsByStrategyID(strategyID *int) ([]StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
	FROM strategy_market_data_assets
	WHERE strategy_id = $1
	`, *strategyID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategyMarketDataAssets := make([]StrategyMarketDataAsset, 0)
	for results.Next() {
		var strategyMarketDataAsset StrategyMarketDataAsset
		results.Scan(
			&strategyMarketDataAsset.ID,
			&strategyMarketDataAsset.StrategyID,
			&strategyMarketDataAsset.BaseAssetID,
			&strategyMarketDataAsset.QuoteAssetID,
			&strategyMarketDataAsset.Name,
			&strategyMarketDataAsset.UUID,
			&strategyMarketDataAsset.AlternateName,
			&strategyMarketDataAsset.StartDate,
			&strategyMarketDataAsset.EndDate,
			&strategyMarketDataAsset.Ticker,
			&strategyMarketDataAsset.Description,
			&strategyMarketDataAsset.SourceID,
			&strategyMarketDataAsset.FrequencyID,
			&strategyMarketDataAsset.CreatedBy,
			&strategyMarketDataAsset.CreatedAt,
			&strategyMarketDataAsset.UpdatedBy,
			&strategyMarketDataAsset.UpdatedAt,
		)

		strategyMarketDataAssets = append(strategyMarketDataAssets, strategyMarketDataAsset)
	}
	return strategyMarketDataAssets, nil
}

func GetStartAndEndDateDiffStrategyMarketDataAssets(diffInDate int) ([]StrategyMarketDataAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
		id,
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
	FROM strategy_market_data_assets
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	strategyMarketDataAssets := make([]StrategyMarketDataAsset, 0)
	for results.Next() {
		var strategyMarketDataAsset StrategyMarketDataAsset
		results.Scan(
			&strategyMarketDataAsset.ID,
			&strategyMarketDataAsset.StrategyID,
			&strategyMarketDataAsset.BaseAssetID,
			&strategyMarketDataAsset.QuoteAssetID,
			&strategyMarketDataAsset.Name,
			&strategyMarketDataAsset.UUID,
			&strategyMarketDataAsset.AlternateName,
			&strategyMarketDataAsset.StartDate,
			&strategyMarketDataAsset.EndDate,
			&strategyMarketDataAsset.Ticker,
			&strategyMarketDataAsset.Description,
			&strategyMarketDataAsset.SourceID,
			&strategyMarketDataAsset.FrequencyID,
			&strategyMarketDataAsset.CreatedBy,
			&strategyMarketDataAsset.CreatedAt,
			&strategyMarketDataAsset.UpdatedBy,
			&strategyMarketDataAsset.UpdatedAt,
		)

		strategyMarketDataAssets = append(strategyMarketDataAssets, strategyMarketDataAsset)
	}
	return strategyMarketDataAssets, nil
}

func UpdateStrategyMarketDataAsset(strategyMarketDataAsset StrategyMarketDataAsset) error {
	// if the strategyMarketDataAsset id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if strategyMarketDataAsset.ID == nil || *strategyMarketDataAsset.ID == 0 {
		return errors.New("strategyMarketDataAsset has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE strategy_market_data_assets SET 
		strategy_id=$1, 
		base_asset_id=$2, 
		quote_asset_id=$3, 
		name=$4, 
		uuid=$5, 
		alternate_name=$6, 
		start_date=$7, 
		end_date=$8, 
		ticker=$9, 
		description=$10, 
		source_id=$11, 
		frequency_id=$12, 
		updated_by=$13, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$14`,
		strategyMarketDataAsset.StrategyID,    //1
		strategyMarketDataAsset.BaseAssetID,   //2
		strategyMarketDataAsset.QuoteAssetID,  //3
		strategyMarketDataAsset.Name,          //4
		strategyMarketDataAsset.UUID,          //5
		strategyMarketDataAsset.AlternateName, //6
		strategyMarketDataAsset.StartDate,     //7
		strategyMarketDataAsset.EndDate,       //8
		strategyMarketDataAsset.Ticker,        //9
		strategyMarketDataAsset.Description,   //10
		strategyMarketDataAsset.SourceID,      //11
		strategyMarketDataAsset.FrequencyID,   //12
		strategyMarketDataAsset.UpdatedBy,     //13
		strategyMarketDataAsset.ID)            //14
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertStrategyMarketDataAsset(strategyMarketDataAsset StrategyMarketDataAsset) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO strategy_market_data_assets 
	(
		strategy_id,
		base_asset_id,
		quote_asset_id,
		name,
		uuid,
		alternate_name,
		start_date,
		end_date,
		ticker,
		description,
		source_id ,
		frequency_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at
		) VALUES (
			$1,
			$2, 
			$3, 
			$4, 
			uuid_generate_v4(), 
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
		&strategyMarketDataAsset.StrategyID,   //1
		&strategyMarketDataAsset.BaseAssetID,  //2
		&strategyMarketDataAsset.QuoteAssetID, //3
		&strategyMarketDataAsset.Name,         //4
		// &strategyMarketDataAsset.UUID,                            //5
		&strategyMarketDataAsset.AlternateName, //5
		&strategyMarketDataAsset.StartDate,     //6
		&strategyMarketDataAsset.EndDate,       //7
		&strategyMarketDataAsset.Ticker,        //8
		&strategyMarketDataAsset.Description,   //9
		&strategyMarketDataAsset.SourceID,      //10
		&strategyMarketDataAsset.FrequencyID,   //11
		&strategyMarketDataAsset.CreatedBy,     //12
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertStrategyMarketDataAssets(strategyMarketDataAssets []StrategyMarketDataAsset) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range strategyMarketDataAssets {
		strategyMarketDataAsset := strategyMarketDataAssets[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(strategyMarketDataAsset.UUID)
		row := []interface{}{
			*strategyMarketDataAsset.StrategyID,   //1
			*strategyMarketDataAsset.BaseAssetID,  //2
			*strategyMarketDataAsset.QuoteAssetID, //3
			strategyMarketDataAsset.Name,          //4
			uuidString,                            //5
			strategyMarketDataAsset.AlternateName, //6
			&strategyMarketDataAsset.StartDate,    //7
			&strategyMarketDataAsset.EndDate,      //8
			strategyMarketDataAsset.Ticker,        //9
			strategyMarketDataAsset.Description,   //10
			*strategyMarketDataAsset.SourceID,     //11
			*strategyMarketDataAsset.FrequencyID,  //12
			strategyMarketDataAsset.CreatedBy,     //13
			&now,                                  //14
			strategyMarketDataAsset.CreatedBy,     //15
			&now,                                  //16
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"strategy_market_data_assets"},
		[]string{
			"strategy_id",    //1
			"base_asset_id",  //2
			"quote_asset_id", //3
			"name",           //4
			"uuid",           //5
			"alternate_name", //6
			"start_date",     //7
			"end_date",       //8
			"ticker",         //9
			"description",    //10
			"source_id",      //11
			"frequency_id",   //12
			"created_by",     //13
			"created_at",     //14
			"updated_by",     //15
			"updated_at",     //16
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
