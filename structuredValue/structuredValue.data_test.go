package structuredvalue

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
	"id",                       //1
	"uuid",                     //2
	"name",                     //3
	"alternate_name",           //4
	"structured_value_type_id", //5
	"created_by",               //6
	"created_at",               //7
	"updated_by",               //8
	"updated_at",               //9
}
var DBColumnsInsertStructuredValues = []string{
	"uuid",                     //1
	"name",                     //2
	"alternate_name",           //3
	"structured_value_type_id", //4
	"created_by",               //5
	"created_at",               //6
	"updated_by",               //7
	"updated_at",               //8
}

var TestData1 = StructuredValue{
	ID:                    utils.Ptr[int](1),                      //1
	UUID:                  "01ef85e8-2c26-441e-8c7f-71d79518ad72", //2
	Name:                  "Crypto",                               //3
	AlternateName:         "Crypto",                               //4
	StructuredValueTypeID: utils.Ptr[int](1),                      //5
	CreatedBy:             "SYSTEM",                               //6
	CreatedAt:             utils.SampleCreatedAtTime,              //7
	UpdatedBy:             "SYSTEM",                               //8
	UpdatedAt:             utils.SampleCreatedAtTime,              //9

}

var TestData2 = StructuredValue{
	ID:                    utils.Ptr[int](3),                      //1
	UUID:                  "4f0d5402-7a7c-402d-a7fc-c56a02b13e03", //2
	Name:                  "Defi",                                 //3
	AlternateName:         "Defi",                                 //4
	StructuredValueTypeID: utils.Ptr[int](2),                      //5
	CreatedBy:             "SYSTEM",                               //6
	CreatedAt:             utils.SampleCreatedAtTime,              //7
	UpdatedBy:             "SYSTEM",                               //8
	UpdatedAt:             utils.SampleCreatedAtTime,              //9
}
var TestAllData = []StructuredValue{TestData1, TestData2}

func AddStructuredValueToMockRows(mock pgxmock.PgxPoolIface, dataList []StructuredValue) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                    //1
			data.UUID,                  //2
			data.Name,                  //3
			data.AlternateName,         //4
			data.StructuredValueTypeID, //5
			data.CreatedBy,             //6
			data.CreatedAt,             //7
			data.UpdatedBy,             //8
			data.UpdatedAt,             //9
		)
	}
	return rows
}

func TestGetStructuredValue(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []StructuredValue{targetData}
	structuredValueID := targetData.ID
	mockRows := AddStructuredValueToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WithArgs(*structuredValueID).WillReturnRows(mockRows)
	foundStructuredValue, err := GetStructuredValue(mock, structuredValueID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStructuredValue", err)
	}
	if cmp.Equal(*foundStructuredValue, targetData) == false {
		t.Errorf("Expected StructuredValue From Method GetStructuredValue: %v is different from actual %v", foundStructuredValue, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WithArgs(structuredValueID).WillReturnRows(noRows)
	foundStructuredValue, err := GetStructuredValue(mock, &structuredValueID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStructuredValue", err)
	}
	if foundStructuredValue != nil {
		t.Errorf("Expected StructuredValue From Method GetStructuredValue: to be empty but got this: %v", foundStructuredValue)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueID := -1
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WithArgs(structuredValueID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStructuredValue, err := GetStructuredValue(mock, &structuredValueID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValue", err)
	}
	if foundStructuredValue != nil {
		t.Errorf("Expected StructuredValue From Method GetStructuredValue: to be empty but got this: %v", foundStructuredValue)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WithArgs(structuredValueID).WillReturnRows(differentModelRows)
	foundStructuredValue, err := GetStructuredValue(mock, &structuredValueID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValue", err)
	}
	if foundStructuredValue != nil {
		t.Errorf("Expected StructuredValue From Method GetStructuredValue: to be empty but got this: %v", foundStructuredValue)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStructuredValue(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	structuredValueID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM structured_values").WithArgs(*structuredValueID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveStructuredValue(mock, structuredValueID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveStructuredValue", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStructuredValueOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	structuredValueID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveStructuredValue(mock, &structuredValueID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStructuredValueOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM structured_values").WithArgs(structuredValueID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveStructuredValue(mock, &structuredValueID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueByStructuredValueTypeIDList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStructuredValueToMockRows(mock, dataList)
	structuredValueTypeID := TestData1.StructuredValueTypeID
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WithArgs(*structuredValueTypeID).WillReturnRows(mockRows)
	foundStructuredValueList, err := GetStructuredValueByStructuredValueTypeIDList(mock, structuredValueTypeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStructuredValueByStructuredValueTypeIDList", err)
	}
	for i, sourceStructuredValue := range dataList {
		if cmp.Equal(sourceStructuredValue, foundStructuredValueList[i]) == false {
			t.Errorf("Expected StructuredValue From Method GetStructuredValueByStructuredValueTypeIDList: %v is different from actual %v", sourceStructuredValue, foundStructuredValueList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueByStructuredValueTypeIDListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueTypeID := -1
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WithArgs(structuredValueTypeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStructuredValueList, err := GetStructuredValueByStructuredValueTypeIDList(mock, &structuredValueTypeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueByStructuredValueTypeIDList", err)
	}
	if len(foundStructuredValueList) != 0 {
		t.Errorf("Expected From Method GetStructuredValueByStructuredValueTypeIDList: to be empty but got this: %v", foundStructuredValueList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueByStructuredValueTypeIDListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	structuredValueTypeID := 1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WithArgs(structuredValueTypeID).WillReturnRows(differentModelRows)
	foundStructuredValue, err := GetStructuredValueByStructuredValueTypeIDList(mock, &structuredValueTypeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueByStructuredValueTypeIDList", err)
	}
	if foundStructuredValue != nil {
		t.Errorf("Expected StructuredValue From Method GetStructuredValueByStructuredValueTypeIDList: to be empty but got this: %v", foundStructuredValue)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStructuredValueToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WillReturnRows(mockRows)
	ids := []int{1}
	foundStructuredValueList, err := GetStructuredValueList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetStructuredValueList", err)
	}
	for i, sourceStructuredValue := range dataList {
		if cmp.Equal(sourceStructuredValue, foundStructuredValueList[i]) == false {
			t.Errorf("Expected StructuredValue From Method GetStructuredValueList: %v is different from actual %v", sourceStructuredValue, foundStructuredValueList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1}
	foundStructuredValueList, err := GetStructuredValueList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueList", err)
	}
	if len(foundStructuredValueList) != 0 {
		t.Errorf("Expected From Method GetStructuredValueList: to be empty but got this: %v", foundStructuredValueList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WillReturnRows(differentModelRows)
	ids := []int{1}
	foundStructuredValue, err := GetStructuredValueList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueList", err)
	}
	if foundStructuredValue != nil {
		t.Errorf("Expected StructuredValue From Method GetStructuredValueList: to be empty but got this: %v", foundStructuredValue)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateStructuredValue(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE structured_values").WithArgs(
		targetData.Name,                  //1
		targetData.AlternateName,         //2
		targetData.StructuredValueTypeID, //3
		targetData.UpdatedBy,             //4
		targetData.ID,                    //5
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateStructuredValue(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateStructuredValue", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStructuredValueOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateStructuredValue(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStructuredValueOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateStructuredValue(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStructuredValueOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE structured_values").WithArgs(
		targetData.Name,                  //1
		targetData.AlternateName,         //2
		targetData.StructuredValueTypeID, //3
		targetData.UpdatedBy,             //4
		targetData.ID,                    //5
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateStructuredValue(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValue(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO structured_values").WithArgs(
		targetData.Name,                  //1
		targetData.AlternateName,         //2
		targetData.StructuredValueTypeID, //3
		targetData.CreatedBy,             //4
		targetData.CreatedBy,             //5
	).WillReturnRows(pgxmock.NewRows([]string{"structured_value_id"}).AddRow(1))
	mock.ExpectCommit()
	structuredValueID, err := InsertStructuredValue(mock, &targetData)
	if structuredValueID < 0 {
		t.Fatalf("structuredValueID should not be negative ID: %d", structuredValueID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertStructuredValue", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValueOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertStructuredValue(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValueOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO structured_values").WithArgs(
		targetData.Name,                  //1
		targetData.AlternateName,         //2
		targetData.StructuredValueTypeID, //3
		targetData.CreatedBy,             //4
		targetData.CreatedBy,             //5
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	structuredValueID, err := InsertStructuredValue(mock, &targetData)
	if structuredValueID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error structuredValueID: %d", structuredValueID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValueOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO structured_values").WithArgs(
		targetData.Name,                  //1
		targetData.AlternateName,         //2
		targetData.StructuredValueTypeID, //3
		targetData.CreatedBy,             //4
		targetData.CreatedBy,             //5
	).WillReturnRows(pgxmock.NewRows([]string{"structured_value_id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	structuredValueID, err := InsertStructuredValue(mock, &targetData)
	if structuredValueID >= 0 {
		t.Fatalf("Expecting -1 for structuredValueID because of error structuredValueID: %d", structuredValueID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValues(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"structured_values"}, DBColumnsInsertStructuredValues)
	targetData := TestAllData
	err = InsertStructuredValues(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertStructuredValues", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStructuredValuesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"structured_values"}, DBColumnsInsertStructuredValues).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertStructuredValues(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddStructuredValueToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"alternate_name = 'test", "structured_value_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WillReturnRows(mockRows)
	foundStructuredValueList, err := GetStructuredValueListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetStructuredValueListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundStructuredValueList[i]) == false {
			t.Errorf("Expected sourceData From Method GetStructuredValueListByPagination: %v is different from actual %v", sourceData, foundStructuredValueList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStructuredValueList, err := GetStructuredValueListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueListByPagination", err)
	}
	if len(foundStructuredValueList) != 0 {
		t.Errorf("Expected From Method GetStructuredValueListByPagination: to be empty but got this: %v", foundStructuredValueList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStructuredValueListByPaginationForCollectRowsErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM structured_values").WillReturnRows(differentModelRows)
	foundStructuredValueList, err := GetStructuredValueListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStructuredValueListByPagination", err)
	}
	if len(foundStructuredValueList) != 0 {
		t.Errorf("Expected From Method GetStructuredValueListByPagination: to be empty but got this: %v", foundStructuredValueList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalStructuredValueCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM structured_values").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalStructuredValueCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalStructuredValueCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalStructuredValueCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM structured_values").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalStructuredValueCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalStructuredValueCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalStructuredValueCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
