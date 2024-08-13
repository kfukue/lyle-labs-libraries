package asset

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
)

func TestGetAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	dataList := []Asset{targetData}
	mockRows := AddAssetToMockRows(mock, dataList)
	assetID := 1
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE id = ?").WithArgs(assetID).WillReturnRows(mockRows)
	foundAsset, err := GetAsset(mock, &assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAsset", err)
	}
	if cmp.Equal(*foundAsset, targetData) == false {
		t.Errorf("Expected Asset From Method GetAsset: %v is different from actual %v", foundAsset, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	noRows := pgxmock.NewRows(DBColumns)
	assetID := 1
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE id = ?").WithArgs(assetID).WillReturnRows(noRows)
	foundAsset, err := GetAsset(mock, &assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAsset", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAsset: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	invalidID := -1
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE id = ?").WithArgs(invalidID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAsset, err := GetAsset(mock, &invalidID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAsset", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAsset: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(assetID).WillReturnRows(differentModelRows)
	foundAsset, err := GetAsset(mock, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAsset", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected foundAsset From Method GetAsset: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByTicker(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Asset{targetData}
	testTicker := targetData.Ticker
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE ticker = ?").WithArgs(testTicker).WillReturnRows(mockRows)
	foundAsset, err := GetAssetByTicker(mock, testTicker)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetByTicker", err)
	}
	if cmp.Equal(*foundAsset, targetData) == false {
		t.Errorf("Expected Asset From Method GetAssetByTicker: %v is different from actual %v", foundAsset, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByTickerForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testTicker := "Fake-Ticker"
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE ticker = ?").WithArgs(testTicker).WillReturnRows(noRows)
	foundAsset, err := GetAssetByTicker(mock, testTicker)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetByTicker", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAssetByTicker: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByTickerForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testTicker := "Fake-Ticker"
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE ticker = ?").WithArgs(testTicker).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAsset, err := GetAssetByTicker(mock, testTicker)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetByTicker", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAssetByTicker: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByTickerForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testTicker := "Fake-Ticker"

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(testTicker).WillReturnRows(differentModelRows)
	foundAsset, err := GetAssetByTicker(mock, testTicker)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetByTicker", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected foundAsset From Method GetAssetByTicker: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByContractAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Asset{targetData}
	testContractAddress := targetData.ContractAddress
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE contract_address = ?").WithArgs(testContractAddress).WillReturnRows(mockRows)
	foundAsset, err := GetAssetByContractAddress(mock, testContractAddress)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetByContractAddress", err)
	}
	if cmp.Equal(*foundAsset, targetData) == false {
		t.Errorf("Expected Asset From Method GetAssetByContractAddress: %v is different from actual %v", foundAsset, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByContractAddressForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testContractAddress := "Fake-ContractAddress"
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE contract_address = ?").WithArgs(testContractAddress).WillReturnRows(noRows)
	foundAsset, err := GetAssetByContractAddress(mock, testContractAddress)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetByContractAddress", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAssetByContractAddress: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByContractAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testContractAddress := "Fake-ContractAddress"
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE contract_address = ?").WithArgs(testContractAddress).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAsset, err := GetAssetByContractAddress(mock, testContractAddress)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetByContractAddress", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAssetByContractAddress: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByContractAddressForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testContractAddress := "Fake-Ticker"

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(testContractAddress).WillReturnRows(differentModelRows)
	foundAsset, err := GetAssetByContractAddress(mock, testContractAddress)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetByContractAddress", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected foundAsset From Method GetAssetByContractAddress: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByCusip(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Asset{targetData}
	testCusip := targetData.Cusip
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE cusip = ?").WithArgs(testCusip).WillReturnRows(mockRows)
	foundAsset, err := GetAssetByCusip(mock, testCusip)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetByCusip", err)
	}
	if cmp.Equal(*foundAsset, targetData) == false {
		t.Errorf("Expected Asset From Method GetAssetByCusip: %v is different from actual %v", foundAsset, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByCusipForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testCusip := "Fake-Cusip"
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE cusip = ?").WithArgs(testCusip).WillReturnRows(noRows)
	foundAsset, err := GetAssetByCusip(mock, testCusip)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetByCusip", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAssetByCusip: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByCusipForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testCusip := "Fake-Cusip"
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE cusip = ?").WithArgs(testCusip).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAsset, err := GetAssetByCusip(mock, testCusip)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetByCusip", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAssetByCusip: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByCusipForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testCusip := "Fake-Ticker"

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(testCusip).WillReturnRows(differentModelRows)
	foundAsset, err := GetAssetByCusip(mock, testCusip)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetByCusip", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected foundAsset From Method GetAssetByCusip: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestGetAssetByBaseAndQuoteID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Asset{targetData}
	baseAssetID := targetData.BaseAssetID
	quoteAssetID := targetData.QuoteAssetID
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*baseAssetID, *quoteAssetID).WillReturnRows(mockRows)
	foundAsset, err := GetAssetByBaseAndQuoteID(mock, baseAssetID, quoteAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetByBaseAndQuoteID", err)
	}
	if cmp.Equal(*foundAsset, targetData) == false {
		t.Errorf("Expected Asset From Method GetAssetByBaseAndQuoteID: %v is different from actual %v", foundAsset, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByBaseAndQuoteIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := utils.Ptr[int](1331)
	quoteAssetID := utils.Ptr[int](1111)
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*baseAssetID, *quoteAssetID).WillReturnRows(noRows)
	foundAsset, err := GetAssetByBaseAndQuoteID(mock, baseAssetID, quoteAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetByBaseAndQuoteID", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAssetByBaseAndQuoteID: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByBaseAndQuoteIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := utils.Ptr[int](-1)
	quoteAssetID := utils.Ptr[int](-1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*baseAssetID, *quoteAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAsset, err := GetAssetByBaseAndQuoteID(mock, baseAssetID, quoteAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetByBaseAndQuoteID", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAssetByBaseAndQuoteID: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByBaseAndQuoteIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := utils.Ptr[int](-1)
	quoteAssetID := utils.Ptr[int](-1)

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*baseAssetID, *quoteAssetID).WillReturnRows(differentModelRows)
	foundAsset, err := GetAssetByBaseAndQuoteID(mock, baseAssetID, quoteAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetByBaseAndQuoteID", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected foundAsset From Method GetAssetByBaseAndQuoteID: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethImportAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{TestData1, TestData2}
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE import_geth = TRUE").WillReturnRows(mockRows)
	foundAssets, err := GetGethImportAssets(mock)

	if err != nil {
		t.Fatalf("an error '%s' in GetGethImportAssets", err)
	}
	testAssets := TestAllData
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetGethImportAssets: %v is different from actual %v", foundAsset, testAssets[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethImportAssetsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE import_geth = TRUE").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetGethImportAssets(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethImportAssets", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetGethImportAssets: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethImportAssetsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WillReturnRows(differentModelRows)
	foundAssets, err := GetGethImportAssets(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethImportAssets", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected foundAssets From Method GetGethImportAssets: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM assets WHERE id = ?").WithArgs(*targetData.ID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveAsset(mock, targetData.ID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveAsset(mock, &taxID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	invalidID := utils.Ptr[int](-1)
	targetData.ID = invalidID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM assets WHERE id = ?").WithArgs(*invalidID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveAsset(mock, invalidID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCurrentTradingAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{TestData1, TestData2}
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM public.get_current_assets").WillReturnRows(mockRows)
	foundAssets, err := GetCurrentTradingAssets(mock)

	if err != nil {
		t.Fatalf("an error '%s' in GetCurrentTradingAssets", err)
	}
	testAssets := TestAllData
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetCurrentTradingAssets: %v is different from actual %v", foundAsset, testAssets[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCurrentTradingAssetsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM public.get_current_assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetCurrentTradingAssets(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetCurrentTradingAssets", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetCurrentTradingAssets: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCurrentTradingAssetsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM public.get_current_assets").WillReturnRows(differentModelRows)
	foundAssets, err := GetCurrentTradingAssets(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetCurrentTradingAssets", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected foundAssets From Method GetCurrentTradingAssets: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCryptoAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{TestData1, TestData2}
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM assets where asset_type_id = 1").WillReturnRows(mockRows)
	foundAssets, err := GetCryptoAssets(mock)

	if err != nil {
		t.Fatalf("an error '%s' in GetCryptoAssets", err)
	}
	testAssets := TestAllData
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetCryptoAssets: %v is different from actual %v", foundAsset, testAssets[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCryptoAssetsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM assets where asset_type_id = 1").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetCryptoAssets(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetCryptoAssets", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetCryptoAssets: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCryptoAssetsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) assets where asset_type_id = 1").WillReturnRows(differentModelRows)
	foundAssets, err := GetCryptoAssets(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetCryptoAssets", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected foundAssets From Method GetCryptoAssets: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

// asset with sources
func TestGetAssetsByAssetTypeAndSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllDataAssetWithSources
	mockRows := AddAssetWithSourcesToMockRows(mock, dataList)
	assetTypeID := dataAssetWithSources1.AssetTypeID
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*assetTypeID, *sourceID).WillReturnRows(mockRows)
	foundAssets, err := GetAssetsByAssetTypeAndSource(mock, assetTypeID, sourceID, excludeIgnoreMarketData)

	if err != nil {
		t.Fatalf("an error '%s' in GetAssetsByAssetTypeAndSource", err)
	}
	testAssets := TestAllDataAssetWithSources
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetAssetsByAssetTypeAndSource: %v is different from actual %v", foundAsset, testAssets[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetsByAssetTypeAndSourceForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetTypeID := dataAssetWithSources1.AssetTypeID
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := false
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*assetTypeID, *sourceID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetAssetsByAssetTypeAndSource(mock, assetTypeID, sourceID, excludeIgnoreMarketData)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetsByAssetTypeAndSource", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetAssetsByAssetTypeAndSource: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetsByAssetTypeAndSourceForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetTypeID := dataAssetWithSources1.AssetTypeID
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := false
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*assetTypeID, *sourceID).WillReturnRows(differentModelRows)
	foundAssets, err := GetAssetsByAssetTypeAndSource(mock, assetTypeID, sourceID, excludeIgnoreMarketData)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetsByAssetTypeAndSource", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected foundAssets From Method GetAssetsByAssetTypeAndSource: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCryptoAssetsBySourceId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllDataAssetWithSources
	mockRows := AddAssetWithSourcesToMockRows(mock, dataList)
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*sourceID).WillReturnRows(mockRows)
	foundAssets, err := GetCryptoAssetsBySourceId(mock, sourceID, excludeIgnoreMarketData)

	if err != nil {
		t.Fatalf("an error '%s' in GetCryptoAssetsBySourceId", err)
	}
	testAssets := TestAllDataAssetWithSources
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetCryptoAssetsBySourceId: %v is different from actual %v", foundAsset, testAssets[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCryptoAssetsBySourceIdForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := false
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*sourceID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetCryptoAssetsBySourceId(mock, sourceID, excludeIgnoreMarketData)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetCryptoAssetsBySourceId", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetCryptoAssetsBySourceId: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCryptoAssetsBySourceIdForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := false
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*sourceID).WillReturnRows(differentModelRows)
	foundAssets, err := GetCryptoAssetsBySourceId(mock, sourceID, excludeIgnoreMarketData)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetCryptoAssetsBySourceId", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected foundAssets From Method GetCryptoAssetsBySourceId: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCryptoAssetsBySourceID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllDataAssetWithSources
	mockRows := AddAssetWithSourcesToMockRows(mock, dataList)
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*sourceID).WillReturnRows(mockRows)
	foundAssets, err := GetCryptoAssetsBySourceID(mock, sourceID, excludeIgnoreMarketData)

	if err != nil {
		t.Fatalf("an error '%s' in GetCryptoAssetsBySourceID", err)
	}
	testAssets := TestAllData
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetCryptoAssetsBySourceID: %v is different from actual %v", foundAsset, testAssets[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetCryptoAssetsBySourceIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := false
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*sourceID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetCryptoAssetsBySourceID(mock, sourceID, excludeIgnoreMarketData)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetCryptoAssetsBySourceID", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetCryptoAssetsBySourceID: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetWithSourceByAssetIdAndSourceID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []AssetWithSources{dataAssetWithSources1}
	mockRows := AddAssetWithSourcesToMockRows(mock, dataList)
	assetID := dataAssetWithSources1.Asset.ID
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*assetID, *sourceID).WillReturnRows(mockRows)
	foundAsset, err := GetAssetWithSourceByAssetIdAndSourceID(mock, assetID, sourceID, excludeIgnoreMarketData)

	if err != nil {
		t.Fatalf("an error '%s' in GetAssetWithSourceByAssetIdAndSourceID", err)
	}
	if cmp.Equal(*foundAsset, dataAssetWithSources1) == false {
		t.Errorf("Expected Asset From Method GetAssetWithSourceByAssetIdAndSourceID: %v is different from actual %v", foundAsset, dataAssetWithSources1)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetWithSourceByAssetIdAndSourceIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	noRows := pgxmock.NewRows(DBColumns)
	assetID := dataAssetWithSources1.Asset.ID
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*assetID, *sourceID).WillReturnRows(noRows)
	foundAsset, err := GetAssetWithSourceByAssetIdAndSourceID(mock, assetID, sourceID, excludeIgnoreMarketData)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetWithSourceByAssetIdAndSourceID", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAssetWithSourceByAssetIdAndSourceID: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetWithSourceByAssetIdAndSourceIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := utils.Ptr[int](-1)
	sourceID := utils.Ptr[int](-1)
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*assetID, *sourceID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAsset, err := GetAssetWithSourceByAssetIdAndSourceID(mock, assetID, sourceID, excludeIgnoreMarketData)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetWithSourceByAssetIdAndSourceID", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected Asset From Method GetAssetWithSourceByAssetIdAndSourceID: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetWithSourceByAssetIdAndSourceIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := utils.Ptr[int](-1)
	sourceID := utils.Ptr[int](-1)
	excludeIgnoreMarketData := true
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*assetID, *sourceID).WillReturnRows(differentModelRows)
	foundAsset, err := GetAssetWithSourceByAssetIdAndSourceID(mock, assetID, sourceID, excludeIgnoreMarketData)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetWithSourceByAssetIdAndSourceID", err)
	}
	if foundAsset != nil {
		t.Errorf("Expected foundAsset From Method GetAssetWithSourceByAssetIdAndSourceID: to be empty but got this: %v", foundAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetWithSourceByAssetIdsAndSourceID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllDataAssetWithSources
	assetIDs := []int{*dataAssetWithSources1.Asset.ID, *dataAssetWithSources2.Asset.ID}
	mockRows := AddAssetWithSourcesToMockRows(mock, dataList)
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(pq.Array(assetIDs), *sourceID).WillReturnRows(mockRows)
	foundAssets, err := GetAssetWithSourceByAssetIdsAndSourceID(mock, assetIDs, sourceID, excludeIgnoreMarketData)

	if err != nil {
		t.Fatalf("an error '%s' in GetAssetWithSourceByAssetIdsAndSourceID", err)
	}
	testAssets := TestAllDataAssetWithSources
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetAssetWithSourceByAssetIdsAndSourceID: %v is different from actual %v", foundAsset, testAssets[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetWithSourceByAssetIdsAndSourceIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetIDs := []int{-1, -2}
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := false
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(pq.Array(assetIDs), *sourceID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetAssetWithSourceByAssetIdsAndSourceID(mock, assetIDs, sourceID, excludeIgnoreMarketData)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetWithSourceByAssetIdsAndSourceID", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetAssetWithSourceByAssetIdsAndSourceID: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetWithSourceByAssetIdsAndSourceIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetIDs := []int{-1, -2}
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := false
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(pq.Array(assetIDs), *sourceID).WillReturnRows(differentModelRows)
	foundAssets, err := GetAssetWithSourceByAssetIdsAndSourceID(mock, assetIDs, sourceID, excludeIgnoreMarketData)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetWithSourceByAssetIdsAndSourceID", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected foundAssets From Method GetAssetWithSourceByAssetIdsAndSourceID: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{TestData1, TestData2}
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery(fmt.Sprintf("^SELECT (.+) FROM assets WHERE")).WillReturnRows(mockRows)
	ids := make([]int, 0)
	ids = append(ids, *TestData1.ID)
	ids = append(ids, *TestData2.ID)
	foundAssets, err := GetAssetList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetList", err)
	}
	testAssets := TestAllData
	for i, foundAsset := range foundAssets {
		// t.Logf("i: %d,  \n mock : %v , \n in memory : %v", i, foundAsset, testAssets[i])

		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetAssetList: %v is different from actual %v", foundAsset, testAssets[i])
		}

	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetIDs := []int{-1, -2}
	mock.ExpectQuery("^SELECT (.+) FROM assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetAssetList(mock, assetIDs)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetList", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetAssetList: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetIDs := []int{-1, -2}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WillReturnRows(differentModelRows)
	foundAssets, err := GetAssetList(mock, assetIDs)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetList", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected foundAssets From Method GetAssetList: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetsByChainId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{TestData1}
	chainID := TestData1.ChainID
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE chain_id = ?").WithArgs(*chainID).WillReturnRows(mockRows)
	foundAssets, err := GetAssetsByChainId(mock, chainID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetsByChainId", err)
	}
	testAssets := dataList
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetAssetsByChainId: %v is different from actual %v", foundAsset, testAssets[i])
		}

	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetsByChainIdForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainID := TestData1.ChainID
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE chain_id = ?").WithArgs(*chainID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetAssetsByChainId(mock, chainID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetsByChainId", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetAssetsByChainId: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetsByChainIdForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	chainID := TestData1.ChainID
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE chain_id = ?").WithArgs(*chainID).WillReturnRows(differentModelRows)
	foundAssets, err := GetAssetsByChainId(mock, chainID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetsByChainId", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected foundAssets From Method GetAssetsByChainId: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{TestData1}
	mockRows := AddAssetToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"asset_type_id = 1", "import_geth = TRUE"}
	mock.ExpectQuery("^SELECT (.+) FROM assets").WillReturnRows(mockRows)
	foundAssets, err := GetAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetListByPagination", err)
	}
	testAssets := dataList
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetAssetListByPagination: %v is different from actual %v", foundAsset, testAssets[i])
		}

	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"asset_type_id = 1", "import_geth = TRUE"}
	mock.ExpectQuery("^SELECT (.+) FROM assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetListByPagination", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetAssetListByPagination: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"asset_type_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM assets").WillReturnRows(differentModelRows)
	foundAssets, err := GetAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetListByPagination", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected From Method GetAssetListByPagination: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalAssetCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfAssetsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM assets").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfAssetsExpected))
	numOfAssets, err := GetTotalAssetCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalAssetCount", err)
	}
	if *numOfAssets != numOfAssetsExpected {
		t.Errorf("Expected Asset From Method GetTotalAssetCount: %d is different from actual %d", numOfAssetsExpected, *numOfAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalAssetCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	totalCount, err := GetTotalAssetCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalAssetCount", err)
	}
	if totalCount != nil {
		t.Errorf("Expected From Method GetTotalAssetCount: to be 0 but got this: %d", totalCount)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDefaultQuoteAssetListBySourceID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllDataAssetWithSources
	sourceID := 1
	mockRows := AddAssetWithSourcesToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM get_default_quotes").WithArgs(sourceID).WillReturnRows(mockRows)
	foundAssets, err := GetDefaultQuoteAssetListBySourceID(mock, &sourceID)
	if err != nil {
		t.Fatalf("an error '%s' in GetDefaultQuoteAssetListBySourceID", err)
	}
	testAssets := TestAllDataAssetWithSources
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetDefaultQuoteAssetListBySourceID: %v is different from actual %v", foundAsset, testAssets[i])
		}

	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDefaultQuoteAssetListBySourceIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	mock.ExpectQuery("^SELECT (.+) FROM get_default_quotes").WithArgs(sourceID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetDefaultQuoteAssetListBySourceID(mock, &sourceID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDefaultQuoteAssetListBySourceID", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetDefaultQuoteAssetListBySourceID: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDefaultQuoteAssetListBySourceIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	sourceID := -1
	mock.ExpectQuery("^SELECT (.+) FROM get_default_quotes").WithArgs(sourceID).WillReturnRows(differentModelRows)
	foundAssets, err := GetDefaultQuoteAssetListBySourceID(mock, &sourceID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDefaultQuoteAssetListBySourceID", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected foundAssets From Method GetDefaultQuoteAssetListBySourceID: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDefaultQuoteAssetList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mockRows := AddAssetToMockRows(mock, TestAllData)
	mock.ExpectQuery("^SELECT (.+) FROM get_default_quotes").WillReturnRows(mockRows)
	foundAssets, err := GetDefaultQuoteAssetList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetDefaultQuoteAssetList", err)
	}
	testAssets := TestAllData
	for i, foundAsset := range foundAssets {
		if cmp.Equal(foundAsset, testAssets[i]) == false {
			t.Errorf("Expected Asset From Method GetDefaultQuoteAssetList: %v is different from actual %v", foundAsset, testAssets[i])
		}

	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDefaultQuoteAssetListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM get_default_quotes").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetDefaultQuoteAssetList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDefaultQuoteAssetList", err)
	}
	if len(foundAssets) != 0 {
		t.Errorf("Expected From Method GetDefaultQuoteAssetList: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDefaultQuoteAssetListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM get_default_quotes").WillReturnRows(differentModelRows)
	foundAssets, err := GetDefaultQuoteAssetList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDefaultQuoteAssetList", err)
	}
	if foundAssets != nil {
		t.Errorf("Expected foundAssets From Method GetDefaultQuoteAssetList: to be empty but got this: %v", foundAssets)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE assets").WithArgs(
		targetData.Name,                //1
		targetData.AlternateName,       //2
		targetData.Cusip,               //3
		targetData.Ticker,              //4
		targetData.BaseAssetID,         //5
		targetData.QuoteAssetID,        //6
		targetData.Description,         //7
		targetData.AssetTypeID,         //8
		targetData.UpdatedBy,           //9
		targetData.ChainID,             //10
		targetData.CategoryID,          //11
		targetData.SubCategoryID,       //12
		targetData.IsDefaultQuote,      //13
		targetData.IgnoreMarketData,    //14
		targetData.Decimals,            //15
		targetData.ContractAddress,     //16
		targetData.StartingBlockNumber, //17
		targetData.ImportGeth,          //18
		targetData.ImportGethInitial,   //19
		targetData.ID,                  //20
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateAsset(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAssetOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAssetOnFailure(t *testing.T) {
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
	mock.ExpectExec("^UPDATE assets").WithArgs(
		targetData.Name,                //1
		targetData.AlternateName,       //2
		targetData.Cusip,               //3
		targetData.Ticker,              //4
		targetData.BaseAssetID,         //5
		targetData.QuoteAssetID,        //6
		targetData.Description,         //7
		targetData.AssetTypeID,         //8
		targetData.UpdatedBy,           //9
		targetData.ChainID,             //10
		targetData.CategoryID,          //11
		targetData.SubCategoryID,       //12
		targetData.IsDefaultQuote,      //13
		targetData.IgnoreMarketData,    //14
		targetData.Decimals,            //15
		targetData.ContractAddress,     //16
		targetData.StartingBlockNumber, //17
		targetData.ImportGeth,          //18
		targetData.ImportGethInitial,   //19
		targetData.ID,                  //20
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO assets").WithArgs(
		targetData.Name,                //1
		targetData.AlternateName,       //2
		targetData.Cusip,               //3
		targetData.Ticker,              //4
		targetData.BaseAssetID,         //5
		targetData.QuoteAssetID,        //6
		targetData.Description,         //7
		targetData.AssetTypeID,         //8
		targetData.CreatedBy,           //9
		targetData.ChainID,             //10
		targetData.CategoryID,          //11
		targetData.SubCategoryID,       //12
		targetData.IsDefaultQuote,      //13
		targetData.IgnoreMarketData,    //14
		targetData.Decimals,            //15
		targetData.ContractAddress,     //16
		targetData.StartingBlockNumber, //17
		targetData.ImportGeth,          //18
		targetData.ImportGethInitial,   //19
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	newID, err := InsertAsset(mock, &targetData)
	if newID < 0 {
		t.Fatalf("ID should not be negative ID: %d", newID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO assets").WithArgs(
		targetData.Name,                //1
		targetData.AlternateName,       //2
		targetData.Cusip,               //3
		targetData.Ticker,              //4
		targetData.BaseAssetID,         //5
		targetData.QuoteAssetID,        //6
		targetData.Description,         //7
		targetData.AssetTypeID,         //8
		targetData.CreatedBy,           //9
		targetData.ChainID,             //10
		targetData.CategoryID,          //11
		targetData.SubCategoryID,       //12
		targetData.IsDefaultQuote,      //13
		targetData.IgnoreMarketData,    //14
		targetData.Decimals,            //15
		targetData.ContractAddress,     //16
		targetData.StartingBlockNumber, //17
		targetData.ImportGeth,          //18
		targetData.ImportGethInitial,   //19
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	newID, err := InsertAsset(mock, &targetData)
	if newID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error ID: %d", newID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO assets").WithArgs(
		targetData.Name,                //1
		targetData.AlternateName,       //2
		targetData.Cusip,               //3
		targetData.Ticker,              //4
		targetData.BaseAssetID,         //5
		targetData.QuoteAssetID,        //6
		targetData.Description,         //7
		targetData.AssetTypeID,         //8
		targetData.CreatedBy,           //9
		targetData.ChainID,             //10
		targetData.CategoryID,          //11
		targetData.SubCategoryID,       //12
		targetData.IsDefaultQuote,      //13
		targetData.IgnoreMarketData,    //14
		targetData.Decimals,            //15
		targetData.ContractAddress,     //16
		targetData.StartingBlockNumber, //17
		targetData.ImportGeth,          //18
		targetData.ImportGethInitial,   //19
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	newID, err := InsertAsset(mock, &targetData)
	if newID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error ID: %d", newID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"assets"}, DBColumnsInsertAssets)
	targetData := TestAllData
	err = InsertAssets(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertAssets", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"assets"}, DBColumnsInsertAssets).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertAssets(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
