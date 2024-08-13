package account

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
	"id",               //1
	"uuid",             //2
	"name",             //3
	"alternate_name",   //4
	"address",          //5
	"name_from_source", //6
	"portfolio_id",     //7
	"source_id",        //8
	"account_type_id",  //9
	"description",      //10
	"created_by",       //11
	"created_at",       //12
	"updated_by",       //13
	"updated_at",       //14
	"chain_id",         //15
}

var DBColumnsInsertAccounts = []string{
	"uuid",             //1
	"name",             //2
	"alternate_name",   //3
	"address",          //4
	"name_from_source", //5
	"portfolio_id",     //6
	"source_id",        //7
	"account_type_id",  //8
	"description",      //9
	"created_by",       //10
	"created_at",       //11
	"updated_by",       //12
	"updated_at",       //13
	"chain_id",         //14
}

var TestData1 = Account{
	ID:             utils.Ptr[int](1),
	UUID:           "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:           "0x82069A502461C3f73705A9Cd6d3aB4ef27112345",
	AlternateName:  "ETH Wallet 1234",
	Address:        "0x82069A502461C3f73705A9Cd6d3aB4ef27112345",
	NameFromSource: "0x82069A502461C3f73705A9Cd6d3aB4ef27112345",
	PortfolioID:    nil,
	SourceID:       utils.Ptr[int](6),
	AccountTypeID:  nil,
	Description:    "",
	CreatedBy:      "SYSTEM",
	CreatedAt:      utils.SampleCreatedAtTime,
	UpdatedBy:      "SYSTEM",
	UpdatedAt:      utils.SampleCreatedAtTime,
	ChainID:        utils.Ptr[int](1),
}

var TestData2 = Account{
	ID:             utils.Ptr[int](2),
	UUID:           "98d5545f-4244-4838-234e-d98722ab1c52",
	Name:           "L-gro",
	AlternateName:  "0x4740887F03191E46E597D1A51c749176AA123451",
	Address:        "0x4740887F03191E46E597D1A51c749176AA123451",
	NameFromSource: "0x4740887F03191E46E597D1A51c749176AA123451",
	PortfolioID:    utils.Ptr[int](1),
	SourceID:       nil,
	AccountTypeID:  utils.Ptr[int](6),
	Description:    "",
	CreatedBy:      "test@gmail.com",
	CreatedAt:      utils.SampleCreatedAtTime,
	UpdatedBy:      "test@gmail.com",
	UpdatedAt:      utils.SampleCreatedAtTime,
	ChainID:        utils.Ptr[int](2),
}

var TestAllData = []Account{TestData1, TestData2}

func AddAccountToMockRows(mock pgxmock.PgxPoolIface, dataList []Account) *pgxmock.Rows {
	rows := mock.NewRows(columns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,
			data.UUID,
			data.Name,
			data.AlternateName,
			data.Address,
			data.NameFromSource,
			data.PortfolioID,
			data.SourceID,
			data.AccountTypeID,
			data.Description,
			data.CreatedBy,
			data.CreatedAt,
			data.UpdatedBy,
			data.UpdatedAt,
			data.ChainID,
		)
	}
	return rows
}

func TestGetAccount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	dataList := []Account{targetData}
	mockRows := AddAccountToMockRows(mock, dataList)
	accountID := utils.Ptr[int](1)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE id = ?").WithArgs(*accountID).WillReturnRows(mockRows)
	foundAccount, err := GetAccount(mock, accountID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccount", err)
	}
	if cmp.Equal(*foundAccount, targetData) == false {
		t.Errorf("Expected Account From Method GetAccount: %v is different from actual %v", foundAccount, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	noRows := pgxmock.NewRows(columns)
	accountID := utils.Ptr[int](1)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE id = ?").WithArgs(*accountID).WillReturnRows(noRows)
	foundAccount, err := GetAccount(mock, accountID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccount", err)
	}
	if foundAccount != nil {
		t.Errorf("Expected Account From Method GetAccount: to be empty but got this: %v", foundAccount)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	invalidID := utils.Ptr[int](-1)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE id = ?").WithArgs(*invalidID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAccount, err := GetAccount(mock, invalidID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAccount", err)
	}
	if foundAccount != nil {
		t.Errorf("Expected Account From Method GetAccount: to be empty but got this: %v", foundAccount)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	accountID := -1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM accounts").WithArgs(accountID).WillReturnRows(differentModelRows)
	foundAccount, err := GetAccount(mock, &accountID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAccount", err)
	}
	if foundAccount != nil {
		t.Errorf("Expected foundAccount From Method GetAccount: to be empty but got this: %v", foundAccount)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountByAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Account{targetData}
	testAddress := targetData.Address
	mockRows := AddAccountToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE address = ?").WithArgs(testAddress).WillReturnRows(mockRows)
	foundAccount, err := GetAccountByAddress(mock, testAddress)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccountByAddress", err)
	}
	if cmp.Equal(*foundAccount, targetData) == false {
		t.Errorf("Expected Account From Method GetAccountByAddress: %v is different from actual %v", foundAccount, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountByAddressForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testAddress := "Fake-Address"
	noRows := pgxmock.NewRows(columns)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE address = ?").WithArgs(testAddress).WillReturnRows(noRows)
	foundAccount, err := GetAccountByAddress(mock, testAddress)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccountByAddress", err)
	}
	if foundAccount != nil {
		t.Errorf("Expected Account From Method GetAccountByAddress: to be empty but got this: %v", foundAccount)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountByAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testAddress := "Fake-Address"
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE address = ?").WithArgs(testAddress).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAccount, err := GetAccountByAddress(mock, testAddress)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAccountByAddress", err)
	}
	if foundAccount != nil {
		t.Errorf("Expected Account From Method GetAccountByAddress: to be empty but got this: %v", foundAccount)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountByAddressForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testAddress := "Fake-Address"
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM accounts").WithArgs(testAddress).WillReturnRows(differentModelRows)
	foundAccount, err := GetAccountByAddress(mock, testAddress)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAccountByAddress", err)
	}
	if foundAccount != nil {
		t.Errorf("Expected foundAccount From Method GetAccountByAddress: to be empty but got this: %v", foundAccount)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountByAlternateName(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Account{targetData}
	testAlternateName := targetData.AlternateName
	mockRows := AddAccountToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE alternate_name = ?").WithArgs(testAlternateName).WillReturnRows(mockRows)
	foundAccount, err := GetAccountByAlternateName(mock, testAlternateName)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccountByAlternateName", err)
	}
	if cmp.Equal(*foundAccount, targetData) == false {
		t.Errorf("Expected Account From Method GetAccountByAlternateName: %v is different from actual %v", foundAccount, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestGetAccountByAlternateNameForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testAlternateName := "Fake-AlternateName"
	noRows := pgxmock.NewRows(columns)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE alternate_name = ?").WithArgs(testAlternateName).WillReturnRows(noRows)
	foundAccount, err := GetAccountByAlternateName(mock, testAlternateName)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccountByAlternateName", err)
	}
	if foundAccount != nil {
		t.Errorf("Expected Account From Method GetAccountByAlternateName: to be empty but got this: %v", foundAccount)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountByAlternateNameForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testAlternateName := "Fake-AlternateName"
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE alternate_name = ?").WithArgs(testAlternateName).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAccount, err := GetAccountByAlternateName(mock, testAlternateName)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAccountByAlternateName", err)
	}
	if foundAccount != nil {
		t.Errorf("Expected Account From Method GetAccountByAlternateName: to be empty but got this: %v", foundAccount)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountByAlternateNameForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testAlternateName := "Fake-Address"
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM accounts").WithArgs(testAlternateName).WillReturnRows(differentModelRows)
	foundAccount, err := GetAccountByAlternateName(mock, testAlternateName)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAccountByAlternateName", err)
	}
	if foundAccount != nil {
		t.Errorf("Expected foundAccount From Method GetAccountByAlternateName: to be empty but got this: %v", foundAccount)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAccount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM accounts WHERE id = ?").WithArgs(*targetData.ID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveAccount(mock, targetData.ID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveAccount", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAccountOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetDataID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveAccount(mock, &targetDataID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAccountOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM accounts WHERE id = ?").WithArgs(-1).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveAccount(mock, targetData.ID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Account{TestData1, TestData2}
	mockRows := AddAccountToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM accounts").WillReturnRows(mockRows)
	ids := []int{1, 2}
	foundAccounts, err := GetAccountList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccount", err)
	}
	testAccounts := TestAllData
	for i, foundAccount := range foundAccounts {
		if cmp.Equal(foundAccount, testAccounts[i]) == false {
			t.Errorf("Expected Account From Method GetAccount: %v is different from actual %v", foundAccount, testAccounts[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM accounts").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := make([]int, 0)
	foundAccounts, err := GetAccountList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAccountList", err)
	}
	if len(foundAccounts) != 0 {
		t.Errorf("Expected From Method GetAccountList: to be empty but got this: %v", foundAccounts)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	ids := []int{1, 2}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM accounts").WillReturnRows(differentModelRows)
	foundAccounts, err := GetAccountList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAccountList", err)
	}
	if foundAccounts != nil {
		t.Errorf("Expected foundAccounts From Method GetAccountList: to be empty but got this: %v", foundAccounts)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAccount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE accounts").WithArgs(
		targetData.Name,
		targetData.AlternateName,
		targetData.Address,
		targetData.NameFromSource,
		targetData.PortfolioID,
		targetData.SourceID,
		targetData.AccountTypeID,
		targetData.Description,
		targetData.UpdatedBy,
		targetData.ChainID,
		targetData.ID,
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateAccount(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveAccount", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAccountOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateAccount(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateAccountOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateAccount(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateAccountOnFailure(t *testing.T) {
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
	mock.ExpectExec("^UPDATE accounts").WithArgs(
		targetData.Name,
		targetData.AlternateName,
		targetData.Address,
		targetData.NameFromSource,
		targetData.PortfolioID,
		targetData.SourceID,
		targetData.AccountTypeID,
		targetData.Description,
		targetData.UpdatedBy,
		targetData.ChainID,
		targetData.ID,
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateAccount(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAccount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO accounts").WithArgs(
		targetData.Name,
		targetData.AlternateName,
		targetData.Address,
		targetData.NameFromSource,
		targetData.PortfolioID,
		targetData.SourceID,
		targetData.AccountTypeID,
		targetData.Description,
		targetData.CreatedBy,
		targetData.ChainID,
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	newID, err := InsertAccount(mock, &targetData)
	if newID < 0 {
		t.Fatalf("ID should not be negative ID: %d", newID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertAccount", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAccountOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertAccount(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAccountOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO accounts").WithArgs(
		targetData.Name,
		targetData.AlternateName,
		targetData.Address,
		targetData.NameFromSource,
		targetData.PortfolioID,
		targetData.SourceID,
		targetData.AccountTypeID,
		targetData.Description,
		targetData.CreatedBy,
		targetData.ChainID,
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	newID, err := InsertAccount(mock, &targetData)
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

func TestInsertAccountOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO accounts").WithArgs(
		targetData.Name,
		targetData.AlternateName,
		targetData.Address,
		targetData.NameFromSource,
		targetData.PortfolioID,
		targetData.SourceID,
		targetData.AccountTypeID,
		targetData.Description,
		targetData.CreatedBy,
		targetData.ChainID,
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	newID, err := InsertAccount(mock, &targetData)
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

func TestInsertAccounts(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"accounts"}, DBColumnsInsertAccounts)
	targetData := TestAllData
	err = InsertAccounts(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertAccounts", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAccountsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"accounts"}, DBColumnsInsertAccounts).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertAccounts(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddAccountToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"account_type_id = 1", "chain_id=1"}
	mock.ExpectQuery("^SELECT (.+) FROM accounts").WillReturnRows(mockRows)
	foundAccountList, err := GetAccountListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccountListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundAccountList[i]) == false {
			t.Errorf("Expected sourceData From Method GetAccountListByPagination: %v is different from actual %v", sourceData, foundAccountList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"account_type_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM accounts").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAccountList, err := GetAccountListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAccountListByPagination", err)
	}
	if len(foundAccountList) != 0 {
		t.Errorf("Expected From Method GetAccountListByPagination: to be empty but got this: %v", foundAccountList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAccountListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"account_type_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM accounts").WillReturnRows(differentModelRows)
	foundAccountList, err := GetAccountListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAccountListByPagination", err)
	}
	if foundAccountList != nil {
		t.Errorf("Expected From Method GetAccountListByPagination: to be empty but got this: %v", foundAccountList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalAccountsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM accounts").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalAccountsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalAccountsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalAccountsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMinersForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM accounts").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalAccountsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalAccountsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalAccountsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
