package marketdata

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
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
}
var DBColumnsInsertMarketDataList = []string{
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

var TestData1 = MarketData{
	ID:                utils.Ptr[int](1),
	UUID:              "01ef85e8-2c26-441e-8c7f-71d79518ad72",
	Name:              "Mog Coin ETH",
	AlternateName:     "Mog Coin ETH",
	StartDate:         utils.SampleCreatedAtTime,
	EndDate:           utils.SampleCreatedAtTime,
	AssetID:           utils.Ptr[int](1),
	OpenUSD:           utils.Ptr[float64](float64(5.0)),
	CloseUSD:          utils.Ptr[float64](float64(7.0)),
	HighUSD:           utils.Ptr[float64](float64(10.0)),
	LowUSD:            utils.Ptr[float64](float64(5.0)),
	PriceUSD:          utils.Ptr[float64](float64(5.0)),
	VolumeUSD:         utils.Ptr[float64](float64(5.0)),
	MarketCapUSD:      utils.Ptr[float64](float64(5.0)),
	Ticker:            "MOG",
	Description:       "Defillama url :https://coins.llama.fi/prices/historical/1720396800/coingecko:mog-coin",
	IntervalID:        utils.Ptr[int](5),
	MarketDataTypeID:  utils.Ptr[int](8),
	SourceID:          utils.Ptr[int](3),
	TotalSupply:       utils.Ptr[float64](float64(5.0)),
	MaxSupply:         utils.Ptr[float64](float64(5.0)),
	CirculatingSupply: utils.Ptr[float64](float64(5.0)),
	Sparkline7d:       fakeSparklineData,
	CreatedBy:         "SYSTEM",
	CreatedAt:         utils.SampleCreatedAtTime,
	UpdatedBy:         "SYSTEM",
	UpdatedAt:         utils.SampleCreatedAtTime,
}

var TestData2 = MarketData{
	ID:                utils.Ptr[int](2),
	UUID:              "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",
	Name:              "PEPE",
	AlternateName:     "PEPE",
	StartDate:         utils.SampleCreatedAtTime,
	EndDate:           utils.SampleCreatedAtTime,
	AssetID:           utils.Ptr[int](1),
	OpenUSD:           utils.Ptr[float64](float64(5.0)),
	CloseUSD:          utils.Ptr[float64](float64(7.0)),
	HighUSD:           utils.Ptr[float64](float64(10.0)),
	LowUSD:            utils.Ptr[float64](float64(5.0)),
	PriceUSD:          utils.Ptr[float64](float64(0.0000113)),
	VolumeUSD:         utils.Ptr[float64](float64(5.0)),
	MarketCapUSD:      utils.Ptr[float64](float64(5.0)),
	Ticker:            "PEPE",
	Description:       "Defillama url :https://coins.llama.fi/prices/historical/1719964800/coingecko:pepe",
	IntervalID:        utils.Ptr[int](5),
	MarketDataTypeID:  utils.Ptr[int](8),
	SourceID:          utils.Ptr[int](3),
	TotalSupply:       utils.Ptr[float64](float64(5.0)),
	MaxSupply:         utils.Ptr[float64](float64(5.0)),
	CirculatingSupply: utils.Ptr[float64](float64(5.0)),
	Sparkline7d:       fakeSparklineData,
	CreatedBy:         "SYSTEM",
	CreatedAt:         utils.SampleCreatedAtTime,
	UpdatedBy:         "SYSTEM",
	UpdatedAt:         utils.SampleCreatedAtTime,
}
var TestAllData = []MarketData{TestData1, TestData2}

func AddMarketDataToMockRows(mock pgxmock.PgxPoolIface, dataList []MarketData) *pgxmock.Rows {
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
		)
	}
	return rows
}

func TestGetMarketData(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []MarketData{targetData}
	marketDataID := targetData.ID
	mockRows := AddMarketDataToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WithArgs(*marketDataID).WillReturnRows(mockRows)
	foundMarketData, err := GetMarketData(mock, marketDataID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketData", err)
	}
	if cmp.Equal(*foundMarketData, targetData) == false {
		t.Errorf("Expected MarketData From Method GetMarketData: %v is different from actual %v", foundMarketData, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WithArgs(marketDataID).WillReturnRows(noRows)
	foundMarketData, err := GetMarketData(mock, &marketDataID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketData", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected MarketData From Method GetMarketData: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WithArgs(marketDataID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketData, err := GetMarketData(mock, &marketDataID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketData", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected MarketData From Method GetMarketData: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	marketDataID := -1
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WithArgs(marketDataID).WillReturnRows(differentModelRows)
	foundMarketData, err := GetMarketData(mock, &marketDataID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketData", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected foundMarketData From Method GetMarketData: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketData(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	marketDataID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data").WithArgs(*marketDataID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveMarketData(mock, marketDataID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveMarketData", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketDataOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveMarketData(mock, &taxID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestRemoveMarketDataOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data").WithArgs(marketDataID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveMarketData(mock, &marketDataID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketDataFromBaseAssetBetweenDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	assetID := targetData.AssetID
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	// from marketdataquote.RemoveMarketDataQuoteFromBaseAssetBetweenDates
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data_quotes").WithArgs(*assetID, startDate, endDate).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// from RemoveMarketDataFromBaseAssetBetweenDates
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data").WithArgs(*assetID, startDate, endDate).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveMarketDataFromBaseAssetBetweenDates(mock, assetID, startDate, endDate)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveMarketDataFromBaseAssetBetweenDates", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketDataFromBaseAssetBetweenDatesOnFailureAtRemoveMarketDataQuoteFromBaseAssetBetweenDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 1
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	// from marketdataquote.RemoveMarketDataQuoteFromBaseAssetBetweenDates
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data_quotes").WithArgs(assetID, startDate, endDate).WillReturnError(fmt.Errorf("Error at RemoveMarketDataQuoteFromBaseAssetBetweenDates ID"))
	mock.ExpectRollback()
	err = RemoveMarketDataFromBaseAssetBetweenDates(mock, &assetID, startDate, endDate)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketDataFromBaseAssetBetweenDatesOnFailureAtBeginOfRemoveMarketDataFromBaseAssetBetweenDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 1
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	// from marketdataquote.RemoveMarketDataQuoteFromBaseAssetBetweenDates
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data_quotes").WithArgs(assetID, startDate, endDate).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// from RemoveMarketDataFromBaseAssetBetweenDates
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveMarketDataFromBaseAssetBetweenDates(mock, &assetID, startDate, endDate)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketDataFromBaseAssetBetweenDatesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 1
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	// from marketdataquote.RemoveMarketDataQuoteFromBaseAssetBetweenDates
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data_quotes").WithArgs(assetID, startDate, endDate).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// from RemoveMarketDataFromBaseAssetBetweenDates
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data").WithArgs(assetID, startDate, endDate).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveMarketDataFromBaseAssetBetweenDates(mock, &assetID, startDate, endDate)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []MarketData{TestData1, TestData2}
	mockRows := AddMarketDataToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WillReturnRows(mockRows)
	ids := []int{1, 2}
	foundMarketDataList, err := GetMarketDataList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketDataList", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundMarketDataList[i]) == false {
			t.Errorf("Expected sourceData From Method GetMarketDataList: %v is different from actual %v", sourceData, foundMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	ids := make([]int, 0)
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetMarketDataList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataList", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetMarketDataList: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	ids := make([]int, 0)
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WillReturnRows(differentModelRows)
	foundMarketDataList, err := GetMarketDataList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataList", err)
	}
	if foundMarketDataList != nil {
		t.Errorf("Expected foundMarketDataList From Method GetMarketDataList: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataListByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []MarketData{TestData1, TestData2}
	mockRows := AddMarketDataToMockRows(mock, dataList)
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WithArgs(pq.Array(uuids)).WillReturnRows(mockRows)
	foundMarketDataList, err := GetMarketDataListByUUIDs(mock, uuids)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketDataListByUUIDs", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundMarketDataList[i]) == false {
			t.Errorf("Expected sourceData From Method GetMarketDataListByUUIDs: %v is different from actual %v", sourceData, foundMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataListByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WithArgs(pq.Array(uuids)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetMarketDataListByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataListByUUIDs", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetMarketDataListByUUIDs: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataListByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WithArgs(pq.Array(uuids)).WillReturnRows(differentModelRows)
	foundMarketDataList, err := GetMarketDataListByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataListByUUIDs", err)
	}
	if foundMarketDataList != nil {
		t.Errorf("Expected foundMarketDataList From Method GetMarketDataListByUUIDs: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffMarketDataList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []MarketData{TestData1, TestData2}
	mockRows := AddMarketDataToMockRows(mock, dataList)
	diffInDate := 10
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WithArgs(diffInDate).WillReturnRows(mockRows)
	foundMarketDataList, err := GetStartAndEndDateDiffMarketDataList(mock, &diffInDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndDateDiffMarketDataList", err)
	}
	testMarketDataList := TestAllData
	for i, foundMarketData := range foundMarketDataList {
		if cmp.Equal(foundMarketData, testMarketDataList[i]) == false {
			t.Errorf("Expected MarketData From Method GetStartAndEndDateDiffMarketDataList: %v is different from actual %v", foundMarketData, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffMarketDataListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := 10
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WithArgs(diffInDate).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetStartAndEndDateDiffMarketDataList(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffMarketDataList", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetStartAndEndDateDiffMarketDataList: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffMarketDataListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	diffInDate := 10
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WithArgs(diffInDate).WillReturnRows(differentModelRows)
	foundMarketDataList, err := GetStartAndEndDateDiffMarketDataList(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffMarketDataList", err)
	}
	if foundMarketDataList != nil {
		t.Errorf("Expected foundMarketDataList From Method GetStartAndEndDateDiffMarketDataList: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateMarketData(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE market_data").WithArgs(
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
		targetData.ID,                    //23
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateMarketData(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateMarketData", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateMarketDataOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateMarketData(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateMarketDataOnFailureAtBegin(t *testing.T) {
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
	err = UpdateMarketData(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateMarketDataOnFailure(t *testing.T) {
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
	mock.ExpectExec("^UPDATE market_data").WithArgs(
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
		targetData.ID,                    //23
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateMarketData(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketData(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO market_data").WithArgs(
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
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	marketDataID, err := InsertMarketData(mock, &targetData)
	if marketDataID < 0 {
		t.Fatalf("marketDataID should not be negative ID: %d", marketDataID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertMarketData", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertMarketData(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO market_data").WithArgs(
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
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	marketDataID, err := InsertMarketData(mock, &targetData)
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

func TestInsertMarketDataOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO market_data").WithArgs(
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
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	marketDataID, err := InsertMarketData(mock, &targetData)
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

func TestInsertMarketDataList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"market_data"}, DBColumnsInsertMarketDataList)
	targetData := TestAllData
	err = InsertMarketDataList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertMarketDataList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"market_data"}, DBColumnsInsertMarketDataList).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertMarketDataList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddMarketDataToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"market_data_type_id = 1", "source_id=1"}
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WillReturnRows(mockRows)
	foundMarketDataList, err := GetMarketDataListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketDataListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundMarketDataList[i]) == false {
			t.Errorf("Expected sourceData From Method GetMarketDataListByPagination: %v is different from actual %v", sourceData, foundMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"market_data_type_id = 1", "source_id=1"}
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetMarketDataListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataListByPagination", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetMarketDataListByPagination: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"market_data_type_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM market_data").WillReturnRows(differentModelRows)
	foundMarketDataList, err := GetMarketDataListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataListByPagination", err)
	}
	if foundMarketDataList != nil {
		t.Errorf("Expected From Method GetMarketDataListByPagination: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMarketDataCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM market_data").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalMarketDataCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalMarketDataCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalMarketDataCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM market_data").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalMarketDataCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalMarketDataCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalMarketDataCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
