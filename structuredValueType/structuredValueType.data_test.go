package structuredvaluetype

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
	"created_by",     //5
	"created_at",     //6
	"updated_by",     //7
	"updated_at",     //8
}
var DBColumnsInsertStructuredValueTypes = []string{
	"uuid",           //1
	"name",           //2
	"alternate_name", //3
	"created_by",     //4
	"created_at",     //5
	"updated_by",     //6
	"updated_at",     //7
}

var TestData1 = StructuredValueType{
	ID:            utils.Ptr[int](1),                      //1
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72", //2
	Name:          "Asset Type",                           //3
	AlternateName: "Asset Type",                           //4
	CreatedBy:     "SYSTEM",                               //5
	CreatedAt:     utils.SampleCreatedAtTime,              //6
	UpdatedBy:     "SYSTEM",                               //7
	UpdatedAt:     utils.SampleCreatedAtTime,              //8

}

var TestData2 = StructuredValueType{
	ID:            utils.Ptr[int](3),                      //1
	UUID:          "4f0d5402-7a7c-402d-a7fc-c56a02b13e03", //2
	Name:          "Interval",                             //3
	AlternateName: "Interval",                             //4
	CreatedBy:     "SYSTEM",                               //5
	CreatedAt:     utils.SampleCreatedAtTime,              //6
	UpdatedBy:     "SYSTEM",                               //7
	UpdatedAt:     utils.SampleCreatedAtTime,              //8
}
var TestAllData = []StructuredValueType{TestData1, TestData2}

func AddStructuredValueTypeToMockRows(mock pgxmock.PgxPoolIface, dataList []StructuredValueType) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,            //1
			data.UUID,          //2
			data.Name,          //3
			data.AlternateName, //4
			data.CreatedBy,     //5
			data.CreatedAt,     //6
			data.UpdatedBy,     //7
			data.UpdatedAt,     //8
		)
	}
	return rows
}

func TestGetStructuredValueType(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []StructuredValueType{targetData}
	structuredValueTypeID := targetData.ID
	mockRows := AddStructuredValueTypeToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM structured_value_types").WithArgs(*structuredValueTypeID).WillReturnRows(mockRows)
	foundStructuredValueType, err := GetStructuredValueType(mock, structuredValueTypeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStructuredValueType", err)
	}
	if cmp.Equal(*foundStructuredValueType, targetData) == false {
		t.Errorf("Expected StructuredValueType From Method GetStructuredValueType: %v is different from actual %v", foundStructuredValueType, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueTypeForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueTypeID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM structured_value_types").WithArgs(structuredValueTypeID).WillReturnRows(noRows)
	foundStructuredValueType, err := GetStructuredValueType(mock, &structuredValueTypeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStructuredValueType", err)
	}
	if foundStructuredValueType != nil {
		t.Errorf("Expected StructuredValueType From Method GetStructuredValueType: to be empty but got this: %v", foundStructuredValueType)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueTypeForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueTypeID := -1
	mock.ExpectQuery("^SELECT (.+) FROM structured_value_types").WithArgs(structuredValueTypeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStructuredValueType, err := GetStructuredValueType(mock, &structuredValueTypeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueType", err)
	}
	if foundStructuredValueType != nil {
		t.Errorf("Expected StructuredValueType From Method GetStructuredValueType: to be empty but got this: %v", foundStructuredValueType)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueTypeForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueTypeID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM structured_value_types").WithArgs(structuredValueTypeID).WillReturnRows(differentModelRows)
	foundStructuredValueType, err := GetStructuredValueType(mock, &structuredValueTypeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueType", err)
	}
	if foundStructuredValueType != nil {
		t.Errorf("Expected StructuredValueType From Method GetStructuredValueType: to be empty but got this: %v", foundStructuredValueType)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStructuredValueType(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	structuredValueTypeID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM structured_value_types").WithArgs(*structuredValueTypeID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveStructuredValueType(mock, structuredValueTypeID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveStructuredValueType", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStructuredValueTypeOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueTypeID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveStructuredValueType(mock, &structuredValueTypeID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStructuredValueTypeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueTypeID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM structured_value_types").WithArgs(structuredValueTypeID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveStructuredValueType(mock, &structuredValueTypeID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueTypeList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStructuredValueTypeToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM structured_value_types").WillReturnRows(mockRows)
	ids := []int{1}
	foundStructuredValueTypeList, err := GetStructuredValueTypeList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetStructuredValueTypeList", err)
	}
	for i, sourceStructuredValueType := range dataList {
		if cmp.Equal(sourceStructuredValueType, foundStructuredValueTypeList[i]) == false {
			t.Errorf("Expected StructuredValueType From Method GetStructuredValueTypeList: %v is different from actual %v", sourceStructuredValueType, foundStructuredValueTypeList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueTypeListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM structured_value_types").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1}
	foundStructuredValueTypeList, err := GetStructuredValueTypeList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueTypeList", err)
	}
	if len(foundStructuredValueTypeList) != 0 {
		t.Errorf("Expected From Method GetStructuredValueTypeList: to be empty but got this: %v", foundStructuredValueTypeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueTypeListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM structured_value_types").WillReturnRows(differentModelRows)
	ids := []int{1}
	foundStructuredValueType, err := GetStructuredValueTypeList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueTypeList", err)
	}
	if foundStructuredValueType != nil {
		t.Errorf("Expected StructuredValueType From Method GetStructuredValueTypeList: to be empty but got this: %v", foundStructuredValueType)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateStructuredValueType(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE structured_value_types").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.UpdatedBy,     //3
		targetData.ID,            //4
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateStructuredValueType(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateStructuredValueType", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStructuredValueTypeOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateStructuredValueType(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStructuredValueTypeOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateStructuredValueType(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStructuredValueTypeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE structured_value_types").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.UpdatedBy,     //3
		targetData.ID,            //4
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateStructuredValueType(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValueType(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO structured_value_types").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.CreatedBy,     //3
		targetData.CreatedBy,     //4
	).WillReturnRows(pgxmock.NewRows([]string{"structured_value_type_id"}).AddRow(1))
	mock.ExpectCommit()
	structuredValueTypeID, err := InsertStructuredValueType(mock, &targetData)
	if structuredValueTypeID < 0 {
		t.Fatalf("structuredValueTypeID should not be negative ID: %d", structuredValueTypeID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertStructuredValueType", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValueTypeOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertStructuredValueType(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValueTypeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO structured_value_types").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.CreatedBy,     //3
		targetData.CreatedBy,     //4
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	structuredValueTypeID, err := InsertStructuredValueType(mock, &targetData)
	if structuredValueTypeID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error structuredValueTypeID: %d", structuredValueTypeID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValueTypeOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO structured_value_types").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.CreatedBy,     //3
		targetData.CreatedBy,     //4
	).WillReturnRows(pgxmock.NewRows([]string{"structured_value_type_id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	structuredValueTypeID, err := InsertStructuredValueType(mock, &targetData)
	if structuredValueTypeID >= 0 {
		t.Fatalf("Expecting -1 for structuredValueTypeID because of error structuredValueTypeID: %d", structuredValueTypeID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValueTypes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"structured_value_types"}, DBColumnsInsertStructuredValueTypes)
	targetData := TestAllData
	err = InsertStructuredValueTypes(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertStructuredValueTypes", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValueTypesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"structured_value_types"}, DBColumnsInsertStructuredValueTypes).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertStructuredValueTypes(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueTypeListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddStructuredValueTypeToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"alternate_name = 'test", "structured_value_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM structured_value_types").WillReturnRows(mockRows)
	foundStructuredValueTypeList, err := GetStructuredValueTypeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetStructuredValueTypeListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundStructuredValueTypeList[i]) == false {
			t.Errorf("Expected sourceData From Method GetStructuredValueTypeListByPagination: %v is different from actual %v", sourceData, foundStructuredValueTypeList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueTypeListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"structured_value_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM structured_value_types").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStructuredValueTypeList, err := GetStructuredValueTypeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueTypeListByPagination", err)
	}
	if len(foundStructuredValueTypeList) != 0 {
		t.Errorf("Expected From Method GetStructuredValueTypeListByPagination: to be empty but got this: %v", foundStructuredValueTypeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueTypeListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"structured_value_type_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM structured_value_types").WillReturnRows(differentModelRows)
	foundStructuredValueTypeList, err := GetStructuredValueTypeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueTypeListByPagination", err)
	}
	if len(foundStructuredValueTypeList) != 0 {
		t.Errorf("Expected From Method GetStructuredValueTypeListByPagination: to be empty but got this: %v", foundStructuredValueTypeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalStructuredValueTypeCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM structured_value_types").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalStructuredValueTypeCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalStructuredValueTypeCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalStructuredValueTypeCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM structured_value_types").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalStructuredValueTypeCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalStructuredValueTypeCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalStructuredValueTypeCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
