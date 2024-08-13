package gethlyleaddresses

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/asset"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

func TestCreateOrGetContractAddressFromExistingAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset TestData2 is Contract
	assetData1 := asset.TestData1
	targetData := TestData2
	dataList := []GethAddress{targetData}
	// query 1: get geth address by asset's contract address
	assetContracAddressRow := AddGethAddressToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(assetData1.ContractAddress).WillReturnRows(assetContracAddressRow)
	// return contract Address
	foundGethAddress, err := CreateOrGetContractAddressFromAsset(mock, &assetData1)
	if err != nil {
		t.Fatalf("an error '%s' in CreateOrGetContractAddressFromAsset", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method CreateOrGetContractAddressFromAsset: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetContractAddressFromExistingAssetForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetData1 := asset.TestData2
	invalidContractAddress := "Invalid-Contract-Address"
	assetData1.ContractAddress = invalidContractAddress
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(invalidContractAddress).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethAddress, err := CreateOrGetContractAddressFromAsset(mock, &assetData1)
	if err == nil {
		t.Fatalf("expected an error '%s' in CreateOrGetContractAddressFromAsset", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method CreateOrGetContractAddressFromAsset: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetContractAddressFromNewAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset TestData2 is Contract
	assetData1 := asset.TestData2
	noRows := pgxmock.NewRows(DBColumns)
	// 1st step`: get geth address by asset's contract address
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(assetData1.ContractAddress).WillReturnRows(noRows)
	// 2nd step insert address
	contractName := fmt.Sprintf("Contract : %s", assetData1.Name)
	contractTypeID := utils.CONTRACT_ADDRESS_TYPE_STRUCTURED_VALUE_ID
	newID := 1
	targetData := GethAddress{
		ID:            &newID,
		Name:          contractName,
		AlternateName: contractName,
		AddressStr:    assetData1.ContractAddress,
		AddressTypeID: &contractTypeID,
		CreatedBy:     utils.SYSTEM_NAME,
	}
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_addresses").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.AddressStr,    //4
		targetData.AddressTypeID, //5
		targetData.CreatedBy,     //6
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(newID))
	mock.ExpectCommit()
	// return contract Address
	foundGethAddress, err := CreateOrGetContractAddressFromAsset(mock, &assetData1)
	if err != nil {
		t.Fatalf("an error '%s' in CreateOrGetContractAddressFromAsset", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method CreateOrGetContractAddressFromAsset: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetContractAddressFromNewAssetForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset TestData2 is Contract
	assetData1 := asset.TestData2
	noRows := pgxmock.NewRows(DBColumns)
	// 1st step`: get geth address by asset's contract address
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(assetData1.ContractAddress).WillReturnRows(noRows)
	// 2nd step insert address
	contractName := fmt.Sprintf("Contract : %s", assetData1.Name)
	contractTypeID := utils.CONTRACT_ADDRESS_TYPE_STRUCTURED_VALUE_ID
	newID := 1
	targetData := GethAddress{
		ID:            &newID,
		Name:          contractName,
		AlternateName: contractName,
		AddressStr:    assetData1.ContractAddress,
		AddressTypeID: &contractTypeID,
		CreatedBy:     utils.SYSTEM_NAME,
	}
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
	foundGethAddress, err := CreateOrGetContractAddressFromAsset(mock, &assetData1)
	if err == nil {
		t.Fatalf("expected an error '%s' in CreateOrGetContractAddressFromAsset", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method CreateOrGetContractAddressFromAsset: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

// func CreateOrGetAddress
func TestCreateOrGetAddressFromExistingAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test address
	targetData := TestData1
	dataList := []GethAddress{targetData}
	// query 1: get geth address by asset's contract address
	assetContracAddressRow := AddGethAddressToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(targetData.AddressStr).WillReturnRows(assetContracAddressRow)
	// return contract Address
	foundGethAddress, err := CreateOrGetAddress(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in CreateOrGetAddress", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method CreateOrGetAddress: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetAddressFromExistingAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	invalidAddress := "Invalid-Contract-Address"
	targetData.AddressStr = invalidAddress
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(invalidAddress).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethAddress, err := CreateOrGetAddress(mock, &targetData)
	if err == nil {
		t.Fatalf("expected an error '%s' in CreateOrGetAddress", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method CreateOrGetAddress: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetAddressFromNewAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset TestData1 = EOA
	targetData := TestData1
	noRows := pgxmock.NewRows(DBColumns)
	// 1st step`: get geth address by asset's contract address
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(targetData.AddressStr).WillReturnRows(noRows)
	newID := 1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_addresses").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.AddressStr,    //4
		targetData.AddressTypeID, //5
		targetData.CreatedBy,     //6
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(newID))
	mock.ExpectCommit()
	// return contract Address
	foundGethAddress, err := CreateOrGetAddress(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in CreateOrGetAddress", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method CreateOrGetAddress: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetAddressFromNewAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset TestData1 = EOA
	targetData := TestData1
	noRows := pgxmock.NewRows(DBColumns)
	// 1st step`: get geth address by asset's contract address
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(targetData.AddressStr).WillReturnRows(noRows)
	// 2nd step insert address
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
	foundGethAddress, err := CreateOrGetAddress(mock, &targetData)
	if err == nil {
		t.Fatalf("expected an error '%s' in CreateOrGetAddress", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method CreateOrGetAddress: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

// func CreateOrGetEOAAddress
func TestCreateOrGetEOAAddressFromExistingAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test address
	targetData := TestData1
	dataList := []GethAddress{targetData}
	// query 1: get geth address by asset's contract address
	assetContracAddressRow := AddGethAddressToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(targetData.AddressStr).WillReturnRows(assetContracAddressRow)
	// return contract Address
	foundGethAddress, err := CreateOrGetEOAAddress(mock, targetData.AddressStr)
	if err != nil {
		t.Fatalf("an error '%s' in CreateOrGetEOAAddress", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method CreateOrGetEOAAddress: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetEOAAddressFromExistingAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	invalidAddress := "Invalid-Contract-Address"
	targetData.AddressStr = invalidAddress
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(invalidAddress).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethAddress, err := CreateOrGetEOAAddress(mock, invalidAddress)
	if err == nil {
		t.Fatalf("expected an error '%s' in CreateOrGetEOAAddress", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method CreateOrGetEOAAddress: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetEOAAddressFromNewAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset TestData1 = EOA
	targetData := TestData1
	addressStr := targetData.AddressStr
	noRows := pgxmock.NewRows(DBColumns)
	// 1st step`: get geth address by asset's contract address
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(addressStr).WillReturnRows(noRows)
	newID := 1
	targetData.ID = &newID
	// change the time to 0
	zeroTime := time.Time{}
	targetData.CreatedAt = zeroTime
	targetData.UpdatedAt = zeroTime
	// reassign the targetData with limited fields
	newData := GethAddress{
		ID:            &newID,
		Name:          targetData.Name,          //1
		AlternateName: targetData.AlternateName, //2
		Description:   targetData.Description,   //3
		AddressStr:    targetData.AddressStr,    //4
		AddressTypeID: targetData.AddressTypeID, //5
		CreatedBy:     targetData.CreatedBy,     //6
	}
	targetData = newData
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_addresses").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.AddressStr,    //4
		targetData.AddressTypeID, //5
		targetData.CreatedBy,     //6
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(newID))
	mock.ExpectCommit()
	// return contract Address
	foundGethAddress, err := CreateOrGetEOAAddress(mock, addressStr)
	if err != nil {
		t.Fatalf("an error '%s' in CreateOrGetEOAAddress", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method CreateOrGetEOAAddress: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetEOAAddressFromNewAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset TestData1 is EOA
	targetData := TestData1
	addressStr := targetData.AddressStr
	noRows := pgxmock.NewRows(DBColumns)
	// 1st step`: get geth address by asset's contract address
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(addressStr).WillReturnRows(noRows)
	// 2nd step insert address
	newID := 1
	targetData.ID = &newID
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
	foundGethAddress, err := CreateOrGetEOAAddress(mock, addressStr)
	if err == nil {
		t.Fatalf("expected an error '%s' in CreateOrGetEOAAddress", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method CreateOrGetEOAAddress: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

// func CreateOrGetContractAddress
func TestCreateOrGetContractAddressFromExistingAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test address
	targetData := TestData2
	dataList := []GethAddress{targetData}
	// query 1: get geth address by asset's contract address
	assetContracAddressRow := AddGethAddressToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(targetData.AddressStr).WillReturnRows(assetContracAddressRow)
	// return contract Address
	foundGethAddress, err := CreateOrGetContractAddress(mock, targetData.AddressStr)
	if err != nil {
		t.Fatalf("an error '%s' in CreateOrGetContractAddress", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method CreateOrGetContractAddress: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetContractAddressFromExistingAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	invalidAddress := "Invalid-Contract-Address"
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(invalidAddress).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethAddress, err := CreateOrGetContractAddress(mock, invalidAddress)
	if err == nil {
		t.Fatalf("expected an error '%s' in CreateOrGetContractAddress", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method CreateOrGetContractAddress: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetContractAddressFromNewAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset TestData2 is Contract
	targetData := TestData2
	addressStr := targetData.AddressStr
	noRows := pgxmock.NewRows(DBColumns)
	// 1st step`: get geth address by asset's contract address
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(addressStr).WillReturnRows(noRows)
	newID := 1
	targetData.ID = &newID
	// change the time to 0
	zeroTime := time.Time{}
	targetData.CreatedAt = zeroTime
	targetData.UpdatedAt = zeroTime
	// reassign the targetData with limited fields
	newData := GethAddress{
		ID:            &newID,
		Name:          targetData.Name,          //1
		AlternateName: targetData.AlternateName, //2
		Description:   targetData.Description,   //3
		AddressStr:    targetData.AddressStr,    //4
		AddressTypeID: targetData.AddressTypeID, //5
		CreatedBy:     targetData.CreatedBy,     //6
	}
	targetData = newData
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_addresses").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.AddressStr,    //4
		targetData.AddressTypeID, //5
		targetData.CreatedBy,     //6
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(newID))
	mock.ExpectCommit()
	// return contract Address
	foundGethAddress, err := CreateOrGetContractAddress(mock, addressStr)
	if err != nil {
		t.Fatalf("an error '%s' in CreateOrGetContractAddress", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method CreateOrGetContractAddress: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateOrGetContractAddressFromNewAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset TestData2 is Contract
	targetData := TestData2
	addressStr := targetData.AddressStr
	noRows := pgxmock.NewRows(DBColumns)
	// 1st step`: get geth address by asset's contract address
	mock.ExpectQuery("^SELECT (.+) FROM geth_addresses").WithArgs(addressStr).WillReturnRows(noRows)
	// 2nd step insert address
	newID := 1
	targetData.ID = &newID
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
	foundGethAddress, err := CreateOrGetContractAddress(mock, addressStr)
	if err == nil {
		t.Fatalf("expected an error '%s' in CreateOrGetContractAddress", err)
	}
	if foundGethAddress != nil {
		t.Errorf("Expected GethAddress From Method CreateOrGetContractAddress: to be empty but got this: %v", foundGethAddress)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

// CreateGethAddress
func TestCreateGethAddressFromNewConrtractAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset testDat2 is Contract
	targetData := TestData2
	addressStr := targetData.AddressStr
	// reassign the targetData with limited fields
	newData := GethAddress{
		Name:          targetData.Name,
		AlternateName: targetData.AlternateName,
		AddressStr:    targetData.AddressStr,
		AddressTypeID: targetData.AddressTypeID,
		CreatedBy:     utils.SYSTEM_NAME,
	}
	targetData = newData
	isEOA := false
	// return contract Address
	foundGethAddress, err := CreateGethAddress(addressStr, isEOA)
	// reassign the UUID as this is random
	targetData.UUID = foundGethAddress.UUID
	if err != nil {
		t.Fatalf("an error '%s' in CreateGethAddress", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method CreateGethAddress: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestCreateGethAddressFromNewEOAAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	// test asset testDat1 is EOA
	targetData := TestData1
	addressStr := targetData.AddressStr
	// reassign the targetData with limited fields
	newData := GethAddress{
		Name:          targetData.Name,
		AlternateName: targetData.AlternateName,
		AddressStr:    targetData.AddressStr,
		AddressTypeID: targetData.AddressTypeID,
		CreatedBy:     utils.SYSTEM_NAME,
	}
	targetData = newData
	isEOA := true
	// return contract Address
	foundGethAddress, err := CreateGethAddress(addressStr, isEOA)
	// reassign the UUID as this is random
	targetData.UUID = foundGethAddress.UUID
	if err != nil {
		t.Fatalf("an error '%s' in CreateGethAddress", err)
	}
	if cmp.Equal(*foundGethAddress, targetData) == false {
		t.Errorf("Expected GethAddress From Method CreateGethAddress: %v is different from actual %v", foundGethAddress, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
