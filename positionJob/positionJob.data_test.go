package positionjob

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
	"position_id",        //1
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
var DBColumnsInsertPositionJobs = []string{
	"position_id",        //1
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

var TestData1 = PositionJob{
	PositionID:       utils.Ptr[int](1),                      //1
	JobID:            utils.Ptr[int](1),                      //2
	UUID:             "01ef85e8-2c26-441e-8c7f-71d79518ad72", //3
	Name:             "Sample Position Job 1",                //4
	AlternateName:    "Sample Position Job 1",                //5
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

var TestData2 = PositionJob{
	PositionID:       utils.Ptr[int](2),                      //1
	JobID:            utils.Ptr[int](2),                      //2
	UUID:             "4f0d5402-7a7c-402d-a7fc-c56a02b13e03", //2
	Name:             "Sample Position Job 2",                //3
	AlternateName:    "Sample Position Job 2",                //4
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
var TestAllData = []PositionJob{TestData1, TestData2}

func AddPositionJobToMockRows(mock pgxmock.PgxPoolIface, dataList []PositionJob) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.PositionID,       //1
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

func TestGetPositionJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []PositionJob{targetData}
	positionID := targetData.PositionID
	jobID := targetData.JobID
	mockRows := AddPositionJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(*positionID, *jobID).WillReturnRows(mockRows)
	foundPositionJob, err := GetPositionJob(mock, positionID, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPositionJob", err)
	}
	if cmp.Equal(*foundPositionJob, targetData) == false {
		t.Errorf("Expected PositionJob From Method GetPositionJob: %v is different from actual %v", foundPositionJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionID := 999
	jobID := 1111
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(positionID, jobID).WillReturnRows(noRows)
	foundPositionJob, err := GetPositionJob(mock, &positionID, &jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPositionJob", err)
	}
	if foundPositionJob != nil {
		t.Errorf("Expected PositionJob From Method GetPositionJob: to be empty but got this: %v", foundPositionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionID := -1
	jobID := 1111
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(positionID, jobID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPositionJob, err := GetPositionJob(mock, &positionID, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionJob", err)
	}
	if foundPositionJob != nil {
		t.Errorf("Expected PositionJob From Method GetPositionJob: to be empty but got this: %v", foundPositionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionID := -1
	jobID := 1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(positionID, jobID).WillReturnRows(differentModelRows)
	foundPositionJob, err := GetPositionJob(mock, &positionID, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionJob", err)
	}
	if foundPositionJob != nil {
		t.Errorf("Expected PositionJob From Method GetPositionJob: to be empty but got this: %v", foundPositionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []PositionJob{targetData}
	positionJobUUID := targetData.UUID
	mockRows := AddPositionJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(positionJobUUID).WillReturnRows(mockRows)
	foundPositionJob, err := GetPositionJobByUUID(mock, positionJobUUID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPositionJobByUUID", err)
	}
	if cmp.Equal(*foundPositionJob, targetData) == false {
		t.Errorf("Expected PositionJob From Method GetPositionJobByUUID: %v is different from actual %v", foundPositionJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobByUUIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionJobUUID := "no-row-uuid"
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(positionJobUUID).WillReturnRows(noRows)
	foundPositionJob, err := GetPositionJobByUUID(mock, positionJobUUID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPositionJobByUUID", err)
	}
	if foundPositionJob != nil {
		t.Errorf("Expected PositionJob From Method GetPositionJobByUUID: to be empty but got this: %v", foundPositionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobByUUIDForSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionJobUUID := "row-error-uuid"
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(positionJobUUID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPositionJob, err := GetPositionJobByUUID(mock, positionJobUUID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionJobByUUID", err)
	}
	if foundPositionJob != nil {
		t.Errorf("Expected PositionJob From Method GetPositionJobByUUID: to be empty but got this: %v", foundPositionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobByUUIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionJobUUID := "row-different-model-uuid"

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(positionJobUUID).WillReturnRows(differentModelRows)
	foundPositionJob, err := GetPositionJobByUUID(mock, positionJobUUID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionJobByUUID", err)
	}
	if foundPositionJob != nil {
		t.Errorf("Expected PositionJob From Method GetPositionJobByUUID: to be empty but got this: %v", foundPositionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	UUIDList := []string{TestData1.UUID, TestData2.UUID}

	mockRows := AddPositionJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(pq.Array(UUIDList)).WillReturnRows(mockRows)
	foundPositionJobList, err := GetPositionJobsByUUIDs(mock, UUIDList)
	if err != nil {
		t.Fatalf("an error '%s' in GetPositionJobsByUUIDs", err)
	}
	for i, sourcePositionJob := range dataList {
		if cmp.Equal(sourcePositionJob, foundPositionJobList[i]) == false {
			t.Errorf("Expected PositionJob From Method GetPositionJobsByUUIDs: %v is different from actual %v", sourcePositionJob, foundPositionJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(pq.Array(UUIDList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPositionJobList, err := GetPositionJobsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionJobsByUUIDs", err)
	}
	if len(foundPositionJobList) != 0 {
		t.Errorf("Expected From Method GetPositionJobsByUUIDs: to be empty but got this: %v", foundPositionJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobsByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WithArgs(pq.Array(UUIDList)).WillReturnRows(differentModelRows)
	foundPositionJob, err := GetPositionJobsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionJobsByUUIDs", err)
	}
	if foundPositionJob != nil {
		t.Errorf("Expected PositionJob From Method GetPositionJobsByUUIDs: to be empty but got this: %v", foundPositionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePositionJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	positionID := targetData.PositionID
	jobID := targetData.JobID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM position_jobs").WithArgs(*positionID, *jobID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemovePositionJob(mock, positionID, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in RemovePositionJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePositionJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionID := -1
	jobID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemovePositionJob(mock, &positionID, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePositionJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionID := -1
	jobID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM position_jobs").WithArgs(positionID, jobID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemovePositionJob(mock, &positionID, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePositionJobByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	positionJobUUID := targetData.UUID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM position_jobs").WithArgs(positionJobUUID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemovePositionJobByUUID(mock, positionJobUUID)
	if err != nil {
		t.Fatalf("an error '%s' in RemovePositionJobByUUID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePositionJobByUUIDOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionJobUUID := "Fail-at-begining"
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemovePositionJobByUUID(mock, positionJobUUID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePositionJobByUUIDOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionJobUUID := "Fail-at-end"
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM position_jobs").WithArgs(positionJobUUID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemovePositionJobByUUID(mock, positionJobUUID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddPositionJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WillReturnRows(mockRows)
	foundPositionJobList, err := GetPositionJobList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetPositionJobList", err)
	}
	for i, sourcePositionJob := range dataList {
		if cmp.Equal(sourcePositionJob, foundPositionJobList[i]) == false {
			t.Errorf("Expected PositionJob From Method GetPositionJobList: %v is different from actual %v", sourcePositionJob, foundPositionJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPositionJobList, err := GetPositionJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionJobList", err)
	}
	if len(foundPositionJobList) != 0 {
		t.Errorf("Expected From Method GetPositionJobList: to be empty but got this: %v", foundPositionJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WillReturnRows(differentModelRows)
	foundPositionJob, err := GetPositionJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionJobList", err)
	}
	if foundPositionJob != nil {
		t.Errorf("Expected PositionJob From Method GetPositionJobList: to be empty but got this: %v", foundPositionJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdatePositionJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE position_jobs").WithArgs(
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
		targetData.PositionID,     //13
		targetData.JobID,          //14
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdatePositionJob(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdatePositionJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePositionJobOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.PositionID = nil
	err = UpdatePositionJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePositionJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.JobID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdatePositionJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePositionJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.PositionID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE position_jobs").WithArgs(
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
		targetData.PositionID,     //13
		targetData.JobID,          //14
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdatePositionJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdatePositionJobByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE position_jobs").WithArgs(
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
	err = UpdatePositionJobByUUID(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdatePositionJobByUUID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePositionJobByUUIDOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.PositionID = nil
	err = UpdatePositionJobByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePositionJobByUUIDOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.JobID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdatePositionJobByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePositionJobByUUIDOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.UUID = ""
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE position_jobs").WithArgs(
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
	err = UpdatePositionJobByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPositionJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.UUID = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO position_jobs").WithArgs(
		targetData.PositionID,       //1
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
	).WillReturnRows(pgxmock.NewRows([]string{"position_id", "job_id"}).AddRow(1, 1))
	mock.ExpectCommit()
	positionID, jobID, err := InsertPositionJob(mock, &targetData)
	if positionID < 0 {
		t.Fatalf("positionID should not be negative ID: %d", positionID)
	}
	if jobID < 0 {
		t.Fatalf("jobID should not be negative ID: %d", jobID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertPositionJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPositionJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.PositionID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertPositionJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPositionJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO position_jobs").WithArgs(
		targetData.PositionID,       //1
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
	positionID, jobID, err := InsertPositionJob(mock, &targetData)
	if positionID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error positionID: %d", positionID)
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

func TestInsertPositionJobOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO position_jobs").WithArgs(
		targetData.PositionID,       //1
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
	).WillReturnRows(pgxmock.NewRows([]string{"position_id", "job_id"}).AddRow(-1, -1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	positionID, jobID, err := InsertPositionJob(mock, &targetData)
	if positionID >= 0 {
		t.Fatalf("Expecting -1 for positionID because of error positionID: %d", positionID)
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

func TestInsertPositionJobs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"position_jobs"}, DBColumnsInsertPositionJobs)
	targetData := TestAllData
	err = InsertPositionJobs(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertPositionJobs", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPositionJobsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"position_jobs"}, DBColumnsInsertPositionJobs).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertPositionJobs(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddPositionJobToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"base_asset_id = 1", "frequency_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WillReturnRows(mockRows)
	foundPositionJobList, err := GetPositionJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetPositionJobListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundPositionJobList[i]) == false {
			t.Errorf("Expected sourceData From Method GetPositionJobListByPagination: %v is different from actual %v", sourceData, foundPositionJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPositionJobList, err := GetPositionJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionJobListByPagination", err)
	}
	if len(foundPositionJobList) != 0 {
		t.Errorf("Expected From Method GetPositionJobListByPagination: to be empty but got this: %v", foundPositionJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionJobListByPaginationForCollectRowsErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM position_jobs").WillReturnRows(differentModelRows)
	foundPositionJobList, err := GetPositionJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionJobListByPagination", err)
	}
	if len(foundPositionJobList) != 0 {
		t.Errorf("Expected From Method GetPositionJobListByPagination: to be empty but got this: %v", foundPositionJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalPositionJobsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM position_jobs").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalPositionJobsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalPositionJobsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalPositionJobsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM position_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalPositionJobsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalPositionJobsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalPositionJobsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
