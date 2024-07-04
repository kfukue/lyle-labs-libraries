package gethlylejobs

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
	"id",                 //1
	"uuid",               //2
	"name",               //3
	"alternate_name",     //4
	"start_date",         //5
	"end_date",           //6
	"description",        //7
	"status_id",          //8
	"job_category_id",    //9
	"import_type_id",     //10
	"chain_id",           //11
	"start_block_number", //12
	"end_block_number",   //13
	"created_by",         //14
	"created_at",         //15
	"updated_by",         //16
	"updated_at",         //17
	"asset_id",           //18
}
var DBColumnsInsertGethProcessJobList = []string{
	"uuid",               //1
	"name",               //2
	"alternate_name",     //3
	"start_date",         //4
	"end_date",           //5
	"description",        //6
	"status_id",          //7
	"job_category_id",    //8
	"import_type_id",     //9
	"chain_id",           //10
	"start_block_number", //11
	"end_block_number",   //12
	"created_by",         //13
	"created_at",         //14
	"updated_by",         //15
	"updated_at",         //16
	"asset_id",           //17
}

var TestData1 = GethProcessJob{
	ID:               utils.Ptr[int](1),
	UUID:             "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:             "Asset ID : 535, Name : PEPE, Import All Swaps For ERC20",
	AlternateName:    "Asset ID : 535, Name : PEPE, Import All Swaps For ERC20",
	StartDate:        utils.SampleCreatedAtTime,
	EndDate:          utils.SampleCreatedAtTime,
	Description:      "",
	StatusID:         utils.Ptr[int](utils.SUCCESS_STRUCTURED_VALUE_ID),
	JobCategoryID:    utils.Ptr[int](utils.EOD_JOB_CATEGORY_STRUCTURED_VALUE_ID),
	ImportTypeID:     utils.Ptr[int](utils.MINER_IMPORT_TYPE_STRUCTURED_VALUE_TYPE_ID),
	ChainID:          utils.Ptr[int](1),
	StartBlockNumber: utils.Ptr[uint64](17046105),
	EndBlockNumber:   utils.Ptr[uint64](17046305),
	CreatedBy:        "SYSTEM",
	CreatedAt:        utils.SampleCreatedAtTime,
	UpdatedBy:        "SYSTEM",
	UpdatedAt:        utils.SampleCreatedAtTime,
	AssetID:          utils.Ptr[int](1),
}

var TestData2 = GethProcessJob{
	ID:               utils.Ptr[int](2),
	UUID:             "880607ab-2833-4ad7-a231-b983a61cad34",
	Name:             "Calculate Cost Basis For ERC20; Asset ID : 539, Name : HAM",
	AlternateName:    "Calculate Cost Basis For ERC20; Asset ID : 539, Name : HAM",
	StartDate:        utils.SampleCreatedAtTime,
	EndDate:          utils.SampleCreatedAtTime,
	Description:      "",
	StatusID:         utils.Ptr[int](utils.SUCCESS_STRUCTURED_VALUE_ID),
	JobCategoryID:    utils.Ptr[int](utils.EOD_JOB_CATEGORY_STRUCTURED_VALUE_ID),
	ImportTypeID:     utils.Ptr[int](utils.MINER_IMPORT_TYPE_STRUCTURED_VALUE_TYPE_ID),
	ChainID:          utils.Ptr[int](1),
	StartBlockNumber: utils.Ptr[uint64](18887528),
	EndBlockNumber:   utils.Ptr[uint64](18887528),
	CreatedBy:        "SYSTEM",
	CreatedAt:        utils.SampleCreatedAtTime,
	UpdatedBy:        "SYSTEM",
	UpdatedAt:        utils.SampleCreatedAtTime,
	AssetID:          utils.Ptr[int](1),
}
var TestAllData = []GethProcessJob{TestData1, TestData2}

func AddGethProcessJobToMockRows(mock pgxmock.PgxPoolIface, dataList []GethProcessJob) *pgxmock.Rows {
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
			data.JobCategoryID,    //9
			data.ImportTypeID,     //10
			data.ChainID,          //11
			data.StartBlockNumber, //12
			data.EndBlockNumber,   //13
			data.CreatedBy,        //14
			data.CreatedAt,        //15
			data.UpdatedBy,        //16
			data.UpdatedAt,        //17
			data.AssetID,          //18
		)
	}
	return rows
}

func TestGetGethProcessJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethProcessJob{targetData}
	gethProcessJobID := targetData.ID
	mockRows := AddGethProcessJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_jobs").WithArgs(*gethProcessJobID).WillReturnRows(mockRows)
	foundGethProcessJob, err := GetGethProcessJob(mock, gethProcessJobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessJob", err)
	}
	if cmp.Equal(*foundGethProcessJob, targetData) == false {
		t.Errorf("Expected GethProcessJob From Method GetGethProcessJob: %v is different from actual %v", foundGethProcessJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessJobID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_jobs").WithArgs(gethProcessJobID).WillReturnRows(noRows)
	foundGethProcessJob, err := GetGethProcessJob(mock, &gethProcessJobID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessJob", err)
	}
	if foundGethProcessJob != nil {
		t.Errorf("Expected GethProcessJob From Method GetGethProcessJob: to be empty but got this: %v", foundGethProcessJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessJobID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_jobs").WithArgs(gethProcessJobID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethProcessJob, err := GetGethProcessJob(mock, &gethProcessJobID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessJob", err)
	}
	if foundGethProcessJob != nil {
		t.Errorf("Expected GethProcessJob From Method GetGethProcessJob: to be empty but got this: %v", foundGethProcessJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLatestGethProcessJobByImportTypeIDAndAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethProcessJob{targetData}
	importTypeID := targetData.ImportTypeID
	assetID := targetData.AssetID
	mockRows := AddGethProcessJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_jobs").WithArgs(*importTypeID, *assetID, utils.SUCCESS_STRUCTURED_VALUE_ID).WillReturnRows(mockRows)
	foundGethProcessJob, err := GetLatestGethProcessJobByImportTypeIDAndAssetID(mock, importTypeID, assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetLatestGethProcessJobByImportTypeIDAndAssetID", err)
	}
	if cmp.Equal(*foundGethProcessJob, targetData) == false {
		t.Errorf("Expected GethProcessJob From Method GetLatestGethProcessJobByImportTypeIDAndAssetID: %v is different from actual %v", foundGethProcessJob, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLatestGethProcessJobByImportTypeIDAndAssetIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	importTypeID := 1
	assetID := 10
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_jobs").WithArgs(importTypeID, assetID, utils.SUCCESS_STRUCTURED_VALUE_ID).WillReturnRows(noRows)
	foundGethProcessJob, err := GetLatestGethProcessJobByImportTypeIDAndAssetID(mock, &importTypeID, &assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetLatestGethProcessJobByImportTypeIDAndAssetID", err)
	}
	if foundGethProcessJob != nil {
		t.Errorf("Expected GethProcessJob From Method GetLatestGethProcessJobByImportTypeIDAndAssetID: to be empty but got this: %v", foundGethProcessJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLatestGethProcessJobByImportTypeIDAndAssetIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	importTypeID := -1
	assetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_jobs").WithArgs(importTypeID, assetID, utils.SUCCESS_STRUCTURED_VALUE_ID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethProcessJob, err := GetLatestGethProcessJobByImportTypeIDAndAssetID(mock, &importTypeID, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetLatestGethProcessJobByImportTypeIDAndAssetID", err)
	}
	if foundGethProcessJob != nil {
		t.Errorf("Expected GethProcessJob From Method GetLatestGethProcessJobByImportTypeIDAndAssetID: to be empty but got this: %v", foundGethProcessJob)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethProcessJob{TestData1, TestData2}
	mockRows := AddGethProcessJobToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_jobs").WillReturnRows(mockRows)
	foundGethProcessJobs, err := GetGethProcessJobList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessJobList", err)
	}
	testGethProcessJobs := TestAllData
	for i, foundGethProcessJob := range foundGethProcessJobs {
		if cmp.Equal(foundGethProcessJob, testGethProcessJobs[i]) == false {
			t.Errorf("Expected GethProcessJob From Method GetGethProcessJobList: %v is different from actual %v", foundGethProcessJob, testGethProcessJobs[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethProcessJobs, err := GetGethProcessJobList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessJobList", err)
	}
	if len(foundGethProcessJobs) != 0 {
		t.Errorf("Expected From Method GetGethProcessJobs: to be empty but got this: %v", foundGethProcessJobs)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethProcessJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	gethProcessJobID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_process_jobs").WithArgs(*gethProcessJobID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethProcessJob(mock, gethProcessJobID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethProcessJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethProcessJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessJobID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_process_jobs").WithArgs(gethProcessJobID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethProcessJob(mock, &gethProcessJobID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethProcessJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_process_jobs").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.StartDate,        //3
		targetData.EndDate,          //4
		targetData.Description,      //5
		targetData.StatusID,         //6
		targetData.JobCategoryID,    //7
		targetData.ImportTypeID,     //8
		targetData.ChainID,          //9
		targetData.StartBlockNumber, //10
		targetData.EndBlockNumber,   //11
		targetData.UpdatedBy,        //12
		targetData.AssetID,          //13
		targetData.ID,               //14

	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethProcessJob(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethProcessJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethProcessJobOnFailureAtBegin(t *testing.T) {
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
	err = UpdateGethProcessJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethProcessJobOnFailure(t *testing.T) {
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
	mock.ExpectExec("^UPDATE geth_process_jobs").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.StartDate,        //3
		targetData.EndDate,          //4
		targetData.Description,      //5
		targetData.StatusID,         //6
		targetData.JobCategoryID,    //7
		targetData.ImportTypeID,     //8
		targetData.ChainID,          //9
		targetData.StartBlockNumber, //10
		targetData.EndBlockNumber,   //11
		targetData.UpdatedBy,        //12
		targetData.AssetID,          //13
		targetData.ID,               //14
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethProcessJob(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJob(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_process_jobs").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.StartDate,        //3
		targetData.EndDate,          //4
		targetData.Description,      //5
		targetData.StatusID,         //6
		targetData.JobCategoryID,    //7
		targetData.ImportTypeID,     //8
		targetData.ChainID,          //9
		targetData.StartBlockNumber, //10
		targetData.EndBlockNumber,   //11
		targetData.CreatedBy,        //12
		targetData.AssetID,          //13
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	gethProcessJobID, err := InsertGethProcessJob(mock, &targetData)
	if gethProcessJobID < 0 {
		t.Fatalf("gethProcessJobID should not be negative ID: %d", gethProcessJobID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethProcessJob", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJobOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_process_jobs").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.StartDate,        //3
		targetData.EndDate,          //4
		targetData.Description,      //5
		targetData.StatusID,         //6
		targetData.JobCategoryID,    //7
		targetData.ImportTypeID,     //8
		targetData.ChainID,          //9
		targetData.StartBlockNumber, //10
		targetData.EndBlockNumber,   //11
		targetData.CreatedBy,        //12
		targetData.AssetID,          //13
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethProcessJobID, err := InsertGethProcessJob(mock, &targetData)
	if gethProcessJobID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethProcessJobID: %d", gethProcessJobID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJobOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_process_jobs").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.StartDate,        //3
		targetData.EndDate,          //4
		targetData.Description,      //5
		targetData.StatusID,         //6
		targetData.JobCategoryID,    //7
		targetData.ImportTypeID,     //8
		targetData.ChainID,          //9
		targetData.StartBlockNumber, //10
		targetData.EndBlockNumber,   //11
		targetData.CreatedBy,        //12
		targetData.AssetID,          //136
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethProcessJobID, err := InsertGethProcessJob(mock, &targetData)
	if gethProcessJobID >= 0 {
		t.Fatalf("Expecting -1 for gethProcessJobID because of error gethProcessJobID: %d", gethProcessJobID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJobList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_process_jobs"}, DBColumnsInsertGethProcessJobList)
	targetData := TestAllData
	err = InsertGethProcessJobList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethProcessJobList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJobListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_process_jobs"}, DBColumnsInsertGethProcessJobList).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertGethProcessJobList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethProcessJobToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_jobs").WillReturnRows(mockRows)
	foundChains, err := GetGethProcessJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessJobListByPagination", err)
	}
	testChains := dataList
	for i, foundChain := range foundChains {
		if cmp.Equal(foundChain, testChains[i]) == false {
			t.Errorf("Expected Chain From Method GetGethProcessJobListByPagination: %v is different from actual %v", foundChain, testChains[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChains, err := GetGethProcessJobListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessJobListByPagination", err)
	}
	if len(foundChains) != 0 {
		t.Errorf("Expected From Method GetGethProcessJobListByPagination: to be empty but got this: %v", foundChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethProcessJobCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_process_jobs").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalGethProcessJobCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalGethProcessJobCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalGethProcessJobCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethProcessJobCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_process_jobs").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalGethProcessJobCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalGethProcessJobCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalGethProcessJobCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
