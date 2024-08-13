package strategy

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
	"id",               //1
	"uuid",             //2
	"name",             //3
	"alternate_name",   //4
	"start_date",       //5
	"end_date",         //6
	"description",      //7
	"strategy_type_id", //8
	"created_by",       //9
	"created_at",       //10
	"updated_by",       //11
	"updated_at",       //12
}
var DBColumnsInsertStrategies = []string{
	"uuid",             //1
	"name",             //2
	"alternate_name",   //3
	"start_date",       //4
	"end_date",         //5
	"description",      //6
	"strategy_type_id", //7
	"created_by",       //8
	"created_at",       //9
	"updated_by",       //10
	"updated_at",       //11
}

var TestData1 = Strategy{
	ID:             utils.Ptr[int](1),                      //1
	UUID:           "01ef85e8-2c26-441e-8c7f-71d79518ad72", //2
	Name:           "Velodrome",                            //3
	AlternateName:  "Velodrome",                            //4
	StartDate:      utils.SampleCreatedAtTime,              //5
	EndDate:        utils.SampleCreatedAtTime,              //6
	Description:    "",                                     //7
	StrategyTypeID: utils.Ptr[int](1),                      //8
	CreatedBy:      "SYSTEM",                               //9
	CreatedAt:      utils.SampleCreatedAtTime,              //10
	UpdatedBy:      "SYSTEM",                               //11
	UpdatedAt:      utils.SampleCreatedAtTime,              //12

}

var TestData2 = Strategy{
	ID:             utils.Ptr[int](3),                                        //1
	UUID:           "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                   //2
	Name:           "Vector Finance",                                         //3
	AlternateName:  "Vector Finance",                                         //4
	StartDate:      utils.SampleCreatedAtTime,                                //5
	EndDate:        utils.SampleCreatedAtTime,                                //6
	Description:    "Reinvest the proceeds back to the vector finance pools", //7
	StrategyTypeID: utils.Ptr[int](2),                                        //8
	CreatedBy:      "SYSTEM",                                                 //9
	CreatedAt:      utils.SampleCreatedAtTime,                                //10
	UpdatedBy:      "SYSTEM",                                                 //11
	UpdatedAt:      utils.SampleCreatedAtTime,                                //12
}
var TestAllData = []Strategy{TestData1, TestData2}

func AddStrategyToMockRows(mock pgxmock.PgxPoolIface, dataList []Strategy) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,             //1
			data.UUID,           //2
			data.Name,           //3
			data.AlternateName,  //4
			data.StartDate,      //5
			data.EndDate,        //6
			data.Description,    //7
			data.StrategyTypeID, //8
			data.CreatedBy,      //9
			data.CreatedAt,      //10
			data.UpdatedBy,      //11
			data.UpdatedAt,      //12
		)
	}
	return rows
}

func TestGetStrategy(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Strategy{targetData}
	strategyID := targetData.ID
	mockRows := AddStrategyToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WithArgs(*strategyID).WillReturnRows(mockRows)
	foundStrategy, err := GetStrategy(mock, strategyID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategy", err)
	}
	if cmp.Equal(*foundStrategy, targetData) == false {
		t.Errorf("Expected Strategy From Method GetStrategy: %v is different from actual %v", foundStrategy, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WithArgs(strategyID).WillReturnRows(noRows)
	foundStrategy, err := GetStrategy(mock, &strategyID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategy", err)
	}
	if foundStrategy != nil {
		t.Errorf("Expected Strategy From Method GetStrategy: to be empty but got this: %v", foundStrategy)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -1
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WithArgs(strategyID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategy, err := GetStrategy(mock, &strategyID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategy", err)
	}
	if foundStrategy != nil {
		t.Errorf("Expected Strategy From Method GetStrategy: to be empty but got this: %v", foundStrategy)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WithArgs(strategyID).WillReturnRows(differentModelRows)
	foundStrategy, err := GetStrategy(mock, &strategyID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategy", err)
	}
	if foundStrategy != nil {
		t.Errorf("Expected Strategy From Method GetStrategy: to be empty but got this: %v", foundStrategy)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStrategy(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	strategyID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM strategies").WithArgs(*strategyID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveStrategy(mock, strategyID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveStrategy", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStrategyOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	strategyID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveStrategy(mock, &strategyID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStrategyOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM strategies").WithArgs(strategyID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveStrategy(mock, &strategyID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategies(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStrategyToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WillReturnRows(mockRows)
	ids := []int{1}
	foundStrategyList, err := GetStrategies(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategies", err)
	}
	for i, sourceStrategy := range dataList {
		if cmp.Equal(sourceStrategy, foundStrategyList[i]) == false {
			t.Errorf("Expected Strategy From Method GetStrategies: %v is different from actual %v", sourceStrategy, foundStrategyList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategiesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1}
	foundStrategyList, err := GetStrategies(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategies", err)
	}
	if len(foundStrategyList) != 0 {
		t.Errorf("Expected From Method GetStrategies: to be empty but got this: %v", foundStrategyList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategiesForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WillReturnRows(differentModelRows)
	ids := []int{1}
	foundStrategy, err := GetStrategies(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategies", err)
	}
	if foundStrategy != nil {
		t.Errorf("Expected Strategy From Method GetStrategies: to be empty but got this: %v", foundStrategy)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategiesByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStrategyToMockRows(mock, dataList)
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WithArgs(pq.Array(UUIDList)).WillReturnRows(mockRows)
	foundStrategyList, err := GetStrategiesByUUIDs(mock, UUIDList)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategiesByUUIDs", err)
	}
	for i, sourceStrategy := range dataList {
		if cmp.Equal(sourceStrategy, foundStrategyList[i]) == false {
			t.Errorf("Expected Strategy From Method GetStrategiesByUUIDs: %v is different from actual %v", sourceStrategy, foundStrategyList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategiesByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WithArgs(pq.Array(UUIDList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyList, err := GetStrategiesByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategiesByUUIDs", err)
	}
	if len(foundStrategyList) != 0 {
		t.Errorf("Expected From Method GetStrategiesByUUIDs: to be empty but got this: %v", foundStrategyList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategiesByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WithArgs(pq.Array(UUIDList)).WillReturnRows(differentModelRows)
	foundStrategy, err := GetStrategiesByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategiesByUUIDs", err)
	}
	if foundStrategy != nil {
		t.Errorf("Expected Strategy From Method GetStrategiesByUUIDs: to be empty but got this: %v", foundStrategy)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStrategies(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStrategyToMockRows(mock, dataList)
	diffInDate := 2
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WithArgs(diffInDate).WillReturnRows(mockRows)
	foundStrategyList, err := GetStartAndEndDateDiffStrategies(mock, &diffInDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndDateDiffStrategies", err)
	}
	for i, sourceStrategy := range dataList {
		if cmp.Equal(sourceStrategy, foundStrategyList[i]) == false {
			t.Errorf("Expected Strategy From Method GetStartAndEndDateDiffStrategies: %v is different from actual %v", sourceStrategy, foundStrategyList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStrategiesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := -1
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WithArgs(diffInDate).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyList, err := GetStartAndEndDateDiffStrategies(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffStrategies", err)
	}
	if len(foundStrategyList) != 0 {
		t.Errorf("Expected From Method GetStartAndEndDateDiffStrategies: to be empty but got this: %v", foundStrategyList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStrategiesForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := 1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WithArgs(diffInDate).WillReturnRows(differentModelRows)
	foundStrategy, err := GetStartAndEndDateDiffStrategies(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffStrategies", err)
	}
	if foundStrategy != nil {
		t.Errorf("Expected Strategy From Method GetStartAndEndDateDiffStrategies: to be empty but got this: %v", foundStrategy)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateStrategy(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE strategies").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.Description,    //5
		targetData.StrategyTypeID, //6
		targetData.UpdatedBy,      //7
		targetData.ID,             //8
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateStrategy(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateStrategy", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStrategyOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateStrategy(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStrategyOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateStrategy(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStrategyOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE strategies").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.Description,    //5
		targetData.StrategyTypeID, //6
		targetData.UpdatedBy,      //7
		targetData.ID,             //8
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateStrategy(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategy(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO strategies").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.Description,    //5
		targetData.StrategyTypeID, //6
		targetData.CreatedBy,      //7
	).WillReturnRows(pgxmock.NewRows([]string{"strategy_id"}).AddRow(1))
	mock.ExpectCommit()
	strategyID, err := InsertStrategy(mock, &targetData)
	if strategyID < 0 {
		t.Fatalf("strategyID should not be negative ID: %d", strategyID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertStrategy", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertStrategy(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO strategies").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.Description,    //5
		targetData.StrategyTypeID, //6
		targetData.CreatedBy,      //7
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	strategyID, err := InsertStrategy(mock, &targetData)
	if strategyID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error strategyID: %d", strategyID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategyOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO strategies").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.Description,    //5
		targetData.StrategyTypeID, //6
		targetData.CreatedBy,      //7
	).WillReturnRows(pgxmock.NewRows([]string{"strategy_id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	strategyID, err := InsertStrategy(mock, &targetData)
	if strategyID >= 0 {
		t.Fatalf("Expecting -1 for strategyID because of error strategyID: %d", strategyID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategies(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"strategies"}, DBColumnsInsertStrategies)
	targetData := TestAllData
	err = InsertStrategies(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertStrategies", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStrategiesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"strategies"}, DBColumnsInsertStrategies).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertStrategies(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddStrategyToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"pool_id = 1", "parent_strategy_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WillReturnRows(mockRows)
	foundStrategyList, err := GetStrategyListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetStrategyListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundStrategyList[i]) == false {
			t.Errorf("Expected sourceData From Method GetStrategyListByPagination: %v is different from actual %v", sourceData, foundStrategyList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"pool_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStrategyList, err := GetStrategyListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyListByPagination", err)
	}
	if len(foundStrategyList) != 0 {
		t.Errorf("Expected From Method GetStrategyListByPagination: to be empty but got this: %v", foundStrategyList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStrategyListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"pool_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM strategies").WillReturnRows(differentModelRows)
	foundStrategyList, err := GetStrategyListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStrategyListByPagination", err)
	}
	if len(foundStrategyList) != 0 {
		t.Errorf("Expected From Method GetStrategyListByPagination: to be empty but got this: %v", foundStrategyList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalStrategiesCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM strategies").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalStrategiesCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalStrategiesCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalStrategiesCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM strategies").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalStrategiesCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalStrategiesCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalStrategiesCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
