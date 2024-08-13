package gethlylejobstopics

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
	"id",                  //1
	"geth_process_job_id", //2
	"uuid",                //3
	"name",                //4
	"alternate_name",      //5
	"description",         //6
	"status_id",           //7
	"topic_str",           //8
	"created_by",          //9
	"created_at",          //10
	"updated_by",          //11
	"updated_at",          //12
}
var DBColumnsInsertGethProcessJobTopicList = []string{
	"geth_process_job_id", //1
	"uuid",                //2
	"name",                //3
	"alternate_name",      //4
	"description",         //5
	"status_id",           //6
	"topic_str",           //7
	"created_by",          //8
	"created_at",          //9
	"updated_by",          //10
	"updated_at",          //11
}

var TestData1 = GethProcessJobTopic{
	ID:            utils.Ptr[int](1),
	UUID:          "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:          "Import Swaps Using Liquidity Pool ID : 15, Name : PEPE/WETH Uniswap V2, SwapSig : PEPE/WETH Uniswap V2",
	AlternateName: "Asset ID : 535, Name : PEPE, Import All Swaps For ERC20",
	Description:   "Original Start Block : 17046105, Current Block : 18768759",
	StatusID:      utils.Ptr[int](utils.SUCCESS_STRUCTURED_VALUE_ID),
	TopicStr:      "Swap(address,uint256,uint256,uint256,uint256,address)",
	CreatedBy:     "SYSTEM",
	CreatedAt:     utils.SampleCreatedAtTime,
	UpdatedBy:     "SYSTEM",
	UpdatedAt:     utils.SampleCreatedAtTime,
}

var TestData2 = GethProcessJobTopic{
	ID:            utils.Ptr[int](2),
	UUID:          "880607ab-2833-4ad7-a231-b983a61cad34",
	Name:          "Import Swaps Using Liquidity Pool ID : 14, Name : PEPE/WETH Uniswap V3, SwapSig : PEPE/WETH Uniswap V3",
	AlternateName: "Asset ID : 535, Name : PEPE, Import All Swaps For ERC20",
	Description:   "Original Start Block : 17046105, Current Block : 18768759",
	StatusID:      utils.Ptr[int](utils.SUCCESS_STRUCTURED_VALUE_ID),
	TopicStr:      "Swap(address,address,int256,int256,uint160,uint128,int24)",
	CreatedBy:     "SYSTEM",
	CreatedAt:     utils.SampleCreatedAtTime,
	UpdatedBy:     "SYSTEM",
	UpdatedAt:     utils.SampleCreatedAtTime,
}
var TestAllData = []GethProcessJobTopic{TestData1, TestData2}

func AddGethProcessJobTopicToMockRows(mock pgxmock.PgxPoolIface, dataList []GethProcessJobTopic) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,               //1
			data.GethProcessJobID, //2
			data.UUID,             //3
			data.Name,             //4
			data.AlternateName,    //5
			data.Description,      //6
			data.StatusID,         //7
			data.TopicStr,         //8
			data.CreatedBy,        //9
			data.CreatedAt,        //10
			data.UpdatedBy,        //11
			data.UpdatedAt,        //12
		)
	}
	return rows
}

func TestGetGethProcessJobTopic(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethProcessJobTopic{targetData}
	gethProcessJobTopicID := targetData.ID
	mockRows := AddGethProcessJobTopicToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_job_topics").WithArgs(*gethProcessJobTopicID).WillReturnRows(mockRows)
	foundGethProcessJobTopic, err := GetGethProcessJobTopic(mock, gethProcessJobTopicID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessJobTopic", err)
	}
	if cmp.Equal(*foundGethProcessJobTopic, targetData) == false {
		t.Errorf("Expected GethProcessJobTopic From Method GetGethProcessJobTopic: %v is different from actual %v", foundGethProcessJobTopic, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobTopicForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessJobTopicID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_job_topics").WithArgs(gethProcessJobTopicID).WillReturnRows(noRows)
	foundGethProcessJobTopic, err := GetGethProcessJobTopic(mock, &gethProcessJobTopicID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessJobTopic", err)
	}
	if foundGethProcessJobTopic != nil {
		t.Errorf("Expected GethProcessJobTopic From Method GetGethProcessJobTopic: to be empty but got this: %v", foundGethProcessJobTopic)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobTopicForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessJobTopicID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_job_topics").WithArgs(gethProcessJobTopicID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethProcessJobTopic, err := GetGethProcessJobTopic(mock, &gethProcessJobTopicID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessJobTopic", err)
	}
	if foundGethProcessJobTopic != nil {
		t.Errorf("Expected GethProcessJobTopic From Method GetGethProcessJobTopic: to be empty but got this: %v", foundGethProcessJobTopic)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobTopicForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessJobTopicID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_job_topics").WithArgs(gethProcessJobTopicID).WillReturnRows(differentModelRows)
	foundGethProcessJobTopic, err := GetGethProcessJobTopic(mock, &gethProcessJobTopicID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessJobTopic", err)
	}
	if foundGethProcessJobTopic != nil {
		t.Errorf("Expected foundGethProcessJobTopic From Method GetGethProcessJobTopic: to be empty but got this: %v", foundGethProcessJobTopic)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobTopicList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethProcessJobTopic{TestData1, TestData2}
	mockRows := AddGethProcessJobTopicToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_job_topics").WillReturnRows(mockRows)
	foundGethProcessJobTopics, err := GetGethProcessJobTopicList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessJobTopicList", err)
	}
	testGethProcessJobTopics := TestAllData
	for i, foundGethProcessJobTopic := range foundGethProcessJobTopics {
		if cmp.Equal(foundGethProcessJobTopic, testGethProcessJobTopics[i]) == false {
			t.Errorf("Expected GethProcessJobTopic From Method GetGethProcessJobTopicList: %v is different from actual %v", foundGethProcessJobTopic, testGethProcessJobTopics[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobTopicListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_job_topics").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethProcessJobTopics, err := GetGethProcessJobTopicList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessJobTopicList", err)
	}
	if len(foundGethProcessJobTopics) != 0 {
		t.Errorf("Expected From Method GetGethProcessJobTopics: to be empty but got this: %v", foundGethProcessJobTopics)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethProcessJobTopic(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	gethProcessJobTopicID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_process_job_topics").WithArgs(*gethProcessJobTopicID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethProcessJobTopic(mock, gethProcessJobTopicID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethProcessJobTopic", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethProcessJobTopicOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessJobTopicID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveGethProcessJobTopic(mock, &gethProcessJobTopicID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethProcessJobTopicOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethProcessJobTopicID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_process_job_topics").WithArgs(gethProcessJobTopicID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethProcessJobTopic(mock, &gethProcessJobTopicID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethProcessJobTopic(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_process_job_topics").WithArgs(
		targetData.GethProcessJobID, //1
		targetData.Name,             //2
		targetData.AlternateName,    //3
		targetData.Description,      //4
		targetData.StatusID,         //5
		targetData.TopicStr,         //6
		targetData.UpdatedBy,        //7
		targetData.ID,               //8
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethProcessJobTopic(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethProcessJobTopic", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethProcessJobTopicOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateGethProcessJobTopic(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethProcessJobTopicOnFailureAtBegin(t *testing.T) {
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
	err = UpdateGethProcessJobTopic(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethProcessJobTopicOnFailure(t *testing.T) {
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
	mock.ExpectExec("^UPDATE geth_process_job_topics").WithArgs(
		targetData.GethProcessJobID, //1
		targetData.Name,             //2
		targetData.AlternateName,    //3
		targetData.Description,      //4
		targetData.StatusID,         //5
		targetData.TopicStr,         //6
		targetData.UpdatedBy,        //7
		targetData.ID,               //8
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethProcessJobTopic(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJobTopic(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_process_job_topics").WithArgs(
		targetData.GethProcessJobID, //1
		targetData.Name,             //2
		targetData.AlternateName,    //3
		targetData.Description,      //4
		targetData.StatusID,         //5
		targetData.TopicStr,         //6
		targetData.CreatedBy,        //7
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	gethProcessJobTopicID, err := InsertGethProcessJobTopic(mock, &targetData)
	if gethProcessJobTopicID < 0 {
		t.Fatalf("GethProcessJobTopicID should not be negative ID: %d", gethProcessJobTopicID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethProcessJobTopic", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJobTopicOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertGethProcessJobTopic(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJobTopicOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_process_job_topics").WithArgs(
		targetData.GethProcessJobID, //1
		targetData.Name,             //2
		targetData.AlternateName,    //3
		targetData.Description,      //4
		targetData.StatusID,         //5
		targetData.TopicStr,         //6
		targetData.CreatedBy,        //7
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethProcessJobTopicID, err := InsertGethProcessJobTopic(mock, &targetData)
	if gethProcessJobTopicID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error GethProcessJobTopicID: %d", gethProcessJobTopicID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJobTopicOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_process_job_topics").WithArgs(
		targetData.GethProcessJobID, //1
		targetData.Name,             //2
		targetData.AlternateName,    //3
		targetData.Description,      //4
		targetData.StatusID,         //5
		targetData.TopicStr,         //6
		targetData.CreatedBy,        //7
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethProcessJobTopicID, err := InsertGethProcessJobTopic(mock, &targetData)
	if gethProcessJobTopicID >= 0 {
		t.Fatalf("Expecting -1 for GethProcessJobTopicID because of error GethProcessJobTopicID: %d", gethProcessJobTopicID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJobTopicList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_process_job_topics"}, DBColumnsInsertGethProcessJobTopicList)
	targetData := TestAllData
	err = InsertGethProcessJobTopicList(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethProcessJobTopicList", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethProcessJobTopicListOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_process_job_topics"}, DBColumnsInsertGethProcessJobTopicList).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertGethProcessJobTopicList(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobTopicListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethProcessJobTopicToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"geth_process_job_id = 1", "name='test'"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_job_topics").WillReturnRows(mockRows)
	foundGethProcessJobTopicList, err := GetGethProcessJobTopicListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethProcessJobTopicListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundGethProcessJobTopicList[i]) == false {
			t.Errorf("Expected sourceData From Method GetGethProcessJobTopicListByPagination: %v is different from actual %v", sourceData, foundGethProcessJobTopicList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobTopicListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"geth_process_job_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_job_topics").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethProcessJobTopicList, err := GetGethProcessJobTopicListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessJobTopicListByPagination", err)
	}
	if len(foundGethProcessJobTopicList) != 0 {
		t.Errorf("Expected From Method GetGethProcessJobTopicListByPagination: to be empty but got this: %v", foundGethProcessJobTopicList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethProcessJobTopicListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"geth_process_job_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_process_job_topics").WillReturnRows(differentModelRows)
	foundGethProcessJobTopicList, err := GetGethProcessJobTopicListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethProcessJobTopicListByPagination", err)
	}
	if foundGethProcessJobTopicList != nil {
		t.Errorf("Expected From Method GetGethProcessJobTopicListByPagination: to be empty but got this: %v", foundGethProcessJobTopicList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethProcessJobTopicCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_process_job_topics").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalGethProcessJobTopicCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalGethProcessJobTopicCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalGethProcessJobTopicCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethProcessJobTopicCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_process_job_topics").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalGethProcessJobTopicCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalGethProcessJobTopicCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalGethProcessJobTopicCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
