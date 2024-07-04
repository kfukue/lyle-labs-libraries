package gethlylevlogjobs

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
	"id",                  //1
	"geth_process_job_id", //2
	"uuid",                //3
	"name",                //4
	"alternate_name",      //5
	"start_date",          //6
	"end_date",            //7
	"description",         //8
	"status_id",           //9
	"job_category_id",     //10
	"asset_id",            //11
	"chain_id",            //12
	"txn_hash",            //13
	"address_id",          //14
	"block_number",        //15
	"index_number",        //16
	"topics_str",          //17
	"created_by",          //18
	"created_at",          //19
	"updated_by",          //20
	"updated_at",          //21
}
var DBColumnsInsertGethProcessVlogJobList = []string{
	"geth_process_job_id", //1
	"uuid",                //2
	"name",                //3
	"alternate_name",      //4
	"start_date",          //5
	"end_date",            //6
	"description",         //7
	"status_id",           //8
	"job_category_id",     //9
	"asset_id",            //10
	"chain_id",            //11
	"txn_hash",            //12
	"address_id",          //13
	"block_number",        //14
	"index_number",        //15
	"topics_str",          //16
	"created_by",          //17
	"created_at",          //18
	"updated_by",          //19
	"updated_at",          //20
}

var TestData1 = GethProcessVlogJob{
	ID:             utils.Ptr[int](1),
	UUID:           "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:           "Import Swaps Using Liquidity Pool ID : 15, Name : PEPE/WETH Uniswap V2, SwapSig : PEPE/WETH Uniswap V2",
	AlternateName:  "Asset ID : 535, Name : PEPE, Import All Swaps For ERC20",
	StartDate:      utils.SampleCreatedAtTime,
	EndDate:        utils.SampleCreatedAtTime,
	Description:    "Original Start Block : 17046105, Current Block : 18768759",
	StatusID:       utils.Ptr[int](utils.SUCCESS_STRUCTURED_VALUE_ID),
	JobCategoryID:  utils.Ptr[int](utils.EOD_JOB_CATEGORY_STRUCTURED_VALUE_ID),
	AssetID:        utils.Ptr[int](1),
	ChainID:        utils.Ptr[int](2),
	TxnHash:        "0xf3d7dc2631d5d6612eb0b2b2dd0121eda3bd7e07da6527ee5e4a27c95a679456",
	AddressID:      utils.Ptr[int](1),
	BlockNumber:    utils.Ptr[uint64](20228646),
	IndexNumber:    utils.Ptr[uint](1),
	TopicsStrArray: []string{"Swap(address,address,int256,int256,uint160,uint128,int24)"},
	CreatedBy:      "SYSTEM",
	CreatedAt:      utils.SampleCreatedAtTime,
	UpdatedBy:      "SYSTEM",
	UpdatedAt:      utils.SampleCreatedAtTime,
}

var TestData2 = GethProcessVlogJob{
	ID:             utils.Ptr[int](2),
	UUID:           "880607ab-2833-4ad7-a231-b983a61cad34",
	Name:           "Import Swaps Using Liquidity Pool ID : 14, Name : PEPE/WETH Uniswap V3, SwapSig : PEPE/WETH Uniswap V3",
	AlternateName:  "Asset ID : 535, Name : PEPE, Import All Swaps For ERC20",
	StartDate:      utils.SampleCreatedAtTime,
	EndDate:        utils.SampleCreatedAtTime,
	Description:    "Original Start Block : 17046105, Current Block : 18768759",
	StatusID:       utils.Ptr[int](utils.SUCCESS_STRUCTURED_VALUE_ID),
	JobCategoryID:  utils.Ptr[int](utils.LIVE_JOB_CATEGORY_STRUCTURED_VALUE_ID),
	AssetID:        utils.Ptr[int](1),
	ChainID:        utils.Ptr[int](2),
	TxnHash:        "0x924342c099f079ad23b73081801a2eb186c5567a51a3fc01586111e5f922f603",
	AddressID:      utils.Ptr[int](1),
	BlockNumber:    utils.Ptr[uint64](20228646),
	IndexNumber:    utils.Ptr[uint](1),
	TopicsStrArray: []string{"Swap(address,address,int256,int256,uint160,uint128,int24)"},
	CreatedBy:      "SYSTEM",
	CreatedAt:      utils.SampleCreatedAtTime,
	UpdatedBy:      "SYSTEM",
	UpdatedAt:      utils.SampleCreatedAtTime,
}
var TestAllData = []GethProcessVlogJob{TestData1, TestData2}

func AddGethProcessVlogJobToMockRows(mock pgxmock.PgxPoolIface, dataList []GethProcessVlogJob) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,               //1
			data.GethProcessJobID, //2
			data.UUID,             //3
			data.Name,             //4
			data.AlternateName,    //5
			data.StartDate,        //6
			data.EndDate,          //7
			data.Description,      //8
			data.StatusID,         //9
			data.JobCategoryID,    //10
			data.AssetID,          //11
			data.ChainID,          //12
			data.TxnHash,          //13
			data.AddressID,        //14
			data.BlockNumber,      //15
			data.IndexNumber,      //16
			data.TopicsStrArray,   //17
			data.CreatedBy,        //18
			data.CreatedAt,        //19
			data.UpdatedBy,        //20
			data.UpdatedAt,        //21
		)
	}
	return rows
}

func TestGetGethProcessVlogJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethProcessVlogJob{targetData}
	gethProcessVlogJobID := targetData.ID
	mockRows := AddGethProcessVlogJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_vlog_jobs").WithArgs(*gethProcessVlogJobID).WillReturnRows(mockRows)
	foundGethProcessVlogJob, err := GetGethProcessVlogJob(mock, gethProcessVlogJobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessVlogJob", err)
	}
	if cmp.Equal(*foundGethProcessVlogJob, targetData) == false {
		t.Errorf("Expected GethProcessVlogJob From Method GetGethProcessVlogJob: %v is different from actual %v", foundGethProcessVlogJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessVlogJobForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessVlogJobID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_vlog_jobs").WithArgs(gethProcessVlogJobID).WillReturnRows(noRows)
	foundGethProcessVlogJob, err := GetGethProcessVlogJob(mock, &gethProcessVlogJobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessVlogJob", err)
	}
	if foundGethProcessVlogJob != nil {
		t.Errorf("Expected GethProcessVlogJob From Method GetGethProcessVlogJob: to be empty but got this: %v", foundGethProcessVlogJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessVlogJobForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessVlogJobID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_vlog_jobs").WithArgs(gethProcessVlogJobID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethProcessVlogJob, err := GetGethProcessVlogJob(mock, &gethProcessVlogJobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessVlogJob", err)
	}
	if foundGethProcessVlogJob != nil {
		t.Errorf("Expected GethProcessVlogJob From Method GetGethProcessVlogJob: to be empty but got this: %v", foundGethProcessVlogJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessVlogJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethProcessVlogJob{TestData1, TestData2}
	mockRows := AddGethProcessVlogJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_vlog_jobs").WillReturnRows(mockRows)
	foundGethProcessVlogJobs, err := GetGethProcessVlogJobList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessVlogJobList", err)
	}
	testGethProcessVlogJobs := TestAllData
	for i, foundGethProcessVlogJob := range foundGethProcessVlogJobs {
		if cmp.Equal(foundGethProcessVlogJob, testGethProcessVlogJobs[i]) == false {
			t.Errorf("Expected GethProcessVlogJob From Method GetGethProcessVlogJobList: %v is different from actual %v", foundGethProcessVlogJob, testGethProcessVlogJobs[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessVlogJobListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_vlog_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethProcessVlogJobs, err := GetGethProcessVlogJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessVlogJobList", err)
	}
	if len(foundGethProcessVlogJobs) != 0 {
		t.Errorf("Expected From Method GetGethProcessVlogJobs: to be empty but got this: %v", foundGethProcessVlogJobs)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethProcessVlogJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	gethProcessVlogJobID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_process_vlog_jobs").WithArgs(*gethProcessVlogJobID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethProcessVlogJob(mock, gethProcessVlogJobID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethProcessVlogJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethProcessVlogJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessVlogJobID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_process_vlog_jobs").WithArgs(gethProcessVlogJobID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethProcessVlogJob(mock, &gethProcessVlogJobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethProcessVlogJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_process_vlog_jobs").WithArgs(
		targetData.GethProcessJobID,         //1
		targetData.Name,                     //2
		targetData.AlternateName,            //3
		targetData.StartDate,                //4
		targetData.EndDate,                  //5
		targetData.Description,              //6
		targetData.StatusID,                 //7
		targetData.JobCategoryID,            //8
		targetData.AssetID,                  //9
		targetData.ChainID,                  //10
		targetData.TxnHash,                  //11
		targetData.AddressID,                //12
		targetData.BlockNumber,              //13
		targetData.IndexNumber,              //14
		pq.Array(targetData.TopicsStrArray), //15
		targetData.UpdatedBy,                //16
		targetData.ID,                       //17
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethProcessVlogJob(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethProcessVlogJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethProcessVlogJobOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethProcessVlogJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethProcessVlogJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_process_vlog_jobs").WithArgs(
		targetData.GethProcessJobID,         //1
		targetData.Name,                     //2
		targetData.AlternateName,            //3
		targetData.StartDate,                //4
		targetData.EndDate,                  //5
		targetData.Description,              //6
		targetData.StatusID,                 //7
		targetData.JobCategoryID,            //8
		targetData.AssetID,                  //9
		targetData.ChainID,                  //10
		targetData.TxnHash,                  //11
		targetData.AddressID,                //12
		targetData.BlockNumber,              //13
		targetData.IndexNumber,              //14
		pq.Array(targetData.TopicsStrArray), //15
		targetData.UpdatedBy,                //16
		targetData.ID,                       //17
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethProcessVlogJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessVlogJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_process_vlog_jobs").WithArgs(
		targetData.GethProcessJobID, //1
		targetData.Name,             //2
		targetData.AlternateName,    //3
		targetData.StartDate,        //4
		targetData.EndDate,          //5
		targetData.Description,      //6
		targetData.StatusID,         //7
		targetData.JobCategoryID,    //8
		targetData.AssetID,          //9
		targetData.ChainID,          //10
		targetData.TxnHash,          //11
		targetData.AddressID,        //12
		targetData.BlockNumber,      //13
		targetData.IndexNumber,      //14
		targetData.TopicsStrArray,   //15
		targetData.CreatedBy,        //16
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	gethProcessVlogJobID, err := InsertGethProcessVlogJob(mock, &targetData)
	if gethProcessVlogJobID < 0 {
		t.Fatalf("GethProcessVlogJobID should not be negative ID: %d", gethProcessVlogJobID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethProcessVlogJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessVlogJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_process_vlog_jobs").WithArgs(
		targetData.GethProcessJobID, //1
		targetData.Name,             //2
		targetData.AlternateName,    //3
		targetData.StartDate,        //4
		targetData.EndDate,          //5
		targetData.Description,      //6
		targetData.StatusID,         //7
		targetData.JobCategoryID,    //8
		targetData.AssetID,          //9
		targetData.ChainID,          //10
		targetData.TxnHash,          //11
		targetData.AddressID,        //12
		targetData.BlockNumber,      //13
		targetData.IndexNumber,      //14
		targetData.TopicsStrArray,   //15
		targetData.CreatedBy,        //16
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethProcessVlogJobID, err := InsertGethProcessVlogJob(mock, &targetData)
	if gethProcessVlogJobID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error GethProcessVlogJobID: %d", gethProcessVlogJobID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessVlogJobOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_process_vlog_jobs").WithArgs(
		targetData.GethProcessJobID, //1
		targetData.Name,             //2
		targetData.AlternateName,    //3
		targetData.StartDate,        //4
		targetData.EndDate,          //5
		targetData.Description,      //6
		targetData.StatusID,         //7
		targetData.JobCategoryID,    //8
		targetData.AssetID,          //9
		targetData.ChainID,          //10
		targetData.TxnHash,          //11
		targetData.AddressID,        //12
		targetData.BlockNumber,      //13
		targetData.IndexNumber,      //14
		targetData.TopicsStrArray,   //15
		targetData.CreatedBy,        //16
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethProcessVlogJobID, err := InsertGethProcessVlogJob(mock, &targetData)
	if gethProcessVlogJobID >= 0 {
		t.Fatalf("Expecting -1 for GethProcessVlogJobID because of error GethProcessVlogJobID: %d", gethProcessVlogJobID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessVlogJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_process_vlog_jobs"}, DBColumnsInsertGethProcessVlogJobList)
	targetData := TestAllData
	err = InsertGethProcessVlogJobList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethProcessVlogJobList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessVlogJobListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_process_vlog_jobs"}, DBColumnsInsertGethProcessVlogJobList).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertGethProcessVlogJobList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessVlogJobListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethProcessVlogJobToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_vlog_jobs").WillReturnRows(mockRows)
	foundChains, err := GetGethProcessVlogJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessVlogJobListByPagination", err)
	}
	testChains := dataList
	for i, foundChain := range foundChains {
		if cmp.Equal(foundChain, testChains[i]) == false {
			t.Errorf("Expected Chain From Method GetGethProcessVlogJobListByPagination: %v is different from actual %v", foundChain, testChains[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessVlogJobListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_vlog_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChains, err := GetGethProcessVlogJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessVlogJobListByPagination", err)
	}
	if len(foundChains) != 0 {
		t.Errorf("Expected From Method GetGethProcessVlogJobListByPagination: to be empty but got this: %v", foundChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethProcessVlogJobCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_process_vlog_jobs").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalGethProcessVlogJobCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalGethProcessVlogJobCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalGethProcessVlogJobCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethProcessVlogJobCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_process_vlog_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalGethProcessVlogJobCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalGethProcessVlogJobCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalGethProcessVlogJobCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
