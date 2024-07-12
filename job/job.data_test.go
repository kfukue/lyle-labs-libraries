package job

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

var DBColumns = []string{
	"id",                 //1
	"uuid",               //2
	"name",               //3
	"alternate_name",     //4
	"start_date",         //5
	"end_date",           //6
	"description",        //7
	"status_id",          //8
	"response_status",    //9
	"request_url",        //10
	"request_body",       //11
	"request_method",     //12
	"response_data",      //13
	"response_data_json", //14
	"job_category_id",    //15
	"created_by",         //16
	"created_at",         //17
	"updated_by",         //18
	"updated_at",         //19
}
var DBColumnsInsertJobs = []string{
	"uuid",               //1
	"name",               //2
	"alternate_name",     //3
	"start_date",         //4
	"end_date",           //5
	"description",        //6
	"status_id",          //7
	"response_status",    //8
	"request_url",        //9
	"request_body",       //10
	"request_method",     //11
	"response_data",      //12
	"response_data_json", //13
	"job_category_id",    //14
	"created_by",         //15
	"created_at",         //16
	"updated_by",         //17
	"updated_at",         //18
}

var TestData1 = Job{
	ID:               utils.Ptr[int](1),                                                                                                   //1
	UUID:             "01ef85e8-2c26-441e-8c7f-71d79518ad72",                                                                              //2
	Name:             "EOD Market Data Import - min date : 2021-01-01 00:00:00 max date : 2022-03-02 00:00:00",                            //3
	AlternateName:    "EOD Market Data Import - min date : 2021-01-01 00:00:00 max date : 2022-03-02 00:00:00",                            //4
	StartDate:        utils.SampleCreatedAtTime,                                                                                           //5
	EndDate:          utils.SampleCreatedAtTime,                                                                                           //6
	Description:      "ImportMarketDataPrices: Getting coingecko market data between min : 2021-01-01 00:00:00 max : 2022-03-02 00:00:00", //7
	StatusID:         utils.Ptr[int](52),                                                                                                  //8
	ResponseStatus:   "Success",                                                                                                           //9
	RequestUrl:       "https://wwww.google.com",                                                                                           //10
	RequestBody:      "Test Body",                                                                                                         //11
	RequestMethod:    "Test Method Get",                                                                                                   //12
	ResponseData:     "Test Data",                                                                                                         //13
	ResponseDataJson: nil,                                                                                                                 //14
	JobCategoryID:    utils.Ptr[int](1),                                                                                                   //15
	CreatedBy:        "SYSTEM",                                                                                                            //16
	CreatedAt:        utils.SampleCreatedAtTime,                                                                                           //17
	UpdatedBy:        "SYSTEM",                                                                                                            //18
	UpdatedAt:        utils.SampleCreatedAtTime,                                                                                           //19

}

var TestData2 = Job{
	ID:               utils.Ptr[int](2),                                                                                  //1
	UUID:             "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                                                             //2
	Name:             "Live Market Data Import - current date : 2022-03-08 11:24:08",                                     //3
	AlternateName:    "Live Market Data Import - current date : 2022-03-08 11:24:08",                                     //4
	StartDate:        utils.SampleCreatedAtTime,                                                                          //5
	EndDate:          utils.SampleCreatedAtTime,                                                                          //6
	Description:      "ImportLatestLiveMarketDataPrices: Getting live coingecko market data as of : 2022-03-08 11:24:08", //7
	StatusID:         utils.Ptr[int](52),                                                                                 //8
	ResponseStatus:   "Success",                                                                                          //9
	RequestUrl:       "https://wwww.google.com",                                                                          //10
	RequestBody:      "Test Body",                                                                                        //11
	RequestMethod:    "Test Method Get",                                                                                  //12
	ResponseData:     "Test Data",                                                                                        //13
	ResponseDataJson: nil,                                                                                                //14
	JobCategoryID:    utils.Ptr[int](1),                                                                                  //15
	CreatedBy:        "SYSTEM",                                                                                           //16
	CreatedAt:        utils.SampleCreatedAtTime,                                                                          //17
	UpdatedBy:        "SYSTEM",                                                                                           //18
	UpdatedAt:        utils.SampleCreatedAtTime,                                                                          //19
}
var TestAllData = []Job{TestData1, TestData2}

func AddJobToMockRows(mock pgxmock.PgxPoolIface, dataList []Job) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,               //1
			data.UUID,             //2
			data.Name,             //3
			data.AlternateName,    //4
			data.StartDate,        //5
			data.EndDate,          //6
			data.Description,      //7
			data.StatusID,         //8
			data.ResponseStatus,   //9
			data.RequestUrl,       //10
			data.RequestBody,      //11
			data.RequestMethod,    //12
			data.ResponseData,     //13
			data.ResponseDataJson, //14
			data.JobCategoryID,    //15
			data.CreatedBy,        //16
			data.CreatedAt,        //17
			data.UpdatedBy,        //18
			data.UpdatedAt,        //19
		)
	}
	return rows
}

func TestGetJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Job{targetData}
	jobID := targetData.ID
	mockRows := AddJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WithArgs(*jobID).WillReturnRows(mockRows)
	foundJob, err := GetJob(mock, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetJob", err)
	}
	if cmp.Equal(*foundJob, targetData) == false {
		t.Errorf("Expected Job From Method GetJob: %v is different from actual %v", foundJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetJobForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	jobID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WithArgs(jobID).WillReturnRows(noRows)
	foundJob, err := GetJob(mock, &jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetJob", err)
	}
	if foundJob != nil {
		t.Errorf("Expected Job From Method GetJob: to be empty but got this: %v", foundJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetJobForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	jobID := -1
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WithArgs(jobID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundJob, err := GetJob(mock, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetJob", err)
	}
	if foundJob != nil {
		t.Errorf("Expected Job From Method GetJob: to be empty but got this: %v", foundJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	jobID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM jobs").WithArgs(*jobID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveJob(mock, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	jobID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM jobs").WithArgs(jobID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveJob(mock, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WillReturnRows(mockRows)
	ids := []int{}
	foundJobList, err := GetJobList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetJobList", err)
	}
	for i, sourceJob := range dataList {
		if cmp.Equal(sourceJob, foundJobList[i]) == false {
			t.Errorf("Expected Job From Method GetJobList: %v is different from actual %v", sourceJob, foundJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetJobListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{}
	foundJobList, err := GetJobList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetJobList", err)
	}
	if len(foundJobList) != 0 {
		t.Errorf("Expected From Method GetJobList: to be empty but got this: %v", foundJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetJobListByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddJobToMockRows(mock, dataList)
	uuidList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WithArgs(pq.Array(uuidList)).WillReturnRows(mockRows)
	foundJobList, err := GetJobListByUUIDs(mock, uuidList)
	if err != nil {
		t.Fatalf("an error '%s' in GetJobListByUUIDs", err)
	}
	for i, sourceJob := range dataList {
		if cmp.Equal(sourceJob, foundJobList[i]) == false {
			t.Errorf("Expected Job From Method GetJobListByUUIDs: %v is different from actual %v", sourceJob, foundJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetJobListByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuidList := []string{"test", "test2"}
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WithArgs(pq.Array(uuidList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundJobList, err := GetJobListByUUIDs(mock, uuidList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetJobListByUUIDs", err)
	}
	if len(foundJobList) != 0 {
		t.Errorf("Expected From Method GetJobListByUUIDs: to be empty but got this: %v", foundJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddJobToMockRows(mock, dataList)
	diffInDate := 10
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WithArgs(diffInDate).WillReturnRows(mockRows)
	foundJobList, err := GetStartAndEndDateDiffJobList(mock, &diffInDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndDateDiffJobList", err)
	}
	for i, sourceJob := range dataList {
		if cmp.Equal(sourceJob, foundJobList[i]) == false {
			t.Errorf("Expected Job From Method GetStartAndEndDateDiffJobList: %v is different from actual %v", sourceJob, foundJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffJobListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := -1
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WithArgs(diffInDate).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundJobList, err := GetStartAndEndDateDiffJobList(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffJobList", err)
	}
	if len(foundJobList) != 0 {
		t.Errorf("Expected From Method GetStartAndEndDateDiffJobList: to be empty but got this: %v", foundJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE jobs").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.StartDate.Format(utils.LayoutPostgres), //3
		targetData.EndDate.Format(utils.LayoutPostgres),   //4
		targetData.Description,                            //5
		targetData.StatusID,                               //6
		targetData.ResponseStatus,                         //7
		targetData.RequestUrl,                             //8
		targetData.RequestBody,                            //9
		targetData.RequestMethod,                          //10
		targetData.ResponseData,                           //11
		targetData.ResponseDataJson,                       //12
		targetData.JobCategoryID,                          //13
		targetData.UpdatedBy,                              //14
		targetData.ID,                                     //15
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateJob(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE jobs").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.StartDate.Format(utils.LayoutPostgres), //3
		targetData.EndDate.Format(utils.LayoutPostgres),   //4
		targetData.Description,                            //5
		targetData.StatusID,                               //6
		targetData.ResponseStatus,                         //7
		targetData.RequestUrl,                             //8
		targetData.RequestBody,                            //9
		targetData.RequestMethod,                          //10
		targetData.ResponseData,                           //11
		targetData.ResponseDataJson,                       //12
		targetData.JobCategoryID,                          //13
		targetData.UpdatedBy,                              //14
		targetData.ID,                                     //15
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO jobs").WithArgs(
		targetData.UUID,          //1
		targetData.Name,          //2
		targetData.AlternateName, //3
		targetData.StartDate.Format(utils.LayoutPostgres), //4
		targetData.EndDate.Format(utils.LayoutPostgres),   //5
		targetData.Description,                            //6
		targetData.StatusID,                               //7
		targetData.ResponseStatus,                         //8
		targetData.RequestUrl,                             //9
		targetData.RequestBody,                            //10
		targetData.RequestMethod,                          //11
		targetData.ResponseData,                           //12
		targetData.ResponseDataJson,                       //13
		targetData.JobCategoryID,                          //14
		targetData.CreatedBy,                              //15
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit()
	jobID, newUUID, err := InsertJob(mock, &targetData)
	if jobID < 0 {
		t.Fatalf("jobID should not be negative ID: %d", jobID)
	}
	if newUUID == "" {
		t.Fatalf("newUUID should not be empty string: %s", newUUID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.StatusID = nil
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO jobs").WithArgs(
		targetData.UUID,          //1
		targetData.Name,          //2
		targetData.AlternateName, //3
		targetData.StartDate.Format(utils.LayoutPostgres), //4
		targetData.EndDate.Format(utils.LayoutPostgres),   //5
		targetData.Description,                            //6
		targetData.StatusID,                               //7
		targetData.ResponseStatus,                         //8
		targetData.RequestUrl,                             //9
		targetData.RequestBody,                            //10
		targetData.RequestMethod,                          //11
		targetData.ResponseData,                           //12
		targetData.ResponseDataJson,                       //13
		targetData.JobCategoryID,                          //14
		targetData.CreatedBy,                              //15
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	jobID, newUUID, err := InsertJob(mock, &targetData)
	if jobID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error jobID: %d", jobID)
	}
	if newUUID != "" {
		t.Fatalf("on failure newUUID should be empty string: %s", newUUID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertJobOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO jobs").WithArgs(
		targetData.UUID,          //1
		targetData.Name,          //2
		targetData.AlternateName, //3
		targetData.StartDate.Format(utils.LayoutPostgres), //4
		targetData.EndDate.Format(utils.LayoutPostgres),   //5
		targetData.Description,                            //6
		targetData.StatusID,                               //7
		targetData.ResponseStatus,                         //8
		targetData.RequestUrl,                             //9
		targetData.RequestBody,                            //10
		targetData.RequestMethod,                          //11
		targetData.ResponseData,                           //12
		targetData.ResponseDataJson,                       //13
		targetData.JobCategoryID,                          //14
		targetData.CreatedBy,                              //15
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	jobID, newUUID, err := InsertJob(mock, &targetData)
	if jobID >= 0 {
		t.Fatalf("Expecting -1 for jobID because of error jobID: %d", jobID)
	}
	if newUUID != "" {
		t.Fatalf("on failure newUUID should be empty string: %s", newUUID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"jobs"}, DBColumnsInsertJobs)
	targetData := TestAllData
	err = InsertJobList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertJobList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertJobListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"jobs"}, DBColumnsInsertJobs).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertJobList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetJobListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddJobToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"status_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WillReturnRows(mockRows)
	foundChains, err := GetJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetJobListByPagination", err)
	}
	testChains := dataList
	for i, foundChain := range foundChains {
		if cmp.Equal(foundChain, testChains[i]) == false {
			t.Errorf("Expected Chain From Method GetJobListByPagination: %v is different from actual %v", foundChain, testChains[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetJobListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"status_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChains, err := GetJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetJobListByPagination", err)
	}
	if len(foundChains) != 0 {
		t.Errorf("Expected From Method GetJobListByPagination: to be empty but got this: %v", foundChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM jobs").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
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
