package account

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var columns = []string{"id", "uuid", "name", "alternate_name", "address", "name_from_source", "portfolio_id", "source_id", "account_type_id", "description", "created_by", "created_at", "updated_by", "updated_at", "chain_id"}

var data1 = Account{
	ID:             utils.Ptr[int](1),
	UUID:           "8806079a-1833-4ad7-a834-b983a61c7b39",
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

var data2 = Account{
	ID:             utils.Ptr[int](2),
	UUID:           "98d5505f-1244-4838-8f3e-d98722ab1c52",
	Name:           "L-gro",
	AlternateName:  "0x4740887F03191E46E597D1A51c749176AA12345",
	Address:        "0x4740887F03191E46E597D1A51c749176AA12345",
	NameFromSource: "0x4740887F03191E46E597D1A51c749176AA12345",
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

var allData = []Account{data1, data2}

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
	targetData := data1
	dataList := []Account{targetData}
	mockRows := AddAccountToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE id = ?").WithArgs(1).WillReturnRows(mockRows)
	accountID := 1
	acct, err := GetAccount(mock, accountID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccount", err)
	}
	if cmp.Equal(*acct, targetData) == false {
		t.Errorf("Expected Account From Method GetAccount: %v is different from actual %v", acct, targetData)
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
	targetData := data2
	dataList := []Account{targetData}
	testAddress := targetData.Address
	mockRows := AddAccountToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE address = ?").WithArgs(testAddress).WillReturnRows(mockRows)
	acct, err := GetAccountByAddress(mock, testAddress)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccountByAddress", err)
	}
	if cmp.Equal(*acct, targetData) == false {
		t.Errorf("Expected Account From Method GetAccountByAddress: %v is different from actual %v", acct, targetData)
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
	targetData := data2
	dataList := []Account{targetData}
	testAlternateName := targetData.AlternateName
	mockRows := AddAccountToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE alternate_name = ?").WithArgs(testAlternateName).WillReturnRows(mockRows)
	acct, err := GetAccountByAlternateName(mock, testAlternateName)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccountByAlternateName", err)
	}
	if cmp.Equal(*acct, targetData) == false {
		t.Errorf("Expected Account From Method GetAccountByAlternateName: %v is different from actual %v", acct, targetData)
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
	targetData := data1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM accounts WHERE id = ?").WithArgs(*targetData.ID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveAccount(mock, *targetData.ID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveAccount", err)
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
	targetData := data1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM accounts WHERE id = ?").WithArgs(-1).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveAccount(mock, *targetData.ID)
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
	dataList := []Account{data1, data2}
	mockRows := AddAccountToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM accounts").WillReturnRows(mockRows)
	ids := make([]int, 0)
	accts, err := GetAccountList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetAccount", err)
	}
	// createdAt, err := time.Parse(utils.LayoutRFC3339, "2020-12-10 22:57:32.478544")
	testAccounts := allData
	// t.Logf("mock : %v \n", accts)
	for i, acct := range accts {
		// t.Logf("i: %d,  \n mock : %v , \n in memory : %v", i, acct, testAccounts[i])

		if cmp.Equal(acct, testAccounts[i]) == false {
			t.Errorf("Expected Account From Method GetAccount: %v is different from actual %v", acct, testAccounts[i])
		}

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
	targetData := data1
	targetData.Name = "Updated Name"
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
	err = UpdateAccount(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveAccount", err)
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
	targetData := data1
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
	err = UpdateAccount(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
