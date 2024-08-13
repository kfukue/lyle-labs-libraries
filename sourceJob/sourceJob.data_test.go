package sourcejob

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"source_id",   //1
	"job_id",      //2
	"uuid",        //3
	"description", //4
	"created_by",  //5
	"created_at",  //6
	"updated_by",  //7
	"updated_at",  //8
}
var DBColumnsInsertSourceJobs = []string{
	"source_id",   //1
	"job_id",      //2
	"uuid",        //3
	"description", //4
	"created_by",  //5
	"created_at",  //6
	"updated_by",  //7
	"updated_at",  //8
}

var TestData1 = SourceJob{
	SourceID:    utils.Ptr[int](1),                      //1
	JobID:       utils.Ptr[int](1),                      //2
	UUID:        "01ef85e8-2c26-441e-8c7f-71d79518ad72", //3
	Description: "",                                     //4
	CreatedBy:   "SYSTEM",                               //5
	CreatedAt:   utils.SampleCreatedAtTime,              //6
	UpdatedBy:   "SYSTEM",                               //7
	UpdatedAt:   utils.SampleCreatedAtTime,              //8

}

var TestData2 = SourceJob{
	SourceID:    utils.Ptr[int](2),                      //1
	JobID:       utils.Ptr[int](2),                      //2
	UUID:        "4f0d5402-7a7c-402d-a7fc-c56a02b13e03", //3
	Description: "",                                     //4
	CreatedBy:   "SYSTEM",                               //5
	CreatedAt:   utils.SampleCreatedAtTime,              //6
	UpdatedBy:   "SYSTEM",                               //7
	UpdatedAt:   utils.SampleCreatedAtTime,              //8
}
var TestAllData = []SourceJob{TestData1, TestData2}

func AddSourceJobToMockRows(mock pgxmock.PgxPoolIface, dataList []SourceJob) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.SourceID,    //1
			data.JobID,       //2
			data.UUID,        //3
			data.Description, //4
			data.CreatedBy,   //5
			data.CreatedAt,   //6
			data.UpdatedBy,   //7
			data.UpdatedAt,   //8
		)
	}
	return rows
}

func TestGetSourceJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []SourceJob{targetData}
	sourceID := targetData.SourceID
	jobID := targetData.JobID
	mockRows := AddSourceJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WithArgs(*sourceID, *jobID).WillReturnRows(mockRows)
	foundSourceJob, err := GetSourceJob(mock, sourceID, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetSourceJob", err)
	}
	if cmp.Equal(*foundSourceJob, targetData) == false {
		t.Errorf("Expected SourceJob From Method GetSourceJob: %v is different from actual %v", foundSourceJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := 999
	jobID := 1111
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WithArgs(sourceID, jobID).WillReturnRows(noRows)
	foundSourceJob, err := GetSourceJob(mock, &sourceID, &jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetSourceJob", err)
	}
	if foundSourceJob != nil {
		t.Errorf("Expected SourceJob From Method GetSourceJob: to be empty but got this: %v", foundSourceJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	jobID := 1111
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WithArgs(sourceID, jobID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundSourceJob, err := GetSourceJob(mock, &sourceID, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceJob", err)
	}
	if foundSourceJob != nil {
		t.Errorf("Expected SourceJob From Method GetSourceJob: to be empty but got this: %v", foundSourceJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	jobID := 1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WithArgs(sourceID, jobID).WillReturnRows(differentModelRows)
	foundSourceJob, err := GetSourceJob(mock, &sourceID, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceJob", err)
	}
	if foundSourceJob != nil {
		t.Errorf("Expected SourceJob From Method GetSourceJob: to be empty but got this: %v", foundSourceJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobBySourceID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []SourceJob{targetData}
	sourceID := targetData.SourceID
	mockRows := AddSourceJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WithArgs(*sourceID).WillReturnRows(mockRows)
	foundSourceJob, err := GetSourceJobBySourceID(mock, sourceID)
	if err != nil {
		t.Fatalf("an error '%s' in GetSourceJobBySourceID", err)
	}
	if cmp.Equal(*foundSourceJob, targetData) == false {
		t.Errorf("Expected SourceJob From Method GetSourceJobBySourceID: %v is different from actual %v", foundSourceJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobBySourceIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WithArgs(sourceID).WillReturnRows(noRows)
	foundSourceJob, err := GetSourceJobBySourceID(mock, &sourceID)
	if err != nil {
		t.Fatalf("an error '%s' in GetSourceJobBySourceID", err)
	}
	if foundSourceJob != nil {
		t.Errorf("Expected SourceJob From Method GetSourceJobBySourceID: to be empty but got this: %v", foundSourceJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobBySourceIDForSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WithArgs(sourceID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundSourceJob, err := GetSourceJobBySourceID(mock, &sourceID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceJobBySourceID", err)
	}
	if foundSourceJob != nil {
		t.Errorf("Expected SourceJob From Method GetSourceJobBySourceID: to be empty but got this: %v", foundSourceJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobBySourceIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -999

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WithArgs(sourceID).WillReturnRows(differentModelRows)
	foundSourceJob, err := GetSourceJobBySourceID(mock, &sourceID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceJobBySourceID", err)
	}
	if foundSourceJob != nil {
		t.Errorf("Expected SourceJob From Method GetSourceJobBySourceID: to be empty but got this: %v", foundSourceJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveSourceJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	sourceID := targetData.SourceID
	jobID := targetData.JobID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM source_jobs").WithArgs(*sourceID, *jobID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveSourceJob(mock, sourceID, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveSourceJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveSourceJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	jobID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveSourceJob(mock, &sourceID, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveSourceJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	jobID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM source_jobs").WithArgs(sourceID, jobID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveSourceJob(mock, &sourceID, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddSourceJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WillReturnRows(mockRows)
	foundSourceJobList, err := GetSourceJobList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetSourceJobList", err)
	}
	for i, sourceSourceJob := range dataList {
		if cmp.Equal(sourceSourceJob, foundSourceJobList[i]) == false {
			t.Errorf("Expected SourceJob From Method GetSourceJobList: %v is different from actual %v", sourceSourceJob, foundSourceJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundSourceJobList, err := GetSourceJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceJobList", err)
	}
	if len(foundSourceJobList) != 0 {
		t.Errorf("Expected From Method GetSourceJobList: to be empty but got this: %v", foundSourceJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WillReturnRows(differentModelRows)
	foundSourceJob, err := GetSourceJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceJobList", err)
	}
	if foundSourceJob != nil {
		t.Errorf("Expected SourceJob From Method GetSourceJobList: to be empty but got this: %v", foundSourceJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateSourceJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE source_jobs").WithArgs(
		targetData.Description, //1
		targetData.UpdatedBy,   //2
		targetData.SourceID,    //3
		targetData.JobID,       //4
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateSourceJob(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateSourceJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateSourceJobOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.SourceID = nil
	err = UpdateSourceJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateSourceJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.JobID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateSourceJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateSourceJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.SourceID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE source_jobs").WithArgs(
		targetData.Description, //1
		targetData.UpdatedBy,   //2
		targetData.SourceID,    //3
		targetData.JobID,       //4
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateSourceJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSourceJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.UUID = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO source_jobs").WithArgs(
		targetData.SourceID,    //1
		targetData.JobID,       //2
		targetData.Description, //3
		targetData.CreatedBy,   //4
	).WillReturnRows(pgxmock.NewRows([]string{"source_id", "job_id"}).AddRow(1, 1))
	mock.ExpectCommit()
	sourceID, jobID, err := InsertSourceJob(mock, &targetData)
	if sourceID < 0 {
		t.Fatalf("sourceID should not be negative ID: %d", sourceID)
	}
	if jobID < 0 {
		t.Fatalf("jobID should not be negative ID: %d", jobID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertSourceJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSourceJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.SourceID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertSourceJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSourceJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO source_jobs").WithArgs(
		targetData.SourceID,    //1
		targetData.JobID,       //2
		targetData.Description, //3
		targetData.CreatedBy,   //4
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	sourceID, jobID, err := InsertSourceJob(mock, &targetData)
	if sourceID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error sourceID: %d", sourceID)
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

func TestInsertSourceJobOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO source_jobs").WithArgs(
		targetData.SourceID,    //1
		targetData.JobID,       //2
		targetData.Description, //3
		targetData.CreatedBy,   //4
	).WillReturnRows(pgxmock.NewRows([]string{"source_id", "job_id"}).AddRow(-1, -1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	sourceID, jobID, err := InsertSourceJob(mock, &targetData)
	if sourceID >= 0 {
		t.Fatalf("Expecting -1 for sourceID because of error sourceID: %d", sourceID)
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

func TestInsertSourceJobs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"source_jobs"}, DBColumnsInsertSourceJobs)
	targetData := TestAllData
	err = InsertSourceJobs(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertSourceJobs", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSourceJobsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"source_jobs"}, DBColumnsInsertSourceJobs).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertSourceJobs(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddSourceJobToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"source_id = 1", "job_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WillReturnRows(mockRows)
	foundSourceJobList, err := GetSourceJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetSourceJobListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundSourceJobList[i]) == false {
			t.Errorf("Expected sourceData From Method GetSourceJobListByPagination: %v is different from actual %v", sourceData, foundSourceJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"source_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundSourceJobList, err := GetSourceJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceJobListByPagination", err)
	}
	if len(foundSourceJobList) != 0 {
		t.Errorf("Expected From Method GetSourceJobListByPagination: to be empty but got this: %v", foundSourceJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceJobListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"source_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM source_jobs").WillReturnRows(differentModelRows)
	foundSourceJobList, err := GetSourceJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceJobListByPagination", err)
	}
	if len(foundSourceJobList) != 0 {
		t.Errorf("Expected From Method GetSourceJobListByPagination: to be empty but got this: %v", foundSourceJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalSourceJobsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM source_jobs").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalSourceJobsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalSourceJobsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalSourceJobsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM source_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalSourceJobsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalSourceJobsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalSourceJobsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
