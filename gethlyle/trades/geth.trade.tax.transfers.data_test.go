package gethlyletrades

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumnsGethTradeTaxTransfers = []string{
	"geth_trade_id",    //1
	"geth_transfer_id", //2
	"tax_id",           //3
	"uuid",             //4
	"name",             //5
	"alternate_name",   //6
	"description",      //7
	"created_by",       //8
	"created_at",       //9
	"updated_by",       //10
	"updated_at",       //11
}
var DBColumnsInsertGethTradeTaxTransfers = []string{
	"geth_trade_id",    //1
	"geth_transfer_id", //2
	"tax_id",           //3
	"uuid",             //4
	"name",             //5
	"alternate_name",   //6
	"description",      //7
	"created_by",       //8
	"created_at",       //9
	"updated_by",       //10
	"updated_at",       //11
}

var TestData1GethTradeTaxTransfer = GethTradeTaxTransfer{
	GethTradeID:    utils.Ptr[int](1),
	GethTransferID: utils.Ptr[int](2),
	TaxID:          utils.Ptr[int](1),
	UUID:           "01ef85e8-2c26-441e-8c7f-71d79518ad72",
	Name:           "APU-ETH/WETH",
	AlternateName:  "APU-ETH/WETH",
	Description:    "Imported by Geth Dex Analyzer",
	CreatedBy:      "SYSTEM",
	CreatedAt:      utils.SampleCreatedAtTime,
	UpdatedBy:      "SYSTEM",
	UpdatedAt:      utils.SampleCreatedAtTime,
}

var TestData2GethTradeTaxTransfer = GethTradeTaxTransfer{
	GethTradeID:    utils.Ptr[int](2),
	GethTransferID: utils.Ptr[int](3),
	TaxID:          utils.Ptr[int](2),
	UUID:           "bd45f8f9-84e5-4052-a895-09f2e5993622",
	Name:           "PEPE/WETH",
	AlternateName:  "PEPE/WETH",
	Description:    "Imported by Geth Dex Analyzer",
	CreatedBy:      "SYSTEM",
	CreatedAt:      utils.SampleCreatedAtTime,
	UpdatedBy:      "SYSTEM",
	UpdatedAt:      utils.SampleCreatedAtTime,
}
var TestAllDataGethTradeTaxTransfer = []GethTradeTaxTransfer{TestData1GethTradeTaxTransfer, TestData2GethTradeTaxTransfer}

func AddGethTradeTaxTransferToMockRows(mock pgxmock.PgxPoolIface, dataList []GethTradeTaxTransfer) *pgxmock.Rows {
	rows := mock.NewRows(DBColumnsGethTradeTaxTransfers)
	for _, data := range dataList {
		rows.AddRow(
			data.GethTradeID,    //1
			data.GethTransferID, //2
			data.TaxID,          //3
			data.UUID,           //4
			data.Name,           //5
			data.AlternateName,  //6
			data.Description,    //7
			data.CreatedBy,      //8
			data.CreatedAt,      //9
			data.UpdatedBy,      //10
			data.UpdatedAt,      //11
		)
	}
	return rows
}

func TestGetAllGethTradeTaxTransfersByTradeID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTradeTaxTransfer{TestData1GethTradeTaxTransfer}
	mockRows := AddGethTradeTaxTransferToMockRows(mock, dataList)
	gethTradeID := TestData1GethTradeTaxTransfer.GethTradeID
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(*gethTradeID).WillReturnRows(mockRows)
	foundGethTradeTaxTransferList, err := GetAllGethTradeTaxTransfersByTradeID(mock, gethTradeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllGethTradeTaxTransfersByTradeID", err)
	}
	testMarketDataList := dataList
	for i, foundGethTradeTaxTransfer := range foundGethTradeTaxTransferList {
		if cmp.Equal(foundGethTradeTaxTransfer, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTradeTaxTransfer From Method GetAllGethTradeTaxTransfersByTradeID: %v is different from actual %v", foundGethTradeTaxTransfer, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTradeTaxTransfersByTradeIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := -999
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(gethTradeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeTaxTransferList, err := GetAllGethTradeTaxTransfersByTradeID(mock, &gethTradeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethTradeTaxTransfersByTradeID", err)
	}
	if len(foundGethTradeTaxTransferList) != 0 {
		t.Errorf("Expected From Method GetAllGethTradeTaxTransfersByTradeID: to be empty but got this: %v", foundGethTradeTaxTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeTaxTransfer(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2GethTradeTaxTransfer
	dataList := []GethTradeTaxTransfer{targetData}
	gethTradeID := targetData.GethTradeID
	gethTransferID := targetData.GethTransferID
	mockRows := AddGethTradeTaxTransferToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_transfers").WithArgs(*gethTradeID, *gethTransferID).WillReturnRows(mockRows)
	foundGethTradeTaxTransfer, err := GetGethTradeTaxTransfer(mock, gethTradeID, gethTransferID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeTaxTransfer", err)
	}
	if cmp.Equal(*foundGethTradeTaxTransfer, targetData) == false {
		t.Errorf("Expected GethTradeTaxTransfer From Method GetGethTradeTaxTransfer: %v is different from actual %v", foundGethTradeTaxTransfer, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeTaxTransferForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := 999
	gethTransferID := 999
	noRows := pgxmock.NewRows(DBColumnsGethTradeTaxTransfers)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_transfers").WithArgs(gethTradeID, gethTransferID).WillReturnRows(noRows)
	foundGethTradeTaxTransfer, err := GetGethTradeTaxTransfer(mock, &gethTradeID, &gethTransferID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeTaxTransfer", err)
	}
	if foundGethTradeTaxTransfer != nil {
		t.Errorf("Expected GethTradeTaxTransfer From Method GetGethTradeTaxTransfer: to be empty but got this: %v", foundGethTradeTaxTransfer)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeTaxTransferForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := -1
	gethTransferID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_transfers").WithArgs(gethTradeID, gethTransferID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeTaxTransfer, err := GetGethTradeTaxTransfer(mock, &gethTradeID, &gethTransferID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeTaxTransfer", err)
	}
	if foundGethTradeTaxTransfer != nil {
		t.Errorf("Expected GethTradeTaxTransfer From Method GetGethTradeTaxTransfer: to be empty but got this: %v", foundGethTradeTaxTransfer)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestRemoveGethTradeTaxTransfer(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeTaxTransfer
	gethTradeID := targetData.GethTradeID
	gethTransferID := targetData.GethTransferID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_trade_transfers").WithArgs(*gethTradeID, *gethTransferID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethTradeTaxTransfer(mock, gethTradeID, gethTransferID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethTradeTaxTransfer", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTradeTaxTransferOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := -1
	gethTransferID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_trade_transfers").WithArgs(gethTradeID, gethTransferID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethTradeTaxTransfer(mock, &gethTradeID, &gethTransferID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeTaxTransferList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTradeTaxTransfer{TestData1GethTradeTaxTransfer, TestData2GethTradeTaxTransfer}
	mockRows := AddGethTradeTaxTransferToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_transfers").WillReturnRows(mockRows)
	gethTradeIds := []int{}
	gethTransferIDs := []int{}
	foundGethTradeTaxTransferList, err := GetGethTradeTaxTransferList(mock, gethTradeIds, gethTransferIDs)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeTaxTransferList", err)
	}
	testMarketDataList := TestAllDataGethTradeTaxTransfer
	for i, foundGethTradeTaxTransfer := range foundGethTradeTaxTransferList {
		if cmp.Equal(foundGethTradeTaxTransfer, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTradeTaxTransfer From Method GetGethTradeTaxTransferList: %v is different from actual %v", foundGethTradeTaxTransfer, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeTaxTransferListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_transfers").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	gethTradeIds := []int{}
	gethTransferIDs := []int{}
	foundGethTradeTaxTransferList, err := GetGethTradeTaxTransferList(mock, gethTradeIds, gethTransferIDs)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeTaxTransferList", err)
	}
	if len(foundGethTradeTaxTransferList) != 0 {
		t.Errorf("Expected From Method GetGethTradeTaxTransferList: to be empty but got this: %v", foundGethTradeTaxTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTradeTaxTransfer(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeTaxTransfer
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_trade_transfers").WithArgs(
		targetData.TaxID,          //1
		targetData.Name,           //2
		targetData.AlternateName,  //3
		targetData.Description,    //4
		targetData.UpdatedBy,      //5
		targetData.GethTradeID,    //6
		targetData.GethTransferID, //7
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethTradeTaxTransfer(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethTradeTaxTransfer", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTradeTaxTransferOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeTaxTransfer
	// name can't be nil
	targetData.Name = ""
	targetData.GethTradeID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethTradeTaxTransfer(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTradeTaxTransferOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeTaxTransfer
	// name can't be nil
	targetData.Name = ""
	targetData.GethTradeID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_trade_transfers").WithArgs(
		targetData.TaxID,          //1
		targetData.Name,           //2
		targetData.AlternateName,  //3
		targetData.Description,    //4
		targetData.UpdatedBy,      //5
		targetData.GethTradeID,    //6
		targetData.GethTransferID, //7
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethTradeTaxTransfer(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeTaxTransfer(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeTaxTransfer
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_trade_transfers").WithArgs(
		targetData.GethTradeID,    //1
		targetData.GethTransferID, //2
		targetData.TaxID,          //3
		targetData.Name,           //4
		targetData.AlternateName,  //5
		targetData.Description,    //6
		targetData.CreatedBy,      //7
	).WillReturnRows(pgxmock.NewRows([]string{"geth_trade_id", "geth_transfer_id"}).AddRow(1, 2))
	mock.ExpectCommit()
	gethTradeID, gethTransferID, err := InsertGethTradeTaxTransfer(mock, &targetData)
	if gethTradeID < 0 {
		t.Fatalf("gethTradeID should not be negative ID: %d", gethTradeID)
	}
	if gethTransferID < 0 {
		t.Fatalf("gethTransferID should not be negative ID: %d", gethTransferID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTradeTaxTransfer", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeTaxTransferOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeTaxTransfer
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_trade_transfers").WithArgs(
		targetData.GethTradeID,    //1
		targetData.GethTransferID, //2
		targetData.TaxID,          //3
		targetData.Name,           //4
		targetData.AlternateName,  //5
		targetData.Description,    //6
		targetData.CreatedBy,      //7
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethTradeID, gethTransferID, err := InsertGethTradeTaxTransfer(mock, &targetData)
	if gethTradeID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethTradeID: %d", gethTradeID)
	}
	if gethTransferID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethTransferID: %d", gethTransferID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeTaxTransferOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeTaxTransfer
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_trade_transfers").WithArgs(
		targetData.GethTradeID,    //1
		targetData.GethTransferID, //2
		targetData.TaxID,          //3
		targetData.Name,           //4
		targetData.AlternateName,  //5
		targetData.Description,    //6
		targetData.CreatedBy,      //7
	).WillReturnRows(pgxmock.NewRows([]string{"geth_trade_id", "geth_transfer_id"}).AddRow(1, 2))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethTradeID, gethTransferID, err := InsertGethTradeTaxTransfer(mock, &targetData)
	if gethTradeID >= 0 {
		t.Fatalf("Expecting -1 for gethTradeID because of error gethTradeID: %d", gethTradeID)
	}
	if gethTransferID >= 0 {
		t.Fatalf("Expecting -1 for gethTransferID because of error gethTransferID: %d", gethTransferID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeTaxTransfers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_trade_transfers"}, DBColumnsInsertGethTradeTaxTransfers)
	targetData := TestAllDataGethTradeTaxTransfer
	err = InsertGethTradeTaxTransfers(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTradeTaxTransfers", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeTaxTransfersOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_trade_transfers"}, DBColumnsInsertGethTradeTaxTransfers).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllDataGethTradeTaxTransfer
	err = InsertGethTradeTaxTransfers(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeTaxTransferListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllDataGethTradeTaxTransfer
	mockRows := AddGethTradeTaxTransferToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"tax_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_transfers").WillReturnRows(mockRows)
	foundGethTradeTaxTransferList, err := GetGethTradeTaxTransferListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeTaxTransferListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundGethTradeTaxTransferList[i]) == false {
			t.Errorf("Expected sourceData From Method GetGethTradeTaxTransferListByPagination: %v is different from actual %v", sourceData, foundGethTradeTaxTransferList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeTaxTransferListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_transfers").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeTaxTransferList, err := GetGethTradeTaxTransferListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeTaxTransferListByPagination", err)
	}
	if len(foundGethTradeTaxTransferList) != 0 {
		t.Errorf("Expected From Method GetGethTradeTaxTransferListByPagination: to be empty but got this: %v", foundGethTradeTaxTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethTradeTaxTransferCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_trade_transfers").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalGethTradeTaxTransferCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalGethTradeTaxTransferCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalGethTradeTaxTransferCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethTradeTaxTransferCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_trade_transfers").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalGethTradeTaxTransferCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalGethTradeTaxTransferCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalGethTradeTaxTransferCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
