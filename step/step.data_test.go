package step

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
	"id",             //1
	"pool_id",        //2
	"parent_step_id", //3
	"uuid",           //4
	"name",           //5
	"alternate_name", //6
	"start_date",     //7
	"end_date",       //8
	"description",    //9
	"action_type_id", //10
	"function_name",  //11
	"step_order",     //12
	"created_by",     //13
	"created_at",     //14
	"updated_by",     //15
	"updated_at",     //16
}
var DBColumnsInsertSteps = []string{
	"pool_id",        //1
	"parent_step_id", //2
	"uuid",           //3
	"name",           //4
	"alternate_name", //5
	"start_date",     //6
	"end_date",       //7
	"description",    //8
	"action_type_id", //9
	"function_name",  //10
	"step_order",     //11
	"created_by",     //12
	"created_at",     //13
	"updated_by",     //14
	"updated_at",     //15
}

var TestData1 = Step{
	ID:            utils.Ptr[int](1),                      //1
	PoolID:        utils.Ptr[int](1),                      //2
	ParentStepId:  nil,                                    //3
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72", //4
	Name:          "VELO/OP-LP-Claim",                     //5
	AlternateName: "VELO/OP-LP-Claim",                     //6
	StartDate:     utils.SampleCreatedAtTime,              //7
	EndDate:       utils.SampleCreatedAtTime,              //8
	Description:   "Claim Reward from VELO/OP LP.",        //9
	ActionTypeID:  utils.Ptr[int](1),                      //10
	FunctionName:  "claimVeloOpLpFromVelodrome",           //11
	StepOrder:     utils.Ptr[int](1),                      //12
	CreatedBy:     "SYSTEM",                               //13
	CreatedAt:     utils.SampleCreatedAtTime,              //14
	UpdatedBy:     "SYSTEM",                               //15
	UpdatedAt:     utils.SampleCreatedAtTime,              //16

}

var TestData2 = Step{
	ID:            utils.Ptr[int](3),                                        //1
	PoolID:        utils.Ptr[int](2),                                        //2
	ParentStepId:  utils.Ptr[int](2),                                        //3
	UUID:          "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                   //4
	Name:          "XPTP - 3 - Reinvest",                                    //5
	AlternateName: "XPTP - 3 - Reinvest",                                    //6
	StartDate:     utils.SampleCreatedAtTime,                                //7
	EndDate:       utils.SampleCreatedAtTime,                                //8
	Description:   "Reinvest the proceeds back to the vector finance pools", //9
	ActionTypeID:  utils.Ptr[int](71),                                       //10
	FunctionName:  "swapXPTPRewards",                                        //11
	StepOrder:     utils.Ptr[int](3),                                        //12
	CreatedBy:     "SYSTEM",                                                 //13
	CreatedAt:     utils.SampleCreatedAtTime,                                //14
	UpdatedBy:     "SYSTEM",                                                 //15
	UpdatedAt:     utils.SampleCreatedAtTime,                                //16
}
var TestAllData = []Step{TestData1, TestData2}

func AddStepToMockRows(mock pgxmock.PgxPoolIface, dataList []Step) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,            //1
			data.PoolID,        //2
			data.ParentStepId,  //3
			data.UUID,          //4
			data.Name,          //5
			data.AlternateName, //6
			data.StartDate,     //7
			data.EndDate,       //8
			data.Description,   //9
			data.ActionTypeID,  //10
			data.FunctionName,  //11
			data.StepOrder,     //12
			data.CreatedBy,     //13
			data.CreatedAt,     //14
			data.UpdatedBy,     //15
			data.UpdatedAt,     //16
		)
	}
	return rows
}

func TestGetStep(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Step{targetData}
	stepID := targetData.ID
	mockRows := AddStepToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(*stepID).WillReturnRows(mockRows)
	foundStep, err := GetStep(mock, stepID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStep", err)
	}
	if cmp.Equal(*foundStep, targetData) == false {
		t.Errorf("Expected Step From Method GetStep: %v is different from actual %v", foundStep, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	stepID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(stepID).WillReturnRows(noRows)
	foundStep, err := GetStep(mock, &stepID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStep", err)
	}
	if foundStep != nil {
		t.Errorf("Expected Step From Method GetStep: to be empty but got this: %v", foundStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	stepID := -1
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(stepID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStep, err := GetStep(mock, &stepID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStep", err)
	}
	if foundStep != nil {
		t.Errorf("Expected Step From Method GetStep: to be empty but got this: %v", foundStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	stepID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(stepID).WillReturnRows(differentModelRows)
	foundStep, err := GetStep(mock, &stepID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStep", err)
	}
	if foundStep != nil {
		t.Errorf("Expected Step From Method GetStep: to be empty but got this: %v", foundStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStep(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	stepID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM steps").WithArgs(*stepID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveStep(mock, stepID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveStep", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStepOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	stepID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveStep(mock, &stepID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStepOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	stepID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM steps").WithArgs(stepID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveStep(mock, &stepID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetSteps(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStepToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM steps").WillReturnRows(mockRows)
	ids := []int{1}
	foundStepList, err := GetSteps(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetSteps", err)
	}
	for i, sourceStep := range dataList {
		if cmp.Equal(sourceStep, foundStepList[i]) == false {
			t.Errorf("Expected Step From Method GetSteps: %v is different from actual %v", sourceStep, foundStepList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM steps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1}
	foundStepList, err := GetSteps(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSteps", err)
	}
	if len(foundStepList) != 0 {
		t.Errorf("Expected From Method GetSteps: to be empty but got this: %v", foundStepList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM steps").WillReturnRows(differentModelRows)
	ids := []int{1}
	foundStep, err := GetSteps(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetSteps", err)
	}
	if foundStep != nil {
		t.Errorf("Expected Step From Method GetSteps: to be empty but got this: %v", foundStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStepToMockRows(mock, dataList)
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(pq.Array(UUIDList)).WillReturnRows(mockRows)
	foundStepList, err := GetStepsByUUIDs(mock, UUIDList)
	if err != nil {
		t.Fatalf("an error '%s' in GetStepsByUUIDs", err)
	}
	for i, sourceStep := range dataList {
		if cmp.Equal(sourceStep, foundStepList[i]) == false {
			t.Errorf("Expected Step From Method GetStepsByUUIDs: %v is different from actual %v", sourceStep, foundStepList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(pq.Array(UUIDList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStepList, err := GetStepsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepsByUUIDs", err)
	}
	if len(foundStepList) != 0 {
		t.Errorf("Expected From Method GetStepsByUUIDs: to be empty but got this: %v", foundStepList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepsByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(pq.Array(UUIDList)).WillReturnRows(differentModelRows)
	foundStep, err := GetStepsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepsByUUIDs", err)
	}
	if foundStep != nil {
		t.Errorf("Expected Step From Method GetStepsByUUIDs: to be empty but got this: %v", foundStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepsFromPoolID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStepToMockRows(mock, dataList)
	poolID := TestData1.PoolID
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(*poolID).WillReturnRows(mockRows)
	foundStepList, err := GetStepsFromPoolID(mock, poolID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStepsFromPoolID", err)
	}
	for i, sourceStep := range dataList {
		if cmp.Equal(sourceStep, foundStepList[i]) == false {
			t.Errorf("Expected Step From Method GetStepsFromPoolID: %v is different from actual %v", sourceStep, foundStepList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepsFromPoolIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	poolID := -1
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(poolID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStepList, err := GetStepsFromPoolID(mock, &poolID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepsFromPoolID", err)
	}
	if len(foundStepList) != 0 {
		t.Errorf("Expected From Method GetStepsFromPoolID: to be empty but got this: %v", foundStepList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepsFromPoolIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	poolID := 1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(poolID).WillReturnRows(differentModelRows)
	foundStep, err := GetStepsFromPoolID(mock, &poolID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepsFromPoolID", err)
	}
	if foundStep != nil {
		t.Errorf("Expected Step From Method GetStepsFromPoolID: to be empty but got this: %v", foundStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffSteps(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStepToMockRows(mock, dataList)
	diffInDate := 2
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(diffInDate).WillReturnRows(mockRows)
	foundStepList, err := GetStartAndEndDateDiffSteps(mock, &diffInDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndDateDiffSteps", err)
	}
	for i, sourceStep := range dataList {
		if cmp.Equal(sourceStep, foundStepList[i]) == false {
			t.Errorf("Expected Step From Method GetStartAndEndDateDiffSteps: %v is different from actual %v", sourceStep, foundStepList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStepsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := -1
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(diffInDate).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStepList, err := GetStartAndEndDateDiffSteps(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffSteps", err)
	}
	if len(foundStepList) != 0 {
		t.Errorf("Expected From Method GetStartAndEndDateDiffSteps: to be empty but got this: %v", foundStepList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStepsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := 1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM steps").WithArgs(diffInDate).WillReturnRows(differentModelRows)
	foundStep, err := GetStartAndEndDateDiffSteps(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffSteps", err)
	}
	if foundStep != nil {
		t.Errorf("Expected Step From Method GetStartAndEndDateDiffSteps: to be empty but got this: %v", foundStep)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateStep(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE steps").WithArgs(
		targetData.PoolID,        //1
		targetData.ParentStepId,  //2
		targetData.UUID,          //3
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.StartDate,     //6
		targetData.EndDate,       //7
		targetData.Description,   //8
		targetData.ActionTypeID,  //9
		targetData.FunctionName,  //10
		targetData.StepOrder,     //11
		targetData.UpdatedBy,     //12
		targetData.ID,            //13
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateStep(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateStep", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStepOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateStep(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStepOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateStep(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStepOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE steps").WithArgs(
		targetData.PoolID,        //1
		targetData.ParentStepId,  //2
		targetData.UUID,          //3
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.StartDate,     //6
		targetData.EndDate,       //7
		targetData.Description,   //8
		targetData.ActionTypeID,  //9
		targetData.FunctionName,  //10
		targetData.StepOrder,     //11
		targetData.UpdatedBy,     //12
		targetData.ID,            //13
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateStep(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStep(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO steps").WithArgs(
		targetData.PoolID,        //1
		targetData.ParentStepId,  //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.StartDate,     //5
		targetData.EndDate,       //6
		targetData.Description,   //7
		targetData.ActionTypeID,  //8
		targetData.FunctionName,  //9
		targetData.StepOrder,     //10
		targetData.CreatedBy,     //11
	).WillReturnRows(pgxmock.NewRows([]string{"step_id"}).AddRow(1))
	mock.ExpectCommit()
	stepID, err := InsertStep(mock, &targetData)
	if stepID < 0 {
		t.Fatalf("stepID should not be negative ID: %d", stepID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertStep", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStepOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertStep(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStepOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO steps").WithArgs(
		targetData.PoolID,        //1
		targetData.ParentStepId,  //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.StartDate,     //5
		targetData.EndDate,       //6
		targetData.Description,   //7
		targetData.ActionTypeID,  //8
		targetData.FunctionName,  //9
		targetData.StepOrder,     //10
		targetData.CreatedBy,     //11
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	stepID, err := InsertStep(mock, &targetData)
	if stepID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error stepID: %d", stepID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStepOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO steps").WithArgs(
		targetData.PoolID,        //1
		targetData.ParentStepId,  //2
		targetData.Name,          //3
		targetData.AlternateName, //4
		targetData.StartDate,     //5
		targetData.EndDate,       //6
		targetData.Description,   //7
		targetData.ActionTypeID,  //8
		targetData.FunctionName,  //9
		targetData.StepOrder,     //10
		targetData.CreatedBy,     //11
	).WillReturnRows(pgxmock.NewRows([]string{"step_id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	stepID, err := InsertStep(mock, &targetData)
	if stepID >= 0 {
		t.Fatalf("Expecting -1 for stepID because of error stepID: %d", stepID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertSteps(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"steps"}, DBColumnsInsertSteps)
	targetData := TestAllData
	err = InsertSteps(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertSteps", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStepsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"steps"}, DBColumnsInsertSteps).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertSteps(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddStepToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"pool_id = 1", "parent_step_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM steps").WillReturnRows(mockRows)
	foundStepList, err := GetStepListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetStepListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundStepList[i]) == false {
			t.Errorf("Expected sourceData From Method GetStepListByPagination: %v is different from actual %v", sourceData, foundStepList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM steps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStepList, err := GetStepListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepListByPagination", err)
	}
	if len(foundStepList) != 0 {
		t.Errorf("Expected From Method GetStepListByPagination: to be empty but got this: %v", foundStepList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepListByPaginationForCollectRowsErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM steps").WillReturnRows(differentModelRows)
	foundStepList, err := GetStepListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepListByPagination", err)
	}
	if len(foundStepList) != 0 {
		t.Errorf("Expected From Method GetStepListByPagination: to be empty but got this: %v", foundStepList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalStepsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM steps").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalStepsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalStepsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalStepsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM steps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalStepsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalStepsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalStepsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
