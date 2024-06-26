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

var columns = []string{
	"id",
	"uuid",
	"name",
	"alternate_name",
	"cusip",
	"ticker",
	"base_asset_id",
	"quote_asset_id",
	"description",
	"asset_type_id",
	"created_by",
	"created_at",
	"updated_by",
	"updated_at",
	"chain_id",
	"category_id",
	"sub_category_id",
	"is_default_quote",
	"ignore_market_data",
	"decimals",
	"contract_address",
	"starting_block_number",
	"import_geth",
	"import_geth_initial",
}

var columnsAssetWithSources = []string{
	"id",
	"uuid",
	"name",
	"alternate_name",
	"cusip",
	"ticker",
	"base_asset_id",
	"quote_asset_id",
	"description",
	"asset_type_id",
	"created_by",
	"created_at",
	"updated_by",
	"updated_at",
	"chain_id",
	"category_id",
	"sub_category_id",
	"is_default_quote",
	"ignore_market_data",
	"decimals",
	"contract_address",
	"starting_block_number",
	"import_geth",
	"import_geth_initial",
	"assetSources.source_id",
	"assetSources.source_identifier",
}

var data1 = Asset{
	ID:                  utils.Ptr[int](1),
	UUID:                "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:                "ETHER",
	AlternateName:       "Ether",
	Cusip:               "",
	Ticker:              "ETH",
	BaseAssetID:         utils.Ptr[int](1),
	QuoteAssetID:        utils.Ptr[int](2),
	Description:         "",
	AssetTypeID:         utils.Ptr[int](1),
	CreatedBy:           "SYSTEM",
	CreatedAt:           utils.SampleCreatedAtTime,
	UpdatedBy:           "SYSTEM",
	UpdatedAt:           utils.SampleCreatedAtTime,
	ChainID:             utils.Ptr[int](2),
	CategoryID:          utils.Ptr[int](27),
	SubCategoryID:       utils.Ptr[int](10),
	IsDefaultQuote:      utils.Ptr[bool](true),
	IgnoreMarketData:    utils.Ptr[bool](false),
	Decimals:            utils.Ptr[int](1),
	ContractAddress:     "SYSTEM",
	StartingBlockNumber: utils.Ptr[uint64](1),
	ImportGeth:          nil,
	ImportGethInitial:   nil,
}

var data2 = Asset{
	ID:                  utils.Ptr[int](2),
	UUID:                "880607ab-2833-4ad7-a231-b983a61c7b334",
	Name:                "BTC",
	AlternateName:       "Bitcoin",
	Cusip:               "",
	Ticker:              "BTC",
	BaseAssetID:         utils.Ptr[int](2),
	QuoteAssetID:        utils.Ptr[int](1),
	Description:         "",
	AssetTypeID:         utils.Ptr[int](1),
	CreatedBy:           "SYSTEM",
	CreatedAt:           utils.SampleCreatedAtTime,
	UpdatedBy:           "SYSTEM",
	UpdatedAt:           utils.SampleCreatedAtTime,
	ChainID:             utils.Ptr[int](1),
	CategoryID:          utils.Ptr[int](28),
	SubCategoryID:       utils.Ptr[int](20),
	IsDefaultQuote:      utils.Ptr[bool](true),
	IgnoreMarketData:    utils.Ptr[bool](false),
	Decimals:            utils.Ptr[int](1),
	ContractAddress:     "SYSTEM",
	StartingBlockNumber: utils.Ptr[uint64](1),
	ImportGeth:          nil,
	ImportGethInitial:   nil,
}
var allData = []Asset{data1, data2}

func AddAssetToMockRows(mock pgxmock.PgxPoolIface, dataList []Asset) *pgxmock.Rows {
	rows := mock.NewRows(columns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,
			data.UUID,
			data.Name,
			data.AlternateName,
			data.Cusip,
			data.Ticker,
			data.BaseAssetID,
			data.QuoteAssetID,
			data.Description,
			data.AssetTypeID,
			data.CreatedBy,
			data.CreatedAt,
			data.UpdatedBy,
			data.UpdatedAt,
			data.ChainID,
			data.CategoryID,
			data.SubCategoryID,
			data.IsDefaultQuote,
			data.IgnoreMarketData,
			data.Decimals,
			data.ContractAddress,
			data.StartingBlockNumber,
			data.ImportGeth,
			data.ImportGethInitial,
		)
	}
	return rows
}

var dataAssetWithSources1 = AssetWithSources{
	Asset:            data1,
	SourceID:         utils.Ptr[int](1),
	SourceIdentifier: "ETH",
}

var dataAssetWithSources2 = AssetWithSources{
	Asset:            data2,
	SourceID:         utils.Ptr[int](1),
	SourceIdentifier: "XBT",
}

var allDataAssetWithSources = []AssetWithSources{dataAssetWithSources1, dataAssetWithSources2}

func AddAssetWithSourcesToMockRows(mock pgxmock.PgxPoolIface, dataList []AssetWithSources) *pgxmock.Rows {
	rows := mock.NewRows(columnsAssetWithSources)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,
			data.UUID,
			data.Name,
			data.AlternateName,
			data.Cusip,
			data.Ticker,
			data.BaseAssetID,
			data.QuoteAssetID,
			data.Description,
			data.AssetTypeID,
			data.CreatedBy,
			data.CreatedAt,
			data.UpdatedBy,
			data.UpdatedAt,
			data.ChainID,
			data.CategoryID,
			data.SubCategoryID,
			data.IsDefaultQuote,
			data.IgnoreMarketData,
			data.Decimals,
			data.ContractAddress,
			data.StartingBlockNumber,
			data.ImportGeth,
			data.ImportGethInitial,
			data.SourceID,
			data.SourceIdentifier,
		)
	}
	return rows
}

func TestGetAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
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
	noRows := pgxmock.NewRows(columns)
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

func TestGetAssetByTicker(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data2
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
	noRows := pgxmock.NewRows(columns)
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

func TestGetAssetByContractAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data2
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
	noRows := pgxmock.NewRows(columns)
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

func TestGetAssetByCusip(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data2
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
	noRows := pgxmock.NewRows(columns)
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

func TestGetAssetByBaseAndQuoteID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data2
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
	noRows := pgxmock.NewRows(columns)
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

func TestGetGethImportAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{data1, data2}
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM assets WHERE import_geth = TRUE").WillReturnRows(mockRows)
	foundAssets, err := GetGethImportAssets(mock)

	if err != nil {
		t.Fatalf("an error '%s' in GetGethImportAssets", err)
	}
	testAssets := allData
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

func TestRemoveAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
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

func TestRemoveAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
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
	dataList := []Asset{data1, data2}
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM public.get_current_assets").WillReturnRows(mockRows)
	foundAssets, err := GetCurrentTradingAssets(mock)

	if err != nil {
		t.Fatalf("an error '%s' in GetCurrentTradingAssets", err)
	}
	testAssets := allData
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

func TestGetCryptoAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{data1, data2}
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM assets where asset_type_id = 1").WillReturnRows(mockRows)
	foundAssets, err := GetCryptoAssets(mock)

	if err != nil {
		t.Fatalf("an error '%s' in GetCryptoAssets", err)
	}
	testAssets := allData
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

// asset with sources
func TestGetAssetsByAssetTypeAndSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := allDataAssetWithSources
	mockRows := AddAssetWithSourcesToMockRows(mock, dataList)
	assetTypeID := dataAssetWithSources1.AssetTypeID
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*assetTypeID, *sourceID).WillReturnRows(mockRows)
	foundAssets, err := GetAssetsByAssetTypeAndSource(mock, assetTypeID, sourceID, excludeIgnoreMarketData)

	if err != nil {
		t.Fatalf("an error '%s' in GetAssetsByAssetTypeAndSource", err)
	}
	testAssets := allDataAssetWithSources
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

func TestGetCryptoAssetsBySourceId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := allDataAssetWithSources
	mockRows := AddAssetWithSourcesToMockRows(mock, dataList)
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*sourceID).WillReturnRows(mockRows)
	foundAssets, err := GetCryptoAssetsBySourceId(mock, sourceID, excludeIgnoreMarketData)

	if err != nil {
		t.Fatalf("an error '%s' in GetCryptoAssetsBySourceId", err)
	}
	testAssets := allDataAssetWithSources
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

func TestGetCryptoAssetsBySourceID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := allDataAssetWithSources
	mockRows := AddAssetWithSourcesToMockRows(mock, dataList)
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(*sourceID).WillReturnRows(mockRows)
	foundAssets, err := GetCryptoAssetsBySourceID(mock, sourceID, excludeIgnoreMarketData)

	if err != nil {
		t.Fatalf("an error '%s' in GetCryptoAssetsBySourceID", err)
	}
	testAssets := allData
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
	noRows := pgxmock.NewRows(columns)
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

func TestGetAssetWithSourceByAssetIdsAndSourceID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := allDataAssetWithSources
	assetIDs := []int{*dataAssetWithSources1.Asset.ID, *dataAssetWithSources2.Asset.ID}
	mockRows := AddAssetWithSourcesToMockRows(mock, dataList)
	sourceID := dataAssetWithSources1.SourceID
	excludeIgnoreMarketData := true
	mock.ExpectQuery("^SELECT (.+) FROM assets").WithArgs(pq.Array(assetIDs), *sourceID).WillReturnRows(mockRows)
	foundAssets, err := GetAssetWithSourceByAssetIdsAndSourceID(mock, assetIDs, sourceID, excludeIgnoreMarketData)

	if err != nil {
		t.Fatalf("an error '%s' in GetAssetWithSourceByAssetIdsAndSourceID", err)
	}
	testAssets := allDataAssetWithSources
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

func TestGetAssetList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{data1, data2}
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery(fmt.Sprintf("^SELECT (.+) FROM assets WHERE")).WillReturnRows(mockRows)
	ids := make([]int, 0)
	ids = append(ids, *data1.ID)
	ids = append(ids, *data2.ID)
	foundAssets, err := GetAssetList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetList", err)
	}
	testAssets := allData
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

func TestGetAssetsByChainId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{data1}
	chainID := data1.ChainID
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
	chainID := data1.ChainID
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

func TestGetAssetListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Asset{data1}
	mockRows := AddAssetToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"id = 1", "import_geth = TRUE"}
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
	filters := []string{"id = 1", "import_geth = TRUE"}
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
	dataList := allData
	sourceID := 1
	mockRows := AddAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM get_default_quotes").WithArgs(sourceID).WillReturnRows(mockRows)
	foundAssets, err := GetDefaultQuoteAssetListBySourceID(mock, &sourceID)
	if err != nil {
		t.Fatalf("an error '%s' in GetDefaultQuoteAssetListBySourceID", err)
	}
	testAssets := allData
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
	chainID := data1.ChainID
	mock.ExpectQuery("^SELECT (.+) FROM get_default_quotes").WithArgs(*chainID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssets, err := GetDefaultQuoteAssetListBySourceID(mock, chainID)
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

func TestGetDefaultQuoteAssetList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mockRows := AddAssetToMockRows(mock, allData)
	mock.ExpectQuery("^SELECT (.+) FROM get_default_quotes").WillReturnRows(mockRows)
	foundAssets, err := GetDefaultQuoteAssetList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetDefaultQuoteAssetList", err)
	}
	testAssets := allData
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

func TestUpdateAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
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

func TestUpdateAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
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
	targetData := data1
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

func TestInsertAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
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
	targetData := data1
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
