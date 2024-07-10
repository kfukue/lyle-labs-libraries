package gethlyleminers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumnsTransactions = []string{
	"miner_id",       //1
	"transaction_id", //2
	"uuid",           //3
	"name",           //4
	"alternate_name", //5
	"description",    //6
	"created_by",     //7
	"created_at",     //8
	"updated_by",     //9
	"updated_at",     //10
}
var DBColumnsInsertGethMinersTransactions = []string{
	"miner_id",       //1
	"transaction_id", //2
	"uuid",           //3
	"name",           //4
	"alternate_name", //5
	"description",    //6
	"created_by",     //7
	"created_at",     //8
	"updated_by",     //9
	"updated_at",     //10
}

var TestData1MinerTransaction = GethMinerTransaction{
	MinerID:       utils.Ptr[int](1),
	TransactionID: utils.Ptr[int](2),
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72",
	Name:          "Meow Miner Transction 1",
	AlternateName: "Meow Miner Transction 1",
	Description:   "",
	CreatedBy:     "SYSTEM",
	CreatedAt:     utils.SampleCreatedAtTime,
	UpdatedBy:     "SYSTEM",
	UpdatedAt:     utils.SampleCreatedAtTime,
}

var TestData2MinerTransaction = GethMinerTransaction{
	MinerID:       utils.Ptr[int](2),
	TransactionID: utils.Ptr[int](3),
	UUID:          "bd45f8f9-84e5-4052-a895-09f2e5993622",
	Name:          "Print The PEPE Transaction 1",
	AlternateName: "Print The PEPE Transaction 1",
	Description:   "",
	CreatedBy:     "SYSTEM",
	CreatedAt:     utils.SampleCreatedAtTime,
	UpdatedBy:     "SYSTEM",
	UpdatedAt:     utils.SampleCreatedAtTime,
}
var TestAllDataMinerTransaction = []GethMinerTransaction{TestData1MinerTransaction, TestData2MinerTransaction}

func AddGethMinerTransactionToMockRows(mock pgxmock.PgxPoolIface, dataList []GethMinerTransaction) *pgxmock.Rows {
	rows := mock.NewRows(DBColumnsTransactions)
	for _, data := range dataList {
		rows.AddRow(
			data.MinerID,       //1
			data.TransactionID, //2
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

func TestGetAllGethMinerTransactionsByMinerID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMinerTransaction{TestData1MinerTransaction}
	mockRows := AddGethMinerTransactionToMockRows(mock, dataList)
	minerID := TestData1MinerTransaction.MinerID
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(*minerID).WillReturnRows(mockRows)
	foundMarketDataList, err := GetAllGethMinerTransactionsByMinerID(mock, minerID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllGethMinerTransactionsByMinerID", err)
	}
	testMarketDataList := dataList
	for i, foundMarketData := range foundMarketDataList {
		if cmp.Equal(foundMarketData, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMinerTransaction From Method GetAllGethMinerTransactionsByMinerID: %v is different from actual %v", foundMarketData, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethMinerTransactionsByMinerIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := -999
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(minerID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetAllGethMinerTransactionsByMinerID(mock, &minerID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethMinerTransactionsByMinerID", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetAllGethMinerTransactionsByMinerID: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinAndMaxDatesFromTransactionsByMinerID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	rows := mock.NewRows([]string{"min_date", "max_date"})
	mockRows := rows.AddRow(&startDate, &endDate)
	minerID := 1
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(minerID).WillReturnRows(mockRows)
	startDateResult, endDateResult, err := GetMinAndMaxDatesFromTransactionsByMinerID(mock, &minerID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMinAndMaxDatesFromTransactionsByMinerID", err)
	}
	if cmp.Equal(startDate, *startDateResult) == false {
		t.Errorf("Expected startDate From Method GetMinAndMaxDatesFromTransactionsByMinerID: %v is different from actual %v", startDate, *startDateResult)
	}
	if cmp.Equal(endDate, *endDateResult) == false {
		t.Errorf("Expected endDate From Method GetMinAndMaxDatesFromTransactionsByMinerID: %v is different from actual %v", endDate, *endDateResult)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinAndMaxDatesFromTransactionsByMinerIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := 99999
	noRows := mock.NewRows([]string{"min_date", "max_date"})
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(minerID).WillReturnRows(noRows)
	startDateResult, endDateResult, err := GetMinAndMaxDatesFromTransactionsByMinerID(mock, &minerID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMinAndMaxDatesFromTransactionsByMinerID", err)
	}
	if startDateResult != nil {
		t.Errorf("startDateResult From Method GetMinAndMaxDatesFromTransactionsByMinerID: to be empty but got this: %v", startDateResult)
	}
	if endDateResult != nil {
		t.Errorf("Expected endDateResult From Method GetMinAndMaxDatesFromTransactionsByMinerID: to be empty but got this: %v", endDateResult)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinAndMaxDatesFromTransactionsByMinerIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := 99999
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(minerID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	startDateResult, endDateResult, err := GetMinAndMaxDatesFromTransactionsByMinerID(mock, &minerID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMinAndMaxDatesFromTransactionsByMinerID", err)
	}
	if startDateResult != nil {
		t.Errorf("Expected startDateResult From Method GetMinAndMaxDatesFromTransactionsByMinerID: to be empty but got this: %v", startDateResult)
	}
	if endDateResult != nil {
		t.Errorf("Expected endDateResult From Method GetMinAndMaxDatesFromTransactionsByMinerID: to be empty but got this: %v", endDateResult)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctAddressesFromGethTransactionsByMinerIDAndBeforeDate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testAddress1 := "0x5A9033Ba22210158EADB84B12Ac05Ca40d7CD01a"
	testAddress2 := "0x0FDda21e2D1088094082dE94335415896Af335f4"
	dataList := []string{testAddress1, testAddress2}
	rows := mock.NewRows([]string{"address"})
	mockRows := rows.AddRow(testAddress1).AddRow(testAddress2)
	minerID := TestData1MinerTransaction.MinerID
	beforeDate := utils.SampleCreatedAtTime
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(*minerID, beforeDate.Format(utils.LayoutPostgres)).WillReturnRows(mockRows)
	foundMarketDataList, err := GetDistinctAddressesFromGethTransactionsByMinerIDAndBeforeDate(mock, minerID, &beforeDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetDistinctAddressesFromGethTransactionsByMinerIDAndBeforeDate", err)
	}
	testMarketDataList := dataList
	for i, foundMarketData := range foundMarketDataList {
		if cmp.Equal(foundMarketData, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMinerTransaction From Method GetDistinctAddressesFromGethTransactionsByMinerIDAndBeforeDate: %v is different from actual %v", foundMarketData, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctAddressesFromGethTransactionsByMinerIDAndBeforeDateForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := -999
	beforeDate := utils.SampleCreatedAtTime
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(minerID, beforeDate.Format(utils.LayoutPostgres)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetDistinctAddressesFromGethTransactionsByMinerIDAndBeforeDate(mock, &minerID, &beforeDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDistinctAddressesFromGethTransactionsByMinerIDAndBeforeDate", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetDistinctAddressesFromGethTransactionsByMinerIDAndBeforeDate: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethMinerTransactionsByTransactionID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMinerTransaction{TestData1MinerTransaction}
	mockRows := AddGethMinerTransactionToMockRows(mock, dataList)
	transactionInputID := TestData1MinerTransaction.TransactionID
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(*transactionInputID).WillReturnRows(mockRows)
	foundMarketDataList, err := GetAllGethMinerTransactionsByTransactionID(mock, transactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllGethMinerTransactionsByTransactionID", err)
	}
	testMarketDataList := dataList
	for i, foundMarketData := range foundMarketDataList {
		if cmp.Equal(foundMarketData, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMinerTransaction From Method GetAllGethMinerTransactionsByTransactionID: %v is different from actual %v", foundMarketData, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethMinerTransactionsByTransactionIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionInputID := -999
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(transactionInputID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetAllGethMinerTransactionsByTransactionID(mock, &transactionInputID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethMinerTransactionsByTransactionID", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetAllGethMinerTransactionsByTransactionID: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2MinerTransaction
	dataList := []GethMinerTransaction{targetData}
	gethMinerID := targetData.MinerID
	transactionInputID := targetData.TransactionID
	mockRows := AddGethMinerTransactionToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(*gethMinerID, *transactionInputID).WillReturnRows(mockRows)
	foundMarketData, err := GetGethMinerTransaction(mock, gethMinerID, transactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMinerTransaction", err)
	}
	if cmp.Equal(*foundMarketData, targetData) == false {
		t.Errorf("Expected GethMinerTransaction From Method GetGethMinerTransaction: %v is different from actual %v", foundMarketData, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerTransactionForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethMinerID := 999
	transactionInputID := 999
	noRows := pgxmock.NewRows(DBColumnsTransactions)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(gethMinerID, transactionInputID).WillReturnRows(noRows)
	foundMarketData, err := GetGethMinerTransaction(mock, &gethMinerID, &transactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMinerTransaction", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected GethMinerTransaction From Method GetGethMinerTransaction: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerTransactionForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethMinerID := -1
	transactionInputID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WithArgs(gethMinerID, transactionInputID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketData, err := GetGethMinerTransaction(mock, &gethMinerID, &transactionInputID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMinerTransaction", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected GethMinerTransaction From Method GetGethMinerTransaction: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMinerTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransaction
	gethMinerID := targetData.MinerID
	transactionInputID := targetData.TransactionID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_miners_transactions").WithArgs(*gethMinerID, *transactionInputID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethMinerTransaction(mock, gethMinerID, transactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethMinerTransaction", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMinerTransactionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethMinerID := -1
	transactionInputID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_miners_transactions").WithArgs(gethMinerID, transactionInputID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethMinerTransaction(mock, &gethMinerID, &transactionInputID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerTransactionList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMinerTransaction{TestData1MinerTransaction, TestData2MinerTransaction}
	mockRows := AddGethMinerTransactionToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WillReturnRows(mockRows)
	minerIDs := []int{}
	transactionInputIDs := []int{}
	foundMarketDataList, err := GetGethMinerTransactionList(mock, minerIDs, transactionInputIDs)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMinerTransactionList", err)
	}
	testMarketDataList := TestAllDataMinerTransaction
	for i, foundMarketData := range foundMarketDataList {
		if cmp.Equal(foundMarketData, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMinerTransaction From Method GetGethMinerTransactionList: %v is different from actual %v", foundMarketData, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerTransactionListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	minerIDs := []int{}
	transactionInputIDs := []int{}
	foundMarketDataList, err := GetGethMinerTransactionList(mock, minerIDs, transactionInputIDs)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMinerTransactionList", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetGethMinerTransactionList: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethMinerTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransaction
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_miners_transactions").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.UpdatedBy,     //4
		targetData.MinerID,       //5
		targetData.TransactionID, //6
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethMinerTransaction(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethMinerTransaction", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethMinerTransactionOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransaction
	// name can't be nil
	targetData.Name = ""
	targetData.MinerID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethMinerTransaction(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethMinerTransactionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransaction
	// name can't be nil
	targetData.Name = ""
	targetData.MinerID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_miners_transactions").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.UpdatedBy,     //4
		targetData.MinerID,       //5
		targetData.TransactionID, //6
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethMinerTransaction(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinerTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransaction
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_miners_transactions").WithArgs(
		targetData.MinerID,       //1
		targetData.TransactionID, //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.Description,   //5
		targetData.CreatedBy,     //6
	).WillReturnRows(pgxmock.NewRows([]string{"id", "transaction_input_id"}).AddRow(1, 2))
	mock.ExpectCommit()
	gethMinerID, transactionInputID, err := InsertGethMinerTransaction(mock, &targetData)
	if gethMinerID < 0 {
		t.Fatalf("gethMinerID should not be negative ID: %d", gethMinerID)
	}
	if transactionInputID < 0 {
		t.Fatalf("transactionInputID should not be negative ID: %d", transactionInputID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethMinerTransaction", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinerTransactionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransaction
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_miners_transactions").WithArgs(
		targetData.MinerID,       //1
		targetData.TransactionID, //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.Description,   //5
		targetData.CreatedBy,     //6
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethMinerID, transactionInputID, err := InsertGethMinerTransaction(mock, &targetData)
	if gethMinerID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethMinerID: %d", gethMinerID)
	}
	if transactionInputID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error transactionInputID: %d", transactionInputID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinerTransactionOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransaction
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_miners_transactions").WithArgs(
		targetData.MinerID,       //1
		targetData.TransactionID, //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.Description,   //5
		targetData.CreatedBy,     //6
	).WillReturnRows(pgxmock.NewRows([]string{"id", "transaction_input_id"}).AddRow(1, 2))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethMinerID, transactionInputID, err := InsertGethMinerTransaction(mock, &targetData)
	if gethMinerID >= 0 {
		t.Fatalf("Expecting -1 for gethMinerID because of error gethMinerID: %d", gethMinerID)
	}
	if transactionInputID >= 0 {
		t.Fatalf("Expecting -1 for transactionInputID because of error transactionInputID: %d", transactionInputID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinersTransactions(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_miners_transactions"}, DBColumnsInsertGethMinersTransactions)
	targetData := TestAllDataMinerTransaction
	err = InsertGethMinersTransactions(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethMinersTransactions", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinersTransactionsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_miners_transactions"}, DBColumnsInsertGethMinersTransactions).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllDataMinerTransaction
	err = InsertGethMinersTransactions(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinerTransactionListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllDataMinerTransaction
	mockRows := AddGethMinerTransactionToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WillReturnRows(mockRows)
	foundChains, err := GetMinerTransactionListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetMinerTransactionListByPagination", err)
	}
	testChains := dataList
	for i, foundChain := range foundChains {
		if cmp.Equal(foundChain, testChains[i]) == false {
			t.Errorf("Expected Chain From Method GetMinerTransactionListByPagination: %v is different from actual %v", foundChain, testChains[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinerTransactionListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChains, err := GetMinerTransactionListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMinerTransactionListByPagination", err)
	}
	if len(foundChains) != 0 {
		t.Errorf("Expected From Method GetMinerTransactionListByPagination: to be empty but got this: %v", foundChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMinerTransactionCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_miners_transactions").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalMinerTransactionCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalMinerTransactionCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalMinerTransactionCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMinerTransactionCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_miners_transactions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalMinerTransactionCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalMinerTransactionCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalMinerTransactionCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
