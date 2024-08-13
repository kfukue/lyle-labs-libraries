package marketdatajob

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
	"market_data_id",     //1
	"jobid",              //2
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
var DBColumnsInsertMarketDataJobs = []string{
	"market_data_id",     //1
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

var TestData1 = MarketDataJob{
	MarketDataID:     utils.Ptr[int](1),                                                                                                   //1
	JobID:            utils.Ptr[int](1),                                                                                                   //2
	UUID:             "01ef85e8-2c26-441e-8c7f-71d79518ad72",                                                                              //3
	Name:             "EOD Market Data Import - min date : 2021-01-01 00:00:00 max date : 2022-03-02 00:00:00",                            //4
	AlternateName:    "EOD Market Data Import - min date : 2021-01-01 00:00:00 max date : 2022-03-02 00:00:00",                            //5
	StartDate:        utils.SampleCreatedAtTime,                                                                                           //6
	EndDate:          utils.SampleCreatedAtTime,                                                                                           //7
	Description:      "ImportMarketDataPrices: Getting coingecko market data between min : 2021-01-01 00:00:00 max : 2022-03-02 00:00:00", //8
	StatusID:         utils.Ptr[int](52),                                                                                                  //9
	ResponseStatus:   "Success",                                                                                                           //10
	RequestUrl:       "https://wwww.google.com",                                                                                           //11
	RequestBody:      "Test Body",                                                                                                         //12
	RequestMethod:    "Test Method Get",                                                                                                   //13
	ResponseData:     "Test Data",                                                                                                         //14
	ResponseDataJson: nil,                                                                                                                 //15
	CreatedBy:        "SYSTEM",                                                                                                            //16
	CreatedAt:        utils.SampleCreatedAtTime,                                                                                           //17
	UpdatedBy:        "SYSTEM",                                                                                                            //18
	UpdatedAt:        utils.SampleCreatedAtTime,                                                                                           //19

}

var TestData2 = MarketDataJob{
	MarketDataID:     utils.Ptr[int](1),                                                                                  //1
	JobID:            utils.Ptr[int](1),                                                                                  //2
	UUID:             "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                                                             //3
	Name:             "Live Market Data Import - current date : 2022-03-08 11:24:08",                                     //4
	AlternateName:    "Live Market Data Import - current date : 2022-03-08 11:24:08",                                     //5
	StartDate:        utils.SampleCreatedAtTime,                                                                          //6
	EndDate:          utils.SampleCreatedAtTime,                                                                          //7
	Description:      "ImportLatestLiveMarketDataPrices: Getting live coingecko market data as of : 2022-03-08 11:24:08", //8
	StatusID:         utils.Ptr[int](52),                                                                                 //9
	ResponseStatus:   "Success",                                                                                          //10
	RequestUrl:       "https://wwww.google.com",                                                                          //11
	RequestBody:      "Test Body",                                                                                        //12
	RequestMethod:    "Test Method Get",                                                                                  //13
	ResponseData:     "Test Data",                                                                                        //14
	ResponseDataJson: nil,                                                                                                //15
	CreatedBy:        "SYSTEM",                                                                                           //16
	CreatedAt:        utils.SampleCreatedAtTime,                                                                          //17
	UpdatedBy:        "SYSTEM",                                                                                           //18
	UpdatedAt:        utils.SampleCreatedAtTime,                                                                          //19
}
var TestAllData = []MarketDataJob{TestData1, TestData2}

func AddMarketDataJobToMockRows(mock pgxmock.PgxPoolIface, dataList []MarketDataJob) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.MarketDataID,     //1
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

func TestGetMarketDataJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []MarketDataJob{targetData}
	marketDataID := targetData.MarketDataID
	jobID := targetData.JobID
	mockRows := AddMarketDataJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WithArgs(*marketDataID, *jobID).WillReturnRows(mockRows)
	foundMarketDataJob, err := GetMarketDataJob(mock, marketDataID, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketDataJob", err)
	}
	if cmp.Equal(*foundMarketDataJob, targetData) == false {
		t.Errorf("Expected MarketDataJob From Method GetMarketDataJob: %v is different from actual %v", foundMarketDataJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := 999
	jobID := 111
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WithArgs(marketDataID, jobID).WillReturnRows(noRows)
	foundMarketDataJob, err := GetMarketDataJob(mock, &marketDataID, &jobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketDataJob", err)
	}
	if foundMarketDataJob != nil {
		t.Errorf("Expected MarketDataJob From Method GetMarketDataJob: to be empty but got this: %v", foundMarketDataJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1
	jobID := -1
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WithArgs(marketDataID, jobID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataJob, err := GetMarketDataJob(mock, &marketDataID, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataJob", err)
	}
	if foundMarketDataJob != nil {
		t.Errorf("Expected MarketDataJob From Method GetMarketDataJob: to be empty but got this: %v", foundMarketDataJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1
	jobID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WithArgs(marketDataID, jobID).WillReturnRows(differentModelRows)
	foundMarketDataJob, err := GetMarketDataJob(mock, &marketDataID, &jobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataJob", err)
	}
	if foundMarketDataJob != nil {
		t.Errorf("Expected MarketDataJob From Method GetMarketDataJob: to be empty but got this: %v", foundMarketDataJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobByMarketDataID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []MarketDataJob{targetData}
	marketDataID := targetData.MarketDataID
	mockRows := AddMarketDataJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WithArgs(*marketDataID).WillReturnRows(mockRows)
	foundMarketDataJob, err := GetMarketDataJobByMarketDataID(mock, marketDataID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketDataJobByMarketDataID", err)
	}
	if cmp.Equal(*foundMarketDataJob, targetData) == false {
		t.Errorf("Expected MarketDataJob From Method GetMarketDataJobByMarketDataID: %v is different from actual %v", foundMarketDataJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobByMarketDataIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WithArgs(marketDataID).WillReturnRows(noRows)
	foundMarketDataJob, err := GetMarketDataJobByMarketDataID(mock, &marketDataID)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketDataJobByMarketDataID", err)
	}
	if foundMarketDataJob != nil {
		t.Errorf("Expected MarketDataJob From Method GetMarketDataJobByMarketDataID: to be empty but got this: %v", foundMarketDataJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobByMarketDataIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WithArgs(marketDataID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataJob, err := GetMarketDataJobByMarketDataID(mock, &marketDataID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataJobByMarketDataID", err)
	}
	if foundMarketDataJob != nil {
		t.Errorf("Expected MarketDataJob From Method GetMarketDataJobByMarketDataID: to be empty but got this: %v", foundMarketDataJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobByMarketDataIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WithArgs(marketDataID).WillReturnRows(differentModelRows)
	foundMarketDataJob, err := GetMarketDataJobByMarketDataID(mock, &marketDataID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataJobByMarketDataID", err)
	}
	if foundMarketDataJob != nil {
		t.Errorf("Expected MarketDataJob From Method GetMarketDataJobByMarketDataID: to be empty but got this: %v", foundMarketDataJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketDataJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	marketDataID := targetData.MarketDataID
	jobID := targetData.JobID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data_jobs").WithArgs(*marketDataID, *jobID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveMarketDataJob(mock, marketDataID, jobID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveMarketDataJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketDataJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.MarketDataID = utils.Ptr[int](-1)
	marketDataID := -1
	jobID := 1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveMarketDataJob(mock, &marketDataID, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveMarketDataJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	marketDataID := -1
	jobID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM market_data_jobs").WithArgs(marketDataID, jobID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveMarketDataJob(mock, &marketDataID, &jobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddMarketDataJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WillReturnRows(mockRows)
	foundMarketDataJobList, err := GetMarketDataJobList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketDataJobList", err)
	}
	for i, sourceMarketDataJob := range dataList {
		if cmp.Equal(sourceMarketDataJob, foundMarketDataJobList[i]) == false {
			t.Errorf("Expected MarketDataJob From Method GetMarketDataJobList: %v is different from actual %v", sourceMarketDataJob, foundMarketDataJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataJobList, err := GetMarketDataJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataJobList", err)
	}
	if len(foundMarketDataJobList) != 0 {
		t.Errorf("Expected From Method GetMarketDataJobList: to be empty but got this: %v", foundMarketDataJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WillReturnRows(differentModelRows)
	foundMarketDataJob, err := GetMarketDataJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataJobList", err)
	}
	if foundMarketDataJob != nil {
		t.Errorf("Expected MarketDataJob From Method GetMarketDataJobList: to be empty but got this: %v", foundMarketDataJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateMarketDataJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE market_data_jobs").WithArgs(
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
		targetData.MarketDataID,   //13
		targetData.JobID,          //14
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateMarketDataJob(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateMarketDataJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateMarketDataJobOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.MarketDataID = nil
	err = UpdateMarketDataJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateMarketDataJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.MarketDataID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateMarketDataJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateMarketDataJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.MarketDataID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE market_data_jobs").WithArgs(
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
		targetData.MarketDataID,   //13
		targetData.JobID,          //14
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateMarketDataJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO market_data_jobs").WithArgs(
		targetData.MarketDataID,  //1
		targetData.JobID,         //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.StartDate.Format(utils.LayoutPostgres), //5
		targetData.EndDate.Format(utils.LayoutPostgres),   //6
		targetData.Description,                            //7
		targetData.StatusID,                               //8
		targetData.ResponseStatus,                         //9
		targetData.RequestUrl,                             //10
		targetData.RequestBody,                            //11
		targetData.RequestMethod,                          //12
		targetData.ResponseData,                           //13
		targetData.ResponseDataJson,                       //14
		targetData.CreatedBy,                              //15
	).WillReturnRows(pgxmock.NewRows([]string{"market_data_id", "job_id"}).AddRow(1, 2))
	mock.ExpectCommit()
	marketDataID, jobID, err := InsertMarketDataJob(mock, &targetData)
	if marketDataID < 0 {
		t.Fatalf("marketDataID should not be negative ID: %d", marketDataID)
	}
	if jobID < 0 {
		t.Fatalf("jobID should not be negative ID: %d", jobID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertMarketDataJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.MarketDataID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertMarketDataJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.StatusID = nil
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO market_data_jobs").WithArgs(
		targetData.MarketDataID,  //1
		targetData.JobID,         //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.StartDate.Format(utils.LayoutPostgres), //5
		targetData.EndDate.Format(utils.LayoutPostgres),   //6
		targetData.Description,                            //7
		targetData.StatusID,                               //8
		targetData.ResponseStatus,                         //9
		targetData.RequestUrl,                             //10
		targetData.RequestBody,                            //11
		targetData.RequestMethod,                          //12
		targetData.ResponseData,                           //13
		targetData.ResponseDataJson,                       //14
		targetData.CreatedBy,                              //15
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	marketDataID, jobID, err := InsertMarketDataJob(mock, &targetData)
	if marketDataID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error marketDataID: %d", marketDataID)
	}
	if jobID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error jobID: %d", jobID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataJobOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO market_data_jobs").WithArgs(
		targetData.MarketDataID,  //1
		targetData.JobID,         //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.StartDate.Format(utils.LayoutPostgres), //5
		targetData.EndDate.Format(utils.LayoutPostgres),   //6
		targetData.Description,                            //7
		targetData.StatusID,                               //8
		targetData.ResponseStatus,                         //9
		targetData.RequestUrl,                             //10
		targetData.RequestBody,                            //11
		targetData.RequestMethod,                          //12
		targetData.ResponseData,                           //13
		targetData.ResponseDataJson,                       //14
		targetData.CreatedBy,                              //15
	).WillReturnRows(pgxmock.NewRows([]string{"market_data_id", "job_id"}).AddRow(1, 2))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	marketDataID, jobID, err := InsertMarketDataJob(mock, &targetData)
	if marketDataID >= 0 {
		t.Fatalf("Expecting -1 for marketDataID because of error marketDataID: %d", marketDataID)
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

func TestInsertMarketDataJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"market_data_jobs"}, DBColumnsInsertMarketDataJobs)
	targetData := TestAllData
	err = InsertMarketDataJobList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertMarketDataJobList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertMarketDataJobListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"market_data_jobs"}, DBColumnsInsertMarketDataJobs).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertMarketDataJobList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddMarketDataJobToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"job_id = 1", "status_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WillReturnRows(mockRows)
	foundMarketDataJobList, err := GetMarketDataJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetMarketDataJobListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundMarketDataJobList[i]) == false {
			t.Errorf("Expected sourceData From Method GetMarketDataJobListByPagination: %v is different from actual %v", sourceData, foundMarketDataJobList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"job_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataJobList, err := GetMarketDataJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataJobListByPagination", err)
	}
	if len(foundMarketDataJobList) != 0 {
		t.Errorf("Expected From Method GetMarketDataJobListByPagination: to be empty but got this: %v", foundMarketDataJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMarketDataJobListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"job_id = -1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM market_data_jobs").WillReturnRows(differentModelRows)
	foundMarketDataJobList, err := GetMarketDataJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetMarketDataJobListByPagination", err)
	}
	if len(foundMarketDataJobList) != 0 {
		t.Errorf("Expected From Method GetMarketDataJobListByPagination: to be empty but got this: %v", foundMarketDataJobList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMarketDataJobsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM market_data_jobs").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalMarketDataJobsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalMarketDataJobsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalMarketDataJobsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM market_data_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalMarketDataJobsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalMarketDataJobsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalMarketDataJobsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
