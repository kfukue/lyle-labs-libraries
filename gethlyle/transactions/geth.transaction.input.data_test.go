package gethlyletransactions

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
	"id",                //1
	"uuid",              //2
	"name",              //3
	"alternate_name",    //4
	"function_name",     //5
	"method_id_str",     //6
	"num_of_parameters", //7
	"description",       //8
	"created_by",        //9
	"created_at",        //10
	"updated_by",        //11
	"updated_at",        //12
}
var DBColumnsInsertGethTransactionInputs = []string{
	"uuid",              //1
	"name",              //2
	"alternate_name",    //3
	"function_name",     //4
	"method_id_str",     //5
	"num_of_parameters", //6
	"description",       //7
	"created_by",        //8
	"created_at",        //9
	"updated_by",        //10
	"updated_at",        //11
}

var TestData1TransactionInput = GethTransactionInput{
	ID:              utils.Ptr[int](1),
	UUID:            "01ef85e8-2c26-441e-8c7f-71d79518ad72",
	Name:            "Meow Miner Compound",
	AlternateName:   "compound",
	FunctionName:    "compound(address ref)",
	MethodIDStr:     "0x284dac23",
	NumOfParameters: utils.Ptr[int](2),
	Description:     "https://snowtrace.io/tx/0x2bfe5a32b4a4685371634f1dbc35515f522cd9ae6d1e400136555e76cc8dbe3e?chainId=43114",
	CreatedBy:       "SYSTEM",
	CreatedAt:       utils.SampleCreatedAtTime,
	UpdatedBy:       "SYSTEM",
	UpdatedAt:       utils.SampleCreatedAtTime,
}

var TestData2TransactionInput = GethTransactionInput{
	ID:              utils.Ptr[int](2),
	UUID:            "bd45f8f9-84e5-4052-a895-09f2e5993622",
	Name:            "Meow Miner Deposit Meow",
	AlternateName:   "depositMEOW",
	FunctionName:    "depositMEOW(uint256 amount, address ref)",
	MethodIDStr:     "0x8e424bae",
	NumOfParameters: utils.Ptr[int](2),
	Description:     "https://snowtrace.io/tx/0xfb243aad98148639224bc41818031ccafec4cc19255fca81ce46161701be5e3f?chainId=43114",
	CreatedBy:       "SYSTEM",
	CreatedAt:       utils.SampleCreatedAtTime,
	UpdatedBy:       "SYSTEM",
	UpdatedAt:       utils.SampleCreatedAtTime,
}
var TestAllDataTransactionInput = []GethTransactionInput{TestData1TransactionInput, TestData2TransactionInput}

func AddGethTransactionInputToMockRows(mock pgxmock.PgxPoolIface, dataList []GethTransactionInput) *pgxmock.Rows {
	rows := mock.NewRows(DBColumnsTransactionInputs)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,              //1
			data.UUID,            //2
			data.Name,            //3
			data.AlternateName,   //4
			data.FunctionName,    //5
			data.MethodIDStr,     //6
			data.NumOfParameters, //7
			data.Description,     //8
			data.CreatedBy,       //9
			data.CreatedAt,       //10
			data.UpdatedBy,       //11
			data.UpdatedAt,       //12
		)
	}
	return rows
}

func TestGetGethTransactionInput(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2TransactionInput
	dataList := []GethTransactionInput{targetData}
	gethTransactionInputID := targetData.ID
	mockRows := AddGethTransactionInputToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs").WithArgs(*gethTransactionInputID).WillReturnRows(mockRows)
	foundGethTransactionInput, err := GetGethTransactionInput(mock, gethTransactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionInput", err)
	}
	if cmp.Equal(*foundGethTransactionInput, targetData) == false {
		t.Errorf("Expected GethTransactionInput From Method GetGethTransactionInput: %v is different from actual %v", foundGethTransactionInput, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionInputForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTransactionInputID := 999
	noRows := pgxmock.NewRows(DBColumnsTransactionInputs)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs").WithArgs(gethTransactionInputID).WillReturnRows(noRows)
	foundGethTransactionInput, err := GetGethTransactionInput(mock, &gethTransactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionInput", err)
	}
	if foundGethTransactionInput != nil {
		t.Errorf("Expected GethTransactionInput From Method GetGethTransactionInput: to be empty but got this: %v", foundGethTransactionInput)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionInputForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTransactionInputID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs").WithArgs(gethTransactionInputID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionInput, err := GetGethTransactionInput(mock, &gethTransactionInputID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionInput", err)
	}
	if foundGethTransactionInput != nil {
		t.Errorf("Expected GethTransactionInput From Method GetGethTransactionInput: to be empty but got this: %v", foundGethTransactionInput)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionInputByFromToAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTransactionInput{TestData1TransactionInput, TestData2TransactionInput}
	mockRows := AddGethTransactionInputToMockRows(mock, dataList)
	fromToAddressID := 1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs").WithArgs(fromToAddressID).WillReturnRows(mockRows)
	foundGethTransactionInputList, err := GetGethTransactionInputByFromToAddress(mock, &fromToAddressID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionInputByFromToAddress", err)
	}
	testMarketDataList := TestAllDataTransactionInput
	for i, foundGethTransactionInput := range foundGethTransactionInputList {
		if cmp.Equal(foundGethTransactionInput, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransactionInput From Method GetGethTransactionInputByFromToAddress: %v is different from actual %v", foundGethTransactionInput, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionInputByFromToAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	fromToAddressID := 1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs").WithArgs(fromToAddressID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionInputList, err := GetGethTransactionInputByFromToAddress(mock, &fromToAddressID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionInputByFromToAddress", err)
	}
	if len(foundGethTransactionInputList) != 0 {
		t.Errorf("Expected From Method GetGethTransactionInputByFromToAddress: to be empty but got this: %v", foundGethTransactionInputList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransactionInput(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1TransactionInput
	transactionInputID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transaction_inputs").WithArgs(*transactionInputID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethTransactionInput(mock, transactionInputID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethTransactionInput", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransactionInputOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionInputID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transaction_inputs").WithArgs(transactionInputID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethTransactionInput(mock, &transactionInputID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionInputList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTransactionInput{TestData1TransactionInput, TestData2TransactionInput}
	mockRows := AddGethTransactionInputToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs").WillReturnRows(mockRows)
	foundGethTransactionInputList, err := GetGethTransactionInputList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionInputList", err)
	}
	testMarketDataList := TestAllDataTransactionInput
	for i, foundGethTransactionInput := range foundGethTransactionInputList {
		if cmp.Equal(foundGethTransactionInput, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransactionInput From Method GetGethTransactionInputList: %v is different from actual %v", foundGethTransactionInput, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionInputListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionInputList, err := GetGethTransactionInputList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionInputList", err)
	}
	if len(foundGethTransactionInputList) != 0 {
		t.Errorf("Expected From Method GetGethTransactionInputList: to be empty but got this: %v", foundGethTransactionInputList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTransactionInput(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1TransactionInput
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transaction_inputs").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.FunctionName,    //3
		targetData.MethodIDStr,     //4
		targetData.NumOfParameters, //5
		targetData.Description,     //6
		targetData.UpdatedBy,       //7
		targetData.ID,              //8
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethTransactionInput(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethTransactionInput", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransactionInputOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1TransactionInput
	// name can't be nil
	targetData.Name = ""
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethTransactionInput(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransactionInputOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1TransactionInput
	// name can't be nil
	targetData.Name = ""
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transaction_inputs").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.FunctionName,    //3
		targetData.MethodIDStr,     //4
		targetData.NumOfParameters, //5
		targetData.Description,     //6
		targetData.UpdatedBy,       //7
		targetData.ID,              //8
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethTransactionInput(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransactionInput(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1TransactionInput
	targetData.Name = "New Name"
	testUUID := "21bc560d-40f9-4f8a-a5e7-f03720fe0e0d"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_transaction_inputs").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.FunctionName,    //3
		targetData.MethodIDStr,     //4
		targetData.NumOfParameters, //5
		targetData.Description,     //6
		targetData.CreatedBy,       //7
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, testUUID))
	mock.ExpectCommit()
	gethMinerID, newUUID, err := InsertGethTransactionInput(mock, &targetData)
	if gethMinerID < 0 {
		t.Fatalf("gethMinerID should not be negative ID: %d", gethMinerID)
	}
	if newUUID == "" {
		t.Fatalf("newUUID should not be empty ID: %s", newUUID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTransactionInput", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransactionInputOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1TransactionInput
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_transaction_inputs").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.FunctionName,    //3
		targetData.MethodIDStr,     //4
		targetData.NumOfParameters, //5
		targetData.Description,     //6
		targetData.CreatedBy,       //7
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethMinerID, newUUID, err := InsertGethTransactionInput(mock, &targetData)
	if gethMinerID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethMinerID: %d", gethMinerID)
	}
	if newUUID != "" {
		t.Fatalf("Expecting empty string for newUUID because of error newUUID: %s", newUUID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransactionInputOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1TransactionInput
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_transaction_inputs").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.FunctionName,    //3
		targetData.MethodIDStr,     //4
		targetData.NumOfParameters, //5
		targetData.Description,     //6
		targetData.CreatedBy,       //7
	).WillReturnRows(pgxmock.NewRows([]string{"id", "transaction_input_id"}).AddRow(1, 2))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethMinerID, newUUID, err := InsertGethTransactionInput(mock, &targetData)
	if gethMinerID >= 0 {
		t.Fatalf("Expecting -1 for gethMinerID because of error gethMinerID: %d", gethMinerID)
	}
	if newUUID != "" {
		t.Fatalf("Expecting empty string for newUUID because of error newUUID: %s", newUUID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransactionInputs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_transaction_inputs"}, DBColumnsInsertGethTransactionInputs)
	targetData := TestAllDataTransactionInput
	err = InsertGethTransactionInputs(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTransactionInputs", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransactionInputsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_transaction_inputs"}, DBColumnsInsertGethTransactionInputs).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllDataTransactionInput
	err = InsertGethTransactionInputs(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionInputListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllDataTransactionInput
	mockRows := AddGethTransactionInputToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs").WillReturnRows(mockRows)
	foundChains, err := GetTransactionInputListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionInputListByPagination", err)
	}
	testChains := dataList
	for i, foundChain := range foundChains {
		if cmp.Equal(foundChain, testChains[i]) == false {
			t.Errorf("Expected Chain From Method GetTransactionInputListByPagination: %v is different from actual %v", foundChain, testChains[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionInputListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChains, err := GetTransactionInputListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionInputListByPagination", err)
	}
	if len(foundChains) != 0 {
		t.Errorf("Expected From Method GetTransactionInputListByPagination: to be empty but got this: %v", foundChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalTransactionInputsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_transaction_inputs").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalTransactionInputsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalTransactionInputsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalTransactionInputsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalTransactionInputsCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_transaction_inputs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalTransactionInputsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalTransactionInputsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalTransactionInputsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionInputByFromMinerID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTransactionInput{TestData1TransactionInput, TestData2TransactionInput}
	mockRows := AddGethTransactionInputToMockRows(mock, dataList)
	minerID := 1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs gti JOIN geth_miners_transaction_inputs gmti").WithArgs(minerID).WillReturnRows(mockRows)
	foundGethTransactionInputList, err := GetGethTransactionInputByFromMinerID(mock, &minerID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionInputByFromMinerID", err)
	}
	testMarketDataList := TestAllDataTransactionInput
	for i, foundGethTransactionInput := range foundGethTransactionInputList {
		if cmp.Equal(foundGethTransactionInput, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransactionInput From Method GetGethTransactionInputByFromMinerID: %v is different from actual %v", foundGethTransactionInput, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionInputByFromMinerIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transaction_inputs gti JOIN geth_miners_transaction_inputs gmti").WithArgs(minerID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionInputList, err := GetGethTransactionInputByFromMinerID(mock, &minerID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionInputByFromMinerID", err)
	}
	if len(foundGethTransactionInputList) != 0 {
		t.Errorf("Expected From Method GetGethTransactionInputByFromMinerID: to be empty but got this: %v", foundGethTransactionInputList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
