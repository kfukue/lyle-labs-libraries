package transaction

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"id",              //1
	"uuid",            //2
	"name",            //3
	"alternate_name",  //4
	"start_date",      //5
	"end_date",        //6
	"description",     //7
	"tx_hash",         //8
	"status_id",       //9
	"from_account_id", //10
	"to_account_id",   //11
	"chain_id",        //12
	"created_by",      //13
	"created_at",      //14
	"updated_by",      //15
	"updated_at",      //16
}
var DBColumnsInsertTransactions = []string{
	"uuid",            //1
	"name",            //2
	"alternate_name",  //3
	"start_date",      //4
	"end_date",        //5
	"description",     //6
	"tx_hash",         //7
	"status_id",       //8
	"from_account_id", //9
	"to_account_id",   //10
	"chain_id",        //11
	"created_by",      //12
	"created_at",      //13
	"updated_by",      //14
	"updated_at",      //15
}

var TestData1 = Transaction{
	ID:            utils.Ptr[int](1),                                                                                                                                  //1
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72",                                                                                                             //2
	Name:          "Date : 2023-05-25 16:15:06, Harvest PTP Strategy : Vector Finance Pool : XPTP",                                                                    //3
	AlternateName: "0xea4e9ded8d172b1b189e0921ddfdd7cd4979101672991f4e9c06472394cc1232",                                                                               //4
	StartDate:     utils.SampleCreatedAtTime,                                                                                                                          //5
	EndDate:       utils.SampleCreatedAtTime,                                                                                                                          //6
	Description:   "&{0xea4e9ded8d172b1b189e0921ddfdd7cd4979101672991f4e9c06472394c12345 0 1200000 25000000000 349 0x423D0FE33031aA4456a17b150804aA57fc151234 false}", //7
	TxHash:        "0xea4e9ded8d172b1b189e0921ddfdd7cd4979101672991f4e9c06472394c31234",                                                                               //8
	StatusID:      utils.Ptr[int](52),                                                                                                                                 //9
	FromAccountID: utils.Ptr[int](12),                                                                                                                                 //10
	ToAccountID:   utils.Ptr[int](233),                                                                                                                                //11
	ChainID:       utils.Ptr[int](13),                                                                                                                                 //12
	CreatedBy:     "SYSTEM",                                                                                                                                           //13
	CreatedAt:     utils.SampleCreatedAtTime,                                                                                                                          //14
	UpdatedBy:     "SYSTEM",                                                                                                                                           //15
	UpdatedAt:     utils.SampleCreatedAtTime,                                                                                                                          //16
}

var TestData2 = Transaction{
	ID:            utils.Ptr[int](2),
	UUID:          "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",
	Name:          "Date : 2023-04-23 16:15:06, Harvest PTP Strategy : Vector Finance Pool : XPTP",                                                                    //3
	AlternateName: "0x52ccb3c3e7d0e9b8c4e834b5844f522a09b09110ec14a0ad3bac8d6e01874578",                                                                               //4
	StartDate:     utils.SampleCreatedAtTime,                                                                                                                          //5
	EndDate:       utils.SampleCreatedAtTime,                                                                                                                          //6
	Description:   "&{0xa1fcb2fbbd01259d469909e2571c6a57f3984d6c9d6e72b27fb7c0474748415 0 1200000 44778882798 320 0x423D0FE33031aA4456a17b150804aA57fc1514748 false}", //7
	TxHash:        "0xdcc37e8e0296c0b2f295b8724e0470049c9d7ebfe09b2a2a782d42c32447458",                                                                                //8
	StatusID:      utils.Ptr[int](52),                                                                                                                                 //9
	FromAccountID: utils.Ptr[int](22),                                                                                                                                 //10
	ToAccountID:   utils.Ptr[int](243),                                                                                                                                //11
	ChainID:       utils.Ptr[int](13),                                                                                                                                 //12
	CreatedBy:     "SYSTEM",                                                                                                                                           //13
	CreatedAt:     utils.SampleCreatedAtTime,                                                                                                                          //14
	UpdatedBy:     "SYSTEM",                                                                                                                                           //15
	UpdatedAt:     utils.SampleCreatedAtTime,                                                                                                                          //16
}
var TestAllData = []Transaction{TestData1, TestData2}

func AddTransactionToMockRows(mock pgxmock.PgxPoolIface, dataList []Transaction) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,            //1
			data.UUID,          //2
			data.Name,          //3
			data.AlternateName, //4
			data.StartDate,     //5
			data.EndDate,       //6
			data.Description,   //7
			data.TxHash,        //8
			data.StatusID,      //9
			data.FromAccountID, //10
			data.ToAccountID,   //11
			data.ChainID,       //12
			data.CreatedBy,     //13
			data.CreatedAt,     //14
			data.UpdatedBy,     //15
			data.UpdatedAt,     //16
		)
	}
	return rows
}

func TestGetTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Transaction{targetData}
	transactionID := targetData.ID
	mockRows := AddTransactionToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WithArgs(*transactionID).WillReturnRows(mockRows)
	foundTransaction, err := GetTransaction(mock, transactionID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransaction", err)
	}
	if cmp.Equal(*foundTransaction, targetData) == false {
		t.Errorf("Expected Transaction From Method GetTransaction: %v is different from actual %v", foundTransaction, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WithArgs(transactionID).WillReturnRows(noRows)
	foundTransaction, err := GetTransaction(mock, &transactionID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransaction", err)
	}
	if foundTransaction != nil {
		t.Errorf("Expected Transaction From Method GetTransaction: to be empty but got this: %v", foundTransaction)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WithArgs(transactionID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransaction, err := GetTransaction(mock, &transactionID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransaction", err)
	}
	if foundTransaction != nil {
		t.Errorf("Expected Transaction From Method GetTransaction: to be empty but got this: %v", foundTransaction)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WithArgs(transactionID).WillReturnRows(differentModelRows)
	foundTransaction, err := GetTransaction(mock, &transactionID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransaction", err)
	}
	if foundTransaction != nil {
		t.Errorf("Expected foundTransaction From Method GetTransaction: to be empty but got this: %v", foundTransaction)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	transactionID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transactions").WithArgs(*transactionID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveTransaction(mock, transactionID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveTransaction", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveTransaction(mock, &transactionID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transactions").WithArgs(transactionID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveTransaction(mock, &transactionID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactions(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddTransactionToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WillReturnRows(mockRows)
	ids := []int{1, 2}
	foundTransactionList, err := GetTransactions(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactions", err)
	}
	testMarketDataList := TestAllData
	for i, foundTransaction := range foundTransactionList {
		if cmp.Equal(foundTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected Transaction From Method GetTransactions: %v is different from actual %v", foundTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WillReturnRows(differentModelRows)
	ids := []int{1, 2}
	foundTransactionList, err := GetTransactions(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactions", err)
	}
	if foundTransactionList != nil {
		t.Errorf("Expected foundTransactionList From Method GetTransactions: to be empty but got this: %v", foundTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1, 2}
	foundTransactionList, err := GetTransactions(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactions", err)
	}
	if len(foundTransactionList) != 0 {
		t.Errorf("Expected From Method GetTransactions: to be empty but got this: %v", foundTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTransactionToMockRows(mock, dataList)
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WithArgs(pq.Array(uuids)).WillReturnRows(mockRows)
	foundTransactionList, err := GetTransactionsByUUIDs(mock, uuids)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionsByUUIDs", err)
	}
	testMarketDataList := TestAllData
	for i, foundTransaction := range foundTransactionList {
		if cmp.Equal(foundTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected Transaction From Method GetTransactionsByUUIDs: %v is different from actual %v", foundTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()

	uuids := []string{"uuid-invalid-1", "uuid-invalid-2"}
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WithArgs(pq.Array(uuids)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionList, err := GetTransactionsByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionsByUUIDs", err)
	}
	if len(foundTransactionList) != 0 {
		t.Errorf("Expected From Method GetTransactionsByUUIDs: to be empty but got this: %v", foundTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionsByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuids := []string{"uuid-invalid-1", "uuid-invalid-2"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WithArgs(pq.Array(uuids)).WillReturnRows(differentModelRows)
	foundTransactionList, err := GetTransactionsByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionsByUUIDs", err)
	}
	if foundTransactionList != nil {
		t.Errorf("Expected foundTransactionList From Method GetTransactionsByUUIDs: to be empty but got this: %v", foundTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffTransactions(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTransactionToMockRows(mock, dataList)
	diffInDate := 2
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WithArgs(diffInDate).WillReturnRows(mockRows)
	foundTransactionList, err := GetStartAndEndDateDiffTransactions(mock, &diffInDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndDateDiffTransactions", err)
	}
	testMarketDataList := TestAllData
	for i, foundTransaction := range foundTransactionList {
		if cmp.Equal(foundTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected Transaction From Method GetStartAndEndDateDiffTransactions: %v is different from actual %v", foundTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffTransactionsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()

	diffInDate := 2
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WithArgs(diffInDate).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionList, err := GetStartAndEndDateDiffTransactions(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffTransactions", err)
	}
	if len(foundTransactionList) != 0 {
		t.Errorf("Expected From Method GetStartAndEndDateDiffTransactions: to be empty but got this: %v", foundTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffTransactionsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := 2
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WithArgs(diffInDate).WillReturnRows(differentModelRows)
	foundTransactionList, err := GetStartAndEndDateDiffTransactions(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffTransactions", err)
	}
	if foundTransactionList != nil {
		t.Errorf("Expected foundTransactionList From Method GetStartAndEndDateDiffTransactions: to be empty but got this: %v", foundTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transactions").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.StartDate,     //3
		targetData.EndDate,       //4
		targetData.Description,   //5
		targetData.TxHash,        //6
		targetData.StatusID,      //7
		targetData.FromAccountID, //8
		targetData.ToAccountID,   //9
		targetData.ChainID,       //10
		targetData.UpdatedBy,     //11
		targetData.ID,            //12
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateTransaction(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateTransaction", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTransactionOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateTransaction(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateTransaction(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transactions").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.StartDate,     //3
		targetData.EndDate,       //4
		targetData.Description,   //5
		targetData.TxHash,        //6
		targetData.StatusID,      //7
		targetData.FromAccountID, //8
		targetData.ToAccountID,   //9
		targetData.ChainID,       //10
		targetData.UpdatedBy,     //11
		targetData.ID,            //12
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateTransaction(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transactions").WithArgs(
		targetData.UUID,          //1
		targetData.Name,          //2
		targetData.AlternateName, //3
		targetData.StartDate,     //4
		targetData.EndDate,       //5
		targetData.Description,   //6
		targetData.TxHash,        //7
		targetData.StatusID,      //8
		targetData.FromAccountID, //9
		targetData.ToAccountID,   //10
		targetData.ChainID,       //11
		targetData.CreatedBy,     //12
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	transactionID, err := InsertTransaction(mock, &targetData)
	if transactionID < 0 {
		t.Fatalf("transactionID should not be negative ID: %d", transactionID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertTransaction", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertTransaction(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestInsertTransactionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.StatusID = nil
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transactions").WithArgs(
		targetData.UUID,          //1
		targetData.Name,          //2
		targetData.AlternateName, //3
		targetData.StartDate,     //4
		targetData.EndDate,       //5
		targetData.Description,   //6
		targetData.TxHash,        //7
		targetData.StatusID,      //8
		targetData.FromAccountID, //9
		targetData.ToAccountID,   //10
		targetData.ChainID,       //11
		targetData.CreatedBy,     //12
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	transactionID, err := InsertTransaction(mock, &targetData)
	if transactionID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error transactionID: %d", transactionID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transactions").WithArgs(
		targetData.UUID,          //1
		targetData.Name,          //2
		targetData.AlternateName, //3
		targetData.StartDate,     //4
		targetData.EndDate,       //5
		targetData.Description,   //6
		targetData.TxHash,        //7
		targetData.StatusID,      //8
		targetData.FromAccountID, //9
		targetData.ToAccountID,   //10
		targetData.ChainID,       //11
		targetData.CreatedBy,     //12
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	transactionID, err := InsertTransaction(mock, &targetData)
	if transactionID >= 0 {
		t.Fatalf("Expecting -1 for transactionID because of error transactionID: %d", transactionID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactions(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"transactions"}, DBColumnsInsertTransactions)
	targetData := TestAllData
	err = InsertTransactions(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertTransactions", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"transactions"}, DBColumnsInsertTransactions).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertTransactions(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionsByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTransactionToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"status_id = 1", "exchange_id=1"}
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WillReturnRows(mockRows)
	foundTransactionList, err := GetTransactionsByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionsByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundTransactionList[i]) == false {
			t.Errorf("Expected sourceData From Method GetTransactionsByPagination: %v is different from actual %v", sourceData, foundTransactionList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionsByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"chain_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionList, err := GetTransactionsByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionsByPagination", err)
	}
	if len(foundTransactionList) != 0 {
		t.Errorf("Expected From Method GetTransactionsByPagination: to be empty but got this: %v", foundTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionsByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"chain_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transactions").WillReturnRows(differentModelRows)
	foundTransactionList, err := GetTransactionsByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionsByPagination", err)
	}
	if foundTransactionList != nil {
		t.Errorf("Expected From Method GetTransactionsByPagination: to be empty but got this: %v", foundTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalTransactionsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM transactions").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalTransactionsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalTransactionsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalTransactionsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM transactions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalTransactionsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalTransactionsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalTransactionsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
