package transactionasset

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
}
var DBColumnsInsertTransactionAssets = []string{
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
}

var TestData1 = TransactionAsset{
	TransactionID:         utils.Ptr[int](1),                      //1
	AssetID:               utils.Ptr[int](1),                      //2
	UUID:                  "01ef85e8-2c26-441e-8c7f-71d79518ad72", //3
	Name:                  "USDC",                                 //4
	AlternateName:         "USDC",                                 //5
	Description:           "",                                     //6
	Quantity:              utils.Ptr[float64](100),                //7
	QuantityUSD:           utils.Ptr[float64](10),                 //8
	MarketDataID:          utils.Ptr[int](1),                      //9
	ManualExchangeRateUSD: utils.Ptr[float64](1),                  //10
	CreatedBy:             "SYSTEM",                               //11
	CreatedAt:             utils.SampleCreatedAtTime,              //12
	UpdatedBy:             "SYSTEM",                               //13
	UpdatedAt:             utils.SampleCreatedAtTime,              //14

}

var TestData2 = TransactionAsset{
	TransactionID:         utils.Ptr[int](1),                      //1
	AssetID:               utils.Ptr[int](1),                      //2
	UUID:                  "4f0d5402-7a7c-402d-a7fc-c56a02b13e03", //3
	Name:                  "BTC",                                  //4
	AlternateName:         "BTC",                                  //5
	Description:           "",                                     //6
	Quantity:              utils.Ptr[float64](100),                //7
	QuantityUSD:           utils.Ptr[float64](100000),             //8
	MarketDataID:          utils.Ptr[int](2),                      //9
	ManualExchangeRateUSD: utils.Ptr[float64](1000),               //10
	CreatedBy:             "SYSTEM",                               //11
	CreatedAt:             utils.SampleCreatedAtTime,              //12
	UpdatedBy:             "SYSTEM",                               //13
	UpdatedAt:             utils.SampleCreatedAtTime,              //14
}
var TestAllData = []TransactionAsset{TestData1, TestData2}

func AddTransactionAssetToMockRows(mock pgxmock.PgxPoolIface, dataList []TransactionAsset) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.TransactionID,         //1
			data.AssetID,               //2
			data.UUID,                  //3
			data.Name,                  //4
			data.AlternateName,         //5
			data.Description,           //6
			data.Quantity,              //7
			data.QuantityUSD,           //8
			data.MarketDataID,          //9
			data.ManualExchangeRateUSD, //10
			data.CreatedBy,             //11
			data.CreatedAt,             //12
			data.UpdatedBy,             //13
			data.UpdatedAt,             //14
		)
	}
	return rows
}

func TestGetTransactionAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []TransactionAsset{targetData}
	transactionID := targetData.TransactionID
	assetID := targetData.AssetID
	mockRows := AddTransactionAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(*transactionID, *assetID).WillReturnRows(mockRows)
	foundTransactionAsset, err := GetTransactionAsset(mock, transactionID, assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionAsset", err)
	}
	if cmp.Equal(*foundTransactionAsset, targetData) == false {
		t.Errorf("Expected TransactionAsset From Method GetTransactionAsset: %v is different from actual %v", foundTransactionAsset, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := 999
	assetID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(transactionID, assetID).WillReturnRows(noRows)
	foundTransactionAsset, err := GetTransactionAsset(mock, &transactionID, &assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionAsset", err)
	}
	if foundTransactionAsset != nil {
		t.Errorf("Expected TransactionAsset From Method GetTransactionAsset: to be empty but got this: %v", foundTransactionAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -999
	assetID := -999
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(transactionID, assetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionAsset, err := GetTransactionAsset(mock, &transactionID, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionAsset", err)
	}
	if foundTransactionAsset != nil {
		t.Errorf("Expected TransactionAsset From Method GetTransactionAsset: to be empty but got this: %v", foundTransactionAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -999
	assetID := -999

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(transactionID, assetID).WillReturnRows(differentModelRows)
	foundTransactionAsset, err := GetTransactionAsset(mock, &transactionID, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionAsset", err)
	}
	if foundTransactionAsset != nil {
		t.Errorf("Expected TransactionAsset From Method GetTransactionAsset: to be empty but got this: %v", foundTransactionAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []TransactionAsset{targetData}
	uuid := targetData.UUID
	mockRows := AddTransactionAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(uuid).WillReturnRows(mockRows)
	foundTransactionAsset, err := GetTransactionAssetByUUID(mock, uuid)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionAssetByUUID", err)
	}
	if cmp.Equal(*foundTransactionAsset, targetData) == false {
		t.Errorf("Expected TransactionAsset From Method GetTransactionAssetByUUID: %v is different from actual %v", foundTransactionAsset, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetByUUIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuid := "invalid-uuid"
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(uuid).WillReturnRows(noRows)
	foundTransactionAsset, err := GetTransactionAssetByUUID(mock, uuid)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionAssetByUUID", err)
	}
	if foundTransactionAsset != nil {
		t.Errorf("Expected TransactionAsset From Method GetTransactionAssetByUUID: to be empty but got this: %v", foundTransactionAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetByUUIDForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuid := "invalid-uuid"
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(uuid).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionAsset, err := GetTransactionAssetByUUID(mock, uuid)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionAssetByUUID", err)
	}
	if foundTransactionAsset != nil {
		t.Errorf("Expected TransactionAsset From Method GetTransactionAssetByUUID: to be empty but got this: %v", foundTransactionAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetByUUIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuid := "invalid-uuid"

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(uuid).WillReturnRows(differentModelRows)
	foundTransactionAsset, err := GetTransactionAssetByUUID(mock, uuid)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionAssetByUUID", err)
	}
	if foundTransactionAsset != nil {
		t.Errorf("Expected TransactionAsset From Method GetTransactionAssetByUUID: to be empty but got this: %v", foundTransactionAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mockRows := AddTransactionAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(pq.Array(uuids)).WillReturnRows(mockRows)
	foundTransactionAssetList, err := GetTransactionAssetsByUUIDs(mock, uuids)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionAssetsByUUIDs", err)
	}
	for i, sourceTransactionAsset := range dataList {
		if cmp.Equal(sourceTransactionAsset, foundTransactionAssetList[i]) == false {
			t.Errorf("Expected TransactionAsset From Method GetTransactionAssetsByUUIDs: %v is different from actual %v", sourceTransactionAsset, foundTransactionAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuids := []string{"test-uuid", "uuid-2"}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(pq.Array(uuids)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionAssetList, err := GetTransactionAssetsByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionAssetsByUUIDs", err)
	}
	if len(foundTransactionAssetList) != 0 {
		t.Errorf("Expected From Method GetTransactionAssetsByUUIDs: to be empty but got this: %v", foundTransactionAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetsByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuids := []string{"test-uuid", "uuid-2"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WithArgs(pq.Array(uuids)).WillReturnRows(differentModelRows)
	foundTransactionAsset, err := GetTransactionAssetsByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionAssetsByUUIDs", err)
	}
	if foundTransactionAsset != nil {
		t.Errorf("Expected TransactionAsset From Method GetTransactionAssetsByUUIDs: to be empty but got this: %v", foundTransactionAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	transactionID := targetData.TransactionID
	assetID := targetData.AssetID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_assets").WithArgs(*transactionID, *assetID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveTransactionAsset(mock, transactionID, assetID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveTransactionAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	assetID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveTransactionAsset(mock, &transactionID, &assetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	assetID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_assets").WithArgs(transactionID, assetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveTransactionAsset(mock, &transactionID, &assetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionAssetByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	uuid := targetData.UUID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_assets").WithArgs(uuid).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveTransactionAssetByUUID(mock, uuid)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveTransactionAssetByUUID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionAssetByUUIDOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuid := "test-uuid"
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveTransactionAssetByUUID(mock, uuid)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionAssetByUUIDOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuid := "test-uuid"
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_assets").WithArgs(uuid).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveTransactionAssetByUUID(mock, uuid)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddTransactionAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WillReturnRows(mockRows)
	foundTransactionAssetList, err := GetTransactionAssets(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionAssets", err)
	}
	for i, sourceTransactionAsset := range dataList {
		if cmp.Equal(sourceTransactionAsset, foundTransactionAssetList[i]) == false {
			t.Errorf("Expected TransactionAsset From Method GetTransactionAssets: %v is different from actual %v", sourceTransactionAsset, foundTransactionAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionAssetList, err := GetTransactionAssets(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionAssets", err)
	}
	if len(foundTransactionAssetList) != 0 {
		t.Errorf("Expected From Method GetTransactionAssets: to be empty but got this: %v", foundTransactionAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WillReturnRows(differentModelRows)
	foundTransactionAsset, err := GetTransactionAssets(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionAssets", err)
	}
	if foundTransactionAsset != nil {
		t.Errorf("Expected TransactionAsset From Method GetTransactionAssets: to be empty but got this: %v", foundTransactionAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTransactionAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_assets").WithArgs(
		targetData.Name,                  //1
		targetData.AlternateName,         //2
		targetData.Description,           //3
		targetData.Quantity,              //4
		targetData.QuantityUSD,           //5
		targetData.MarketDataID,          //6
		targetData.ManualExchangeRateUSD, //7
		targetData.UpdatedBy,             //8
		targetData.TransactionID,         //9
		targetData.AssetID,               //10
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateTransactionAsset(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateTransactionAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionAssetOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.AssetID = nil
	err = UpdateTransactionAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.TransactionID = utils.Ptr[int](1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateTransactionAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.AssetID = utils.Ptr[int](1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_assets").WithArgs(
		targetData.Name,                  //1
		targetData.AlternateName,         //2
		targetData.Description,           //3
		targetData.Quantity,              //4
		targetData.QuantityUSD,           //5
		targetData.MarketDataID,          //6
		targetData.ManualExchangeRateUSD, //7
		targetData.UpdatedBy,             //8
		targetData.TransactionID,         //9
		targetData.AssetID,               //10
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateTransactionAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTransactionAssetByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_assets").WithArgs(
		targetData.Name,                  //1
		targetData.AlternateName,         //2
		targetData.Description,           //3
		targetData.Quantity,              //4
		targetData.QuantityUSD,           //5
		targetData.MarketDataID,          //6
		targetData.ManualExchangeRateUSD, //7
		targetData.UpdatedBy,             //8
		targetData.UUID,                  //9
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateTransactionAssetByUUID(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateTransactionAssetByUUID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionAssetByUUIDOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.AssetID = nil
	err = UpdateTransactionAssetByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionAssetByUUIDOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.TransactionID = utils.Ptr[int](1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateTransactionAssetByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionAssetByUUIDOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.AssetID = utils.Ptr[int](1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_assets").WithArgs(
		targetData.Name,                  //1
		targetData.AlternateName,         //2
		targetData.Description,           //3
		targetData.Quantity,              //4
		targetData.QuantityUSD,           //5
		targetData.MarketDataID,          //6
		targetData.ManualExchangeRateUSD, //7
		targetData.UpdatedBy,             //8
		targetData.UUID,                  //9
	).WillReturnError(fmt.Errorf("Random error"))

	mock.ExpectRollback()
	err = UpdateTransactionAssetByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transaction_assets").WithArgs(
		targetData.TransactionID,         //1
		targetData.AssetID,               //2
		targetData.UUID,                  //3
		targetData.Name,                  //4
		targetData.AlternateName,         //5
		targetData.Description,           //6
		targetData.Quantity,              //7
		targetData.QuantityUSD,           //8
		targetData.MarketDataID,          //9
		targetData.ManualExchangeRateUSD, //10
		targetData.CreatedBy,             //11
	).WillReturnRows(pgxmock.NewRows([]string{"transaction_id", "asset_id"}).AddRow(1, 1))
	mock.ExpectCommit()
	transactionID, assetID, err := InsertTransactionAsset(mock, &targetData)
	if transactionID < 0 {
		t.Fatalf("transactionID should not be negative ID: %d", transactionID)
	}
	if assetID < 0 {
		t.Fatalf("assetID should not be negative ID: %d", assetID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertTransactionAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.AssetID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertTransactionAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transaction_assets").WithArgs(
		targetData.TransactionID,         //1
		targetData.AssetID,               //2
		targetData.UUID,                  //3
		targetData.Name,                  //4
		targetData.AlternateName,         //5
		targetData.Description,           //6
		targetData.Quantity,              //7
		targetData.QuantityUSD,           //8
		targetData.MarketDataID,          //9
		targetData.ManualExchangeRateUSD, //10
		targetData.CreatedBy,             //11
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	transactionID, assetID, err := InsertTransactionAsset(mock, &targetData)
	if transactionID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error transactionID: %d", transactionID)
	}
	if assetID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error assetID: %d", assetID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionAssetOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transaction_assets").WithArgs(
		targetData.TransactionID,         //1
		targetData.AssetID,               //2
		targetData.UUID,                  //3
		targetData.Name,                  //4
		targetData.AlternateName,         //5
		targetData.Description,           //6
		targetData.Quantity,              //7
		targetData.QuantityUSD,           //8
		targetData.MarketDataID,          //9
		targetData.ManualExchangeRateUSD, //10
		targetData.CreatedBy,             //11
	).WillReturnRows(pgxmock.NewRows([]string{"transaction_id", "asset_id"}).AddRow(1, 1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	transactionID, assetID, err := InsertTransactionAsset(mock, &targetData)
	if transactionID >= 0 {
		t.Fatalf("Expecting -1 for transactionID because of error transactionID: %d", transactionID)
	}
	if assetID >= 0 {
		t.Fatalf("Expecting -1 for assetID because of error assetID: %d", assetID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"transaction_assets"}, DBColumnsInsertTransactionAssets)
	targetData := TestAllData
	err = InsertTransactionAssets(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertTransactionAssets", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionAssetsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"transaction_assets"}, DBColumnsInsertTransactionAssets).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertTransactionAssets(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTransactionAssetToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"asset_id = 1", "market_data_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WillReturnRows(mockRows)
	foundTransactionAssetList, err := GetTransactionAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionAssetListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundTransactionAssetList[i]) == false {
			t.Errorf("Expected sourceData From Method GetTransactionAssetListByPagination: %v is different from actual %v", sourceData, foundTransactionAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"market_data_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionAssetList, err := GetTransactionAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionAssetListByPagination", err)
	}
	if len(foundTransactionAssetList) != 0 {
		t.Errorf("Expected From Method GetTransactionAssetListByPagination: to be empty but got this: %v", foundTransactionAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionAssetListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"market_data_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_assets").WillReturnRows(differentModelRows)
	foundTransactionAssetList, err := GetTransactionAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionAssetListByPagination", err)
	}
	if len(foundTransactionAssetList) != 0 {
		t.Errorf("Expected From Method GetTransactionAssetListByPagination: to be empty but got this: %v", foundTransactionAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalTransactionAssetsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM transaction_assets").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalTransactionAssetsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalTransactionAssetsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalTransactionAssetsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM transaction_assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalTransactionAssetsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalTransactionAssetsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalTransactionAssetsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
