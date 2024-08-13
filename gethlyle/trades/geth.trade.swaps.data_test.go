package gethlyletrades

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	gethlyleswaps "github.com/kfukue/lyle-labs-libraries/gethlyle/swaps"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumnsGethTradeSwaps = []string{
	"geth_trade_id",  //1
	"geth_swap_id",   //2
	"uuid",           //3
	"name",           //4
	"alternate_name", //5
	"description",    //6
	"created_by",     //7
	"created_at",     //8
	"updated_by",     //9
	"updated_at",     //10
}
var DBColumnsInsertGethTradeSwaps = []string{
	"geth_trade_id",  //1
	"geth_swap_id",   //2
	"uuid",           //3
	"name",           //4
	"alternate_name", //5
	"description",    //6
	"created_by",     //7
	"created_at",     //8
	"updated_by",     //9
	"updated_at",     //10
}

var TestData1GethTradeSwap = GethTradeSwap{
	GethTradeID:   utils.Ptr[int](1),
	GethSwapID:    utils.Ptr[int](2),
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72",
	Name:          "APU-ETH/WETH",
	AlternateName: "APU-ETH/WETH",
	Description:   "Imported by Geth Dex Analyzer",
	CreatedBy:     "SYSTEM",
	CreatedAt:     utils.SampleCreatedAtTime,
	UpdatedBy:     "SYSTEM",
	UpdatedAt:     utils.SampleCreatedAtTime,
}

var TestData2GethTradeSwap = GethTradeSwap{
	GethTradeID:   utils.Ptr[int](2),
	GethSwapID:    utils.Ptr[int](3),
	UUID:          "bd45f8f9-84e5-4052-a895-09f2e5993622",
	Name:          "PEPE/WETH",
	AlternateName: "PEPE/WETH",
	Description:   "Imported by Geth Dex Analyzer",
	CreatedBy:     "SYSTEM",
	CreatedAt:     utils.SampleCreatedAtTime,
	UpdatedBy:     "SYSTEM",
	UpdatedAt:     utils.SampleCreatedAtTime,
}
var TestAllDataGethTradeSwap = []GethTradeSwap{TestData1GethTradeSwap, TestData2GethTradeSwap}

func AddGethTradeSwapToMockRows(mock pgxmock.PgxPoolIface, dataList []GethTradeSwap) *pgxmock.Rows {
	rows := mock.NewRows(DBColumnsGethTradeSwaps)
	for _, data := range dataList {
		rows.AddRow(
			data.GethTradeID,   //1
			data.GethSwapID,    //2
			data.UUID,          //3
			data.Name,          //4
			data.AlternateName, //5
			data.Description,   //6
			data.CreatedBy,     //7
			data.CreatedAt,     //8
			data.UpdatedBy,     //9
			data.UpdatedAt,     //10
		)
	}
	return rows
}

func TestGetAllGethTradeSwapsByTradeID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTradeSwap{TestData1GethTradeSwap}
	mockRows := AddGethTradeSwapToMockRows(mock, dataList)
	gethTradeID := TestData1GethTradeSwap.GethTradeID
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(*gethTradeID).WillReturnRows(mockRows)
	foundGethTradeSwapList, err := GetAllGethTradeSwapsByTradeID(mock, gethTradeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllGethTradeSwapsByTradeID", err)
	}
	testMarketDataList := dataList
	for i, foundGethTradeSwap := range foundGethTradeSwapList {
		if cmp.Equal(foundGethTradeSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTradeSwap From Method GetAllGethTradeSwapsByTradeID: %v is different from actual %v", foundGethTradeSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTradeSwapsByTradeIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := -999
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(gethTradeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeSwapList, err := GetAllGethTradeSwapsByTradeID(mock, &gethTradeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethTradeSwapsByTradeID", err)
	}
	if len(foundGethTradeSwapList) != 0 {
		t.Errorf("Expected From Method GetAllGethTradeSwapsByTradeID: to be empty but got this: %v", foundGethTradeSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTradeSwapsByTradeIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(gethTradeID).WillReturnRows(differentModelRows)
	foundGethTradeSwapList, err := GetAllGethTradeSwapsByTradeID(mock, &gethTradeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethTradeSwapsByTradeID", err)
	}
	if foundGethTradeSwapList != nil {
		t.Errorf("Expected foundGethTradeSwapList From Method GetAllGethTradeSwapsByTradeID: to be empty but got this: %v", foundGethTradeSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeSwap(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2GethTradeSwap
	dataList := []GethTradeSwap{targetData}
	gethTradeID := targetData.GethTradeID
	gethSwapID := targetData.GethSwapID
	mockRows := AddGethTradeSwapToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WithArgs(*gethTradeID, *gethSwapID).WillReturnRows(mockRows)
	foundGethTradeSwap, err := GetGethTradeSwap(mock, gethSwapID, gethTradeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeSwap", err)
	}
	if cmp.Equal(*foundGethTradeSwap, targetData) == false {
		t.Errorf("Expected GethTradeSwap From Method GetGethTradeSwap: %v is different from actual %v", foundGethTradeSwap, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeSwapForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := 999
	gethSwapID := 999
	noRows := pgxmock.NewRows(DBColumnsGethTradeSwaps)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WithArgs(gethTradeID, gethSwapID).WillReturnRows(noRows)
	foundGethTradeSwap, err := GetGethTradeSwap(mock, &gethSwapID, &gethTradeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeSwap", err)
	}
	if foundGethTradeSwap != nil {
		t.Errorf("Expected GethTradeSwap From Method GetGethTradeSwap: to be empty but got this: %v", foundGethTradeSwap)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeSwapForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := -1
	gethSwapID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WithArgs(gethTradeID, gethSwapID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeSwap, err := GetGethTradeSwap(mock, &gethTradeID, &gethSwapID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeSwap", err)
	}
	if foundGethTradeSwap != nil {
		t.Errorf("Expected GethTradeSwap From Method GetGethTradeSwap: to be empty but got this: %v", foundGethTradeSwap)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeSwapForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := -1
	gethSwapID := -1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WithArgs(gethTradeID, gethSwapID).WillReturnRows(differentModelRows)
	foundGethTradeSwap, err := GetGethTradeSwap(mock, &gethTradeID, &gethSwapID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeSwap", err)
	}
	if foundGethTradeSwap != nil {
		t.Errorf("Expected GethTrade From Method GetGethTradeSwap: to be empty but got this: %v", foundGethTradeSwap)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTradeSwap(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeSwap
	gethTradeID := targetData.GethTradeID
	gethSwapID := targetData.GethSwapID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_trade_swaps").WithArgs(*gethTradeID, *gethSwapID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethTradeSwap(mock, gethTradeID, gethSwapID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethTradeSwap", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTradeSwapOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	gethTradeID := -1
	gethSwapID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveGethTradeSwap(mock, &gethTradeID, &gethSwapID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTradeSwapOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := -1
	gethSwapID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_trade_swaps").WithArgs(gethTradeID, gethSwapID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethTradeSwap(mock, &gethTradeID, &gethSwapID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeSwapList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTradeSwap{TestData1GethTradeSwap, TestData2GethTradeSwap}
	mockRows := AddGethTradeSwapToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WillReturnRows(mockRows)
	gethTradeIds := []int{1, 2}
	gethSwapIDs := []int{1, 2}
	foundGethTradeSwapList, err := GetGethTradeSwapList(mock, gethTradeIds, gethSwapIDs)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeSwapList", err)
	}
	testMarketDataList := TestAllDataGethTradeSwap
	for i, foundGethTradeSwap := range foundGethTradeSwapList {
		if cmp.Equal(foundGethTradeSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTradeSwap From Method GetGethTradeSwapList: %v is different from actual %v", foundGethTradeSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeSwapListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	gethTradeIds := []int{}
	gethSwapIDs := []int{}
	foundGethTradeSwapList, err := GetGethTradeSwapList(mock, gethTradeIds, gethSwapIDs)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeSwapList", err)
	}
	if len(foundGethTradeSwapList) != 0 {
		t.Errorf("Expected From Method GetGethTradeSwapList: to be empty but got this: %v", foundGethTradeSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeSwapListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeIds := []int{}
	gethSwapIDs := []int{}

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WillReturnRows(differentModelRows)
	foundGethTradeList, err := GetGethTradeSwapList(mock, gethTradeIds, gethSwapIDs)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeSwapList", err)
	}
	if foundGethTradeList != nil {
		t.Errorf("Expected foundGethTradeList From Method GetGethTradeSwapList: to be empty but got this: %v", foundGethTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTradeSwap(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeSwap
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_trade_swaps").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.UpdatedBy,     //4
		targetData.GethTradeID,   //5
		targetData.GethSwapID,    //6
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethTradeSwap(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethTradeSwap", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTradeSwapOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeSwap
	// name can't be nil
	targetData.Name = ""
	targetData.GethTradeID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethTradeSwap(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTradeSwapOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeSwap
	targetData.GethSwapID = nil
	err = UpdateGethTradeSwap(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTradeSwapOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeSwap
	// name can't be nil
	targetData.Name = ""
	targetData.GethTradeID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_trade_swaps").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.UpdatedBy,     //4
		targetData.GethTradeID,   //5
		targetData.GethSwapID,    //6
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethTradeSwap(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeSwap(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeSwap
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_trade_swaps").WithArgs(
		targetData.GethTradeID,   //1
		targetData.GethSwapID,    //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.Description,   //5
		targetData.CreatedBy,     //6
	).WillReturnRows(pgxmock.NewRows([]string{"geth_trade_id", "geth_swap_id"}).AddRow(1, 2))
	mock.ExpectCommit()
	gethTradeID, gethSwapID, err := InsertGethTradeSwap(mock, &targetData)
	if gethTradeID < 0 {
		t.Fatalf("gethTradeID should not be negative ID: %d", gethTradeID)
	}
	if gethSwapID < 0 {
		t.Fatalf("gethSwapID should not be negative ID: %d", gethSwapID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTradeSwap", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeSwapOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeSwap
	targetData.GethSwapID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertGethTradeSwap(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeSwapOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeSwap
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_trade_swaps").WithArgs(
		targetData.GethTradeID,   //1
		targetData.GethSwapID,    //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.Description,   //5
		targetData.CreatedBy,     //6
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethTradeID, gethSwapID, err := InsertGethTradeSwap(mock, &targetData)
	if gethTradeID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethTradeID: %d", gethTradeID)
	}
	if gethSwapID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethSwapID: %d", gethSwapID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeSwapOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1GethTradeSwap
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_trade_swaps").WithArgs(
		targetData.GethTradeID,   //1
		targetData.GethSwapID,    //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.Description,   //5
		targetData.CreatedBy,     //6
	).WillReturnRows(pgxmock.NewRows([]string{"geth_trade_id", "geth_swap_id"}).AddRow(1, 2))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethTradeID, gethSwapID, err := InsertGethTradeSwap(mock, &targetData)
	if gethTradeID >= 0 {
		t.Fatalf("Expecting -1 for gethTradeID because of error gethTradeID: %d", gethTradeID)
	}
	if gethSwapID >= 0 {
		t.Fatalf("Expecting -1 for gethSwapID because of error gethSwapID: %d", gethSwapID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeSwaps(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_trade_swaps"}, DBColumnsInsertGethTradeSwaps)
	targetData := TestAllDataGethTradeSwap
	err = InsertGethTradeSwaps(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTradeSwaps", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeSwapsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_trade_swaps"}, DBColumnsInsertGethTradeSwaps).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllDataGethTradeSwap
	err = InsertGethTradeSwaps(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMissingTradesFromSwapsByBaseAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []gethlyleswaps.GethSwap{gethlyleswaps.TestData1, gethlyleswaps.TestData2}
	mockRows := gethlyleswaps.AddGethSwapToMockRows(mock, dataList)
	baseAssetID := 1
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID).WillReturnRows(mockRows)
	foundGethTradeSwapList, err := GetMissingTradesFromSwapsByBaseAssetID(mock, &baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMissingTradesFromSwapsByBaseAssetID", err)
	}
	testMarketDataList := gethlyleswaps.TestAllData
	for i, foundGethTradeSwap := range foundGethTradeSwapList {
		if cmp.Equal(foundGethTradeSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTradeSwap From Method GetMissingTradesFromSwapsByBaseAssetID: %v is different from actual %v", foundGethTradeSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMissingTradesFromSwapsByBaseAssetIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})

	foundGethTradeSwapList, err := GetMissingTradesFromSwapsByBaseAssetID(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMissingTradesFromSwapsByBaseAssetID", err)
	}
	if len(foundGethTradeSwapList) != 0 {
		t.Errorf("Expected From Method GetMissingTradesFromSwapsByBaseAssetID: to be empty but got this: %v", foundGethTradeSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMissingTradesFromSwapsByBaseAssetIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID).WillReturnRows(differentModelRows)
	foundGethTradeSwapList, err := GetMissingTradesFromSwapsByBaseAssetID(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMissingTradesFromSwapsByBaseAssetID", err)
	}
	if foundGethTradeSwapList != nil {
		t.Errorf("Expected foundGethTradeSwapList From Method GetMissingTradesFromSwapsByBaseAssetID: to be empty but got this: %v", foundGethTradeSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMissingTxnHashesFromSwapsByBaseAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	maxBlockNumber := uint64(10000)
	mockRows := mock.NewRows([]string{"txn_hash"}).AddRow(gethlyleswaps.TestData1.TxnHash).AddRow(gethlyleswaps.TestData2.TxnHash)

	baseAssetID := 1
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID, maxBlockNumber).WillReturnRows(mockRows)
	foundGethTradeSwapList, err := GetMissingTxnHashesFromSwapsByBaseAssetID(mock, &baseAssetID, &maxBlockNumber)
	if err != nil {
		t.Fatalf("an error '%s' in GetMissingTxnHashesFromSwapsByBaseAssetID", err)
	}
	results := []string{gethlyleswaps.TestData1.TxnHash, gethlyleswaps.TestData2.TxnHash}
	for i, result := range results {
		if cmp.Equal(result, foundGethTradeSwapList[i]) == false {
			t.Errorf("Expected txnHash From Method GetMissingTxnHashesFromSwapsByBaseAssetID: %v is different from actual %v", result, foundGethTradeSwapList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMissingTxnHashesFromSwapsByBaseAssetIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	maxBlockNumber := uint64(10000)
	baseAssetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID, maxBlockNumber).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeSwapList, err := GetMissingTxnHashesFromSwapsByBaseAssetID(mock, &baseAssetID, &maxBlockNumber)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMissingTxnHashesFromSwapsByBaseAssetID", err)
	}
	if len(foundGethTradeSwapList) != 0 {
		t.Errorf("Expected From Method GetMissingTxnHashesFromSwapsByBaseAssetID: to be empty but got this: %v", foundGethTradeSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinMaxBlocksOfMissingSwapByBaseAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minBlockNumber := uint64(10000)
	maxBlockNumber := uint64(10000)
	mockRows := mock.NewRows([]string{"min_block_number", "max_block_number"}).AddRow(minBlockNumber, maxBlockNumber)

	baseAssetID := 1
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID).WillReturnRows(mockRows)
	foundMinBlockNumber, foundMaxBlockNumber, err := GetMinMaxBlocksOfMissingSwapByBaseAssetID(mock, &baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMinMaxBlocksOfMissingSwapByBaseAssetID", err)
	}
	if cmp.Equal(minBlockNumber, *foundMinBlockNumber) == false {
		t.Errorf("Expected minBlockNumber From Method GetMinMaxBlocksOfMissingSwapByBaseAssetID: %v is different from actual %v", minBlockNumber, foundMinBlockNumber)
	}
	if cmp.Equal(maxBlockNumber, *foundMaxBlockNumber) == false {
		t.Errorf("Expected maxBlockNumber From Method GetMinMaxBlocksOfMissingSwapByBaseAssetID: %v is different from actual %v", maxBlockNumber, foundMaxBlockNumber)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinMaxBlocksOfMissingSwapByBaseAssetIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMinBlockNumber, foundMaxBlockNumber, err := GetMinMaxBlocksOfMissingSwapByBaseAssetID(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMinMaxBlocksOfMissingSwapByBaseAssetID", err)
	}
	if foundMinBlockNumber != nil {
		t.Errorf("Expected From Method GetMinMaxBlocksOfMissingSwapByBaseAssetID: to be nil but got this: %v", foundMinBlockNumber)
	}
	if foundMaxBlockNumber != nil {
		t.Errorf("Expected From Method GetMinMaxBlocksOfMissingSwapByBaseAssetID: to be nil but got this: %v", foundMaxBlockNumber)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetFirstNonProcessedSwapBlockNumberForTrades(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startingBlockNumber := uint64(10000)
	mockRows := mock.NewRows([]string{"starting_block_number"}).AddRow(startingBlockNumber)
	baseAssetID := 1
	mock.ExpectQuery("^WITH max_existing_block_swaps as").WithArgs(baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID).WillReturnRows(mockRows)
	foundStartingBlockNumber, err := GetFirstNonProcessedSwapBlockNumberForTrades(mock, &baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetFirstNonProcessedSwapBlockNumberForTrades", err)
	}
	if cmp.Equal(startingBlockNumber, *foundStartingBlockNumber) == false {
		t.Errorf("Expected startingBlockNumber From Method GetMinMaxBlocksOfMissingSwapByBaseAssetID: %v is different from actual %v", startingBlockNumber, foundStartingBlockNumber)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetFirstNonProcessedSwapBlockNumberForTradesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	mock.ExpectQuery("^WITH max_existing_block_swaps as").WithArgs(baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStartingBlockNumber, err := GetFirstNonProcessedSwapBlockNumberForTrades(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetFirstNonProcessedSwapBlockNumberForTrades", err)
	}
	if foundStartingBlockNumber != nil {
		t.Errorf("Expected From Method foundStartingBlockNumber: to be nil but got this: %v", foundStartingBlockNumber)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeSwapListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllDataGethTradeSwap
	mockRows := AddGethTradeSwapToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"name = 'test", "alternate_name='test'"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WillReturnRows(mockRows)
	foundGethTradeSwapList, err := GetGethTradeSwapListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeSwapListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundGethTradeSwapList[i]) == false {
			t.Errorf("Expected sourceData From Method GetGethTradeSwapListByPagination: %v is different from actual %v", sourceData, foundGethTradeSwapList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeSwapListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"geth_swap_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeSwapList, err := GetGethTradeSwapListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeSwapListByPagination", err)
	}
	if len(foundGethTradeSwapList) != 0 {
		t.Errorf("Expected From Method GetGethTradeSwapListByPagination: to be empty but got this: %v", foundGethTradeSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeSwapListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"geth_swap_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WillReturnRows(differentModelRows)
	foundGethTradeSwapList, err := GetGethTradeSwapListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeSwapListByPagination", err)
	}
	if foundGethTradeSwapList != nil {
		t.Errorf("Expected From Method GetGethTradeSwapListByPagination: to be empty but got this: %v", foundGethTradeSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethTradeSwapCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_trade_swaps").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalGethTradeSwapCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalGethTradeSwapCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalGethTradeSwapCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethTradeSwapCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_trade_swaps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalGethTradeSwapCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalGethTradeSwapCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalGethTradeSwapCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
