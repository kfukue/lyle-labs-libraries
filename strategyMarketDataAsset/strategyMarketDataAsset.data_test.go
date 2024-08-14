package strategymarketdataasset

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
	"id",             //1
	"strategy_id",    //2
	"base_asset_id",  //3
	"quote_asset_id", //4
	"uuid",           //5
	"name",           //6
	"alternate_name", //7
	"start_date",     //8
	"end_date",       //9
	"ticker",         //10
	"description",    //11
	"source_id",      //12
	"frequency_id",   //13
	"created_by",     //14
	"created_at",     //15
	"updated_by",     //16
	"updated_at",     //17
}
var DBColumnsInsertStrategyMarketDataAssets = []string{
	"strategy_id",    //1
	"base_asset_id",  //2
	"quote_asset_id", //3
	"uuid",           //4
	"name",           //5
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
}

var TestData1 = StrategyMarketDataAsset{
	ID:            utils.Ptr[int](1),                      //1
	StrategyID:    utils.Ptr[int](1),                      //2
	BaseAssetID:   utils.Ptr[int](1),                      //3
	QuoteAssetID:  utils.Ptr[int](2),                      //4
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72", //5
	Name:          "VELO/USD",                             //6
	AlternateName: "VELO/USD",                             //7
	StartDate:     utils.SampleCreatedAtTime,              //8
	EndDate:       utils.SampleCreatedAtTime,              //9
	Ticker:        "VELO/USD",                             //10
	Description:   "",                                     //11
	SourceID:      utils.Ptr[int](2),                      //12
	FrequencyID:   utils.Ptr[int](2),                      //13
	CreatedBy:     "SYSTEM",                               //14
	CreatedAt:     utils.SampleCreatedAtTime,              //15
	UpdatedBy:     "SYSTEM",                               //16
	UpdatedAt:     utils.SampleCreatedAtTime,              //17

}

var TestData2 = StrategyMarketDataAsset{
	ID:            utils.Ptr[int](3),                      //1
	StrategyID:    utils.Ptr[int](1),                      //2
	BaseAssetID:   utils.Ptr[int](3),                      //3
	QuoteAssetID:  utils.Ptr[int](2),                      //4
	UUID:          "4f0d5402-7a7c-402d-a7fc-c56a02b13e03", //5
	Name:          "VELO/USD",                             //6
	AlternateName: "VELO/USD",                             //7
	StartDate:     utils.SampleCreatedAtTime,              //8
	EndDate:       utils.SampleCreatedAtTime,              //9
	Ticker:        "OP/USD",                               //10
	Description:   "",                                     //11
	SourceID:      utils.Ptr[int](3),                      //12
	FrequencyID:   utils.Ptr[int](1),                      //13
	CreatedBy:     "SYSTEM",                               //14
	CreatedAt:     utils.SampleCreatedAtTime,              //15
	UpdatedBy:     "SYSTEM",                               //16
	UpdatedAt:     utils.SampleCreatedAtTime,              //17
}
var TestAllData = []StrategyMarketDataAsset{TestData1, TestData2}

func AddStrategyMarketDataAssetToMockRows(mock pgxmock.PgxPoolIface, dataList []StrategyMarketDataAsset) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,            //1
			data.StrategyID,    //2
			data.BaseAssetID,   //3
			data.QuoteAssetID,  //4
			data.UUID,          //5
			data.Name,          //6
			data.AlternateName, //7
			data.StartDate,     //8
			data.EndDate,       //9
			data.Ticker,        //10
			data.Description,   //11
			data.SourceID,      //12
			data.FrequencyID,   //13
			data.CreatedBy,     //14
			data.CreatedAt,     //15
			data.UpdatedBy,     //16
			data.UpdatedAt,     //17
		)
	}
	return rows
}

func TestGetStrategyMarketDataAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []StrategyMarketDataAsset{targetData}
	strategyMarketDataAssetID := targetData.ID
	mockRows := AddStrategyMarketDataAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(*strategyMarketDataAssetID).WillReturnRows(mockRows)
	foundStrategyMarketDataAsset, err := GetStrategyMarketDataAsset(mock, strategyMarketDataAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyMarketDataAsset", err)
	}
	if cmp.Equal(*foundStrategyMarketDataAsset, targetData) == false {
		t.Errorf("Expected StrategyMarketDataAsset From Method GetStrategyMarketDataAsset: %v is different from actual %v", foundStrategyMarketDataAsset, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyMarketDataAssetID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(strategyMarketDataAssetID).WillReturnRows(noRows)
	foundStrategyMarketDataAsset, err := GetStrategyMarketDataAsset(mock, &strategyMarketDataAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyMarketDataAsset", err)
	}
	if foundStrategyMarketDataAsset != nil {
		t.Errorf("Expected StrategyMarketDataAsset From Method GetStrategyMarketDataAsset: to be empty but got this: %v", foundStrategyMarketDataAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyMarketDataAssetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(strategyMarketDataAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyMarketDataAsset, err := GetStrategyMarketDataAsset(mock, &strategyMarketDataAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyMarketDataAsset", err)
	}
	if foundStrategyMarketDataAsset != nil {
		t.Errorf("Expected StrategyMarketDataAsset From Method GetStrategyMarketDataAsset: to be empty but got this: %v", foundStrategyMarketDataAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyMarketDataAssetID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(strategyMarketDataAssetID).WillReturnRows(differentModelRows)
	foundStrategyMarketDataAsset, err := GetStrategyMarketDataAsset(mock, &strategyMarketDataAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyMarketDataAsset", err)
	}
	if foundStrategyMarketDataAsset != nil {
		t.Errorf("Expected StrategyMarketDataAsset From Method GetStrategyMarketDataAsset: to be empty but got this: %v", foundStrategyMarketDataAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStrategyMarketDataAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	strategyMarketDataAssetID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM strategy_market_data_assets").WithArgs(*strategyMarketDataAssetID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveStrategyMarketDataAsset(mock, strategyMarketDataAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveStrategyMarketDataAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStrategyMarketDataAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyMarketDataAssetID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveStrategyMarketDataAsset(mock, &strategyMarketDataAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStrategyMarketDataAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyMarketDataAssetID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM strategy_market_data_assets").WithArgs(strategyMarketDataAssetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveStrategyMarketDataAsset(mock, &strategyMarketDataAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStrategyMarketDataAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WillReturnRows(mockRows)
	ids := []int{1}
	foundStrategyMarketDataAssetList, err := GetStrategyMarketDataAssets(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyMarketDataAssets", err)
	}
	for i, sourceStrategyMarketDataAsset := range dataList {
		if cmp.Equal(sourceStrategyMarketDataAsset, foundStrategyMarketDataAssetList[i]) == false {
			t.Errorf("Expected StrategyMarketDataAsset From Method GetStrategyMarketDataAssets: %v is different from actual %v", sourceStrategyMarketDataAsset, foundStrategyMarketDataAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1}
	foundStrategyMarketDataAssetList, err := GetStrategyMarketDataAssets(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyMarketDataAssets", err)
	}
	if len(foundStrategyMarketDataAssetList) != 0 {
		t.Errorf("Expected From Method GetStrategyMarketDataAssets: to be empty but got this: %v", foundStrategyMarketDataAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WillReturnRows(differentModelRows)
	ids := []int{1}
	foundStrategyMarketDataAsset, err := GetStrategyMarketDataAssets(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyMarketDataAssets", err)
	}
	if foundStrategyMarketDataAsset != nil {
		t.Errorf("Expected StrategyMarketDataAsset From Method GetStrategyMarketDataAssets: to be empty but got this: %v", foundStrategyMarketDataAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStrategyMarketDataAssetToMockRows(mock, dataList)
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(pq.Array(UUIDList)).WillReturnRows(mockRows)
	foundStrategyMarketDataAssetList, err := GetStrategyMarketDataAssetsByUUIDs(mock, UUIDList)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyMarketDataAssetsByUUIDs", err)
	}
	for i, sourceStrategyMarketDataAsset := range dataList {
		if cmp.Equal(sourceStrategyMarketDataAsset, foundStrategyMarketDataAssetList[i]) == false {
			t.Errorf("Expected StrategyMarketDataAsset From Method GetStrategyMarketDataAssetsByUUIDs: %v is different from actual %v", sourceStrategyMarketDataAsset, foundStrategyMarketDataAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(pq.Array(UUIDList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyMarketDataAssetList, err := GetStrategyMarketDataAssetsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyMarketDataAssetsByUUIDs", err)
	}
	if len(foundStrategyMarketDataAssetList) != 0 {
		t.Errorf("Expected From Method GetStrategyMarketDataAssetsByUUIDs: to be empty but got this: %v", foundStrategyMarketDataAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetsByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(pq.Array(UUIDList)).WillReturnRows(differentModelRows)
	foundStrategyMarketDataAsset, err := GetStrategyMarketDataAssetsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyMarketDataAssetsByUUIDs", err)
	}
	if foundStrategyMarketDataAsset != nil {
		t.Errorf("Expected StrategyMarketDataAsset From Method GetStrategyMarketDataAssetsByUUIDs: to be empty but got this: %v", foundStrategyMarketDataAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetsByStrategyID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStrategyMarketDataAssetToMockRows(mock, dataList)
	strategyID := 1
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(strategyID).WillReturnRows(mockRows)
	foundStrategyMarketDataAssetList, err := GetStrategyMarketDataAssetsByStrategyID(mock, &strategyID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyMarketDataAssetsByStrategyID", err)
	}
	for i, sourceStrategyMarketDataAsset := range dataList {
		if cmp.Equal(sourceStrategyMarketDataAsset, foundStrategyMarketDataAssetList[i]) == false {
			t.Errorf("Expected StrategyMarketDataAsset From Method GetStrategyMarketDataAssetsByStrategyID: %v is different from actual %v", sourceStrategyMarketDataAsset, foundStrategyMarketDataAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetsByStrategyIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -1
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(strategyID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyMarketDataAssetList, err := GetStrategyMarketDataAssetsByStrategyID(mock, &strategyID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyMarketDataAssetsByStrategyID", err)
	}
	if len(foundStrategyMarketDataAssetList) != 0 {
		t.Errorf("Expected From Method GetStrategyMarketDataAssetsByStrategyID: to be empty but got this: %v", foundStrategyMarketDataAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetsByStrategyIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := 1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(strategyID).WillReturnRows(differentModelRows)
	foundStrategyMarketDataAsset, err := GetStrategyMarketDataAssetsByStrategyID(mock, &strategyID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyMarketDataAssetsByStrategyID", err)
	}
	if foundStrategyMarketDataAsset != nil {
		t.Errorf("Expected StrategyMarketDataAsset From Method GetStrategyMarketDataAssetsByStrategyID: to be empty but got this: %v", foundStrategyMarketDataAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStrategyMarketDataAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStrategyMarketDataAssetToMockRows(mock, dataList)
	diffInDate := 2
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(diffInDate).WillReturnRows(mockRows)
	foundStrategyMarketDataAssetList, err := GetStartAndEndDateDiffStrategyMarketDataAssets(mock, &diffInDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndDateDiffStrategyMarketDataAssets", err)
	}
	for i, sourceStrategyMarketDataAsset := range dataList {
		if cmp.Equal(sourceStrategyMarketDataAsset, foundStrategyMarketDataAssetList[i]) == false {
			t.Errorf("Expected StrategyMarketDataAsset From Method GetStartAndEndDateDiffStrategyMarketDataAssets: %v is different from actual %v", sourceStrategyMarketDataAsset, foundStrategyMarketDataAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStrategyMarketDataAssetsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := -1
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(diffInDate).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyMarketDataAssetList, err := GetStartAndEndDateDiffStrategyMarketDataAssets(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffStrategyMarketDataAssets", err)
	}
	if len(foundStrategyMarketDataAssetList) != 0 {
		t.Errorf("Expected From Method GetStartAndEndDateDiffStrategyMarketDataAssets: to be empty but got this: %v", foundStrategyMarketDataAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStrategyMarketDataAssetsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := 1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WithArgs(diffInDate).WillReturnRows(differentModelRows)
	foundStrategyMarketDataAsset, err := GetStartAndEndDateDiffStrategyMarketDataAssets(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffStrategyMarketDataAssets", err)
	}
	if foundStrategyMarketDataAsset != nil {
		t.Errorf("Expected StrategyMarketDataAsset From Method GetStartAndEndDateDiffStrategyMarketDataAssets: to be empty but got this: %v", foundStrategyMarketDataAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateStrategyMarketDataAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE strategy_market_data_assets").WithArgs(
		targetData.StrategyID,    //1
		targetData.BaseAssetID,   //2
		targetData.QuoteAssetID,  //3
		targetData.UUID,          //4
		targetData.Name,          //5
		targetData.AlternateName, //6
		targetData.StartDate,     //7
		targetData.EndDate,       //8
		targetData.Ticker,        //9
		targetData.Description,   //10
		targetData.SourceID,      //11
		targetData.FrequencyID,   //12
		targetData.UpdatedBy,     //13
		targetData.ID,            //14
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateStrategyMarketDataAsset(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateStrategyMarketDataAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStrategyMarketDataAssetOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateStrategyMarketDataAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStrategyMarketDataAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateStrategyMarketDataAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStrategyMarketDataAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE strategy_market_data_assets").WithArgs(
		targetData.StrategyID,    //1
		targetData.BaseAssetID,   //2
		targetData.QuoteAssetID,  //3
		targetData.UUID,          //4
		targetData.Name,          //5
		targetData.AlternateName, //6
		targetData.StartDate,     //7
		targetData.EndDate,       //8
		targetData.Ticker,        //9
		targetData.Description,   //10
		targetData.SourceID,      //11
		targetData.FrequencyID,   //12
		targetData.UpdatedBy,     //13
		targetData.ID,            //14
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateStrategyMarketDataAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyMarketDataAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO strategy_market_data_assets").WithArgs(
		targetData.StrategyID,    //1
		targetData.BaseAssetID,   //2
		targetData.QuoteAssetID,  //3
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.StartDate,     //6
		targetData.EndDate,       //7
		targetData.Ticker,        //8
		targetData.Description,   //9
		targetData.SourceID,      //10
		targetData.FrequencyID,   //11
		targetData.CreatedBy,     //12
	).WillReturnRows(pgxmock.NewRows([]string{"strategy_market_data_asset_id"}).AddRow(1))
	mock.ExpectCommit()
	strategyMarketDataAssetID, err := InsertStrategyMarketDataAsset(mock, &targetData)
	if strategyMarketDataAssetID < 0 {
		t.Fatalf("strategyMarketDataAssetID should not be negative ID: %d", strategyMarketDataAssetID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertStrategyMarketDataAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyMarketDataAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertStrategyMarketDataAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyMarketDataAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO strategy_market_data_assets").WithArgs(
		targetData.StrategyID,    //1
		targetData.BaseAssetID,   //2
		targetData.QuoteAssetID,  //3
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.StartDate,     //6
		targetData.EndDate,       //7
		targetData.Ticker,        //8
		targetData.Description,   //9
		targetData.SourceID,      //10
		targetData.FrequencyID,   //11
		targetData.CreatedBy,     //12
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	strategyMarketDataAssetID, err := InsertStrategyMarketDataAsset(mock, &targetData)
	if strategyMarketDataAssetID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error strategyMarketDataAssetID: %d", strategyMarketDataAssetID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyMarketDataAssetOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO strategy_market_data_assets").WithArgs(
		targetData.StrategyID,    //1
		targetData.BaseAssetID,   //2
		targetData.QuoteAssetID,  //3
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.StartDate,     //6
		targetData.EndDate,       //7
		targetData.Ticker,        //8
		targetData.Description,   //9
		targetData.SourceID,      //10
		targetData.FrequencyID,   //11
		targetData.CreatedBy,     //12
	).WillReturnRows(pgxmock.NewRows([]string{"strategy_market_data_asset_id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	strategyMarketDataAssetID, err := InsertStrategyMarketDataAsset(mock, &targetData)
	if strategyMarketDataAssetID >= 0 {
		t.Fatalf("Expecting -1 for strategyMarketDataAssetID because of error strategyMarketDataAssetID: %d", strategyMarketDataAssetID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyMarketDataAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"strategy_market_data_assets"}, DBColumnsInsertStrategyMarketDataAssets)
	targetData := TestAllData
	err = InsertStrategyMarketDataAssets(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertStrategyMarketDataAssets", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyMarketDataAssetsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"strategy_market_data_assets"}, DBColumnsInsertStrategyMarketDataAssets).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertStrategyMarketDataAssets(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddStrategyMarketDataAssetToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"strategy_id = 1", "quote_asset_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WillReturnRows(mockRows)
	foundStrategyMarketDataAssetList, err := GetStrategyMarketDataAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyMarketDataAssetListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundStrategyMarketDataAssetList[i]) == false {
			t.Errorf("Expected sourceData From Method GetStrategyMarketDataAssetListByPagination: %v is different from actual %v", sourceData, foundStrategyMarketDataAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"strategy_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyMarketDataAssetList, err := GetStrategyMarketDataAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyMarketDataAssetListByPagination", err)
	}
	if len(foundStrategyMarketDataAssetList) != 0 {
		t.Errorf("Expected From Method GetStrategyMarketDataAssetListByPagination: to be empty but got this: %v", foundStrategyMarketDataAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyMarketDataAssetListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"strategy_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_market_data_assets").WillReturnRows(differentModelRows)
	foundStrategyMarketDataAssetList, err := GetStrategyMarketDataAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyMarketDataAssetListByPagination", err)
	}
	if len(foundStrategyMarketDataAssetList) != 0 {
		t.Errorf("Expected From Method GetStrategyMarketDataAssetListByPagination: to be empty but got this: %v", foundStrategyMarketDataAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalStrategyMarketDataAssetsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM strategy_market_data_assets").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalStrategyMarketDataAssetsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalStrategyMarketDataAssetsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalStrategyMarketDataAssetsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMinersForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM strategy_market_data_assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalStrategyMarketDataAssetsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalStrategyMarketDataAssetsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalStrategyMarketDataAssetsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
