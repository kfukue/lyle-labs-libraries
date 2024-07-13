package gethlyleaddresses

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

func TestGetGethAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethAddress{targetData}
	exchangeID := targetData.ID
	mockRows := AddGethAddressToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(*exchangeID).WillReturnRows(mockRows)
	foundGethAddress, err := GetGethAddress(mock, exchangeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethAddress", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method GetGethAddress: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	exchangeID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(exchangeID).WillReturnRows(noRows)
	foundGethAddress, err := GetGethAddress(mock, &exchangeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethAddress", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method GetGethAddress: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	exchangeID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(exchangeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethAddress, err := GetGethAddress(mock, &exchangeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethAddress", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method GetGethAddress: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressByAddressStr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethAddress{targetData}
	addressStr := targetData.AddressStr
	mockRows := AddGethAddressToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(addressStr).WillReturnRows(mockRows)
	foundGethAddress, err := GetGethAddressByAddressStr(mock, addressStr)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethAddressByAddressStr", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method GetGethAddressByAddressStr: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressByAddressStrForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	addressStr := "0x1234567894561234567898456121345678987456"
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(addressStr).WillReturnRows(noRows)
	foundGethAddress, err := GetGethAddressByAddressStr(mock, addressStr)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethAddressByAddressStr", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method GetGethAddressByAddressStr: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressByAddressStrForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	addressStr := "0xInvalide"
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(addressStr).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethAddress, err := GetGethAddressByAddressStr(mock, addressStr)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethAddressByAddressStr", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method GetGethAddressByAddressStr: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethAddress{TestData1, TestData2}
	mockRows := AddGethAddressToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WillReturnRows(mockRows)
	foundGethAddresses, err := GetGethAddressList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethAddressList", err)
	}
	testGethAddresses := TestAllData
	for i, foundGethAddress := range foundGethAddresses {
		if cmp.Equal(foundGethAddress, testGethAddresses[i]) == false {
			t.Errorf("Expected GethAddress From Method GetGethAddressList: %v is different from actual %v", foundGethAddress, testGethAddresses[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethAddresses, err := GetGethAddressList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethAddressList", err)
	}
	if len(foundGethAddresses) != 0 {
		t.Errorf("Expected From Method GetGethAddresses: to be empty but got this: %v", foundGethAddresses)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressListByAddressStr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethAddress{TestData1, TestData2}
	mockRows := AddGethAddressToMockRows(mock, dataList)
	addressStrList := []string{TestData1.AddressStr, TestData2.AddressStr}
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(pq.Array(addressStrList)).WillReturnRows(mockRows)
	foundGethAddresses, err := GetGethAddressListByAddressStr(mock, addressStrList)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethAddressListByAddressStr", err)
	}
	testGethAddresses := TestAllData
	for i, foundGethAddress := range foundGethAddresses {
		if cmp.Equal(foundGethAddress, testGethAddresses[i]) == false {
			t.Errorf("Expected GethAddress From Method GetGethAddressListByAddressStr: %v is different from actual %v", foundGethAddress, testGethAddresses[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressListByAddressStrForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	addressStrList := []string{"0x", "0x1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(pq.Array(addressStrList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethAddresses, err := GetGethAddressListByAddressStr(mock, addressStrList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethAddressListByAddressStr", err)
	}
	if len(foundGethAddresses) != 0 {
		t.Errorf("Expected From Method GetGethAddressListByAddressStr: to be empty but got this: %v", foundGethAddresses)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressListByIds(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethAddress{TestData1, TestData2}
	mockRows := AddGethAddressToMockRows(mock, dataList)
	ids := []int{*TestData1.ID, *TestData2.ID}
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(pq.Array(ids)).WillReturnRows(mockRows)
	foundGethAddresses, err := GetGethAddressListByIds(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethAddressListByIds", err)
	}
	testGethAddresses := TestAllData
	for i, foundGethAddress := range foundGethAddresses {
		if cmp.Equal(foundGethAddress, testGethAddresses[i]) == false {
			t.Errorf("Expected GethAddress From Method GetGethAddressListByIds: %v is different from actual %v", foundGethAddress, testGethAddresses[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressListByIdsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	ids := []int{-1, -2}
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(pq.Array(ids)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethAddresses, err := GetGethAddressListByIds(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethAddressListByIds", err)
	}
	if len(foundGethAddresses) != 0 {
		t.Errorf("Expected From Method GetGethAddressListByIds: to be empty but got this: %v", foundGethAddresses)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	exchangeID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_addresses").WithArgs(*exchangeID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethAddress(mock, exchangeID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethAddress", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethAddressOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	exchangeID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_addresses").WithArgs(exchangeID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethAddress(mock, &exchangeID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_addresses").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.AddressStr,    //4
		targetData.AddressTypeID, //5
		targetData.UpdatedBy,     //6
		targetData.ID,            //7

	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethAddress(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethAddress", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethAddressOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethAddress(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethAddressOnFailure(t *testing.T) {
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
	mock.ExpectExec("^UPDATE geth_addresses").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.AddressStr,    //4
		targetData.AddressTypeID, //5
		targetData.UpdatedBy,     //6
		targetData.ID,            //7
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethAddress(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_addresses").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.AddressStr,    //4
		targetData.AddressTypeID, //5
		targetData.CreatedBy,     //6
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	gethAddressID, err := InsertGethAddress(mock, &targetData)
	if gethAddressID < 0 {
		t.Fatalf("gethAddressID should not be negative ID: %d", gethAddressID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethAddress", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethAddressOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_addresses").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.AddressStr,    //4
		targetData.AddressTypeID, //5
		targetData.CreatedBy,     //6
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethAddressID, err := InsertGethAddress(mock, &targetData)
	if gethAddressID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethAddressID: %d", gethAddressID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethAddressOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_addresses").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.AddressStr,    //4
		targetData.AddressTypeID, //5
		targetData.CreatedBy,     //6
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethAddressID, err := InsertGethAddress(mock, &targetData)
	if gethAddressID >= 0 {
		t.Fatalf("Expecting -1 for gethAddressID because of error gethAddressID: %d", gethAddressID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethAddressList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_addresses"}, DBColumnsInsertGethAddressList)
	targetData := TestAllData
	err = InsertGethAddressList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethAddressList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethAddressListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_addresses"}, DBColumnsInsertGethAddressList).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertGethAddressList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethAddressToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"exchange_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WillReturnRows(mockRows)
	foundGethAddressList, err := GetGethAddressListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethAddressListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundGethAddressList[i]) == false {
			t.Errorf("Expected foundGethAddressList From Method GetGethAddressListByPagination: %v is different from actual %v", sourceData, foundGethAddressList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethAddressListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"exchange_type_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethAddressList, err := GetGethAddressListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethAddressListByPagination", err)
	}
	if len(foundGethAddressList) != 0 {
		t.Errorf("Expected From Method GetGethAddressListByPagination: to be empty but got this: %v", foundGethAddressList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethAddressCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_addresses").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalGethAddressCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalGethAddressCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalGethAddressCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethAddressCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_addresses").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalGethAddressCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalGethAddressCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalGethAddressCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
