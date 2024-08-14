package marketdata

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	marketdataquote "github.com/kfukue/lyle-labs-libraries/v2/marketDataQuote"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
	decimal "github.com/shopspring/decimal"
)

func GetMarketData(dbConnPgx utils.PgxIface, marketDataID *int) (*MarketData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
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
	WHERE id = $1`, *marketDataID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	defer row.Close()
	marketData, err := pgx.CollectOneRow(row, pgx.RowToStructByName[MarketData])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &marketData, nil
}

func RemoveMarketData(dbConnPgx utils.PgxIface, marketDataID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethMarketData DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM market_data WHERE id = $1`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *marketDataID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func RemoveMarketDataFromBaseAssetBetweenDates(dbConnPgx utils.PgxIface, assetID *int, startDate, endDate time.Time) error {
	log.Printf("start : %s end : %s", startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres))
	err := marketdataquote.RemoveMarketDataQuoteFromBaseAssetBetweenDates(dbConnPgx, assetID, startDate, endDate)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveMarketDataFromBaseAssetBetweenDates DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE
		FROM market_data 
		WHERE 
			asset_id = $1
			AND start_date BETWEEN $2 and $3`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *assetID, startDate, endDate); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetMarketDataList(dbConnPgx utils.PgxIface, ids []int) ([]MarketData, error) {
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
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	marketDataList, err := pgx.CollectRows(results, pgx.RowToStructByName[MarketData])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return marketDataList, nil
}

func GetMarketDataListByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]MarketData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
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
	marketDataList, err := pgx.CollectRows(results, pgx.RowToStructByName[MarketData])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return marketDataList, nil
}

func GetStartAndEndDateDiffMarketDataList(dbConnPgx utils.PgxIface, diffInDate *int) ([]MarketData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
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
	`, *diffInDate)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	marketDataList, err := pgx.CollectRows(results, pgx.RowToStructByName[MarketData])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return marketDataList, nil
}

func UpdateMarketData(dbConnPgx utils.PgxIface, marketData *MarketData) error {
	// if the marketData id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if marketData.ID == nil || *marketData.ID == 0 {
		return errors.New("marketData has invalid ID")
	}
	layoutPostgres := utils.LayoutPostgres
	startDate := marketData.StartDate
	endDate := marketData.EndDate
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateMarketData DbConn.Begin   %s", err.Error())
		return err
	}
	log.Println(fmt.Sprintf("Updating start: %s, end : %s", startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)))
	sql := `UPDATE market_data SET 
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
		WHERE id=$23`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
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
		marketData.ID);                              //23
	err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertMarketData(dbConnPgx utils.PgxIface, marketData *MarketData) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertMarketData DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	layoutPostgres := utils.LayoutPostgres
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO market_data 
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
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, err
	}
	return int(insertID), nil
}
func InsertMarketDataList(dbConnPgx utils.PgxIface, marketDataList []MarketData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range marketDataList {
		marketData := marketDataList[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(marketData.UUID)

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
	copyCount, err := dbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("InsertMarketDataListManual: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetMarketDataListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]MarketData, error) {
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
	FROM market_data
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
	marketDataList, err := pgx.CollectRows(results, pgx.RowToStructByName[MarketData])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return marketDataList, nil
}

func GetTotalMarketDataCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM market_data`)
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
