package source

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
	"id",             //1
	"uuid",           //2
	"name",           //3
	"alternate_name", //4
	"url",            //5
	"ticker",         //6
	"description",    //7
	"created_by",     //8
	"created_at",     //9
	"updated_by",     //10
	"updated_at",     //11
}
var DBColumnsInsertSources = []string{
	"uuid",           //1
	"name",           //2
	"alternate_name", //3
	"url",            //4
	"ticker",         //5
	"description",    //6
	"created_by",     //7
	"created_at",     //8
	"updated_by",     //9
	"updated_at",     //10
}

var TestData1 = Source{
	ID:            utils.Ptr[int](1),                      //1
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72", //2
	Name:          "Sample Position Job 1",                //3
	AlternateName: "Sample Position Job 1",                //4
	URL:           "https://wwww.coingecko.com",           //5
	Ticker:        "CoingGecko",                           //6
	Description:   "",                                     //7
	CreatedBy:     "SYSTEM",                               //8
	CreatedAt:     utils.SampleCreatedAtTime,              //9
	UpdatedBy:     "SYSTEM",                               //10
	UpdatedAt:     utils.SampleCreatedAtTime,              //11

}

var TestData2 = Source{
	ID:            utils.Ptr[int](2),                      //1
	UUID:          "4f0d5402-7a7c-402d-a7fc-c56a02b13e03", //2
	Name:          "Sample Position Job 2",                //3
	AlternateName: "Sample Position Job 2",                //4
	URL:           "https://wwww.coinbase.com",            //5
	Ticker:        "COIN",                                 //6
	Description:   "",                                     //7
	CreatedBy:     "SYSTEM",                               //8
	CreatedAt:     utils.SampleCreatedAtTime,              //9
	UpdatedBy:     "SYSTEM",                               //10
	UpdatedAt:     utils.SampleCreatedAtTime,              //11
}
var TestAllData = []Source{TestData1, TestData2}

func AddSourceToMockRows(mock pgxmock.PgxPoolIface, dataList []Source) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,            //1
			data.UUID,          //2
			data.Name,          //3
			data.AlternateName, //4
			data.URL,           //5
			data.Ticker,        //6
			data.Description,   //7
			data.CreatedBy,     //8
			data.CreatedAt,     //9
			data.UpdatedBy,     //10
			data.UpdatedAt,     //11
		)
	}
	return rows
}

func TestGetSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Source{targetData}
	sourceID := targetData.ID
	mockRows := AddSourceToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM sources").WithArgs(*sourceID).WillReturnRows(mockRows)
	foundSource, err := GetSource(mock, sourceID)
	if err != nil {
		t.Fatalf("an error '%s' in GetSource", err)
	}
	if cmp.Equal(*foundSource, targetData) == false {
		t.Errorf("Expected Source From Method GetSource: %v is different from actual %v", foundSource, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM sources").WithArgs(sourceID).WillReturnRows(noRows)
	foundSource, err := GetSource(mock, &sourceID)
	if err != nil {
		t.Fatalf("an error '%s' in GetSource", err)
	}
	if foundSource != nil {
		t.Errorf("Expected Source From Method GetSource: to be empty but got this: %v", foundSource)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	mock.ExpectQuery("^SELECT (.+) FROM sources").WithArgs(sourceID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundSource, err := GetSource(mock, &sourceID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSource", err)
	}
	if foundSource != nil {
		t.Errorf("Expected Source From Method GetSource: to be empty but got this: %v", foundSource)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM sources").WithArgs(sourceID).WillReturnRows(differentModelRows)
	foundSource, err := GetSource(mock, &sourceID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSource", err)
	}
	if foundSource != nil {
		t.Errorf("Expected Source From Method GetSource: to be empty but got this: %v", foundSource)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	sourceID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM sources").WithArgs(*sourceID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveSource(mock, sourceID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveSource", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveSourceOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveSource(mock, &sourceID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveSourceOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM sources").WithArgs(sourceID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveSource(mock, &sourceID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	ids := []int{1, 2}
	mockRows := AddSourceToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM sources").WillReturnRows(mockRows)
	foundSourceList, err := GetSourceList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetSourceList", err)
	}
	for i, sourceSource := range dataList {
		if cmp.Equal(sourceSource, foundSourceList[i]) == false {
			t.Errorf("Expected Source From Method GetSourceList: %v is different from actual %v", sourceSource, foundSourceList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	ids := []int{1, 2}
	mock.ExpectQuery("^SELECT (.+) FROM sources").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundSourceList, err := GetSourceList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceList", err)
	}
	if len(foundSourceList) != 0 {
		t.Errorf("Expected From Method GetSourceList: to be empty but got this: %v", foundSourceList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	ids := []int{1, 2}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM sources").WillReturnRows(differentModelRows)
	foundSource, err := GetSourceList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceList", err)
	}
	if foundSource != nil {
		t.Errorf("Expected Source From Method GetSourceList: to be empty but got this: %v", foundSource)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE sources").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.URL,           //3
		targetData.Ticker,        //4
		targetData.Description,   //5
		targetData.UpdatedBy,     //6
		targetData.ID,            //7
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateSource(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateSource", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateSourceOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateSource(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateSourceOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateSource(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateSourceOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE sources").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.URL,           //3
		targetData.Ticker,        //4
		targetData.Description,   //5
		targetData.UpdatedBy,     //6
		targetData.ID,            //7
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateSource(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.UUID = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO sources").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.URL,           //3
		targetData.Ticker,        //4
		targetData.Description,   //5
		targetData.CreatedBy,     //6
		targetData.CreatedBy,     //7
	).WillReturnRows(pgxmock.NewRows([]string{"source_id"}).AddRow(1))
	mock.ExpectCommit()
	sourceID, err := InsertSource(mock, &targetData)
	if sourceID < 0 {
		t.Fatalf("sourceID should not be negative ID: %d", sourceID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertSource", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSourceOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertSource(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSourceOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO sources").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.URL,           //3
		targetData.Ticker,        //4
		targetData.Description,   //5
		targetData.CreatedBy,     //6
		targetData.CreatedBy,     //7
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	sourceID, err := InsertSource(mock, &targetData)
	if sourceID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error sourceID: %d", sourceID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSourceOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO sources").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.URL,           //3
		targetData.Ticker,        //4
		targetData.Description,   //5
		targetData.CreatedBy,     //6
		targetData.CreatedBy,     //7
	).WillReturnRows(pgxmock.NewRows([]string{"source_id"}).AddRow(-1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	sourceID, err := InsertSource(mock, &targetData)
	if sourceID >= 0 {
		t.Fatalf("Expecting -1 for sourceID because of error sourceID: %d", sourceID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSources(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"sources"}, DBColumnsInsertSources)
	targetData := TestAllData
	err = InsertSources(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertSources", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSourcesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"sources"}, DBColumnsInsertSources).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertSources(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddSourceToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"ticker = 'COIN'", "name = 'coingecko'"}
	mock.ExpectQuery("^SELECT (.+) FROM sources").WillReturnRows(mockRows)
	foundSourceList, err := GetSourceListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetSourceListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundSourceList[i]) == false {
			t.Errorf("Expected sourceData From Method GetSourceListByPagination: %v is different from actual %v", sourceData, foundSourceList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"ticker = 'COIN'", "name = 'coingecko'"}
	mock.ExpectQuery("^SELECT (.+) FROM sources").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundSourceList, err := GetSourceListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceListByPagination", err)
	}
	if len(foundSourceList) != 0 {
		t.Errorf("Expected From Method GetSourceListByPagination: to be empty but got this: %v", foundSourceList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSourceListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"ticker = 'COIN'", "name = 'coingecko'"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM sources").WillReturnRows(differentModelRows)
	foundSourceList, err := GetSourceListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSourceListByPagination", err)
	}
	if len(foundSourceList) != 0 {
		t.Errorf("Expected From Method GetSourceListByPagination: to be empty but got this: %v", foundSourceList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalSourcesCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM sources").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalSourcesCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalSourcesCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalSourcesCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM sources").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalSourcesCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalSourcesCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalSourcesCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
