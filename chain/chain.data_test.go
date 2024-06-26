package chain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var columns = []string{
	"id",                 //1
	"uuid",               //2
	"base_asset_id",      //3
	"name",               //4
	"alternate_name",     //5
	"address",            //6
	"chain_type_id",      //7
	"description",        //8
	"created_by",         //9
	"created_at",         //10
	"updated_by",         //11
	"updated_at",         //12
	"rpc_url",            //13
	"chain_id",           //14
	"block_explorer_url", //15
	"rpc_url_dev",        //16
	"rpc_url_prod",       //17
	"rpc_url_archive",    //18
}

var data1 = Chain{
	ID:               utils.Ptr[int](1),
	UUID:             "880607ab-2833-4ad7-a231-b983a61c7b39",
	BaseAssetID:      utils.Ptr[int](2),
	Name:             "Ethereum",
	AlternateName:    "ETH",
	Address:          "",
	ChainTypeID:      utils.Ptr[int](19),
	Description:      "",
	CreatedBy:        "SYSTEM",
	CreatedAt:        utils.SampleCreatedAtTime,
	UpdatedBy:        "SYSTEM",
	UpdatedAt:        utils.SampleCreatedAtTime,
	RpcURL:           "ws://erigon.dappnode:8545",
	ChainID:          utils.Ptr[int](1),
	BlockExplorerURL: "https://etherscan.io/",
	RpcURLDev:        "ws://erigon.dappnode:8545",
	RpcURLProd:       "ws://erigon.dappnode:8545",
	RpcURLArchive:    "ws://erigon.dappnode:8545",
}

var data2 = Chain{
	ID:               utils.Ptr[int](1),
	UUID:             "880607ab-2833-4ad7-a231-b983a61c7b334",
	BaseAssetID:      utils.Ptr[int](495),
	Name:             "Avalanche",
	AlternateName:    "AVAX",
	Address:          "",
	ChainTypeID:      utils.Ptr[int](19),
	Description:      "",
	CreatedBy:        "SYSTEM",
	CreatedAt:        utils.SampleCreatedAtTime,
	UpdatedBy:        "SYSTEM",
	UpdatedAt:        utils.SampleCreatedAtTime,
	RpcURL:           "https://api.avax.network/ext/bc/C/rpc",
	ChainID:          utils.Ptr[int](43114),
	BlockExplorerURL: "https://snowtrace.io",
	RpcURLDev:        "https://api.avax.network/ext/bc/C/rpc",
	RpcURLProd:       "https://api.avax.network/ext/bc/C/rpc",
	RpcURLArchive:    "https://api.avax.network/ext/bc/C/rpc",
}
var allData = []Chain{data1, data2}

func AddChainToMockRows(mock pgxmock.PgxPoolIface, dataList []Chain) *pgxmock.Rows {
	rows := mock.NewRows(columns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,               //1
			data.UUID,             //2
			data.BaseAssetID,      //3
			data.Name,             //4
			data.AlternateName,    //5
			data.Address,          //6
			data.ChainTypeID,      //8
			data.Description,      //7
			data.CreatedBy,        //9
			data.CreatedAt,        //10
			data.UpdatedBy,        //11
			data.UpdatedAt,        //12
			data.RpcURL,           //13
			data.ChainID,          //14
			data.BlockExplorerURL, //15
			data.RpcURLDev,        //16
			data.RpcURLProd,       //17
			data.RpcURLArchive,    //18

		)
	}
	return rows
}

func TestGetChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data2
	dataList := []Chain{targetData}
	chainID := targetData.ChainID
	mockRows := AddChainToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM chains").WithArgs(*chainID).WillReturnRows(mockRows)
	foundChain, err := GetChain(mock, chainID)
	if err != nil {
		t.Fatalf("an error '%s' in GetChain", err)
	}
	if cmp.Equal(*foundChain, targetData) == false {
		t.Errorf("Expected Chain From Method GetChain: %v is different from actual %v", foundChain, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetChainForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainID := 999
	noRows := pgxmock.NewRows(columns)
	mock.ExpectQuery("^SELECT (.+) FROM chains").WithArgs(chainID).WillReturnRows(noRows)
	foundChain, err := GetChain(mock, &chainID)
	if err != nil {
		t.Fatalf("an error '%s' in GetChain", err)
	}
	if foundChain != nil {
		t.Errorf("Expected Chain From Method GetChain: to be empty but got this: %v", foundChain)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetChainForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainID := -1
	mock.ExpectQuery("^SELECT (.+) FROM chains").WithArgs(chainID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChain, err := GetChain(mock, &chainID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetChain", err)
	}
	if foundChain != nil {
		t.Errorf("Expected Chain From Method GetChain: to be empty but got this: %v", foundChain)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetChainByAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data2
	dataList := []Chain{targetData}
	chainAddress := targetData.Address
	mockRows := AddChainToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM chains").WithArgs(chainAddress).WillReturnRows(mockRows)
	foundChain, err := GetChainByAddress(mock, chainAddress)
	if err != nil {
		t.Fatalf("an error '%s' in GetChainByAddress", err)
	}
	if cmp.Equal(*foundChain, targetData) == false {
		t.Errorf("Expected Chain From Method GetChainByAddress: %v is different from actual %v", foundChain, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetChainByAddressForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainAddress := "non-existing-address"
	noRows := pgxmock.NewRows(columns)
	mock.ExpectQuery("^SELECT (.+) FROM chains").WithArgs(chainAddress).WillReturnRows(noRows)
	foundChain, err := GetChainByAddress(mock, chainAddress)
	if err != nil {
		t.Fatalf("an error '%s' in GetChainByAddress", err)
	}
	if foundChain != nil {
		t.Errorf("Expected Chain From Method GetChainByAddress: to be empty but got this: %v", foundChain)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetChainByAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainAddress := "xxx"
	mock.ExpectQuery("^SELECT (.+) FROM chains").WithArgs(chainAddress).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChain, err := GetChainByAddress(mock, chainAddress)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetChainByAddress", err)
	}
	if foundChain != nil {
		t.Errorf("Expected Chain From Method GetChainByAddress: to be empty but got this: %v", foundChain)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestGetChainByAlternateName(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data2
	dataList := []Chain{targetData}
	chainAlternateName := targetData.AlternateName
	mockRows := AddChainToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM chains").WithArgs(chainAlternateName).WillReturnRows(mockRows)
	foundChain, err := GetChainByAlternateName(mock, chainAlternateName)
	if err != nil {
		t.Fatalf("an error '%s' in GetChainByAlternateName", err)
	}
	if cmp.Equal(*foundChain, targetData) == false {
		t.Errorf("Expected Chain From Method GetChainByAlternateName: %v is different from actual %v", foundChain, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetChainByAlternateNameForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainAlternateName := "non-existing-alternate-name"
	noRows := pgxmock.NewRows(columns)
	mock.ExpectQuery("^SELECT (.+) FROM chains").WithArgs(chainAlternateName).WillReturnRows(noRows)
	foundChain, err := GetChainByAlternateName(mock, chainAlternateName)
	if err != nil {
		t.Fatalf("an error '%s' in GetChainByAlternateName", err)
	}
	if foundChain != nil {
		t.Errorf("Expected Chain From Method GetChainByAlternateName: to be empty but got this: %v", foundChain)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetChainByAlternateNameForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainAlternateName := "xxx"
	mock.ExpectQuery("^SELECT (.+) FROM chains").WithArgs(chainAlternateName).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChain, err := GetChainByAlternateName(mock, chainAlternateName)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetChainByAlternateName", err)
	}
	if foundChain != nil {
		t.Errorf("Expected Chain From Method GetChainByAlternateName: to be empty but got this: %v", foundChain)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	chainID := targetData.ChainID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM chains").WithArgs(*chainID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveChain(mock, chainID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveChain", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveChainOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM chains").WithArgs(chainID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveChain(mock, &chainID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetChainList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Chain{data1, data2}
	mockRows := AddChainToMockRows(mock, dataList)
	chainIDs := []int{1, 2}
	mock.ExpectQuery("^SELECT (.+) FROM chains").WillReturnRows(mockRows)
	foundChains, err := GetChainList(mock, chainIDs)
	if err != nil {
		t.Fatalf("an error '%s' in GetChainList", err)
	}
	testChains := allData
	for i, foundChain := range foundChains {
		if cmp.Equal(foundChain, testChains[i]) == false {
			t.Errorf("Expected Chain From Method GetChainList: %v is different from actual %v", foundChain, testChains[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetChainListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainIDs := []int{-1, -2}
	mock.ExpectQuery("^SELECT (.+) FROM chains").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChains, err := GetChainList(mock, chainIDs)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetChainList", err)
	}
	if len(foundChains) != 0 {
		t.Errorf("Expected From Method GetChainList: to be empty but got this: %v", foundChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestGetChainListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := allData
	mockRows := AddChainToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"chain_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM chains").WillReturnRows(mockRows)
	foundChains, err := GetChainListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetChainListByPagination", err)
	}
	testChains := dataList
	for i, foundChain := range foundChains {
		if cmp.Equal(foundChain, testChains[i]) == false {
			t.Errorf("Expected Chain From Method GetChainListByPagination: %v is different from actual %v", foundChain, testChains[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetChainListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"chain_type_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM chains").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChains, err := GetChainListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetChainListByPagination", err)
	}
	if len(foundChains) != 0 {
		t.Errorf("Expected From Method GetChainListByPagination: to be empty but got this: %v", foundChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalChainCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM chains").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalChainCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalChainCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalChainCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalChainCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM chains").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalChainCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalChainCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalChainCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE chains").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.Address,          //3
		targetData.ChainTypeID,      //4
		targetData.Description,      //5
		targetData.UpdatedBy,        //6
		targetData.BaseAssetID,      //7
		targetData.RpcURL,           //8
		targetData.ChainID,          //9
		targetData.BlockExplorerURL, //10
		targetData.RpcURLDev,        //11
		targetData.RpcURLProd,       //12
		targetData.RpcURLArchive,    //13
		targetData.ID,               //14
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateChain(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateChain", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateChainOnFailure(t *testing.T) {
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
	mock.ExpectExec("^UPDATE chains").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.Address,          //3
		targetData.ChainTypeID,      //4
		targetData.Description,      //5
		targetData.UpdatedBy,        //6
		targetData.BaseAssetID,      //7
		targetData.RpcURL,           //8
		targetData.ChainID,          //9
		targetData.BlockExplorerURL, //10
		targetData.RpcURLDev,        //11
		targetData.RpcURLProd,       //12
		targetData.RpcURLArchive,    //13
		targetData.ID,               //14
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateChain(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO chains").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.Address,          //3
		targetData.ChainTypeID,      //4
		targetData.Description,      //5
		targetData.UpdatedBy,        //6
		targetData.BaseAssetID,      //7
		targetData.RpcURL,           //8
		targetData.ChainID,          //9
		targetData.BlockExplorerURL, //10
		targetData.RpcURLDev,        //11
		targetData.RpcURLProd,       //12
		targetData.RpcURLArchive,    //13
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	chainID, err := InsertChain(mock, &targetData)
	if chainID < 0 {
		t.Fatalf("chainID should not be negative ID: %d", chainID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertChain", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertChainOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO chains").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.Address,          //3
		targetData.ChainTypeID,      //4
		targetData.Description,      //5
		targetData.UpdatedBy,        //6
		targetData.BaseAssetID,      //7
		targetData.RpcURL,           //8
		targetData.ChainID,          //9
		targetData.BlockExplorerURL, //10
		targetData.RpcURLDev,        //11
		targetData.RpcURLProd,       //12
		targetData.RpcURLArchive,    //13
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	chainID, err := InsertChain(mock, &targetData)
	if chainID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error chainID: %d", chainID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertChainOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO chains").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.Address,          //3
		targetData.ChainTypeID,      //4
		targetData.Description,      //5
		targetData.UpdatedBy,        //6
		targetData.BaseAssetID,      //7
		targetData.RpcURL,           //8
		targetData.ChainID,          //9
		targetData.BlockExplorerURL, //10
		targetData.RpcURLDev,        //11
		targetData.RpcURLProd,       //12
		targetData.RpcURLArchive,    //13
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	chainID, err := InsertChain(mock, &targetData)
	if chainID >= 0 {
		t.Fatalf("Expecting -1 for chainID because of error chainID: %d", chainID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
