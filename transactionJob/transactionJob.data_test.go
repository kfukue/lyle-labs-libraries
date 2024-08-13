package transactionjob

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
	"transaction_id",     //1
	"job_id",             //2
	"uuid",               //3
	"name",               //4
	"alternate_name",     //5
	"start_date",         //6
	"end_date",           //7
	"description",        //8
	"status_id",          //9
	"response_status",    //10
	"request_url",        //11
	"request_body",       //12
	"request_method",     //13
	"response_data",      //14
	"response_data_json", //15
	"created_by",         //16
	"created_at",         //17
	"updated_by",         //18
	"updated_at",         //19
}
var DBColumnsInsertTransactionJobs = []string{
	"transaction_id",     //1
	"job_id",             //2
	"uuid",               //3
	"name",               //4
	"alternate_name",     //5
	"start_date",         //6
	"end_date",           //7
	"description",        //8
	"status_id",          //9
	"response_status",    //10
	"request_url",        //11
	"request_body",       //12
	"request_method",     //13
	"response_data",      //14
	"response_data_json", //15
	"created_by",         //16
	"created_at",         //17
	"updated_by",         //18
	"updated_at",         //19
}

var TestData1 = TransactionJob{
	TransactionID:    utils.Ptr[int](1),                      //1
	JobID:            utils.Ptr[int](1),                      //2
	UUID:             "01ef85e8-2c26-441e-8c7f-71d79518ad72", //3
	Name:             "Sample Transaction Job 1",             //4
	AlternateName:    "Sample Transaction Job 1",             //5
	StartDate:        utils.SampleCreatedAtTime,              //6
	EndDate:          utils.SampleCreatedAtTime,              //7
	Description:      "",                                     //8
	StatusID:         utils.Ptr[int](52),                     //9
	ResponseStatus:   "Success",                              //10
	RequestUrl:       "https://wwww.google.com",              //11
	RequestBody:      "Test Body",                            //12
	RequestMethod:    "Test Method Get",                      //13
	ResponseData:     "Test Data",                            //14
	ResponseDataJson: nil,                                    //15
	CreatedBy:        "SYSTEM",                               //16
	CreatedAt:        utils.SampleCreatedAtTime,              //17
	UpdatedBy:        "SYSTEM",                               //18
	UpdatedAt:        utils.SampleCreatedAtTime,              //19

}

var TestData2 = TransactionJob{
	TransactionID:    utils.Ptr[int](2),                      //1
	JobID:            utils.Ptr[int](2),                      //2
	UUID:             "4f0d5402-7a7c-402d-a7fc-c56a02b13e03", //2
	Name:             "Sample Transaction Job 2",             //3
	AlternateName:    "Sample Transaction Job 2",             //4
	StartDate:        utils.SampleCreatedAtTime,              //6
	EndDate:          utils.SampleCreatedAtTime,              //7
	Description:      "",                                     //8
	StatusID:         utils.Ptr[int](52),                     //9
	ResponseStatus:   "Success",                              //10
	RequestUrl:       "https://wwww.google.com",              //11
	RequestBody:      "Test Body",                            //12
	RequestMethod:    "Test Method Get",                      //13
	ResponseData:     "Test Data",                            //14
	ResponseDataJson: nil,                                    //15
	CreatedBy:        "SYSTEM",                               //16
	CreatedAt:        utils.SampleCreatedAtTime,              //17
	UpdatedBy:        "SYSTEM",                               //18
	UpdatedAt:        utils.SampleCreatedAtTime,              //19
}
var TestAllData = []TransactionJob{TestData1, TestData2}

func AddTransactionJobToMockRows(mock pgxmock.PgxPoolIface, dataList []TransactionJob) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.TransactionID,    //1
			data.JobID,            //2
			data.UUID,             //3
			data.Name,             //4
			data.AlternateName,    //5
			data.StartDate,        //6
			data.EndDate,          //7
			data.Description,      //8
			data.StatusID,         //9
			data.ResponseStatus,   //10
			data.RequestUrl,       //11
			data.RequestBody,      //12
			data.RequestMethod,    //13
			data.ResponseData,     //14
			data.ResponseDataJson, //15
			data.CreatedBy,        //16
			data.CreatedAt,        //17
			data.UpdatedBy,        //18
			data.UpdatedAt,        //19
		)
	}
	return rows
}

func TestGetTransactionJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []TransactionJob{targetData}
	transactionID := targetData.TransactionID
	jobID := targetData.JobID
	mockRows := AddTransactionJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(*transactionID, *jobID).WillReturnRows(mockRows)
	foundTransactionJob, err := GetTransactionJob(mock, transactionID, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionJob", err)
	}
	if cmp.Equal(*foundTransactionJob, targetData) == false {
		t.Errorf("Expected TransactionJob From Method GetTransactionJob: %v is different from actual %v", foundTransactionJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := 999
	jobID := 1111
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(transactionID, jobID).WillReturnRows(noRows)
	foundTransactionJob, err := GetTransactionJob(mock, &transactionID, &jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionJob", err)
	}
	if foundTransactionJob != nil {
		t.Errorf("Expected TransactionJob From Method GetTransactionJob: to be empty but got this: %v", foundTransactionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	jobID := 1111
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(transactionID, jobID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionJob, err := GetTransactionJob(mock, &transactionID, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionJob", err)
	}
	if foundTransactionJob != nil {
		t.Errorf("Expected TransactionJob From Method GetTransactionJob: to be empty but got this: %v", foundTransactionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	jobID := 1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(transactionID, jobID).WillReturnRows(differentModelRows)
	foundTransactionJob, err := GetTransactionJob(mock, &transactionID, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionJob", err)
	}
	if foundTransactionJob != nil {
		t.Errorf("Expected TransactionJob From Method GetTransactionJob: to be empty but got this: %v", foundTransactionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []TransactionJob{targetData}
	transactionJobUUID := targetData.UUID
	mockRows := AddTransactionJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(transactionJobUUID).WillReturnRows(mockRows)
	foundTransactionJob, err := GetTransactionJobByUUID(mock, transactionJobUUID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionJobByUUID", err)
	}
	if cmp.Equal(*foundTransactionJob, targetData) == false {
		t.Errorf("Expected TransactionJob From Method GetTransactionJobByUUID: %v is different from actual %v", foundTransactionJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobByUUIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionJobUUID := "no-row-uuid"
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(transactionJobUUID).WillReturnRows(noRows)
	foundTransactionJob, err := GetTransactionJobByUUID(mock, transactionJobUUID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionJobByUUID", err)
	}
	if foundTransactionJob != nil {
		t.Errorf("Expected TransactionJob From Method GetTransactionJobByUUID: to be empty but got this: %v", foundTransactionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobByUUIDForSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionJobUUID := "row-error-uuid"
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(transactionJobUUID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionJob, err := GetTransactionJobByUUID(mock, transactionJobUUID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionJobByUUID", err)
	}
	if foundTransactionJob != nil {
		t.Errorf("Expected TransactionJob From Method GetTransactionJobByUUID: to be empty but got this: %v", foundTransactionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobByUUIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionJobUUID := "row-different-model-uuid"

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(transactionJobUUID).WillReturnRows(differentModelRows)
	foundTransactionJob, err := GetTransactionJobByUUID(mock, transactionJobUUID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionJobByUUID", err)
	}
	if foundTransactionJob != nil {
		t.Errorf("Expected TransactionJob From Method GetTransactionJobByUUID: to be empty but got this: %v", foundTransactionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	UUIDList := []string{TestData1.UUID, TestData2.UUID}

	mockRows := AddTransactionJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(pq.Array(UUIDList)).WillReturnRows(mockRows)
	foundTransactionJobList, err := GetTransactionJobsByUUIDs(mock, UUIDList)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionJobsByUUIDs", err)
	}
	for i, sourceTransactionJob := range dataList {
		if cmp.Equal(sourceTransactionJob, foundTransactionJobList[i]) == false {
			t.Errorf("Expected TransactionJob From Method GetTransactionJobsByUUIDs: %v is different from actual %v", sourceTransactionJob, foundTransactionJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(pq.Array(UUIDList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionJobList, err := GetTransactionJobsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionJobsByUUIDs", err)
	}
	if len(foundTransactionJobList) != 0 {
		t.Errorf("Expected From Method GetTransactionJobsByUUIDs: to be empty but got this: %v", foundTransactionJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobsByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WithArgs(pq.Array(UUIDList)).WillReturnRows(differentModelRows)
	foundTransactionJob, err := GetTransactionJobsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionJobsByUUIDs", err)
	}
	if foundTransactionJob != nil {
		t.Errorf("Expected TransactionJob From Method GetTransactionJobsByUUIDs: to be empty but got this: %v", foundTransactionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	transactionID := targetData.TransactionID
	jobID := targetData.JobID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_jobs").WithArgs(*transactionID, *jobID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveTransactionJob(mock, transactionID, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveTransactionJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	jobID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveTransactionJob(mock, &transactionID, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionID := -1
	jobID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_jobs").WithArgs(transactionID, jobID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveTransactionJob(mock, &transactionID, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionJobByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	transactionJobUUID := targetData.UUID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_jobs").WithArgs(transactionJobUUID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveTransactionJobByUUID(mock, transactionJobUUID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveTransactionJobByUUID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionJobByUUIDOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionJobUUID := "Fail-at-begining"
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveTransactionJobByUUID(mock, transactionJobUUID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTransactionJobByUUIDOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionJobUUID := "Fail-at-end"
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM transaction_jobs").WithArgs(transactionJobUUID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveTransactionJobByUUID(mock, transactionJobUUID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddTransactionJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WillReturnRows(mockRows)
	foundTransactionJobList, err := GetTransactionJobList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionJobList", err)
	}
	for i, sourceTransactionJob := range dataList {
		if cmp.Equal(sourceTransactionJob, foundTransactionJobList[i]) == false {
			t.Errorf("Expected TransactionJob From Method GetTransactionJobList: %v is different from actual %v", sourceTransactionJob, foundTransactionJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionJobList, err := GetTransactionJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionJobList", err)
	}
	if len(foundTransactionJobList) != 0 {
		t.Errorf("Expected From Method GetTransactionJobList: to be empty but got this: %v", foundTransactionJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WillReturnRows(differentModelRows)
	foundTransactionJob, err := GetTransactionJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionJobList", err)
	}
	if foundTransactionJob != nil {
		t.Errorf("Expected TransactionJob From Method GetTransactionJobList: to be empty but got this: %v", foundTransactionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTransactionJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_jobs").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.Description,    //5
		targetData.StatusID,       //6
		targetData.ResponseStatus, //7
		targetData.RequestUrl,     //8
		targetData.RequestBody,    //9
		targetData.RequestMethod,  //10
		targetData.ResponseData,   //11
		targetData.UpdatedBy,      //12
		targetData.TransactionID,  //13
		targetData.JobID,          //14
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateTransactionJob(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateTransactionJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionJobOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.TransactionID = nil
	err = UpdateTransactionJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.JobID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateTransactionJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.TransactionID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_jobs").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.Description,    //5
		targetData.StatusID,       //6
		targetData.ResponseStatus, //7
		targetData.RequestUrl,     //8
		targetData.RequestBody,    //9
		targetData.RequestMethod,  //10
		targetData.ResponseData,   //11
		targetData.UpdatedBy,      //12
		targetData.TransactionID,  //13
		targetData.JobID,          //14
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateTransactionJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTransactionJobByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_jobs").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.Description,    //5
		targetData.StatusID,       //6
		targetData.ResponseStatus, //7
		targetData.RequestUrl,     //8
		targetData.RequestBody,    //9
		targetData.RequestMethod,  //10
		targetData.ResponseData,   //11
		targetData.UpdatedBy,      //12
		targetData.UUID,           //13
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateTransactionJobByUUID(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateTransactionJobByUUID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionJobByUUIDOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.TransactionID = nil
	err = UpdateTransactionJobByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionJobByUUIDOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.JobID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateTransactionJobByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTransactionJobByUUIDOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.UUID = ""
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE transaction_jobs").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.Description,    //5
		targetData.StatusID,       //6
		targetData.ResponseStatus, //7
		targetData.RequestUrl,     //8
		targetData.RequestBody,    //9
		targetData.RequestMethod,  //10
		targetData.ResponseData,   //11
		targetData.UpdatedBy,      //12
		targetData.UUID,           //13
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateTransactionJobByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.UUID = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transaction_jobs").WithArgs(
		targetData.TransactionID,    //1
		targetData.JobID,            //2
		targetData.Name,             //3
		targetData.AlternateName,    //4
		targetData.StartDate,        //5
		targetData.EndDate,          //6
		targetData.Description,      //7
		targetData.StatusID,         //8
		targetData.ResponseStatus,   //9
		targetData.RequestUrl,       //10
		targetData.RequestBody,      //11
		targetData.RequestMethod,    //12
		targetData.ResponseData,     //13
		targetData.ResponseDataJson, //14
		targetData.CreatedBy,        //15
	).WillReturnRows(pgxmock.NewRows([]string{"transaction_id", "job_id"}).AddRow(1, 1))
	mock.ExpectCommit()
	transactionID, jobID, err := InsertTransactionJob(mock, &targetData)
	if transactionID < 0 {
		t.Fatalf("transactionID should not be negative ID: %d", transactionID)
	}
	if jobID < 0 {
		t.Fatalf("jobID should not be negative ID: %d", jobID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertTransactionJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.TransactionID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertTransactionJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transaction_jobs").WithArgs(
		targetData.TransactionID,    //1
		targetData.JobID,            //2
		targetData.Name,             //3
		targetData.AlternateName,    //4
		targetData.StartDate,        //5
		targetData.EndDate,          //6
		targetData.Description,      //7
		targetData.StatusID,         //8
		targetData.ResponseStatus,   //9
		targetData.RequestUrl,       //10
		targetData.RequestBody,      //11
		targetData.RequestMethod,    //12
		targetData.ResponseData,     //13
		targetData.ResponseDataJson, //14
		targetData.CreatedBy,        //15
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	transactionID, jobID, err := InsertTransactionJob(mock, &targetData)
	if transactionID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error transactionID: %d", transactionID)
	}
	if jobID >= 0 {
		t.Fatalf("Expecting -1 for jobID because of error jobID: %d", jobID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionJobOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO transaction_jobs").WithArgs(
		targetData.TransactionID,    //1
		targetData.JobID,            //2
		targetData.Name,             //3
		targetData.AlternateName,    //4
		targetData.StartDate,        //5
		targetData.EndDate,          //6
		targetData.Description,      //7
		targetData.StatusID,         //8
		targetData.ResponseStatus,   //9
		targetData.RequestUrl,       //10
		targetData.RequestBody,      //11
		targetData.RequestMethod,    //12
		targetData.ResponseData,     //13
		targetData.ResponseDataJson, //14
		targetData.CreatedBy,        //15
	).WillReturnRows(pgxmock.NewRows([]string{"transaction_id", "job_id"}).AddRow(-1, -1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	transactionID, jobID, err := InsertTransactionJob(mock, &targetData)
	if transactionID >= 0 {
		t.Fatalf("Expecting -1 for transactionID because of error transactionID: %d", transactionID)
	}
	if jobID >= 0 {
		t.Fatalf("Expecting -1 for jobID because of error jobID: %d", jobID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionJobs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"transaction_jobs"}, DBColumnsInsertTransactionJobs)
	targetData := TestAllData
	err = InsertTransactionJobs(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertTransactionJobs", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTransactionJobsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"transaction_jobs"}, DBColumnsInsertTransactionJobs).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertTransactionJobs(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTransactionJobToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"base_asset_id = 1", "frequency_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WillReturnRows(mockRows)
	foundTransactionJobList, err := GetTransactionJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransactionJobListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundTransactionJobList[i]) == false {
			t.Errorf("Expected sourceData From Method GetTransactionJobListByPagination: %v is different from actual %v", sourceData, foundTransactionJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"frequency_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTransactionJobList, err := GetTransactionJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionJobListByPagination", err)
	}
	if len(foundTransactionJobList) != 0 {
		t.Errorf("Expected From Method GetTransactionJobListByPagination: to be empty but got this: %v", foundTransactionJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionJobListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"frequency_id = -1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM transaction_jobs").WillReturnRows(differentModelRows)
	foundTransactionJobList, err := GetTransactionJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransactionJobListByPagination", err)
	}
	if len(foundTransactionJobList) != 0 {
		t.Errorf("Expected From Method GetTransactionJobListByPagination: to be empty but got this: %v", foundTransactionJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalTransactionJobsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM transaction_jobs").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalTransactionJobsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalTransactionJobsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalTransactionJobsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM transaction_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalTransactionJobsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalTransactionJobsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalTransactionJobsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
