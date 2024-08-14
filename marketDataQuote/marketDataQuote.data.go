package marketdataquote

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
)

func GetLatestLiveMarketData(dbConnPgx utils.PgxIface) ([]MarketDataQuoteResults, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
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
		mq.uuid,
		mq.name,
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

func RemoveMarketDataQuoteFromBaseAssetBetweenDates(dbConnPgx utils.PgxIface, assetID *int, startDate, endDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethMarketData DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE
		FROM market_data_quotes
		WHERE market_data_id IN (
			SELECT ID
			FROM market_data
			WHERE asset_id = $1
			AND start_date BETWEEN $2 and $3
		);`

	if _, err := dbConnPgx.Exec(ctx, sql, *assetID, startDate, endDate); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
func GetAllMarketDataFromStrategyID(dbConnPgx utils.PgxIface, strategyID *int) ([]MarketDataQuoteResults, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
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
		mq.uuid,
		mq.name,
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
	ORDER BY m.start_date desc , baseA.ticker`, *strategyID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

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

func InsertMarketDataQuote(dbConnPgx utils.PgxIface, marketDataQuote *MarketDataQuote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertMarketDataQuote DbConn.Begin   %s", err.Error())
		return err
	}
	_, err = dbConnPgx.Query(ctx, `INSERT INTO market_data_quotes 
	(
		"market_data_id",        
		"base_asset_id",         
		"quote_asset_id",        
		"uuid",                  
		"name",                  
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
			current_timestamp at time zone 'UTC',
			$30,
			current_timestamp at time zone 'UTC',
		)`,
		marketDataQuote.MarketDataID,         //1
		marketDataQuote.BaseAssetID,          //2
		marketDataQuote.QuoteAssetID,         //3
		marketDataQuote.UUID,                 //4
		marketDataQuote.Name,                 //5
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
		marketDataQuote.CreatedBy,            //30
	)

	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return err
	}
	return nil
}
func InsertMarketDataQuoteList(dbConnPgx utils.PgxIface, marketDataQuoteList []MarketDataQuote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range marketDataQuoteList {
		marketDataQuote := marketDataQuoteList[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(marketDataQuote.UUID)
		row := []interface{}{
			*marketDataQuote.MarketDataID,                                     //1
			*marketDataQuote.BaseAssetID,                                      //2
			*marketDataQuote.QuoteAssetID,                                     //3
			uuidString,                                                        //4
			marketDataQuote.Name,                                              //5
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
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"market_data_quotes"},
		[]string{
			"market_data_id",         //1
			"base_asset_id",          //2
			"quote_asset_id",         //3
			"uuid",                   //4
			"name",                   //5
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
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertMarketDataQuoteList: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetMarketDataQuoteListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]MarketDataQuote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		market_data_id,
		base_asset_id,
		quote_asset_id ,
  		uuid,
		name,
		alternate_name,
		open,
		close,
		high_24h,
		low_24h,
		price,
		volume,
		market_cap,
		ticker,
		description,
		source_id,
		fully_diluted_valution,
		ath,
		ath_date date,
		atl,
		atl_date date,
		rename to high_24h,
		rename to low_24h,
		price_change_1h,
		price_change_24h,
		price_change_7d,
		price_change_30d,
		price_change_60d,
		price_change_200d,
		price_change_1y,
		created_by,
		created_at,
		updated_by,
		updated_at,
	FROM market_data_quotes
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

	marketDataQuoteList, err := pgx.CollectRows(results, pgx.RowToStructByName[MarketDataQuote])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return marketDataQuoteList, nil
}

func GetTotalMarketDataQuoteCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM market_data_quotes`)
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
