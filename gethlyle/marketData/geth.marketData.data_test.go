package gethlylemarketdata

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/shopspring/decimal"
)

var DBColumns = []string{
	"id",                  //1
	"uuid",                //2
	"name",                //3
	"alternate_name",      //4
	"start_date",          //5
	"end_date",            //6
	"asset_id",            //7
	"open_usd",            //8
	"close_usd",           //9
	"high_usd",            //10
	"low_usd",             //11
	"price_usd",           //12
	"volume_usd",          //13
	"market_cap_usd",      //14
	"ticker",              //15
	"description",         //16
	"interval_id",         //17
	"market_data_type_id", //18
	"source_id",           //19
	"total_supply",        //20
	"max_supply",          //21
	"circulating_supply",  //22
	"sparkline_7d",        //23
	"created_by",          //24
	"created_at",          //25
	"updated_by",          //26
	"updated_at",          //27
	"geth_process_job_id", //28
}
var DBColumnsInsertGethMarketDataList = []string{
	"uuid",                //1
	"name",                //2
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
	"geth_process_job_id", //27
}
var fakeSparklineData = []decimal.Decimal{
	decimal.NewFromFloat(1.0),
	decimal.NewFromFloat(2.0),
	decimal.NewFromFloat(3.0),
	decimal.NewFromFloat(4.0),
	decimal.NewFromFloat(5.0),
	decimal.NewFromFloat(6.0),
	decimal.NewFromFloat(7.0),
}

var TestData1 = GethMarketData{
	ID:                utils.Ptr[int](1),
	UUID:              "01ef85e8-2c26-441e-8c7f-71d79518ad72",
	Name:              "Mog Coin ETH",
	AlternateName:     "Mog Coin ETH",
	StartDate:         utils.SampleCreatedAtTime,
	EndDate:           utils.SampleCreatedAtTime,
	AssetID:           utils.Ptr[int](1),
	OpenUSD:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	CloseUSD:          utils.Ptr[decimal.Decimal](decimal.NewFromFloat(7.0)),
	HighUSD:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(10.0)),
	LowUSD:            utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	PriceUSD:          utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	VolumeUSD:         utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	MarketCapUSD:      utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	Ticker:            "MOG",
	Description:       "Defillama url :https://coins.llama.fi/prices/historical/1720396800/coingecko:mog-coin",
	IntervalID:        utils.Ptr[int](5),
	MarketDataTypeID:  utils.Ptr[int](8),
	SourceID:          utils.Ptr[int](3),
	TotalSupply:       utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	MaxSupply:         utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	CirculatingSupply: utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	Sparkline7d:       fakeSparklineData,
	CreatedBy:         "SYSTEM",
	CreatedAt:         utils.SampleCreatedAtTime,
	UpdatedBy:         "SYSTEM",
	UpdatedAt:         utils.SampleCreatedAtTime,
	GethProcessJobID:  utils.Ptr[int](4571186),
}

var TestData2 = GethMarketData{
	ID:                utils.Ptr[int](2),
	UUID:              "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",
	Name:              "PEPE",
	AlternateName:     "PEPE",
	StartDate:         utils.SampleCreatedAtTime,
	EndDate:           utils.SampleCreatedAtTime,
	AssetID:           utils.Ptr[int](1),
	OpenUSD:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	CloseUSD:          utils.Ptr[decimal.Decimal](decimal.NewFromFloat(7.0)),
	HighUSD:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(10.0)),
	LowUSD:            utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	PriceUSD:          utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.0000113)),
	VolumeUSD:         utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	MarketCapUSD:      utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	Ticker:            "PEPE",
	Description:       "Defillama url :https://coins.llama.fi/prices/historical/1719964800/coingecko:pepe",
	IntervalID:        utils.Ptr[int](5),
	MarketDataTypeID:  utils.Ptr[int](8),
	SourceID:          utils.Ptr[int](3),
	TotalSupply:       utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	MaxSupply:         utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	CirculatingSupply: utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)),
	Sparkline7d:       fakeSparklineData,
	CreatedBy:         "SYSTEM",
	CreatedAt:         utils.SampleCreatedAtTime,
	UpdatedBy:         "SYSTEM",
	UpdatedAt:         utils.SampleCreatedAtTime,
	GethProcessJobID:  utils.Ptr[int](4514258),
}
var TestAllData = []GethMarketData{TestData1, TestData2}

func AddGethMarketDataToMockRows(mock pgxmock.PgxPoolIface, dataList []GethMarketData) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                //1
			data.UUID,              //2
			data.Name,              //3
			data.AlternateName,     //4
			data.StartDate,         //5
			data.EndDate,           //6
			data.AssetID,           //7
			data.OpenUSD,           //8
			data.CloseUSD,          //9
			data.HighUSD,           //10
			data.LowUSD,            //11
			data.PriceUSD,          //12
			data.VolumeUSD,         //13
			data.MarketCapUSD,      //14
			data.Ticker,            //15
			data.Description,       //16
			data.IntervalID,        //17
			data.MarketDataTypeID,  //18
			data.SourceID,          //19
			data.TotalSupply,       //20
			data.MaxSupply,         //21
			data.CirculatingSupply, //22
			data.Sparkline7d,       //23
			data.CreatedBy,         //24
			data.CreatedAt,         //25
			data.UpdatedBy,         //26
			data.UpdatedAt,         //27
			data.GethProcessJobID,  //28
		)
	}
	return rows
}
func TestGetMinAndMaxDatesFromGethMarketByAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	rows := mock.NewRows([]string{"min_date", "max_date"})
	mockRows := rows.AddRow(&startDate, &endDate)
	assetID := 1
	marketDataTypeID := utils.END_OF_DAY_MARKET_DATA_TYPE_STRUCTURED_VALUE_ID
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(assetID, marketDataTypeID).WillReturnRows(mockRows)
	startDateResult, endDateResult, err := GetMinAndMaxDatesFromGethMarketByAssetID(mock, &assetID, &marketDataTypeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMinAndMaxDatesFromGethMarketByAssetID", err)
	}
	if cmp.Equal(startDate, *startDateResult) == false {
		t.Errorf("Expected startDate From Method GetMinAndMaxDatesFromGethMarketByAssetID: %v is different from actual %v", startDate, *startDateResult)
	}
	if cmp.Equal(endDate, *endDateResult) == false {
		t.Errorf("Expected endDate From Method GetMinAndMaxDatesFromGethMarketByAssetID: %v is different from actual %v", endDate, *endDateResult)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinAndMaxDatesFromGethMarketByAssetIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 99999
	marketDataTypeID := utils.END_OF_DAY_MARKET_DATA_TYPE_STRUCTURED_VALUE_ID
	noRows := mock.NewRows([]string{"min_date", "max_date"})
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(assetID, marketDataTypeID).WillReturnRows(noRows)
	startDateResult, endDateResult, err := GetMinAndMaxDatesFromGethMarketByAssetID(mock, &assetID, &marketDataTypeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMinAndMaxDatesFromGethMarketByAssetID", err)
	}
	if startDateResult != nil {
		t.Errorf("startDateResult From Method GetMinAndMaxDatesFromGethMarketByAssetID: to be empty but got this: %v", startDateResult)
	}
	if endDateResult != nil {
		t.Errorf("Expected endDateResult From Method GetMinAndMaxDatesFromGethMarketByAssetID: to be empty but got this: %v", endDateResult)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinAndMaxDatesFromGethMarketByAssetIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 99999
	marketDataTypeID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(assetID, marketDataTypeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	startDateResult, endDateResult, err := GetMinAndMaxDatesFromGethMarketByAssetID(mock, &assetID, &marketDataTypeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMinAndMaxDatesFromGethMarketByAssetID", err)
	}
	if startDateResult != nil {
		t.Errorf("Expected startDateResult From Method GetMinAndMaxDatesFromGethMarketByAssetID: to be empty but got this: %v", startDateResult)
	}
	if endDateResult != nil {
		t.Errorf("Expected endDateResult From Method GetMinAndMaxDatesFromGethMarketByAssetID: to be empty but got this: %v", endDateResult)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketData(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethMarketData{targetData}
	marketDataID := targetData.ID
	mockRows := AddGethMarketDataToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(*marketDataID).WillReturnRows(mockRows)
	foundMarketData, err := GetGethMarketData(mock, marketDataID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMarketData", err)
	}
	if cmp.Equal(*foundMarketData, targetData) == false {
		t.Errorf("Expected GethMarketData From Method GetGethMarketData: %v is different from actual %v", foundMarketData, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(marketDataID).WillReturnRows(noRows)
	foundMarketData, err := GetGethMarketData(mock, &marketDataID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMarketData", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected GethMarketData From Method GetGethMarketData: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(marketDataID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketData, err := GetGethMarketData(mock, &marketDataID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMarketData", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected GethMarketData From Method GetGethMarketData: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(marketDataID).WillReturnRows(differentModelRows)
	foundMarketData, err := GetGethMarketData(mock, &marketDataID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMarketData", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected foundMarketData From Method GetGethMarketData: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataByAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethMarketData{targetData}
	startDate := targetData.StartDate
	startDateStr := startDate.Format(utils.LayoutISO)
	assetID := targetData.AssetID
	marketDataTypeID := targetData.MarketDataTypeID
	mockRows := AddGethMarketDataToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(startDateStr, *assetID, *marketDataTypeID).WillReturnRows(mockRows)
	foundMarketData, err := GetGethMarketDataByAssetID(mock, &startDate, assetID, marketDataTypeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMarketDataByAssetID", err)
	}
	if cmp.Equal(*foundMarketData, targetData) == false {
		t.Errorf("Expected GethMarketData From Method GetGethMarketDataByAssetID: %v is different from actual %v", foundMarketData, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataByAssetIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startDate := utils.SampleCreatedAtTime
	startDateStr := startDate.Format(utils.LayoutISO)
	assetID := 999
	marketDataTypeID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(startDateStr, assetID, marketDataTypeID).WillReturnRows(noRows)
	foundMarketData, err := GetGethMarketDataByAssetID(mock, &startDate, &assetID, &marketDataTypeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMarketDataByAssetID", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected GethMarketData From Method GetGethMarketDataByAssetID: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataByAssetIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startDate := utils.SampleCreatedAtTime
	startDateStr := startDate.Format(utils.LayoutISO)
	assetID := -1
	marketDataTypeID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(startDateStr, assetID, marketDataTypeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketData, err := GetGethMarketDataByAssetID(mock, &startDate, &assetID, &marketDataTypeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMarketDataByAssetID", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected GethMarketData From Method GetGethMarketDataByAssetID: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataByAssetIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startDate := utils.SampleCreatedAtTime
	startDateStr := startDate.Format(utils.LayoutISO)
	marketDataTypeID := -1
	assetID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(startDateStr, assetID, marketDataTypeID).WillReturnRows(differentModelRows)
	foundMarketData, err := GetGethMarketDataByAssetID(mock, &startDate, &assetID, &marketDataTypeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMarketDataByAssetID", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected foundMarketData From Method GetGethMarketDataByAssetID: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketData(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	marketDataID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_market_data").WithArgs(*marketDataID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethMarketData(mock, marketDataID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethMarketData", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveGethMarketData(mock, &marketDataID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_market_data").WithArgs(marketDataID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethMarketData(mock, &marketDataID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataFromBaseAssetBetweenDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	assetID := targetData.AssetID
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_market_data").WithArgs(*assetID, startDate, endDate).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethMarketDataFromBaseAssetBetweenDates(mock, assetID, &startDate, &endDate)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethMarketDataFromBaseAssetBetweenDates", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataFromBaseAssetBetweenDatesOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := -1
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveGethMarketDataFromBaseAssetBetweenDates(mock, &assetID, &startDate, &endDate)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataFromBaseAssetBetweenDatesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := -1
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_market_data").WithArgs(assetID, startDate, endDate).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethMarketDataFromBaseAssetBetweenDates(mock, &assetID, &startDate, &endDate)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataByMarketDataTypeIDFromBaseAssetBetweenDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	assetID := targetData.AssetID
	marketDataTypeID := targetData.MarketDataTypeID
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_market_data").WithArgs(startDate, endDate, *marketDataTypeID, *assetID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethMarketDataByMarketDataTypeIDFromBaseAssetBetweenDates(mock, assetID, marketDataTypeID, &startDate, &endDate)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethMarketDataByMarketDataTypeIDFromBaseAssetBetweenDates", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataByMarketDataTypeIDFromBaseAssetBetweenDatesOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 1
	marketDataTypeID := -1
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveGethMarketDataByMarketDataTypeIDFromBaseAssetBetweenDates(mock, &assetID, &marketDataTypeID, &startDate, &endDate)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataByMarketDataTypeIDFromBaseAssetBetweenDatesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 1
	marketDataTypeID := -1
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_market_data").WithArgs(startDate, endDate, marketDataTypeID, assetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethMarketDataByMarketDataTypeIDFromBaseAssetBetweenDates(mock, &assetID, &marketDataTypeID, &startDate, &endDate)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataByMarketDataTypeIDFromBaseAssetAsOfDate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	assetID := targetData.AssetID
	marketDataTypeID := targetData.MarketDataTypeID
	asOfDate := utils.SampleCreatedAtTime
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_market_data").WithArgs(asOfDate, *marketDataTypeID, *assetID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethMarketDataByMarketDataTypeIDFromBaseAssetAsOfDate(mock, assetID, marketDataTypeID, &asOfDate)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethMarketDataByMarketDataTypeIDFromBaseAssetAsOfDate", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataByMarketDataTypeIDFromBaseAssetAsOfDateOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 1
	marketDataTypeID := -1
	asOfDate := utils.SampleCreatedAtTime
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveGethMarketDataByMarketDataTypeIDFromBaseAssetAsOfDate(mock, &assetID, &marketDataTypeID, &asOfDate)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMarketDataByMarketDataTypeIDFromBaseAssetAsOfDateOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 1
	marketDataTypeID := -1
	asOfDate := utils.SampleCreatedAtTime
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_market_data").WithArgs(asOfDate, marketDataTypeID, assetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethMarketDataByMarketDataTypeIDFromBaseAssetAsOfDate(mock, &assetID, &marketDataTypeID, &asOfDate)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMarketData{TestData1, TestData2}
	mockRows := AddGethMarketDataToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WillReturnRows(mockRows)
	ids := []int{1, 2}
	foundMarketDataList, err := GetGethMarketDataList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMarketDataList", err)
	}
	testMarketDataList := TestAllData
	for i, foundMarketData := range foundMarketDataList {
		if cmp.Equal(foundMarketData, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMarketData From Method GetGethMarketDataList: %v is different from actual %v", foundMarketData, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	ids := []int{1, 2}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WillReturnRows(differentModelRows)
	foundMarketDataList, err := GetGethMarketDataList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMarketDataList", err)
	}
	if foundMarketDataList != nil {
		t.Errorf("Expected foundMarketDataList From Method GetGethMarketDataList: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	ids := make([]int, 0)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetGethMarketDataList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMarketDataList", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetGethMarketDataList: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataListByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMarketData{TestData1, TestData2}
	mockRows := AddGethMarketDataToMockRows(mock, dataList)
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(pq.Array(uuids)).WillReturnRows(mockRows)
	foundMarketDataList, err := GetGethMarketDataListByUUIDs(mock, uuids)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMarketDataListByUUIDs", err)
	}
	testMarketDataList := TestAllData
	for i, foundMarketData := range foundMarketDataList {
		if cmp.Equal(foundMarketData, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMarketData From Method GetGethMarketDataListByUUIDs: %v is different from actual %v", foundMarketData, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataListByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(pq.Array(uuids)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetGethMarketDataListByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMarketDataListByUUIDs", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetGethMarketDataListByUUIDs: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataListByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuids := []string{TestData1.UUID, TestData2.UUID}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(pq.Array(uuids)).WillReturnRows(differentModelRows)
	foundMarketDataList, err := GetGethMarketDataListByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMarketDataListByUUIDs", err)
	}
	if foundMarketDataList != nil {
		t.Errorf("Expected foundMarketDataList From Method GetGethMarketDataListByUUIDs: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffGethMarketDataList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMarketData{TestData1, TestData2}
	mockRows := AddGethMarketDataToMockRows(mock, dataList)
	diffInDate := 10
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(diffInDate).WillReturnRows(mockRows)
	foundMarketDataList, err := GetStartAndEndDateDiffGethMarketDataList(mock, &diffInDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndDateDiffGethMarketDataList", err)
	}
	testMarketDataList := TestAllData
	for i, foundMarketData := range foundMarketDataList {
		if cmp.Equal(foundMarketData, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMarketData From Method GetStartAndEndDateDiffGethMarketDataList: %v is different from actual %v", foundMarketData, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffGethMarketDataListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := 10
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(diffInDate).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetStartAndEndDateDiffGethMarketDataList(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffGethMarketDataList", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetStartAndEndDateDiffGethMarketDataList: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffGethMarketDataListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := 10
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WithArgs(diffInDate).WillReturnRows(differentModelRows)
	foundMarketDataList, err := GetStartAndEndDateDiffGethMarketDataList(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffGethMarketDataList", err)
	}
	if foundMarketDataList != nil {
		t.Errorf("Expected foundMarketDataList From Method GetStartAndEndDateDiffGethMarketDataList: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethMarketData(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_market_data").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.StartDate.Format(utils.LayoutPostgres), //3
		targetData.EndDate.Format(utils.LayoutPostgres),   //4
		targetData.AssetID,               //5
		targetData.OpenUSD,               //6
		targetData.CloseUSD,              //7
		targetData.HighUSD,               //8
		targetData.LowUSD,                //9
		targetData.PriceUSD,              //10
		targetData.VolumeUSD,             //11
		targetData.MarketCapUSD,          //12
		targetData.Ticker,                //13
		targetData.Description,           //14
		targetData.IntervalID,            //15
		targetData.MarketDataTypeID,      //16
		targetData.SourceID,              //17
		targetData.TotalSupply,           //18
		targetData.MaxSupply,             //19
		targetData.CirculatingSupply,     //20
		pq.Array(targetData.Sparkline7d), //21
		targetData.UpdatedBy,             //22
		targetData.GethProcessJobID,      //23
		targetData.ID,                    //24
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethMarketData(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethMarketData", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethMarketDataOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateGethMarketData(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethMarketDataOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethMarketData(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethMarketDataOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_market_data").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.StartDate.Format(utils.LayoutPostgres), //3
		targetData.EndDate.Format(utils.LayoutPostgres),   //4
		targetData.AssetID,               //5
		targetData.OpenUSD,               //6
		targetData.CloseUSD,              //7
		targetData.HighUSD,               //8
		targetData.LowUSD,                //9
		targetData.PriceUSD,              //10
		targetData.VolumeUSD,             //11
		targetData.MarketCapUSD,          //12
		targetData.Ticker,                //13
		targetData.Description,           //14
		targetData.IntervalID,            //15
		targetData.MarketDataTypeID,      //16
		targetData.SourceID,              //17
		targetData.TotalSupply,           //18
		targetData.MaxSupply,             //19
		targetData.CirculatingSupply,     //20
		pq.Array(targetData.Sparkline7d), //21
		targetData.UpdatedBy,             //22
		targetData.GethProcessJobID,      //23
		targetData.ID,                    //24
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethMarketData(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMarketData(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_market_data").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.StartDate.Format(utils.LayoutPostgres), //3
		targetData.EndDate.Format(utils.LayoutPostgres),   //4
		targetData.AssetID,               //5
		targetData.OpenUSD,               //6
		targetData.CloseUSD,              //7
		targetData.HighUSD,               //8
		targetData.LowUSD,                //9
		targetData.PriceUSD,              //10
		targetData.VolumeUSD,             //11
		targetData.MarketCapUSD,          //12
		targetData.Ticker,                //13
		targetData.Description,           //14
		targetData.IntervalID,            //15
		targetData.MarketDataTypeID,      //16
		targetData.SourceID,              //17
		targetData.TotalSupply,           //18
		targetData.MaxSupply,             //19
		targetData.CirculatingSupply,     //20
		pq.Array(targetData.Sparkline7d), //21
		targetData.CreatedBy,             //22
		targetData.GethProcessJobID,      //23
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	marketDataID, err := InsertGethMarketData(mock, &targetData)
	if marketDataID < 0 {
		t.Fatalf("marketDataID should not be negative ID: %d", marketDataID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethMarketData", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMarketDataOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertGethMarketData(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMarketDataOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_market_data").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.StartDate.Format(utils.LayoutPostgres), //3
		targetData.EndDate.Format(utils.LayoutPostgres),   //4
		targetData.AssetID,               //5
		targetData.OpenUSD,               //6
		targetData.CloseUSD,              //7
		targetData.HighUSD,               //8
		targetData.LowUSD,                //9
		targetData.PriceUSD,              //10
		targetData.VolumeUSD,             //11
		targetData.MarketCapUSD,          //12
		targetData.Ticker,                //13
		targetData.Description,           //14
		targetData.IntervalID,            //15
		targetData.MarketDataTypeID,      //16
		targetData.SourceID,              //17
		targetData.TotalSupply,           //18
		targetData.MaxSupply,             //19
		targetData.CirculatingSupply,     //20
		pq.Array(targetData.Sparkline7d), //21
		targetData.CreatedBy,             //22
		targetData.GethProcessJobID,      //23
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	marketDataID, err := InsertGethMarketData(mock, &targetData)
	if marketDataID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error marketDataID: %d", marketDataID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMarketDataOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_market_data").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.StartDate.Format(utils.LayoutPostgres), //3
		targetData.EndDate.Format(utils.LayoutPostgres),   //4
		targetData.AssetID,               //5
		targetData.OpenUSD,               //6
		targetData.CloseUSD,              //7
		targetData.HighUSD,               //8
		targetData.LowUSD,                //9
		targetData.PriceUSD,              //10
		targetData.VolumeUSD,             //11
		targetData.MarketCapUSD,          //12
		targetData.Ticker,                //13
		targetData.Description,           //14
		targetData.IntervalID,            //15
		targetData.MarketDataTypeID,      //16
		targetData.SourceID,              //17
		targetData.TotalSupply,           //18
		targetData.MaxSupply,             //19
		targetData.CirculatingSupply,     //20
		pq.Array(targetData.Sparkline7d), //21
		targetData.CreatedBy,             //22
		targetData.GethProcessJobID,      //23
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	marketDataID, err := InsertGethMarketData(mock, &targetData)
	if marketDataID >= 0 {
		t.Fatalf("Expecting -1 for marketDataID because of error marketDataID: %d", marketDataID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMarketDataList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_market_data"}, DBColumnsInsertGethMarketDataList)
	targetData := TestAllData
	err = InsertGethMarketDataList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethMarketDataList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMarketDataListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_market_data"}, DBColumnsInsertGethMarketDataList).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertGethMarketDataList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethMarketDataToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"asset_id = 1", "source_id=1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WillReturnRows(mockRows)
	foundMarketDataList, err := GetGethMarketDataListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMarketDataListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundMarketDataList[i]) == false {
			t.Errorf("Expected sourceData From Method GetGethMarketDataListByPagination: %v is different from actual %v", sourceData, foundMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"asset_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetGethMarketDataListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMarketDataListByPagination", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetGethMarketDataListByPagination: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMarketDataListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"asset_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_market_data").WillReturnRows(differentModelRows)
	foundMarketDataList, err := GetGethMarketDataListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMarketDataListByPagination", err)
	}
	if foundMarketDataList != nil {
		t.Errorf("Expected From Method GetGethMarketDataListByPagination: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestGetTotalGethMarketDataCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_market_data").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalGethMarketDataCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalGethMarketDataCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalGethMarketDataCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMarketDataCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_market_data").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalGethMarketDataCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalGethMarketDataCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalGethMarketDataCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
