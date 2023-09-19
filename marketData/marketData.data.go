package marketdata

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
	decimal "github.com/shopspring/decimal"
)

func GetMarketData(marketDataID int) (*MarketData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	asset_id,
	open_usd,
	close_usd,
	high_usd,
	low_usd,
	price_usd,
	volume_usd,
	market_cap_usd,
	ticker,
	description,
	interval_id,
	market_data_type_id,
	source_id,

	total_supply,
	max_supply,
	circulating_supply,
	sparkline_7d,

	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	FROM market_data 
	WHERE id = $1`, marketDataID)

	marketData := &MarketData{}
	err := row.Scan(
		&marketData.ID,
		&marketData.UUID,
		&marketData.Name,
		&marketData.AlternateName,
		&marketData.StartDate,
		&marketData.EndDate,
		&marketData.AssetID,
		&marketData.OpenUSD,
		&marketData.CloseUSD,
		&marketData.HighUSD,
		&marketData.LowUSD,
		&marketData.PriceUSD,
		&marketData.VolumeUSD,
		&marketData.MarketCapUSD,
		&marketData.Ticker,
		&marketData.Description,
		&marketData.IntervalID,
		&marketData.MarketDataTypeID,
		&marketData.SourceID,
		&marketData.TotalSupply,
		&marketData.MaxSupply,
		&marketData.CirculatingSupply,
		&marketData.Sparkline7d,
		&marketData.CreatedBy,
		&marketData.CreatedAt,
		&marketData.UpdatedBy,
		&marketData.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return marketData, nil
}

func GetTopTenMarketDatas() ([]MarketData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	asset_id,
	open_usd,
	close_usd,
	high_usd,
	low_usd,
	price_usd,
	volume_usd,
	market_cap_usd,
	ticker,
	description,
	interval_id,
	market_data_type_id,
	source_id,
	total_supply,
	max_supply,
	circulating_supply,
	sparkline_7d,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM market_data 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	marketDataList := make([]MarketData, 0)
	for results.Next() {
		var marketData MarketData
		results.Scan(
			&marketData.ID,
			&marketData.UUID,
			&marketData.Name,
			&marketData.AlternateName,
			&marketData.StartDate,
			&marketData.EndDate,
			&marketData.AssetID,
			&marketData.OpenUSD,
			&marketData.CloseUSD,
			&marketData.HighUSD,
			&marketData.LowUSD,
			&marketData.PriceUSD,
			&marketData.VolumeUSD,
			&marketData.MarketCapUSD,
			&marketData.Ticker,
			&marketData.Description,
			&marketData.IntervalID,
			&marketData.MarketDataTypeID,
			&marketData.SourceID,
			&marketData.TotalSupply,
			&marketData.MaxSupply,
			&marketData.CirculatingSupply,
			&marketData.Sparkline7d,
			&marketData.CreatedBy,
			&marketData.CreatedAt,
			&marketData.UpdatedBy,
			&marketData.UpdatedAt,
		)

		marketDataList = append(marketDataList, marketData)
	}
	return marketDataList, nil
}

func RemoveMarketData(marketDataID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM market_data WHERE id = $1`, marketDataID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveMarketDataFromBaseAssetBetweenDates(assetID int, startDate, endDate time.Time) error {
	log.Println(fmt.Sprintf("start : %s end : %s", startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)))
	err := RemoveMarketDataQuoteFromBaseAssetBetweenDates(assetID, startDate, endDate)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err = database.DbConnPgx.Query(ctx, `DELETE FROM market_data 
		WHERE asset_id = $1
			AND start_date BETWEEN $2 and $3
	`, assetID, startDate, endDate)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveMarketDataQuoteFromBaseAssetBetweenDates(assetID int, startDate, endDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `
		DELETE FROM market_data_quotes
		WHERE market_data_id IN (
			SELECT ID
			FROM market_data
			WHERE asset_id = $1
			AND start_date BETWEEN $2 and $3
		);
	`, assetID, startDate, endDate)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetMarketDataList(ids []int) ([]MarketData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	asset_id,
	open_usd,
	close_usd,
	high_usd,
	low_usd,
	price_usd,
	volume_usd,
	market_cap_usd,
	ticker,
	description,
	interval_id,
	market_data_type_id,
	source_id,
	total_supply,
	max_supply,
	circulating_supply,
	sparkline_7d,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM market_data`
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
	marketDataList := make([]MarketData, 0)
	for results.Next() {
		var marketData MarketData
		results.Scan(
			&marketData.ID,
			&marketData.UUID,
			&marketData.Name,
			&marketData.AlternateName,
			&marketData.StartDate,
			&marketData.EndDate,
			&marketData.AssetID,
			&marketData.OpenUSD,
			&marketData.CloseUSD,
			&marketData.HighUSD,
			&marketData.LowUSD,
			&marketData.PriceUSD,
			&marketData.VolumeUSD,
			&marketData.MarketCapUSD,
			&marketData.Ticker,
			&marketData.Description,
			&marketData.IntervalID,
			&marketData.MarketDataTypeID,
			&marketData.SourceID,
			&marketData.TotalSupply,
			&marketData.MaxSupply,
			&marketData.CirculatingSupply,
			&marketData.Sparkline7d,
			&marketData.CreatedBy,
			&marketData.CreatedAt,
			&marketData.UpdatedBy,
			&marketData.UpdatedAt,
		)

		marketDataList = append(marketDataList, marketData)
	}
	return marketDataList, nil
}

func GetMarketDataListByUUIDs(UUIDList []string) ([]MarketData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	asset_id,
	open_usd,
	close_usd,
	high_usd,
	low_usd,
	price_usd,
	volume_usd,
	market_cap_usd,
	ticker,
	description,
	interval_id,
	market_data_type_id,
	source_id,
	total_supply,
	max_supply,
	circulating_supply,
	sparkline_7d,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM market_data
	WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	marketDataList := make([]MarketData, 0)
	for results.Next() {
		var marketData MarketData
		results.Scan(
			&marketData.ID,
			&marketData.UUID,
			&marketData.Name,
			&marketData.AlternateName,
			&marketData.StartDate,
			&marketData.EndDate,
			&marketData.AssetID,
			&marketData.OpenUSD,
			&marketData.CloseUSD,
			&marketData.HighUSD,
			&marketData.LowUSD,
			&marketData.PriceUSD,
			&marketData.VolumeUSD,
			&marketData.MarketCapUSD,
			&marketData.Ticker,
			&marketData.Description,
			&marketData.IntervalID,
			&marketData.MarketDataTypeID,
			&marketData.SourceID,
			&marketData.TotalSupply,
			&marketData.MaxSupply,
			&marketData.CirculatingSupply,
			&marketData.Sparkline7d,
			&marketData.CreatedBy,
			&marketData.CreatedAt,
			&marketData.UpdatedBy,
			&marketData.UpdatedAt,
		)

		marketDataList = append(marketDataList, marketData)
	}
	return marketDataList, nil
}

func GetStartAndEndDateDiffMarketDataList(diffInDate int) ([]MarketData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	start_date,
	end_date,
	asset_id,
	open_usd,
	close_usd,
	high_usd,
	low_usd,
	price_usd,
	volume_usd,
	market_cap_usd,
	ticker,
	description,
	interval_id,
	market_data_type_id,
	source_id,
	total_supply,
	max_supply,
	circulating_supply,
	sparkline_7d,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM market_data
	WHERE DATE_PART('day', AGE(start_date, end_date)) =$1
	`, diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	marketDataList := make([]MarketData, 0)
	for results.Next() {
		var marketData MarketData
		results.Scan(
			&marketData.ID,
			&marketData.UUID,
			&marketData.Name,
			&marketData.AlternateName,
			&marketData.StartDate,
			&marketData.EndDate,
			&marketData.AssetID,
			&marketData.OpenUSD,
			&marketData.CloseUSD,
			&marketData.HighUSD,
			&marketData.LowUSD,
			&marketData.PriceUSD,
			&marketData.VolumeUSD,
			&marketData.MarketCapUSD,
			&marketData.Ticker,
			&marketData.Description,
			&marketData.IntervalID,
			&marketData.MarketDataTypeID,
			&marketData.SourceID,
			&marketData.TotalSupply,
			&marketData.MaxSupply,
			&marketData.CirculatingSupply,
			&marketData.Sparkline7d,
			&marketData.CreatedBy,
			&marketData.CreatedAt,
			&marketData.UpdatedBy,
			&marketData.UpdatedAt,
		)

		marketDataList = append(marketDataList, marketData)
	}
	return marketDataList, nil
}
func GetLatestLiveMarketData() ([]MarketDataQuoteResults, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
	SELECT 
	m.start_date,
	mj.end_date,
	baseA.name, 
	baseA.ticker,
	quoteA.name,
	quoteA.ticker,
	mq.market_data_id,
	mq.base_asset_id,
	mq.quote_asset_id,
	mq.name,
	mq.uuid,
	mq.alternate_name,
	mq.open,
	mq.close,
	mq.high_24h,
	mq.low_24h,
	mq.price,
	mq.volume,
	mq.market_cap,
	mq.ticker,
	mq.description,
	mq.source_id,
	mq.fully_diluted_valution,
	mq.ath,
	mq.ath_date,
	mq.atl,
	mq.atl_date,
	mq.price_change_1h,
	mq.price_change_24h,
	mq.price_change_7d,
	mq.price_change_30d,
	mq.price_change_60d,
	mq.price_change_200d,
	mq.price_change_1y,
	mq.created_by,
	mq.created_at,
	mq.updated_by,
	mq.updated_at
	FROM market_data_quotes mq
	LEFT JOIN market_data m ON m.id = mq.market_data_id
	LEFT JOIN market_data_jobs mj on m.id = mj.market_data_id
	LEFT JOIN assets baseA on mq.base_asset_id = baseA.id
	LEFT JOIN assets quoteA on mq.quote_asset_id = quoteA.id
	WHERE mq.market_data_id IN (
		SELECT market_data_id FROM market_data_jobs
		where job_id = (
			SELECT ID FROM jobs WHERE
			job_category_id = 56 
			ORDER BY end_date desc
			LIMIT 1
			)
		)
	AND mq.quote_asset_id = 34
	ORDER BY baseA.ticker, m.start_date desc`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	MarketDataQuoteResultsList := make([]MarketDataQuoteResults, 0)
	for results.Next() {
		var MarketDataQuoteResults MarketDataQuoteResults
		results.Scan(
			&MarketDataQuoteResults.StartDate,
			&MarketDataQuoteResults.EndDate,
			&MarketDataQuoteResults.BaseAssetName,
			&MarketDataQuoteResults.BaseAssetTicker,
			&MarketDataQuoteResults.QuoteAssetName,
			&MarketDataQuoteResults.QuoteAssetTicker,
			&MarketDataQuoteResults.MarketDataID,
			&MarketDataQuoteResults.BaseAssetID,
			&MarketDataQuoteResults.QuoteAssetID,
			&MarketDataQuoteResults.UUID,
			&MarketDataQuoteResults.Name,
			&MarketDataQuoteResults.AlternateName,
			&MarketDataQuoteResults.Open,
			&MarketDataQuoteResults.Close,
			&MarketDataQuoteResults.High24h,
			&MarketDataQuoteResults.Low24h,
			&MarketDataQuoteResults.Price,
			&MarketDataQuoteResults.Volume,
			&MarketDataQuoteResults.MarketCap,
			&MarketDataQuoteResults.Ticker,
			&MarketDataQuoteResults.Description,
			&MarketDataQuoteResults.SourceID,
			&MarketDataQuoteResults.FullyDilutedValution,
			&MarketDataQuoteResults.Ath,
			&MarketDataQuoteResults.AthDate,
			&MarketDataQuoteResults.Atl,
			&MarketDataQuoteResults.AtlDate,
			&MarketDataQuoteResults.PriceChange1h,
			&MarketDataQuoteResults.PriceChange24h,
			&MarketDataQuoteResults.PriceChange7d,
			&MarketDataQuoteResults.PriceChange30d,
			&MarketDataQuoteResults.PriceChange60d,
			&MarketDataQuoteResults.PriceChange200d,
			&MarketDataQuoteResults.PriceChange1y,
			&MarketDataQuoteResults.CreatedBy,
			&MarketDataQuoteResults.CreatedAt,
			&MarketDataQuoteResults.UpdatedBy,
			&MarketDataQuoteResults.UpdatedAt,
		)

		MarketDataQuoteResultsList = append(MarketDataQuoteResultsList, MarketDataQuoteResults)
	}
	return MarketDataQuoteResultsList, nil
}

func GetAllMarketDataFromStrategyID(strategyID *int) ([]MarketDataQuoteResults, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
	SELECT 
	m.start_date,
	mj.end_date,
	baseA.name, 
	baseA.ticker,
	quoteA.name,
	quoteA.ticker,
	mq.market_data_id,
	mq.base_asset_id,
	mq.quote_asset_id,
	mq.name,
	mq.uuid,
	mq.alternate_name,
	mq.open,
	mq.close,
	mq.high_24h,
	mq.low_24h,
	mq.price,
	mq.volume,
	mq.market_cap,
	mq.ticker,
	mq.description,
	mq.source_id,
	mq.fully_diluted_valution,
	mq.ath,
	mq.ath_date,
	mq.atl,
	mq.atl_date,
	mq.price_change_1h,
	mq.price_change_24h,
	mq.price_change_7d,
	mq.price_change_30d,
	mq.price_change_60d,
	mq.price_change_200d,
	mq.price_change_1y,
	mq.created_by,
	mq.created_at,
	mq.updated_by,
	mq.updated_at
	FROM market_data_quotes mq
	LEFT JOIN market_data m ON m.id = mq.market_data_id
	LEFT JOIN market_data_jobs mj on m.id = mj.market_data_id
	LEFT JOIN assets baseA on mq.base_asset_id = baseA.id
	LEFT JOIN assets quoteA on mq.quote_asset_id = quoteA.id
	WHERE mq.market_data_id IN (
		SELECT market_data_id 
			FROM market_data_jobs
			WHERE job_id IN (
			SELECT job_id FROM strategy_jobs
			where strategy_id = $1
		)
	)
	-- 	AND mq.quote_asset_id = 34 ignore usd since need ecd ptp / ptp where base is ptp
	ORDER BY m.start_date desc , baseA.ticker`, strategyID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	MarketDataQuoteResultsList := make([]MarketDataQuoteResults, 0)
	for results.Next() {
		var MarketDataQuoteResults MarketDataQuoteResults
		results.Scan(
			&MarketDataQuoteResults.StartDate,
			&MarketDataQuoteResults.EndDate,
			&MarketDataQuoteResults.BaseAssetName,
			&MarketDataQuoteResults.BaseAssetTicker,
			&MarketDataQuoteResults.QuoteAssetName,
			&MarketDataQuoteResults.QuoteAssetTicker,
			&MarketDataQuoteResults.MarketDataID,
			&MarketDataQuoteResults.BaseAssetID,
			&MarketDataQuoteResults.QuoteAssetID,
			&MarketDataQuoteResults.UUID,
			&MarketDataQuoteResults.Name,
			&MarketDataQuoteResults.AlternateName,
			&MarketDataQuoteResults.Open,
			&MarketDataQuoteResults.Close,
			&MarketDataQuoteResults.High24h,
			&MarketDataQuoteResults.Low24h,
			&MarketDataQuoteResults.Price,
			&MarketDataQuoteResults.Volume,
			&MarketDataQuoteResults.MarketCap,
			&MarketDataQuoteResults.Ticker,
			&MarketDataQuoteResults.Description,
			&MarketDataQuoteResults.SourceID,
			&MarketDataQuoteResults.FullyDilutedValution,
			&MarketDataQuoteResults.Ath,
			&MarketDataQuoteResults.AthDate,
			&MarketDataQuoteResults.Atl,
			&MarketDataQuoteResults.AtlDate,
			&MarketDataQuoteResults.PriceChange1h,
			&MarketDataQuoteResults.PriceChange24h,
			&MarketDataQuoteResults.PriceChange7d,
			&MarketDataQuoteResults.PriceChange30d,
			&MarketDataQuoteResults.PriceChange60d,
			&MarketDataQuoteResults.PriceChange200d,
			&MarketDataQuoteResults.PriceChange1y,
			&MarketDataQuoteResults.CreatedBy,
			&MarketDataQuoteResults.CreatedAt,
			&MarketDataQuoteResults.UpdatedBy,
			&MarketDataQuoteResults.UpdatedAt,
		)

		MarketDataQuoteResultsList = append(MarketDataQuoteResultsList, MarketDataQuoteResults)
	}
	return MarketDataQuoteResultsList, nil
}

func UpdateMarketData(marketData MarketData) error {
	// if the marketData id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if marketData.ID == nil || *marketData.ID == 0 {
		return errors.New("marketData has invalid ID")
	}
	layoutPostgres := utils.LayoutPostgres
	startDate := marketData.StartDate
	endDate := marketData.EndDate
	log.Println(fmt.Sprintf("Updating start: %s, end : %s", startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)))
	_, err := database.DbConnPgx.Query(ctx, `UPDATE market_data SET 
		name=$1,  
		alternate_name=$2, 
		start_date =$3,
		end_date =$4,
		asset_id=$5, 
		open_usd=$6, 
		close_usd=$7, 
		high_usd=$8, 
		low_usd=$9, 
		price_usd=$10, 
		volume_usd=$11, 
		market_cap_usd=$12, 
		ticker=$13, 
		description=$14, 
		interval_id=$15, 
		market_data_type_id=$16, 
		source_id=$17, 
		total_supply=$18,
		max_supply=$19,
		circulating_supply=$20,
		sparkline_7d=$21,

		updated_by=$22, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$23`,
		marketData.Name,                             //1
		marketData.AlternateName,                    //2
		marketData.StartDate.Format(layoutPostgres), //3
		marketData.EndDate.Format(layoutPostgres),   //4
		marketData.AssetID,                          //5
		marketData.OpenUSD,                          //6
		marketData.CloseUSD,                         //7
		marketData.HighUSD,                          //8
		marketData.LowUSD,                           //9
		marketData.PriceUSD,                         //10
		marketData.VolumeUSD,                        //11
		marketData.MarketCapUSD,                     //12
		marketData.Ticker,                           //13
		marketData.Description,                      //14
		marketData.IntervalID,                       //15
		marketData.MarketDataTypeID,                 //16
		marketData.SourceID,                         //17
		marketData.TotalSupply,                      //18
		marketData.MaxSupply,                        //19
		marketData.CirculatingSupply,                //20
		pq.Array(marketData.Sparkline7d),            //21
		marketData.UpdatedBy,                        //22
		marketData.ID)                               //23
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertMarketData(marketData MarketData) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	layoutPostgres := utils.LayoutPostgres
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO market_data 
	(
		name,  
		uuid,
		alternate_name, 
		start_date,
		end_date,
		asset_id,
		open_usd,
		close_usd,
		high_usd,
		low_usd,
		price_usd,
		volume_usd,
		market_cap_usd,
		ticker,
		description,
		interval_id,
		market_data_type_id,
		source_id,
		total_supply,
		max_supply,
		circulating_supply,
		sparkline_7d,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES (
			$1,
			uuid_generate_v4(), 
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
			$20,
			$21,
			$22,
			current_timestamp at time zone 'UTC',
			$22,
			current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		marketData.Name,                             //1
		marketData.AlternateName,                    //2
		marketData.StartDate.Format(layoutPostgres), //3
		marketData.EndDate.Format(layoutPostgres),   //4
		marketData.AssetID,                          //5
		marketData.OpenUSD,                          //6
		marketData.CloseUSD,                         //7
		marketData.HighUSD,                          //8
		marketData.LowUSD,                           //9
		marketData.PriceUSD,                         //10
		marketData.VolumeUSD,                        //11
		marketData.MarketCapUSD,                     //12
		marketData.Ticker,                           //13
		marketData.Description,                      //14
		marketData.IntervalID,                       //15
		marketData.MarketDataTypeID,                 //16
		marketData.SourceID,                         //17
		marketData.TotalSupply,                      //18
		marketData.MaxSupply,                        //19
		marketData.CirculatingSupply,                //20
		pq.Array(marketData.Sparkline7d),            //21
		marketData.CreatedBy,                        //22
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertMarketDataListManual(marketDataList []MarketData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range marketDataList {
		marketData := marketDataList[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(marketData.UUID)
		startDate := &pgtype.Date{}
		startDate.Time = marketData.StartDate
		endDate := &pgtype.Date{}
		endDate.Time = marketData.EndDate
		openUSD := utils.ConvertFloatToDecimal(marketData.OpenUSD)
		closeUSD := utils.ConvertFloatToDecimal(marketData.CloseUSD)
		highUSD := utils.ConvertFloatToDecimal(marketData.HighUSD)
		lowUSD := utils.ConvertFloatToDecimal(marketData.LowUSD)
		priceUSD := utils.ConvertFloatToDecimal(marketData.PriceUSD)
		volumeUSD := utils.ConvertFloatToDecimal(marketData.VolumeUSD)
		marketCapUSD := utils.ConvertFloatToDecimal(marketData.MarketCapUSD)
		totalSupply := utils.ConvertFloatToDecimal(marketData.TotalSupply)
		maxSupply := utils.ConvertFloatToDecimal(marketData.MaxSupply)
		circulatingSupply := utils.ConvertFloatToDecimal(marketData.CirculatingSupply)
		sparkline := marketData.Sparkline7d
		sparklineArray := make([]*decimal.Decimal, len(sparkline))
		if len(sparkline) > 0 {
			for i, value := range sparkline {
				x := &value
				sparklineArray[i] = utils.ConvertFloatToDecimal(x)
			}
		}
		nowDate := &pgtype.Date{}
		nowDate.Time = now
		row := []interface{}{
			marketData.Name,              //1
			uuidString,                   //2
			marketData.AlternateName,     //3
			&marketData.StartDate,        //4
			&marketData.EndDate,          //5
			*marketData.AssetID,          //6
			openUSD,                      //7
			closeUSD,                     //8
			highUSD,                      //9
			lowUSD,                       //10
			priceUSD,                     //11
			volumeUSD,                    //12
			marketCapUSD,                 //13
			marketData.Ticker,            //14
			marketData.Description,       //15
			*marketData.IntervalID,       //16
			*marketData.MarketDataTypeID, //17
			*marketData.SourceID,         //18
			totalSupply,                  //19
			maxSupply,                    //20
			circulatingSupply,            //21
			pq.Array(sparklineArray),     //22
			marketData.CreatedBy,         //23
			&now,                         //24
			marketData.CreatedBy,         //25
			&now,                         //26
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB
	// conn, err := database.DbConn.Conn(ctx)
	// if err != nil {
	// 	// handle error from acquiring connection from DB pool
	// }

	// err = conn.Raw(func(driverConn interface{}) error {
	// 	conn := driverConn.(*stdlib.Conn).Conn() // conn is a *pgx.Conn
	// 	// Do pgx specific stuff with conn
	// 	copyCount, err := conn.CopyFrom(context.Background(),
	// 		pgx.Identifier{"market_data"},
	// 		[]string{
	// 			"name",                //1
	// 			"uuid",                //2
	// 			"alternate_name",      //3
	// 			"start_date",          //4
	// 			"end_date",            //5
	// 			"asset_id",            //6
	// 			"open_usd",            //7
	// 			"close_usd",           //8
	// 			"high_usd",            //9
	// 			"low_usd",             //10
	// 			"price_usd",           //11
	// 			"volume_usd",          //12
	// 			"market_cap_usd",      //13
	// 			"ticker",              //14
	// 			"description",         //15
	// 			"interval_id",         //16
	// 			"market_data_type_id", //17
	// 			"source_id",           //18
	// 			"total_supply",        //19
	// 			"max_supply",          //20
	// 			"circulating_supply",  //21
	// 			"sparkline_7d",        //22
	// 			"created_by",          //23
	// 			"created_at",          //24
	// 			"updated_by",          //25
	// 			"updated_at",          //26
	// 		},
	// 		pgx.CopyFromRows(rows),
	// 	)
	// 	log.Println(fmt.Printf("copy count: %d", copyCount))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		// handle error that occurred while using *pgx.Conn
	// 	}
	// 	return nil
	// })
	// _, err := database.DbConnPgx.CopyFrom(
	// 	ctx,
	// 	pgx.Identifier{"market_data"},
	// 	[]string{
	// 		"name",                //1
	// 		"uuid",                //2
	// 		// "alternate_name",      //3
	// 		"start_date",          //4
	// 		"end_date",            //5
	// 		"asset_id",            //6
	// 		"open_usd",            //7
	// 		"close_usd",           //8
	// 		"high_usd",            //9
	// 		"low_usd",             //10
	// 		"price_usd",           //11
	// 		"volume_usd",          //12
	// 		"market_cap_usd",      //13
	// 		"ticker",              //14
	// 		"description",         //15
	// 		"interval_id",         //16
	// 		"market_data_type_id", //17
	// 		"source_id",           //18
	// 		"total_supply",        //19
	// 		"max_supply",          //20
	// 		"circulating_supply",  //21
	// 		// "sparkline_7d",        //22
	// 		"created_by", //23
	// 		"created_at", //24
	// 		"updated_by", //25
	// 		"updated_at", //26
	// 	},
	// 	pgx.CopyFromRows(rows),
	// )
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"market_data"},
		[]string{
			"name",                //1
			"uuid",                //2
			"alternate_name",      //3
			"start_date",          //4
			"end_date",            //5
			"asset_id",            //6
			"open_usd",            //7
			"close_usd",           //8
			"high_usd",            //9
			"low_usd",             //10
			"price_usd",           //11
			"volume_usd",          //12
			"market_cap_usd",      //13
			"ticker",              //14
			"description",         //15
			"interval_id",         //16
			"market_data_type_id", //17
			"source_id",           //18
			"total_supply",        //19
			"max_supply",          //20
			"circulating_supply",  //21
			"sparkline_7d",        //22
			"created_by",          //23
			"created_at",          //24
			"updated_by",          //25
			"updated_at",          //26
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

func InsertMarketDataList(marketDataList []MarketData) error {
	txn, err := database.DbConn.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	layoutPostgres := utils.LayoutPostgres
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn(
		"market_data",
		"name",                //1
		"uuid",                //2
		"alternate_name",      //3
		"start_date",          //4
		"end_date",            //5
		"asset_id",            //6
		"open_usd",            //7
		"close_usd",           //8
		"high_usd",            //9
		"low_usd",             //10
		"price_usd",           //11
		"volume_usd",          //12
		"market_cap_usd",      //13
		"ticker",              //14
		"description",         //15
		"interval_id",         //16
		"market_data_type_id", //17
		"source_id",           //18
		"total_supply",        //19
		"max_supply",          //20
		"circulating_supply",  //21
		"sparkline_7d",        //22
		"created_by",          //23
		"created_at",          //24
		"updated_by",          //25
		"updated_at",          //26
	))
	if err != nil {
		log.Fatal(err)
	}
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	for _, marketData := range marketDataList {
		_, err = stmt.Exec(
			marketData.Name,                             //1
			marketData.UUID,                             //2
			marketData.AlternateName,                    //3
			marketData.StartDate.Format(layoutPostgres), //4
			marketData.EndDate.Format(layoutPostgres),   //5
			marketData.AssetID,                          //6
			marketData.OpenUSD,                          //7
			marketData.CloseUSD,                         //8
			marketData.HighUSD,                          //9
			marketData.LowUSD,                           //10
			marketData.PriceUSD,                         //11
			marketData.VolumeUSD,                        //12
			marketData.MarketCapUSD,                     //13
			marketData.Ticker,                           //14
			marketData.Description,                      //15
			marketData.IntervalID,                       //16
			marketData.MarketDataTypeID,                 //17
			marketData.SourceID,                         //18
			marketData.TotalSupply,                      //19
			marketData.MaxSupply,                        //20
			marketData.CirculatingSupply,                //21
			pq.Array(marketData.Sparkline7d),            //22
			marketData.CreatedBy,                        //23
			now.Format(layoutPostgres),                  //24
			marketData.CreatedBy,                        //25
			now.Format(layoutPostgres),                  //26
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func InsertMarketDataQuote(marketDataQuote MarketDataQuote) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	layoutPostgres := utils.LayoutPostgres
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO market_data_quotes 
	(
		"market_data_id",         //1
		"base_asset_id",          //2
		"quote_asset_id",         //3
		"name",                   //4
		"uuid",                   //5
		"alternate_name",         //6
		"open",                   //7
		"close",                  //8
		"high_24h",               //9
		"low_24h",                //10
		"price",                  //11
		"volume",                 //12
		"market_cap",             //13
		"ticker",                 //14
		"description",            //15
		"source_id",              //16
		"fully_diluted_valution", //17
		"ath",                    //18
		"ath_date",               //19
		"atl",                    //20
		"atl_date",               //21
		"price_change_1h",        //22
		"price_change_24h",       //23
		"price_change_7d",        //24
		"price_change_30d",       //25
		"price_change_60d",       //26
		"price_change_200d",      //27
		"price_change_1y",        //28
		"created_by",             //29
		"created_at",             //30
		"updated_by",             //31
		"updated_at",             //32
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
			$13,
			$14,
			$15,
			$16,
			$17,
			$18,
			$19,
			$20,
			$21,
			$22,
			$23,
			$24,
			$25,
			$26,
			$27,
			$28,
			$29,
			$30,
			$31,
			$32,
		)
		RETURNING id`,
		marketDataQuote.MarketDataID,         //1
		marketDataQuote.BaseAssetID,          //2
		marketDataQuote.QuoteAssetID,         //3
		marketDataQuote.Name,                 //4
		marketDataQuote.UUID,                 //5
		marketDataQuote.AlternateName,        //6
		marketDataQuote.Open,                 //7
		marketDataQuote.Close,                //8
		marketDataQuote.High24h,              //9
		marketDataQuote.Low24h,               //10
		marketDataQuote.Price,                //11
		marketDataQuote.Volume,               //12
		marketDataQuote.MarketCap,            //13
		marketDataQuote.Ticker,               //14
		marketDataQuote.Description,          //15
		marketDataQuote.SourceID,             //16
		marketDataQuote.FullyDilutedValution, //17
		marketDataQuote.Ath,                  //18
		marketDataQuote.AthDate,              //19
		marketDataQuote.Atl,                  //20
		marketDataQuote.AtlDate,              //21
		marketDataQuote.PriceChange1h,        //22
		marketDataQuote.PriceChange24h,       //23
		marketDataQuote.PriceChange7d,        //24
		marketDataQuote.PriceChange30d,       //25
		marketDataQuote.PriceChange60d,       //26
		marketDataQuote.PriceChange200d,      //27
		marketDataQuote.PriceChange1y,        //28
		marketDataQuote.CreatedBy,            //29
		now.Format(layoutPostgres),           //30
		marketDataQuote.CreatedBy,            //31
		now.Format(layoutPostgres),           //32
	).Scan(&insertID)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func InsertMarketDataQuoteListManual(marketDataQuoteList []MarketDataQuote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range marketDataQuoteList {
		marketDataQuote := marketDataQuoteList[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(marketDataQuote.UUID)
		row := []interface{}{
			*marketDataQuote.MarketDataID,                                     //1
			*marketDataQuote.BaseAssetID,                                      //2
			*marketDataQuote.QuoteAssetID,                                     //3
			marketDataQuote.Name,                                              //4
			uuidString,                                                        //5
			marketDataQuote.AlternateName,                                     //6
			utils.ConvertFloatToDecimal(marketDataQuote.Open),                 //7
			utils.ConvertFloatToDecimal(marketDataQuote.Close),                //8
			utils.ConvertFloatToDecimal(marketDataQuote.High24h),              //9
			utils.ConvertFloatToDecimal(marketDataQuote.Low24h),               //10
			utils.ConvertFloatToDecimal(marketDataQuote.Price),                //11
			utils.ConvertFloatToDecimal(marketDataQuote.Volume),               //12
			utils.ConvertFloatToDecimal(marketDataQuote.MarketCap),            //13
			marketDataQuote.Ticker,                                            //14
			marketDataQuote.Description,                                       //15
			*marketDataQuote.SourceID,                                         //16
			utils.ConvertFloatToDecimal(marketDataQuote.FullyDilutedValution), //17
			utils.ConvertFloatToDecimal(marketDataQuote.Ath),                  //18
			&marketDataQuote.AthDate,                                          //19
			utils.ConvertFloatToDecimal(marketDataQuote.Atl),                  //20
			&marketDataQuote.AtlDate,                                          //21
			utils.ConvertFloatToDecimal(marketDataQuote.PriceChange1h),        //22
			utils.ConvertFloatToDecimal(marketDataQuote.PriceChange24h),       //23
			utils.ConvertFloatToDecimal(marketDataQuote.PriceChange7d),        //24
			utils.ConvertFloatToDecimal(marketDataQuote.PriceChange30d),       //25
			utils.ConvertFloatToDecimal(marketDataQuote.PriceChange60d),       //26
			utils.ConvertFloatToDecimal(marketDataQuote.PriceChange200d),      //27
			utils.ConvertFloatToDecimal(marketDataQuote.PriceChange1y),        //28
			marketDataQuote.CreatedBy,                                         //29
			&now,                                                              //30
			marketDataQuote.CreatedBy,                                         //31
			&now,                                                              //32
		}
		rows = append(rows, row)
	}
	// Given db is a *sql.DB
	// conn, err := database.DbConn.Conn(ctx)
	// if err != nil {
	// 	// handle error from acquiring connection from DB pool
	// }

	// err = conn.Raw(func(driverConn interface{}) error {
	// 	conn := driverConn.(*stdlib.Conn).Conn() // conn is a *pgx.Conn
	// 	// Do pgx specific stuff with conn
	// 	copyCount, err := conn.CopyFrom(context.Background(),
	// 		pgx.Identifier{"market_data_quotes"},
	// 		[]string{
	// 			"market_data_id",
	// 			"base_asset_id",
	// 			"quote_asset_id",
	// 			"name",
	// 			"uuid",
	// 			"alternate_name",
	// 			"open",
	// 			"close",
	// 			"high_24h",
	// 			"low_24h",
	// 			"price",
	// 			"volume",
	// 			"market_cap",
	// 			"ticker",
	// 			"description",
	// 			"source_id",
	// 			"fully_diluted_valution",
	// 			"ath",
	// 			"ath_date",
	// 			"atl",
	// 			"atl_date",
	// 			"price_change_1h",
	// 			"price_change_24h",
	// 			"price_change_7d",
	// 			"price_change_30d",
	// 			"price_change_60d",
	// 			"price_change_200d",
	// 			"price_change_1y",
	// 			"created_by",
	// 			"created_at",
	// 			"updated_by",
	// 			"updated_at",
	// 		},
	// 		pgx.CopyFromRows(rows),
	// 	)
	// 	log.Println(fmt.Printf("copy count: %d", copyCount))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		// handle error that occurred while using *pgx.Conn
	// 	}
	// 	return nil
	// })

	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"market_data_quotes"},
		[]string{
			"market_data_id",
			"base_asset_id",
			"quote_asset_id",
			"name",
			"uuid",
			"alternate_name",
			"open",
			"close",
			"high_24h",
			"low_24h",
			"price",
			"volume",
			"market_cap",
			"ticker",
			"description",
			"source_id",
			"fully_diluted_valution",
			"ath",
			"ath_date",
			"atl",
			"atl_date",
			"price_change_1h",
			"price_change_24h",
			"price_change_7d",
			"price_change_30d",
			"price_change_60d",
			"price_change_200d",
			"price_change_1y",
			"created_by",
			"created_at",
			"updated_by",
			"updated_at",
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}

	// if err != nil {
	// 	// handle error that occurred while using *pgx.Conn
	// }
	return nil
}
func InsertMarketDataQuoteList(marketDataQuoteList []MarketDataQuote) error {
	txn, err := database.DbConn.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	layoutPostgres := utils.LayoutPostgres
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn(
		"market_data_quotes",
		"market_data_id",         //1
		"base_asset_id",          //2
		"quote_asset_id",         //3
		"name",                   //4
		"uuid",                   //5
		"alternate_name",         //6
		"open",                   //7
		"close",                  //8
		"high_24h",               //9
		"low_24h",                //10
		"price",                  //11
		"volume",                 //12
		"market_cap",             //13
		"ticker",                 //14
		"description",            //15
		"source_id",              //16
		"fully_diluted_valution", //17
		"ath",                    //18
		"ath_date",               //19
		"atl",                    //20
		"atl_date",               //21
		"price_change_1h",        //22
		"price_change_24h",       //23
		"price_change_7d",        //24
		"price_change_30d",       //25
		"price_change_60d",       //26
		"price_change_200d",      //27
		"price_change_1y",        //28
		"created_by",             //29
		"created_at",             //30
		"updated_by",             //31
		"updated_at",             //32
	))
	if err != nil {
		log.Fatal(err)
	}
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	for _, marketDataQuote := range marketDataQuoteList {
		_, err = stmt.Exec(
			marketDataQuote.MarketDataID,         //1
			marketDataQuote.BaseAssetID,          //2
			marketDataQuote.QuoteAssetID,         //3
			marketDataQuote.Name,                 //4
			marketDataQuote.UUID,                 //5
			marketDataQuote.AlternateName,        //6
			marketDataQuote.Open,                 //7
			marketDataQuote.Close,                //8
			marketDataQuote.High24h,              //9
			marketDataQuote.Low24h,               //10
			marketDataQuote.Price,                //11
			marketDataQuote.Volume,               //12
			marketDataQuote.MarketCap,            //13
			marketDataQuote.Ticker,               //14
			marketDataQuote.Description,          //15
			marketDataQuote.SourceID,             //16
			marketDataQuote.FullyDilutedValution, //17
			marketDataQuote.Ath,                  //18
			marketDataQuote.AthDate,              //19
			marketDataQuote.Atl,                  //20
			marketDataQuote.AtlDate,              //21
			marketDataQuote.PriceChange1h,        //22
			marketDataQuote.PriceChange24h,       //23
			marketDataQuote.PriceChange7d,        //24
			marketDataQuote.PriceChange30d,       //25
			marketDataQuote.PriceChange60d,       //26
			marketDataQuote.PriceChange200d,      //27
			marketDataQuote.PriceChange1y,        //28
			marketDataQuote.CreatedBy,            //29
			now.Format(layoutPostgres),           //30
			marketDataQuote.CreatedBy,            //31
			now.Format(layoutPostgres),           //32
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
