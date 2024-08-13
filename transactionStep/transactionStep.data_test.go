package transactionstep

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
	"transaction_id", //1
	"step_id",        //2
	"uuid",           //3
	"name",           //4
	"alternate_name", //5
	"description",    //6
	"created_by",     //7
	"created_at",     //8
	"updated_by",     //9
	"updated_at",     //10
}
var DBColumnsInsertTransactionSteps = []string{
	"transaction_id", //1
	"step_id",        //2
	"uuid",           //3
	"name",           //4
	"alternate_name", //5
	"description",    //6
	"created_by",     //7
	"created_at",     //8
	"updated_by",     //9
	"updated_at",     //10
}

var TestData1 = TransactionStep{
	TransactionID: utils.Ptr[int](1),                                                                            //1
	StepID:        utils.Ptr[int](1),                                                                            //2
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72",                                                       //3
	Name:          "Date : 2022-09-19 16:15:05, Harvest PTP Strategy : Vector Finance Pool : XPTP",              //4
	AlternateName: "0x5ab87fb9469455d2b5430aaf53a8067ebb41d591196491b97134c7eaa51477458",                        //5
	Description:   "Chain : Avalanche; Tx: 0x5ab87fb9469455d2b5430aaf53a8067ebb41d591196491b97134c7eaa51477458", //6
	CreatedBy:     "SYSTEM",                                                                                     //7
	CreatedAt:     utils.SampleCreatedAtTime,                                                                    //8
	UpdatedBy:     "SYSTEM",                                                                                     //9
	UpdatedAt:     utils.SampleCreatedAtTime,                                                                    //10

}

var TestData2 = TransactionStep{
	TransactionID: utils.Ptr[int](2),                                                                         //1
	StepID:        utils.Ptr[int](2),                                                                         //2
	UUID:          "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                                                    //3
	Name:          "Date : 2022-08-04 21:52:03, Harvest VELO/OP LP Reward : Velodrome Pool : VELO/OP LP",     //4
	AlternateName: "0x936ba1095ce1ea3edc861f94a3b1a1dcad9bb8c913cc1d4e5b1f7f6f8a94715",                       //5
	Description:   "Chain : Optimism; Tx: 0x936ba1095ce1ea3edc861f94a3b1a1dcad9bb8c913cc1d4e5b1f7f6f8a94715", //6
	CreatedBy:     "SYSTEM",                                                                                  //7
	CreatedAt:     utils.SampleCreatedAtTime,                                                                 //8
	UpdatedBy:     "SYSTEM",                                                                                  //9
	UpdatedAt:     utils.SampleCreatedAtTime,                                                                 //10
}
var TestAllData = []TransactionStep{TestData1, TestData2}

func AddTransactionStepToMockRows(mock pgxmock.PgxPoolIface, dataList []TransactionStep) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.TransactionID, //1
			data.StepID,        //2
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

func TestGetTransactionStep(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []TransactionStep{targetData}
	transactionID := targetData.TransactionID
	stepID := targetData.StepID
	mockRows := AddTransactionStepToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(*transactionID, *stepID).WillReturnRows(mockRows)
	foundTransactionStep, err := GetTransactionStep(mock, transactionID, stepID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionStep", err)
	}
	if cmp.Equal(*foundTransactionStep, targetData) == false {
		t.Errorf("Expected TransactionStep From Method GetTransactionStep: %v is different from actual %v", foundTransactionStep, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := 999
	stepID := 1111
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(transactionID, stepID).WillReturnRows(noRows)
	foundTransactionStep, err := GetTransactionStep(mock, &transactionID, &stepID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionStep", err)
	}
	if foundTransactionStep != nil {
		t.Errorf("Expected TransactionStep From Method GetTransactionStep: to be empty but got this: %v", foundTransactionStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	stepID := 1111
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(transactionID, stepID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionStep, err := GetTransactionStep(mock, &transactionID, &stepID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionStep", err)
	}
	if foundTransactionStep != nil {
		t.Errorf("Expected TransactionStep From Method GetTransactionStep: to be empty but got this: %v", foundTransactionStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	stepID := 1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(transactionID, stepID).WillReturnRows(differentModelRows)
	foundTransactionStep, err := GetTransactionStep(mock, &transactionID, &stepID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionStep", err)
	}
	if foundTransactionStep != nil {
		t.Errorf("Expected TransactionStep From Method GetTransactionStep: to be empty but got this: %v", foundTransactionStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []TransactionStep{targetData}
	transactionStepUUID := targetData.UUID
	mockRows := AddTransactionStepToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(transactionStepUUID).WillReturnRows(mockRows)
	foundTransactionStep, err := GetTransactionStepByUUID(mock, transactionStepUUID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionStepByUUID", err)
	}
	if cmp.Equal(*foundTransactionStep, targetData) == false {
		t.Errorf("Expected TransactionStep From Method GetTransactionStepByUUID: %v is different from actual %v", foundTransactionStep, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepByUUIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionStepUUID := "no-row-uuid"
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(transactionStepUUID).WillReturnRows(noRows)
	foundTransactionStep, err := GetTransactionStepByUUID(mock, transactionStepUUID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionStepByUUID", err)
	}
	if foundTransactionStep != nil {
		t.Errorf("Expected TransactionStep From Method GetTransactionStepByUUID: to be empty but got this: %v", foundTransactionStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepByUUIDForSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionStepUUID := "row-error-uuid"
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(transactionStepUUID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionStep, err := GetTransactionStepByUUID(mock, transactionStepUUID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionStepByUUID", err)
	}
	if foundTransactionStep != nil {
		t.Errorf("Expected TransactionStep From Method GetTransactionStepByUUID: to be empty but got this: %v", foundTransactionStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepByUUIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionStepUUID := "row-different-model-uuid"

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(transactionStepUUID).WillReturnRows(differentModelRows)
	foundTransactionStep, err := GetTransactionStepByUUID(mock, transactionStepUUID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionStepByUUID", err)
	}
	if foundTransactionStep != nil {
		t.Errorf("Expected TransactionStep From Method GetTransactionStepByUUID: to be empty but got this: %v", foundTransactionStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	UUIDList := []string{TestData1.UUID, TestData2.UUID}

	mockRows := AddTransactionStepToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(pq.Array(UUIDList)).WillReturnRows(mockRows)
	foundTransactionStepList, err := GetTransactionStepsByUUIDs(mock, UUIDList)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionStepsByUUIDs", err)
	}
	for i, sourceTransactionStep := range dataList {
		if cmp.Equal(sourceTransactionStep, foundTransactionStepList[i]) == false {
			t.Errorf("Expected TransactionStep From Method GetTransactionStepsByUUIDs: %v is different from actual %v", sourceTransactionStep, foundTransactionStepList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(pq.Array(UUIDList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionStepList, err := GetTransactionStepsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionStepsByUUIDs", err)
	}
	if len(foundTransactionStepList) != 0 {
		t.Errorf("Expected From Method GetTransactionStepsByUUIDs: to be empty but got this: %v", foundTransactionStepList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepsByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WithArgs(pq.Array(UUIDList)).WillReturnRows(differentModelRows)
	foundTransactionStep, err := GetTransactionStepsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionStepsByUUIDs", err)
	}
	if foundTransactionStep != nil {
		t.Errorf("Expected TransactionStep From Method GetTransactionStepsByUUIDs: to be empty but got this: %v", foundTransactionStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionStep(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	transactionID := targetData.TransactionID
	stepID := targetData.StepID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_steps").WithArgs(*transactionID, *stepID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveTransactionStep(mock, transactionID, stepID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveTransactionStep", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionStepOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	stepID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveTransactionStep(mock, &transactionID, &stepID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionStepOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	stepID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_steps").WithArgs(transactionID, stepID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveTransactionStep(mock, &transactionID, &stepID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionStepByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	transactionStepUUID := targetData.UUID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_steps").WithArgs(transactionStepUUID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveTransactionStepByUUID(mock, transactionStepUUID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveTransactionStepByUUID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionStepByUUIDOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionStepUUID := "Fail-at-begining"
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveTransactionStepByUUID(mock, transactionStepUUID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionStepByUUIDOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionStepUUID := "Fail-at-end"
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_steps").WithArgs(transactionStepUUID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveTransactionStepByUUID(mock, transactionStepUUID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddTransactionStepToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WillReturnRows(mockRows)
	foundTransactionStepList, err := GetTransactionStepList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionStepList", err)
	}
	for i, sourceTransactionStep := range dataList {
		if cmp.Equal(sourceTransactionStep, foundTransactionStepList[i]) == false {
			t.Errorf("Expected TransactionStep From Method GetTransactionStepList: %v is different from actual %v", sourceTransactionStep, foundTransactionStepList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionStepList, err := GetTransactionStepList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionStepList", err)
	}
	if len(foundTransactionStepList) != 0 {
		t.Errorf("Expected From Method GetTransactionStepList: to be empty but got this: %v", foundTransactionStepList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WillReturnRows(differentModelRows)
	foundTransactionStep, err := GetTransactionStepList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionStepList", err)
	}
	if foundTransactionStep != nil {
		t.Errorf("Expected TransactionStep From Method GetTransactionStepList: to be empty but got this: %v", foundTransactionStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTransactionStep(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_steps").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.UpdatedBy,     //4
		targetData.TransactionID, //5
		targetData.StepID,        //6
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateTransactionStep(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateTransactionStep", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionStepOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.TransactionID = nil
	err = UpdateTransactionStep(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionStepOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.StepID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateTransactionStep(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionStepOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.TransactionID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_steps").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.UpdatedBy,     //4
		targetData.TransactionID, //5
		targetData.StepID,        //6
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateTransactionStep(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTransactionStepByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_steps").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.UpdatedBy,     //4
		targetData.UUID,          //5
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateTransactionStepByUUID(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateTransactionStepByUUID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionStepByUUIDOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.TransactionID = nil
	err = UpdateTransactionStepByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionStepByUUIDOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.StepID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateTransactionStepByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionStepByUUIDOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.UUID = ""
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_steps").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.Description,   //3
		targetData.UpdatedBy,     //4
		targetData.UUID,          //5
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateTransactionStepByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionStep(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transaction_steps").WithArgs(
		targetData.TransactionID, //1
		targetData.StepID,        //2
		targetData.UUID,          //3
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.Description,   //6
		targetData.CreatedBy,     //7
	).WillReturnRows(pgxmock.NewRows([]string{"transaction_id", "step_id"}).AddRow(1, 1))
	mock.ExpectCommit()
	transactionID, stepID, err := InsertTransactionStep(mock, &targetData)
	if transactionID < 0 {
		t.Fatalf("transactionID should not be negative ID: %d", transactionID)
	}
	if stepID < 0 {
		t.Fatalf("stepID should not be negative ID: %d", stepID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertTransactionStep", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionStepOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.TransactionID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertTransactionStep(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionStepOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transaction_steps").WithArgs(
		targetData.TransactionID, //1
		targetData.StepID,        //2
		targetData.UUID,          //3
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.Description,   //6
		targetData.CreatedBy,     //7
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	transactionID, stepID, err := InsertTransactionStep(mock, &targetData)
	if transactionID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error transactionID: %d", transactionID)
	}
	if stepID >= 0 {
		t.Fatalf("Expecting -1 for stepID because of error stepID: %d", stepID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionStepOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transaction_steps").WithArgs(
		targetData.TransactionID, //1
		targetData.StepID,        //2
		targetData.UUID,          //3
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.Description,   //6
		targetData.CreatedBy,     //7
	).WillReturnRows(pgxmock.NewRows([]string{"transaction_id", "step_id"}).AddRow(-1, -1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	transactionID, stepID, err := InsertTransactionStep(mock, &targetData)
	if transactionID >= 0 {
		t.Fatalf("Expecting -1 for transactionID because of error transactionID: %d", transactionID)
	}
	if stepID >= 0 {
		t.Fatalf("Expecting -1 for stepID because of error stepID: %d", stepID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionSteps(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"transaction_steps"}, DBColumnsInsertTransactionSteps)
	targetData := TestAllData
	err = InsertTransactionSteps(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertTransactionSteps", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionStepsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"transaction_steps"}, DBColumnsInsertTransactionSteps).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertTransactionSteps(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTransactionStepToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"step_id = 1", "transaction_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WillReturnRows(mockRows)
	foundTransactionStepList, err := GetTransactionStepListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionStepListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundTransactionStepList[i]) == false {
			t.Errorf("Expected sourceData From Method GetTransactionStepListByPagination: %v is different from actual %v", sourceData, foundTransactionStepList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"transaction_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionStepList, err := GetTransactionStepListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionStepListByPagination", err)
	}
	if len(foundTransactionStepList) != 0 {
		t.Errorf("Expected From Method GetTransactionStepListByPagination: to be empty but got this: %v", foundTransactionStepList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionStepListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"transaction_id = -1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_steps").WillReturnRows(differentModelRows)
	foundTransactionStepList, err := GetTransactionStepListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionStepListByPagination", err)
	}
	if len(foundTransactionStepList) != 0 {
		t.Errorf("Expected From Method GetTransactionStepListByPagination: to be empty but got this: %v", foundTransactionStepList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalTransactionStepsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM transaction_steps").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalTransactionStepsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalTransactionStepsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalTransactionStepsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM transaction_steps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalTransactionStepsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalTransactionStepsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalTransactionStepsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
