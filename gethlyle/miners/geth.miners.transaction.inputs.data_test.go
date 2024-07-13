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

var DBColumnsTransactionInputs = []string{
	"miner_id",             //1
	"transaction_input_id", //2
	"uuid",                 //3
	"name",                 //4
	"alternate_name",       //5
	"description",          //6
	"created_by",           //7
	"created_at",           //8
	"updated_by",           //9
	"updated_at",           //10
}
var DBColumnsInsertGethMinersTransactionInputs = []string{
	"miner_id",             //1
	"transaction_input_id", //2
	"uuid",                 //3
	"name",                 //4
	"alternate_name",       //5
	"description",          //6
	"created_by",           //7
	"created_at",           //8
	"updated_by",           //9
	"updated_at",           //10
}

var TestData1MinerTransactionInput = GethMinerTransactionInput{
	MinerID:            utils.Ptr[int](1),
	TransactionInputID: utils.Ptr[int](2),
	UUID:               "01ef85e8-2c26-441e-8c7f-71d79518ad72",
	Name:               "Meow Miner Compound",
	AlternateName:      "compound",
	Description:        "",
	CreatedBy:          "SYSTEM",
	CreatedAt:          utils.SampleCreatedAtTime,
	UpdatedBy:          "SYSTEM",
	UpdatedAt:          utils.SampleCreatedAtTime,
}

var TestData2MinerTransactionInput = GethMinerTransactionInput{
	MinerID:            utils.Ptr[int](1),
	TransactionInputID: utils.Ptr[int](2),
	UUID:               "bd45f8f9-84e5-4052-a895-09f2e5993622",
	Name:               "Meow Miner Deposit Meow",
	AlternateName:      "depositMEOW",
	Description:        "",
	CreatedBy:          "SYSTEM",
	CreatedAt:          utils.SampleCreatedAtTime,
	UpdatedBy:          "SYSTEM",
	UpdatedAt:          utils.SampleCreatedAtTime,
}
var TestAllDataMinerTransactionInput = []GethMinerTransactionInput{TestData1MinerTransactionInput, TestData2MinerTransactionInput}

func AddGethMinerTransactionInputToMockRows(mock pgxmock.PgxPoolIface, dataList []GethMinerTransactionInput) *pgxmock.Rows {
	rows := mock.NewRows(DBColumnsTransactionInputs)
	for _, data := range dataList {
		rows.AddRow(
			data.MinerID,            //1
			data.TransactionInputID, //2
			data.UUID,               //3
			data.Name,               //4
			data.AlternateName,      //5
			data.Description,        //6
			data.CreatedBy,          //7
			data.CreatedAt,          //8
			data.UpdatedBy,          //9
			data.UpdatedAt,          //10
		)
	}
	return rows
}

func TestGetAllGethMinerTransactionInputsByMinerID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMinerTransactionInput{TestData1MinerTransactionInput}
	mockRows := AddGethMinerTransactionInputToMockRows(mock, dataList)
	minerID := TestData1MinerTransactionInput.MinerID
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WithArgs(*minerID).WillReturnRows(mockRows)
	foundGethMinerTransactionInputList, err := GetAllGethMinerTransactionInputsByMinerID(mock, minerID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllGethMinerTransactionInputsByMinerID", err)
	}
	testMarketDataList := dataList
	for i, foundGethMinerTransactionInput := range foundGethMinerTransactionInputList {
		if cmp.Equal(foundGethMinerTransactionInput, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMinerTransactionInput From Method GetAllGethMinerTransactionInputsByMinerID: %v is different from actual %v", foundGethMinerTransactionInput, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethMinerTransactionInputsByMinerIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := -999
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WithArgs(minerID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethMinerTransactionInputList, err := GetAllGethMinerTransactionInputsByMinerID(mock, &minerID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethMinerTransactionInputsByMinerID", err)
	}
	if len(foundGethMinerTransactionInputList) != 0 {
		t.Errorf("Expected From Method GetAllGethMinerTransactionInputsByMinerID: to be empty but got this: %v", foundGethMinerTransactionInputList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethMinerTransactionInputsByTransactionInputID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMinerTransactionInput{TestData1MinerTransactionInput}
	mockRows := AddGethMinerTransactionInputToMockRows(mock, dataList)
	transactionInputID := TestData1MinerTransactionInput.TransactionInputID
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WithArgs(*transactionInputID).WillReturnRows(mockRows)
	foundGethMinerTransactionInputList, err := GetAllGethMinerTransactionInputsByTransactionInputID(mock, transactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllGethMinerTransactionInputsByTransactionInputID", err)
	}
	testMarketDataList := dataList
	for i, foundGethMinerTransactionInput := range foundGethMinerTransactionInputList {
		if cmp.Equal(foundGethMinerTransactionInput, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMinerTransactionInput From Method GetAllGethMinerTransactionInputsByTransactionInputID: %v is different from actual %v", foundGethMinerTransactionInput, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethMinerTransactionInputsByTransactionInputIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionInputID := -999
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WithArgs(transactionInputID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethMinerTransactionInputList, err := GetAllGethMinerTransactionInputsByTransactionInputID(mock, &transactionInputID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethMinerTransactionInputsByTransactionInputID", err)
	}
	if len(foundGethMinerTransactionInputList) != 0 {
		t.Errorf("Expected From Method GetAllGethMinerTransactionInputsByTransactionInputID: to be empty but got this: %v", foundGethMinerTransactionInputList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerTransactionInput(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2MinerTransactionInput
	dataList := []GethMinerTransactionInput{targetData}
	gethMinerID := targetData.MinerID
	transactionInputID := targetData.TransactionInputID
	mockRows := AddGethMinerTransactionInputToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WithArgs(*gethMinerID, *transactionInputID).WillReturnRows(mockRows)
	foundGethMinerTransactionInput, err := GetGethMinerTransactionInput(mock, gethMinerID, transactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMinerTransactionInput", err)
	}
	if cmp.Equal(*foundGethMinerTransactionInput, targetData) == false {
		t.Errorf("Expected GethMinerTransactionInput From Method GetGethMinerTransactionInput: %v is different from actual %v", foundGethMinerTransactionInput, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerTransactionInputForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethMinerID := 999
	transactionInputID := 999
	noRows := pgxmock.NewRows(DBColumnsTransactionInputs)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WithArgs(gethMinerID, transactionInputID).WillReturnRows(noRows)
	foundGethMinerTransactionInput, err := GetGethMinerTransactionInput(mock, &gethMinerID, &transactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMinerTransactionInput", err)
	}
	if foundGethMinerTransactionInput != nil {
		t.Errorf("Expected GethMinerTransactionInput From Method GetGethMinerTransactionInput: to be empty but got this: %v", foundGethMinerTransactionInput)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerTransactionInputForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethMinerID := -1
	transactionInputID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WithArgs(gethMinerID, transactionInputID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethMinerTransactionInput, err := GetGethMinerTransactionInput(mock, &gethMinerID, &transactionInputID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMinerTransactionInput", err)
	}
	if foundGethMinerTransactionInput != nil {
		t.Errorf("Expected GethMinerTransactionInput From Method GetGethMinerTransactionInput: to be empty but got this: %v", foundGethMinerTransactionInput)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMinerTransactionInput(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransactionInput
	gethMinerID := targetData.MinerID
	transactionInputID := targetData.TransactionInputID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_miners_transaction_inputs").WithArgs(*gethMinerID, *transactionInputID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethMinerTransactionInput(mock, gethMinerID, transactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethMinerTransactionInput", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMinerTransactionInputOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethMinerID := -1
	transactionInputID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_miners_transaction_inputs").WithArgs(gethMinerID, transactionInputID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethMinerTransactionInput(mock, &gethMinerID, &transactionInputID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerTransactionInputList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMinerTransactionInput{TestData1MinerTransactionInput, TestData2MinerTransactionInput}
	mockRows := AddGethMinerTransactionInputToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WillReturnRows(mockRows)
	minerIDs := []int{}
	transactionInputIDs := []int{}
	foundGethMinerTransactionInputList, err := GetGethMinerTransactionInputList(mock, minerIDs, transactionInputIDs)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMinerTransactionInputList", err)
	}
	testMarketDataList := TestAllDataMinerTransactionInput
	for i, foundGethMinerTransactionInput := range foundGethMinerTransactionInputList {
		if cmp.Equal(foundGethMinerTransactionInput, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMinerTransactionInput From Method GetGethMinerTransactionInputList: %v is different from actual %v", foundGethMinerTransactionInput, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerTransactionInputListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	minerIDs := []int{}
	transactionInputIDs := []int{}
	foundGethMinerTransactionInputList, err := GetGethMinerTransactionInputList(mock, minerIDs, transactionInputIDs)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMinerTransactionInputList", err)
	}
	if len(foundGethMinerTransactionInputList) != 0 {
		t.Errorf("Expected From Method GetGethMinerTransactionInputList: to be empty but got this: %v", foundGethMinerTransactionInputList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethMinerTransactionInput(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransactionInput
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_miners_transaction_inputs").WithArgs(
		targetData.Name,               //1
		targetData.AlternateName,      //2
		targetData.Description,        //3
		targetData.UpdatedBy,          //4
		targetData.MinerID,            //5
		targetData.TransactionInputID, //6
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethMinerTransactionInput(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethMinerTransactionInput", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethMinerTransactionInputOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransactionInput
	// name can't be nil
	targetData.Name = ""
	targetData.MinerID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethMinerTransactionInput(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethMinerTransactionInputOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransactionInput
	// name can't be nil
	targetData.Name = ""
	targetData.MinerID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_miners_transaction_inputs").WithArgs(
		targetData.Name,               //1
		targetData.AlternateName,      //2
		targetData.Description,        //3
		targetData.UpdatedBy,          //4
		targetData.MinerID,            //5
		targetData.TransactionInputID, //6
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethMinerTransactionInput(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinerTransactionInput(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransactionInput
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_miners_transaction_inputs").WithArgs(
		targetData.MinerID,            //1
		targetData.TransactionInputID, //2
		targetData.Name,               //3
		targetData.AlternateName,      //4
		targetData.Description,        //5
		targetData.CreatedBy,          //6
	).WillReturnRows(pgxmock.NewRows([]string{"miner_id", "transaction_input_id"}).AddRow(1, 2))
	mock.ExpectCommit()
	gethMinerID, transactionInputID, err := InsertGethMinerTransactionInput(mock, &targetData)
	if gethMinerID < 0 {
		t.Fatalf("gethMinerID should not be negative ID: %d", gethMinerID)
	}
	if transactionInputID < 0 {
		t.Fatalf("transactionInputID should not be negative ID: %d", transactionInputID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethMinerTransactionInput", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinerTransactionInputOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransactionInput
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_miners_transaction_inputs").WithArgs(
		targetData.MinerID,            //1
		targetData.TransactionInputID, //2
		targetData.Name,               //3
		targetData.AlternateName,      //4
		targetData.Description,        //5
		targetData.CreatedBy,          //6
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethMinerID, transactionInputID, err := InsertGethMinerTransactionInput(mock, &targetData)
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

func TestInsertGethMinerTransactionInputOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinerTransactionInput
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_miners_transaction_inputs").WithArgs(
		targetData.MinerID,            //1
		targetData.TransactionInputID, //2
		targetData.Name,               //3
		targetData.AlternateName,      //4
		targetData.Description,        //5
		targetData.CreatedBy,          //6
	).WillReturnRows(pgxmock.NewRows([]string{"miner_id", "transaction_input_id"}).AddRow(1, 2))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethMinerID, transactionInputID, err := InsertGethMinerTransactionInput(mock, &targetData)
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

func TestInsertGethMinersTransactionInputs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_miners_transaction_inputs"}, DBColumnsInsertGethMinersTransactionInputs)
	targetData := TestAllDataMinerTransactionInput
	err = InsertGethMinersTransactionInputs(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethMinersTransactionInputs", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinersTransactionInputsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_miners_transaction_inputs"}, DBColumnsInsertGethMinersTransactionInputs).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllDataMinerTransactionInput
	err = InsertGethMinersTransactionInputs(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinerTransactionInputListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllDataMinerTransactionInput
	mockRows := AddGethMinerTransactionInputToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WillReturnRows(mockRows)
	foundGethMinerTransactionInputList, err := GetMinerTransactionInputListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetMinerTransactionInputListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundGethMinerTransactionInputList[i]) == false {
			t.Errorf("Expected sourceData From Method GetMinerTransactionInputListByPagination: %v is different from actual %v", sourceData, foundGethMinerTransactionInputList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinerTransactionInputListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transaction_inputs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethMinerTransactionInputList, err := GetMinerTransactionInputListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMinerTransactionInputListByPagination", err)
	}
	if len(foundGethMinerTransactionInputList) != 0 {
		t.Errorf("Expected From Method GetMinerTransactionInputListByPagination: to be empty but got this: %v", foundGethMinerTransactionInputList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMinerTransactionInputCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_miners_transaction_inputs").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalMinerTransactionInputCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalMinerTransactionInputCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalMinerTransactionInputCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMinerTransactionInputCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_miners_transaction_inputs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalMinerTransactionInputCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalMinerTransactionInputCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalMinerTransactionInputCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
