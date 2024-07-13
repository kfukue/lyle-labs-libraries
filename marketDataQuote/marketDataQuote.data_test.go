package marketdataquote

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
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
}
var DBColumnsInsertMarketDataQuoteList = []string{
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
}
var fakeSparklineData = []float64{
	float64(1.0),
	float64(2.0),
	float64(3.0),
	float64(4.0),
	float64(5.0),
	float64(6.0),
	float64(7.0),
}

var TestData1 = MarketDataQuote{
	MarketDataID:         utils.Ptr[int](1),                                         //1
	BaseAssetID:          utils.Ptr[int](1),                                         //2
	QuoteAssetID:         utils.Ptr[int](1),                                         //3
	UUID:                 "01ef85e8-2c26-441e-8c7f-71d79518ad72",                    //4
	Name:                 "Basic Attention Token/United States Dollar",              //5
	AlternateName:        "basic-attention-token/usd - 2021-01-06",                  //6
	Open:                 utils.Ptr[float64](float64(0.23220798385492783)),          //7
	Close:                utils.Ptr[float64](float64(0.33)),                         //8
	High24h:              utils.Ptr[float64](float64(0.66)),                         //9
	Low24h:               utils.Ptr[float64](float64(0.11)),                         //10
	Price:                utils.Ptr[float64](float64(0.11)),                         //11
	Volume:               utils.Ptr[float64](float64(171848480.40188548)),           //12
	MarketCap:            utils.Ptr[float64](float64(342801219.6935223)),            //13
	Ticker:               "usd",                                                     //14
	Description:          "Market Data UUID : 9bc45fd6-c60d-4252-a3cc-06b807b7cd19", //15
	SourceID:             utils.Ptr[int](3),                                         //16
	FullyDilutedValution: utils.Ptr[float64](float64(342801219)),                    //17
	Ath:                  utils.Ptr[float64](float64(0.8888)),                       //18
	AthDate:              utils.SampleCreatedAtTime,                                 //19
	Atl:                  utils.Ptr[float64](float64(0.111)),                        //20
	AtlDate:              utils.SampleCreatedAtTime,                                 //21
	PriceChange1h:        utils.Ptr[float64](float64(0.55)),                         //22
	PriceChange24h:       utils.Ptr[float64](float64(0.111)),                        //23
	PriceChange7d:        utils.Ptr[float64](float64(0.3)),                          //24
	PriceChange30d:       utils.Ptr[float64](float64(0.222)),                        //25
	PriceChange60d:       utils.Ptr[float64](float64(0.44)),                         //26
	PriceChange200d:      utils.Ptr[float64](float64(0.55)),                         //27
	PriceChange1y:        utils.Ptr[float64](float64(0.44)),                         //28
	CreatedBy:            "SYSTEM",                                                  //29
	CreatedAt:            utils.SampleCreatedAtTime,                                 //30
	UpdatedBy:            "SYSTEM",                                                  //31
	UpdatedAt:            utils.SampleCreatedAtTime,                                 //32
}

var TestData2 = MarketDataQuote{
	MarketDataID:         utils.Ptr[int](2),                                         //1
	BaseAssetID:          utils.Ptr[int](2),                                         //2
	QuoteAssetID:         utils.Ptr[int](1),                                         //3
	UUID:                 "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                    //4
	Name:                 "Augur/United States Dollar",                              //5
	AlternateName:        "augur/usd - 2021-01-07",                                  //6
	Open:                 utils.Ptr[float64](float64(18.932096109991285)),           //7
	Close:                utils.Ptr[float64](float64(18.33)),                        //8
	High24h:              utils.Ptr[float64](float64(18.66)),                        //9
	Low24h:               utils.Ptr[float64](float64(18.11)),                        //10
	Price:                utils.Ptr[float64](float64(18.11)),                        //11
	Volume:               utils.Ptr[float64](float64(30389278.116494928)),           //12
	MarketCap:            utils.Ptr[float64](float64(106742851.67464235)),           //13
	Ticker:               "usd",                                                     //14
	Description:          "Market Data UUID : b03077f6-29c8-4d52-9f6f-ea1f70be00a4", //15
	SourceID:             utils.Ptr[int](3),                                         //16
	FullyDilutedValution: utils.Ptr[float64](float64(30389278.116494928)),           //17
	Ath:                  utils.Ptr[float64](float64(21)),                           //18
	AthDate:              utils.SampleCreatedAtTime,                                 //19
	Atl:                  utils.Ptr[float64](float64(1.44)),                         //20
	AtlDate:              utils.SampleCreatedAtTime,                                 //21
	PriceChange1h:        utils.Ptr[float64](float64(0.55)),                         //22
	PriceChange24h:       utils.Ptr[float64](float64(0.111)),                        //23
	PriceChange7d:        utils.Ptr[float64](float64(1.3)),                          //24
	PriceChange30d:       utils.Ptr[float64](float64(2.222)),                        //25
	PriceChange60d:       utils.Ptr[float64](float64(3.44)),                         //26
	PriceChange200d:      utils.Ptr[float64](float64(5.55)),                         //27
	PriceChange1y:        utils.Ptr[float64](float64(10.44)),                        //28
	CreatedBy:            "SYSTEM",                                                  //29
	CreatedAt:            utils.SampleCreatedAtTime,                                 //30
	UpdatedBy:            "SYSTEM",                                                  //31
	UpdatedAt:            utils.SampleCreatedAtTime,                                 //32
}
var TestAllData = []MarketDataQuote{TestData1, TestData2}

func AddMarketDataQuoteToMockRows(mock pgxmock.PgxPoolIface, dataList []MarketDataQuote) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.MarketDataID,         //1
			data.BaseAssetID,          //2
			data.QuoteAssetID,         //3
			data.UUID,                 //4
			data.Name,                 //5
			data.AlternateName,        //6
			data.Open,                 //7
			data.Close,                //8
			data.High24h,              //9
			data.Low24h,               //10
			data.Price,                //11
			data.Volume,               //12
			data.MarketCap,            //13
			data.Ticker,               //14
			data.Description,          //15
			data.SourceID,             //16
			data.FullyDilutedValution, //17
			data.Ath,                  //18
			data.AthDate,              //19
			data.Atl,                  //20
			data.AtlDate,              //21
			data.PriceChange1h,        //22
			data.PriceChange24h,       //23
			data.PriceChange7d,        //24
			data.PriceChange30d,       //25
			data.PriceChange60d,       //26
			data.PriceChange200d,      //27
			data.PriceChange1y,        //28
			data.CreatedBy,            //29
			data.CreatedAt,            //30
			data.UpdatedBy,            //31
			data.UpdatedAt,            //32
		)
	}
	return rows
}

var DBColumnsMarketDataQuoteResults = []string{
	"start_date",             //33
	"end_date",               //34
	"base_asset_name",        //35
	"base_asset_ticker",      //36
	"quote_asset_name",       //37
	"quote_asset_ticker",     //38
	"market_data_quotes_id",  //1
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
	"updated_at",             //32 end of MarketDataQuote

}

func AddMarketDataQuoteResultsToMockRows(mock pgxmock.PgxPoolIface, dataList []MarketDataQuoteResults) *pgxmock.Rows {
	rows := mock.NewRows(DBColumnsMarketDataQuoteResults)
	for _, data := range dataList {
		rows.AddRow(
			data.StartDate,                            //33
			data.EndDate,                              //34
			data.BaseAssetName,                        //35
			data.BaseAssetTicker,                      //36
			data.QuoteAssetName,                       //37
			data.QuoteAssetTicker,                     //38
			data.MarketDataQuote.MarketDataID,         //1
			data.MarketDataQuote.BaseAssetID,          //2
			data.MarketDataQuote.QuoteAssetID,         //3
			data.MarketDataQuote.UUID,                 //4
			data.MarketDataQuote.Name,                 //5
			data.MarketDataQuote.AlternateName,        //6
			data.MarketDataQuote.Open,                 //7
			data.MarketDataQuote.Close,                //8
			data.MarketDataQuote.High24h,              //9
			data.MarketDataQuote.Low24h,               //10
			data.MarketDataQuote.Price,                //11
			data.MarketDataQuote.Volume,               //12
			data.MarketDataQuote.MarketCap,            //13
			data.MarketDataQuote.Ticker,               //14
			data.MarketDataQuote.Description,          //15
			data.MarketDataQuote.SourceID,             //16
			data.MarketDataQuote.FullyDilutedValution, //17
			data.MarketDataQuote.Ath,                  //18
			data.MarketDataQuote.AthDate,              //19
			data.MarketDataQuote.Atl,                  //20
			data.MarketDataQuote.AtlDate,              //21
			data.MarketDataQuote.PriceChange1h,        //22
			data.MarketDataQuote.PriceChange24h,       //23
			data.MarketDataQuote.PriceChange7d,        //24
			data.MarketDataQuote.PriceChange30d,       //25
			data.MarketDataQuote.PriceChange60d,       //26
			data.MarketDataQuote.PriceChange200d,      //27
			data.MarketDataQuote.PriceChange1y,        //28
			data.MarketDataQuote.CreatedBy,            //29
			data.MarketDataQuote.CreatedAt,            //30
			data.MarketDataQuote.UpdatedBy,            //31
			data.MarketDataQuote.UpdatedAt,            //32

		)
	}
	return rows
}

func TestGetLatestLiveMarketData(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testData1MarketDataQuoteResults := MarketDataQuoteResults{}
	testData1MarketDataQuoteResults.MarketDataQuote = TestData1
	testData1MarketDataQuoteResults.StartDate = utils.SampleCreatedAtTime
	testData1MarketDataQuoteResults.EndDate = utils.SampleCreatedAtTime
	testData1MarketDataQuoteResults.BaseAssetName = "Basic Attention Token"
	testData1MarketDataQuoteResults.BaseAssetTicker = "BAT"
	testData1MarketDataQuoteResults.QuoteAssetName = "United States Dollar"
	testData1MarketDataQuoteResults.QuoteAssetTicker = "usd"
	testData1MarketDataQuoteResults2 := MarketDataQuoteResults{}
	testData1MarketDataQuoteResults2.MarketDataQuote = TestData2
	testData1MarketDataQuoteResults2.StartDate = utils.SampleCreatedAtTime
	testData1MarketDataQuoteResults2.EndDate = utils.SampleCreatedAtTime
	testData1MarketDataQuoteResults2.BaseAssetName = "Augur"
	testData1MarketDataQuoteResults2.BaseAssetTicker = "REP"
	testData1MarketDataQuoteResults2.QuoteAssetName = "United States Dollar"
	testData1MarketDataQuoteResults2.QuoteAssetTicker = "usd"
	dataList := []MarketDataQuoteResults{testData1MarketDataQuoteResults, testData1MarketDataQuoteResults2}
	mockRows := AddMarketDataQuoteResultsToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM market_data_quotes").WillReturnRows(mockRows)
	foundMarketDataQuoteResultsList, err := GetLatestLiveMarketData(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetLatestLiveMarketData", err)
	}
	for i, sourceMarketDataQuoteResults := range dataList {
		if cmp.Equal(sourceMarketDataQuoteResults, foundMarketDataQuoteResultsList[i]) == false {
			t.Errorf("Expected MarketDataQuoteResults From Method GetLatestLiveMarketData: %v is different from actual %v", sourceMarketDataQuoteResults, foundMarketDataQuoteResultsList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLatestLiveMarketDataForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM market_data_quotes").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataQuoteResultsList, err := GetLatestLiveMarketData(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetLatestLiveMarketData", err)
	}
	if len(foundMarketDataQuoteResultsList) != 0 {
		t.Errorf("Expected From Method GetLatestLiveMarketData: to be empty but got this: %v", foundMarketDataQuoteResultsList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketDataQuoteFromBaseAssetBetweenDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	assetID := targetData.BaseAssetID
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	// from marketdataquote.RemoveMarketDataQuoteQuoteFromBaseAssetBetweenDates
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data_quotes").WithArgs(*assetID, startDate, endDate).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveMarketDataQuoteFromBaseAssetBetweenDates(mock, assetID, startDate, endDate)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveMarketDataQuoteFromBaseAssetBetweenDates", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketDataQuoteFromBaseAssetBetweenDatesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 1
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	// from marketdataquote.RemoveMarketDataQuoteQuoteFromBaseAssetBetweenDates
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data_quotes").WithArgs(assetID, startDate, endDate).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveMarketDataQuoteFromBaseAssetBetweenDates(mock, &assetID, startDate, endDate)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllMarketDataFromStrategyID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testData1MarketDataQuoteResults := MarketDataQuoteResults{}
	testData1MarketDataQuoteResults.MarketDataQuote = TestData1
	testData1MarketDataQuoteResults.StartDate = utils.SampleCreatedAtTime
	testData1MarketDataQuoteResults.EndDate = utils.SampleCreatedAtTime
	testData1MarketDataQuoteResults.BaseAssetName = "Basic Attention Token"
	testData1MarketDataQuoteResults.BaseAssetTicker = "BAT"
	testData1MarketDataQuoteResults.QuoteAssetName = "United States Dollar"
	testData1MarketDataQuoteResults.QuoteAssetTicker = "usd"
	testData1MarketDataQuoteResults2 := MarketDataQuoteResults{}
	testData1MarketDataQuoteResults2.MarketDataQuote = TestData2
	testData1MarketDataQuoteResults2.StartDate = utils.SampleCreatedAtTime
	testData1MarketDataQuoteResults2.EndDate = utils.SampleCreatedAtTime
	testData1MarketDataQuoteResults2.BaseAssetName = "Augur"
	testData1MarketDataQuoteResults2.BaseAssetTicker = "REP"
	testData1MarketDataQuoteResults2.QuoteAssetName = "United States Dollar"
	testData1MarketDataQuoteResults2.QuoteAssetTicker = "usd"
	dataList := []MarketDataQuoteResults{testData1MarketDataQuoteResults, testData1MarketDataQuoteResults2}
	mockRows := AddMarketDataQuoteResultsToMockRows(mock, dataList)
	strategyID := 1
	mock.ExpectQuery("^SELECT (.+) FROM market_data_quotes").WithArgs(strategyID).WillReturnRows(mockRows)
	foundMarketDataQuoteResultsList, err := GetAllMarketDataFromStrategyID(mock, &strategyID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllMarketDataFromStrategyID", err)
	}
	for i, sourceMarketDataQuoteResults := range dataList {
		if cmp.Equal(sourceMarketDataQuoteResults, foundMarketDataQuoteResultsList[i]) == false {
			t.Errorf("Expected MarketDataQuoteResults From Method GetAllMarketDataFromStrategyID: %v is different from actual %v", sourceMarketDataQuoteResults, foundMarketDataQuoteResultsList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllMarketDataFromStrategyIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := 10
	mock.ExpectQuery("^SELECT (.+) FROM market_data_quotes").WithArgs(strategyID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataQuoteList, err := GetAllMarketDataFromStrategyID(mock, &strategyID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllMarketDataFromStrategyID", err)
	}
	if len(foundMarketDataQuoteList) != 0 {
		t.Errorf("Expected From Method GetAllMarketDataFromStrategyID: to be empty but got this: %v", foundMarketDataQuoteList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataQuote(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO market_data_quotes").WithArgs(
		targetData.MarketDataID,         //1
		targetData.BaseAssetID,          //2
		targetData.QuoteAssetID,         //3
		targetData.UUID,                 //4
		targetData.Name,                 //5
		targetData.AlternateName,        //6
		targetData.Open,                 //7
		targetData.Close,                //8
		targetData.High24h,              //9
		targetData.Low24h,               //10
		targetData.Price,                //11
		targetData.Volume,               //12
		targetData.MarketCap,            //13
		targetData.Ticker,               //14
		targetData.Description,          //15
		targetData.SourceID,             //16
		targetData.FullyDilutedValution, //17
		targetData.Ath,                  //18
		targetData.AthDate,              //19
		targetData.Atl,                  //20
		targetData.AtlDate,              //21
		targetData.PriceChange1h,        //22
		targetData.PriceChange24h,       //23
		targetData.PriceChange7d,        //24
		targetData.PriceChange30d,       //25
		targetData.PriceChange60d,       //26
		targetData.PriceChange200d,      //27
		targetData.PriceChange1y,        //28
		targetData.CreatedBy,            //29
		targetData.CreatedBy,            //30
	).WillReturnRows(pgxmock.NewRows([]string{""}).AddRow(1))
	mock.ExpectCommit()
	err = InsertMarketDataQuote(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertMarketDataQuote", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataQuoteOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO market_data_quotes").WithArgs(
		targetData.MarketDataID,         //1
		targetData.BaseAssetID,          //2
		targetData.QuoteAssetID,         //3
		targetData.UUID,                 //4
		targetData.Name,                 //5
		targetData.AlternateName,        //6
		targetData.Open,                 //7
		targetData.Close,                //8
		targetData.High24h,              //9
		targetData.Low24h,               //10
		targetData.Price,                //11
		targetData.Volume,               //12
		targetData.MarketCap,            //13
		targetData.Ticker,               //14
		targetData.Description,          //15
		targetData.SourceID,             //16
		targetData.FullyDilutedValution, //17
		targetData.Ath,                  //18
		targetData.AthDate,              //19
		targetData.Atl,                  //20
		targetData.AtlDate,              //21
		targetData.PriceChange1h,        //22
		targetData.PriceChange24h,       //23
		targetData.PriceChange7d,        //24
		targetData.PriceChange30d,       //25
		targetData.PriceChange60d,       //26
		targetData.PriceChange200d,      //27
		targetData.PriceChange1y,        //28
		targetData.CreatedBy,            //29
		targetData.CreatedBy,            //30
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	err = InsertMarketDataQuote(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataQuoteOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO market_data_quotes").WithArgs(
		targetData.MarketDataID,         //1
		targetData.BaseAssetID,          //2
		targetData.QuoteAssetID,         //3
		targetData.UUID,                 //4
		targetData.Name,                 //5
		targetData.AlternateName,        //6
		targetData.Open,                 //7
		targetData.Close,                //8
		targetData.High24h,              //9
		targetData.Low24h,               //10
		targetData.Price,                //11
		targetData.Volume,               //12
		targetData.MarketCap,            //13
		targetData.Ticker,               //14
		targetData.Description,          //15
		targetData.SourceID,             //16
		targetData.FullyDilutedValution, //17
		targetData.Ath,                  //18
		targetData.AthDate,              //19
		targetData.Atl,                  //20
		targetData.AtlDate,              //21
		targetData.PriceChange1h,        //22
		targetData.PriceChange24h,       //23
		targetData.PriceChange7d,        //24
		targetData.PriceChange30d,       //25
		targetData.PriceChange60d,       //26
		targetData.PriceChange200d,      //27
		targetData.PriceChange1y,        //28
		targetData.CreatedBy,            //29
		targetData.CreatedBy,            //30
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	err = InsertMarketDataQuote(mock, &targetData)

	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataQuoteList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"market_data_quotes"}, DBColumnsInsertMarketDataQuoteList)
	targetData := TestAllData
	err = InsertMarketDataQuoteList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertMarketDataQuoteList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataQuoteListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"market_data_quotes"}, DBColumnsInsertMarketDataQuoteList).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertMarketDataQuoteList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataQuoteListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddMarketDataQuoteToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"quote_asset_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM market_data_quotes").WillReturnRows(mockRows)
	foundMarketDataQuoteList, err := GetMarketDataQuoteListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketDataQuoteListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundMarketDataQuoteList[i]) == false {
			t.Errorf("Expected sourceData From Method GetMarketDataQuoteListByPagination: %v is different from actual %v", sourceData, foundMarketDataQuoteList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataQuoteListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"quote_asset_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM market_data_quotes").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataQuoteList, err := GetMarketDataQuoteListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataQuoteListByPagination", err)
	}
	if len(foundMarketDataQuoteList) != 0 {
		t.Errorf("Expected From Method GetMarketDataQuoteListByPagination: to be empty but got this: %v", foundMarketDataQuoteList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMarketDataQuoteCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM market_data_quotes").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalMarketDataQuoteCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalMarketDataQuoteCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalMarketDataQuoteCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMarketDataQuoteCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM market_data_quotes").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalMarketDataQuoteCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalMarketDataQuoteCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalMarketDataQuoteCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
