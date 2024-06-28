package dexjob

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

var columns = []string{
	"id",                 //1
	"job_id",             //2
	"uuid",               //3
	"name",               //4
	"alternate_name",     //5
	"start_date",         //6
	"end_date",           //7
	"description",        //8
	"status_id",          //9
	"chain_id",           //10
	"exchange_id",        //11
	"transaction_hashes", //12
	"created_by",         //13
	"created_at",         //14
	"updated_by",         //15
	"updated_at",         //16
}

var columnsInsertList = []string{
	"job_id",             //1
	"uuid",               //2
	"name",               //3
	"alternate_name",     //4
	"start_date",         //5
	"end_date",           //6
	"description",        //7
	"status_id",          //8
	"chain_id",           //9
	"exchange_id",        //10
	"transaction_hashes", //11
	"created_by",         //12
	"created_at",         //13
	"updated_by",         //14
	"updated_at",         //15
}
var testTxns = []string{"0x1706fb8bf07d31852bbb0e5d1c8b0378c60b87a1fdccc36eab706603d67522d4", "0x0dc5e228f2520f74abfab4a97867dbf54e5bfc73e5a1d2a68aa79420ae1dd611"}

var data1 = DexTxnJob{
	ID:                utils.Ptr[int](1),
	JobID:             utils.Ptr[int](1),
	UUID:              "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:              "Test Run Job 1",
	AlternateName:     "Testing",
	StartDate:         utils.SampleCreatedAtTime,
	EndDate:           utils.SampleCreatedAtTime,
	Description:       "",
	StatusID:          utils.Ptr[int](1),
	ChainID:           utils.Ptr[int](1),
	ExchangeID:        utils.Ptr[int](1),
	TransactionHashes: testTxns,
	CreatedBy:         "SYSTEM",
	CreatedAt:         utils.SampleCreatedAtTime,
	UpdatedBy:         "SYSTEM",
	UpdatedAt:         utils.SampleCreatedAtTime,
}

var data2 = DexTxnJob{
	ID:                utils.Ptr[int](2),
	JobID:             utils.Ptr[int](3),
	UUID:              "880607ab-2833-4ad7-a231-b983a61cad34",
	Name:              "Test Run Job 2",
	AlternateName:     "Testing 2",
	StartDate:         utils.SampleCreatedAtTime,
	EndDate:           utils.SampleCreatedAtTime,
	Description:       "",
	StatusID:          utils.Ptr[int](2),
	ChainID:           utils.Ptr[int](3),
	ExchangeID:        utils.Ptr[int](5),
	TransactionHashes: testTxns,
	CreatedBy:         "SYSTEM",
	CreatedAt:         utils.SampleCreatedAtTime,
	UpdatedBy:         "SYSTEM",
	UpdatedAt:         utils.SampleCreatedAtTime,
}
var allData = []DexTxnJob{data1, data2}

func AddDexTxnJobToMockRows(mock pgxmock.PgxPoolIface, dataList []DexTxnJob) *pgxmock.Rows {
	rows := mock.NewRows(columns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                //1
			data.JobID,             //2
			data.UUID,              //3
			data.Name,              //4
			data.AlternateName,     //5
			data.StartDate,         //6
			data.EndDate,           //8
			data.Description,       //7
			data.StatusID,          //9
			data.ChainID,           //10
			data.ExchangeID,        //11
			data.TransactionHashes, //12
			data.CreatedBy,         //13
			data.CreatedAt,         //14
			data.UpdatedBy,         //15
			data.UpdatedAt,         //16

		)
	}
	return rows
}

func TestGetDexTxnJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data2
	dataList := []DexTxnJob{targetData}
	dexTxnID := targetData.ID
	mockRows := AddDexTxnJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM dex_txn_jobs").WithArgs(*dexTxnID).WillReturnRows(mockRows)
	foundDexTxnJob, err := GetDexTxnJob(mock, dexTxnID)
	if err != nil {
		t.Fatalf("an error '%s' in GetDexTxnJob", err)
	}
	if cmp.Equal(*foundDexTxnJob, targetData) == false {
		t.Errorf("Expected DexTxnJob From Method GetDexTxnJob: %v is different from actual %v", foundDexTxnJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDexTxnJobForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dexTxnID := 999
	noRows := pgxmock.NewRows(columns)
	mock.ExpectQuery("^SELECT (.+) FROM dex_txn_jobs").WithArgs(dexTxnID).WillReturnRows(noRows)
	foundDexTxnJob, err := GetDexTxnJob(mock, &dexTxnID)
	if err != nil {
		t.Fatalf("an error '%s' in GetDexTxnJob", err)
	}
	if foundDexTxnJob != nil {
		t.Errorf("Expected DexTxnJob From Method GetDexTxnJob: to be empty but got this: %v", foundDexTxnJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDexTxnJobForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dexTxnID := -1
	mock.ExpectQuery("^SELECT (.+) FROM dex_txn_jobs").WithArgs(dexTxnID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundDexTxnJob, err := GetDexTxnJob(mock, &dexTxnID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDexTxnJob", err)
	}
	if foundDexTxnJob != nil {
		t.Errorf("Expected DexTxnJob From Method GetDexTxnJob: to be empty but got this: %v", foundDexTxnJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDexTxnJobByJobId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []DexTxnJob{data1, data2}
	mockRows := AddDexTxnJobToMockRows(mock, dataList)
	jobID := 2
	mock.ExpectQuery("^SELECT (.+) FROM dex_txn_jobs").WithArgs(jobID).WillReturnRows(mockRows)
	foundDexTxnJobs, err := GetDexTxnJobByJobId(mock, &jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetDexTxnJobByJobId", err)
	}
	testDexTxnJobs := allData
	for i, foundDexTxnJob := range foundDexTxnJobs {
		if cmp.Equal(foundDexTxnJob, testDexTxnJobs[i]) == false {
			t.Errorf("Expected DexTxnJob From Method GetDexTxnJobByJobId: %v is different from actual %v", foundDexTxnJob, testDexTxnJobs[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDexTxnJobByJobIdForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	jobID := -1
	mock.ExpectQuery("^SELECT (.+) FROM dex_txn_jobs").WithArgs(jobID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundDexTxnJobs, err := GetDexTxnJobByJobId(mock, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDexTxnJobByJobId", err)
	}
	if len(foundDexTxnJobs) != 0 {
		t.Errorf("Expected From Method GetDexTxnJobByJobId: to be empty but got this: %v", foundDexTxnJobs)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDexTxnJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []DexTxnJob{data1, data2}
	mockRows := AddDexTxnJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM dex_txn_jobs").WillReturnRows(mockRows)
	foundDexTxnJobs, err := GetDexTxnJobList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetDexTxnJobList", err)
	}
	testDexTxnJobs := allData
	for i, foundDexTxnJob := range foundDexTxnJobs {
		if cmp.Equal(foundDexTxnJob, testDexTxnJobs[i]) == false {
			t.Errorf("Expected DexTxnJob From Method GetDexTxnJobList: %v is different from actual %v", foundDexTxnJob, testDexTxnJobs[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDexTxnJobListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM dex_txn_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundDexTxnJobs, err := GetDexTxnJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDexTxnJobList", err)
	}
	if len(foundDexTxnJobs) != 0 {
		t.Errorf("Expected From Method GetDexTxnJobList: to be empty but got this: %v", foundDexTxnJobs)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveDexTxnJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	dexTxnID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM dex_txn_jobs").WithArgs(*dexTxnID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveDexTxnJob(mock, dexTxnID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveDexTxnJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveDexTxnJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dexTxnID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM dex_txn_jobs").WithArgs(dexTxnID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveDexTxnJob(mock, &dexTxnID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateDexTxnJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE dex_txn_jobs").WithArgs(
		targetData.Name,                        //1
		targetData.AlternateName,               //2
		targetData.StartDate,                   //3
		targetData.EndDate,                     //4
		targetData.Description,                 //5
		targetData.StatusID,                    //6
		targetData.ChainID,                     //7
		targetData.ExchangeID,                  //8
		pq.Array(targetData.TransactionHashes), //9
		targetData.UpdatedBy,                   //10
		targetData.ID,                          //11
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateDexTxnJob(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateDexTxnJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateDexTxnJobOnFailure(t *testing.T) {
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
	mock.ExpectExec("^UPDATE dex_txn_jobs").WithArgs(
		targetData.Name,                        //1
		targetData.AlternateName,               //2
		targetData.StartDate,                   //3
		targetData.EndDate,                     //4
		targetData.Description,                 //5
		targetData.StatusID,                    //6
		targetData.ChainID,                     //7
		targetData.ExchangeID,                  //8
		pq.Array(targetData.TransactionHashes), //9
		targetData.UpdatedBy,                   //10
		targetData.ID,                          //11
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateDexTxnJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertDexTxnJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO dex_txn_jobs").WithArgs(
		targetData.JobID,                       //1
		targetData.Name,                        //2
		targetData.AlternateName,               //3
		targetData.StartDate,                   //4
		targetData.EndDate,                     //5
		targetData.Description,                 //6
		targetData.StatusID,                    //7
		targetData.ChainID,                     //8
		targetData.ExchangeID,                  //9
		pq.Array(targetData.TransactionHashes), //10
		targetData.CreatedBy,                   //11
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	chainID, err := InsertDexTxnJob(mock, &targetData)
	if chainID < 0 {
		t.Fatalf("chainID should not be negative ID: %d", chainID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertDexTxnJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertDexTxnJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO dex_txn_jobs").WithArgs(
		targetData.JobID,                       //1
		targetData.Name,                        //2
		targetData.AlternateName,               //3
		targetData.StartDate,                   //4
		targetData.EndDate,                     //5
		targetData.Description,                 //6
		targetData.StatusID,                    //7
		targetData.ChainID,                     //8
		targetData.ExchangeID,                  //9
		pq.Array(targetData.TransactionHashes), //10
		targetData.CreatedBy,                   //11
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	chainID, err := InsertDexTxnJob(mock, &targetData)
	if chainID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error chainID: %d", chainID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertDexTxnJobOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO dex_txn_jobs").WithArgs(
		targetData.JobID,                       //1
		targetData.Name,                        //2
		targetData.AlternateName,               //3
		targetData.StartDate,                   //4
		targetData.EndDate,                     //5
		targetData.Description,                 //6
		targetData.StatusID,                    //7
		targetData.ChainID,                     //8
		targetData.ExchangeID,                  //9
		pq.Array(targetData.TransactionHashes), //10
		targetData.CreatedBy,                   //11
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	chainID, err := InsertDexTxnJob(mock, &targetData)
	if chainID >= 0 {
		t.Fatalf("Expecting -1 for chainID because of error chainID: %d", chainID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertDexTxnJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"dex_txn_jobs"}, columnsInsertList)
	targetData := allData
	err = InsertDexTxnJobList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertDexTxnJobList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertDexTxnJobListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"dex_txn_jobs"}, columnsInsertList).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := allData
	err = InsertDexTxnJobList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
