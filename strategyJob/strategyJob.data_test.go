package strategyjob

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"strategy_id",        //1
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
var DBColumnsInsertStrategyJobList = []string{
	"strategy_id",        //1
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

var TestData1 = StrategyJob{
	StrategyID:       utils.Ptr[int](1),                                                                //1
	JobID:            utils.Ptr[int](1),                                                                //2
	UUID:             "01ef85e8-2c26-441e-8c7f-71d79518ad72",                                           //3
	Name:             "RetrievePriceFromStrategyMarketDataAssets - current date : 2022-06-30 16:10:04", //4
	AlternateName:    "RetrievePriceFromStrategyMarketDataAssets - current date : 2022-06-30 16:10:04", //5
	StartDate:        utils.SampleCreatedAtTime,                                                        //6
	EndDate:          utils.SampleCreatedAtTime,                                                        //7
	Description:      "",                                                                               //8
	StatusID:         utils.Ptr[int](52),                                                               //9
	ResponseStatus:   "Success",                                                                        //10
	RequestUrl:       "https://wwww.google.com",                                                        //11
	RequestBody:      "Test Body",                                                                      //12
	RequestMethod:    "Test Method Get",                                                                //13
	ResponseData:     "Test Data",                                                                      //14
	ResponseDataJson: nil,                                                                              //15
	CreatedBy:        "SYSTEM",                                                                         //16
	CreatedAt:        utils.SampleCreatedAtTime,                                                        //17
	UpdatedBy:        "SYSTEM",                                                                         //18
	UpdatedAt:        utils.SampleCreatedAtTime,                                                        //19

}

var TestData2 = StrategyJob{
	StrategyID:       utils.Ptr[int](2),                                                                                           //1
	JobID:            utils.Ptr[int](2),                                                                                           //2
	UUID:             "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                                                                      //2
	Name:             "RetrievePriceFromStrategyMarketDataAssets: Getting live coingecko market data as of : 2022-07-16 16:10:05", //3
	AlternateName:    "RetrievePriceFromStrategyMarketDataAssets: Getting live coingecko market data as of : 2022-07-16 16:10:05", //4
	StartDate:        utils.SampleCreatedAtTime,                                                                                   //6
	EndDate:          utils.SampleCreatedAtTime,                                                                                   //7
	Description:      "",                                                                                                          //8
	StatusID:         utils.Ptr[int](52),                                                                                          //9
	ResponseStatus:   "Success",                                                                                                   //10
	RequestUrl:       "https://wwww.google.com",                                                                                   //11
	RequestBody:      "Test Body",                                                                                                 //12
	RequestMethod:    "Test Method Get",                                                                                           //13
	ResponseData:     "Test Data",                                                                                                 //14
	ResponseDataJson: nil,                                                                                                         //15
	CreatedBy:        "SYSTEM",                                                                                                    //16
	CreatedAt:        utils.SampleCreatedAtTime,                                                                                   //17
	UpdatedBy:        "SYSTEM",                                                                                                    //18
	UpdatedAt:        utils.SampleCreatedAtTime,                                                                                   //19
}
var TestAllData = []StrategyJob{TestData1, TestData2}

func AddStrategyJobToMockRows(mock pgxmock.PgxPoolIface, dataList []StrategyJob) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.StrategyID,       //1
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

func TestGetStrategyJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []StrategyJob{targetData}
	strategyID := targetData.StrategyID
	jobID := targetData.JobID
	mockRows := AddStrategyJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WithArgs(*strategyID, *jobID).WillReturnRows(mockRows)
	foundStrategyJob, err := GetStrategyJob(mock, strategyID, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyJob", err)
	}
	if cmp.Equal(*foundStrategyJob, targetData) == false {
		t.Errorf("Expected StrategyJob From Method GetStrategyJob: %v is different from actual %v", foundStrategyJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := 999
	jobID := 1111
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WithArgs(strategyID, jobID).WillReturnRows(noRows)
	foundStrategyJob, err := GetStrategyJob(mock, &strategyID, &jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyJob", err)
	}
	if foundStrategyJob != nil {
		t.Errorf("Expected StrategyJob From Method GetStrategyJob: to be empty but got this: %v", foundStrategyJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -1
	jobID := 1111
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WithArgs(strategyID, jobID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyJob, err := GetStrategyJob(mock, &strategyID, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyJob", err)
	}
	if foundStrategyJob != nil {
		t.Errorf("Expected StrategyJob From Method GetStrategyJob: to be empty but got this: %v", foundStrategyJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -1
	jobID := 1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WithArgs(strategyID, jobID).WillReturnRows(differentModelRows)
	foundStrategyJob, err := GetStrategyJob(mock, &strategyID, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyJob", err)
	}
	if foundStrategyJob != nil {
		t.Errorf("Expected StrategyJob From Method GetStrategyJob: to be empty but got this: %v", foundStrategyJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobByStrategyID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []StrategyJob{targetData}
	strategyID := targetData.StrategyID
	mockRows := AddStrategyJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WithArgs(*strategyID).WillReturnRows(mockRows)
	foundStrategyJob, err := GetStrategyJobByStrategyID(mock, strategyID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyJobByStrategyID", err)
	}
	if cmp.Equal(*foundStrategyJob, targetData) == false {
		t.Errorf("Expected StrategyJob From Method GetStrategyJobByStrategyID: %v is different from actual %v", foundStrategyJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobByStrategyIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WithArgs(strategyID).WillReturnRows(noRows)
	foundStrategyJob, err := GetStrategyJobByStrategyID(mock, &strategyID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyJobByStrategyID", err)
	}
	if foundStrategyJob != nil {
		t.Errorf("Expected StrategyJob From Method GetStrategyJobByStrategyID: to be empty but got this: %v", foundStrategyJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobByStrategyIDForSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -1
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WithArgs(strategyID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyJob, err := GetStrategyJobByStrategyID(mock, &strategyID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyJobByStrategyID", err)
	}
	if foundStrategyJob != nil {
		t.Errorf("Expected StrategyJob From Method GetStrategyJobByStrategyID: to be empty but got this: %v", foundStrategyJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobByStrategyIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -11
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WithArgs(strategyID).WillReturnRows(differentModelRows)
	foundStrategyJob, err := GetStrategyJobByStrategyID(mock, &strategyID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyJobByStrategyID", err)
	}
	if foundStrategyJob != nil {
		t.Errorf("Expected StrategyJob From Method GetStrategyJobByStrategyID: to be empty but got this: %v", foundStrategyJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStrategyJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	strategyID := targetData.StrategyID
	jobID := targetData.JobID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM strategy_jobs").WithArgs(*strategyID, *jobID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveStrategyJob(mock, strategyID, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveStrategyJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStrategyJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -1
	jobID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveStrategyJob(mock, &strategyID, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStrategyJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -1
	jobID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM strategy_jobs").WithArgs(strategyID, jobID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveStrategyJob(mock, &strategyID, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStrategyJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WillReturnRows(mockRows)
	foundStrategyJobList, err := GetStrategyJobList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyJobList", err)
	}
	for i, sourceStrategyJob := range dataList {
		if cmp.Equal(sourceStrategyJob, foundStrategyJobList[i]) == false {
			t.Errorf("Expected StrategyJob From Method GetStrategyJobList: %v is different from actual %v", sourceStrategyJob, foundStrategyJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyJobList, err := GetStrategyJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyJobList", err)
	}
	if len(foundStrategyJobList) != 0 {
		t.Errorf("Expected From Method GetStrategyJobList: to be empty but got this: %v", foundStrategyJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WillReturnRows(differentModelRows)
	foundStrategyJob, err := GetStrategyJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyJobList", err)
	}
	if foundStrategyJob != nil {
		t.Errorf("Expected StrategyJob From Method GetStrategyJobList: to be empty but got this: %v", foundStrategyJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateStrategyJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE strategy_jobs").WithArgs(
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
		targetData.UpdatedBy,                              //13
		targetData.StrategyID,                             //14
		targetData.JobID,                                  //15
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateStrategyJob(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateStrategyJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStrategyJobOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.StrategyID = nil
	err = UpdateStrategyJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStrategyJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.JobID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateStrategyJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStrategyJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.StrategyID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE strategy_jobs").WithArgs(
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
		targetData.UpdatedBy,                              //13
		targetData.StrategyID,                             //14
		targetData.JobID,                                  //15
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateStrategyJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.UUID = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO strategy_jobs").WithArgs(
		targetData.StrategyID,       //1
		targetData.JobID,            //2
		targetData.UUID,             //3
		targetData.Name,             //4
		targetData.AlternateName,    //5
		targetData.StartDate,        //6
		targetData.EndDate,          //7
		targetData.Description,      //8
		targetData.StatusID,         //9
		targetData.ResponseStatus,   //10
		targetData.RequestUrl,       //11
		targetData.RequestBody,      //12
		targetData.RequestMethod,    //13
		targetData.ResponseData,     //14
		targetData.ResponseDataJson, //15
		targetData.CreatedBy,        //16
	).WillReturnRows(pgxmock.NewRows([]string{"strategy_id", "job_id"}).AddRow(1, 1))
	mock.ExpectCommit()
	strategyID, jobID, err := InsertStrategyJob(mock, &targetData)
	if strategyID < 0 {
		t.Fatalf("strategyID should not be negative ID: %d", strategyID)
	}
	if jobID < 0 {
		t.Fatalf("jobID should not be negative ID: %d", jobID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertStrategyJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.StrategyID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertStrategyJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO strategy_jobs").WithArgs(
		targetData.StrategyID,       //1
		targetData.JobID,            //2
		targetData.UUID,             //3
		targetData.Name,             //4
		targetData.AlternateName,    //5
		targetData.StartDate,        //6
		targetData.EndDate,          //7
		targetData.Description,      //8
		targetData.StatusID,         //9
		targetData.ResponseStatus,   //10
		targetData.RequestUrl,       //11
		targetData.RequestBody,      //12
		targetData.RequestMethod,    //13
		targetData.ResponseData,     //14
		targetData.ResponseDataJson, //15
		targetData.CreatedBy,        //16
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	strategyID, jobID, err := InsertStrategyJob(mock, &targetData)
	if strategyID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error strategyID: %d", strategyID)
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

func TestInsertStrategyJobOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO strategy_jobs").WithArgs(
		targetData.StrategyID,       //1
		targetData.JobID,            //2
		targetData.UUID,             //3
		targetData.Name,             //4
		targetData.AlternateName,    //5
		targetData.StartDate,        //6
		targetData.EndDate,          //7
		targetData.Description,      //8
		targetData.StatusID,         //9
		targetData.ResponseStatus,   //10
		targetData.RequestUrl,       //11
		targetData.RequestBody,      //12
		targetData.RequestMethod,    //13
		targetData.ResponseData,     //14
		targetData.ResponseDataJson, //15
		targetData.CreatedBy,        //16
	).WillReturnRows(pgxmock.NewRows([]string{"strategy_id", "job_id"}).AddRow(-1, -1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	strategyID, jobID, err := InsertStrategyJob(mock, &targetData)
	if strategyID >= 0 {
		t.Fatalf("Expecting -1 for strategyID because of error strategyID: %d", strategyID)
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

func TestInsertStrategyJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"strategy_jobs"}, DBColumnsInsertStrategyJobList)
	targetData := TestAllData
	err = InsertStrategyJobList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertStrategyJobList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyJobListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"strategy_jobs"}, DBColumnsInsertStrategyJobList).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertStrategyJobList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddStrategyJobToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"base_asset_id = 1", "frequency_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WillReturnRows(mockRows)
	foundStrategyJobList, err := GetStrategyJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyJobListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundStrategyJobList[i]) == false {
			t.Errorf("Expected sourceData From Method GetStrategyJobListByPagination: %v is different from actual %v", sourceData, foundStrategyJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyJobList, err := GetStrategyJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyJobListByPagination", err)
	}
	if len(foundStrategyJobList) != 0 {
		t.Errorf("Expected From Method GetStrategyJobListByPagination: to be empty but got this: %v", foundStrategyJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyJobListByPaginationForCollectRowsErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM strategy_jobs").WillReturnRows(differentModelRows)
	foundStrategyJobList, err := GetStrategyJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyJobListByPagination", err)
	}
	if len(foundStrategyJobList) != 0 {
		t.Errorf("Expected From Method GetStrategyJobListByPagination: to be empty but got this: %v", foundStrategyJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalStrategyJobsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM strategy_jobs").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalStrategyJobsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalStrategyJobsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalStrategyJobsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM strategy_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalStrategyJobsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalStrategyJobsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalStrategyJobsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
