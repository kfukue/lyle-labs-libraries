package stepasset

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
	"step_id",          //2
	"asset_id",         //3
	"swap_asset_id",    //4
	"target_pool_id",   //5
	"uuid",             //6
	"name",             //7
	"alternate_name",   //8
	"start_date",       //9
	"end_date",         //10
	"description",      //11
	"action_parameter", //12
	"created_by",       //13
	"created_at",       //14
	"updated_by",       //15
	"updated_at",       //16
}
var DBColumnsInsertStepAssets = []string{
	"step_id",          //1
	"asset_id",         //2
	"swap_asset_id",    //3
	"target_pool_id",   //4
	"uuid",             //5
	"name",             //6
	"alternate_name",   //7
	"start_date",       //8
	"end_date",         //9
	"description",      //10
	"action_parameter", //11
	"created_by",       //12
	"created_at",       //13
	"updated_by",       //14
	"updated_at",       //15
}

var TestData1 = StepAsset{
	ID:              utils.Ptr[int](1),                                   //1
	StepID:          utils.Ptr[int](1),                                   //2
	AssetID:         utils.Ptr[int](1),                                   //3
	SwapAssetID:     utils.Ptr[int](1),                                   //4
	TargetPoolID:    utils.Ptr[int](1),                                   //5
	UUID:            "01ef85e8-2c26-441e-8c7f-71d79518ad72",              //6
	Name:            "Swap QI to USDC",                                   //7
	AlternateName:   "Swap QI to USDC",                                   //8
	StartDate:       utils.SampleCreatedAtTime,                           //9
	EndDate:         utils.SampleCreatedAtTime,                           //10
	Description:     "Swap reward to USDT using action parameter (in %)", //11
	ActionParameter: utils.Ptr[float64](50),                              //12
	CreatedBy:       "SYSTEM",                                            //13
	CreatedAt:       utils.SampleCreatedAtTime,                           //14
	UpdatedBy:       "SYSTEM",                                            //15
	UpdatedAt:       utils.SampleCreatedAtTime,                           //16

}

var TestData2 = StepAsset{
	ID:              utils.Ptr[int](3),                                        //1
	StepID:          utils.Ptr[int](1),                                        //2
	AssetID:         utils.Ptr[int](1),                                        //3
	SwapAssetID:     utils.Ptr[int](1),                                        //4
	TargetPoolID:    utils.Ptr[int](1),                                        //5
	UUID:            "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                   //6
	Name:            "Reinvest USDT To USDT Pool Vector",                      //7
	AlternateName:   "Reinvest USDT To USDT Pool Vector",                      //8
	StartDate:       utils.SampleCreatedAtTime,                                //9
	EndDate:         utils.SampleCreatedAtTime,                                //10
	Description:     "Reinvest the USDT into Vector Finance stable coin pool", //11
	ActionParameter: utils.Ptr[float64](100),                                  //12
	CreatedBy:       "SYSTEM",                                                 //13
	CreatedAt:       utils.SampleCreatedAtTime,                                //14
	UpdatedBy:       "SYSTEM",                                                 //15
	UpdatedAt:       utils.SampleCreatedAtTime,                                //16
}
var TestAllData = []StepAsset{TestData1, TestData2}

func AddStepAssetToMockRows(mock pgxmock.PgxPoolIface, dataList []StepAsset) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,              //1
			data.StepID,          //2
			data.AssetID,         //3
			data.SwapAssetID,     //4
			data.TargetPoolID,    //5
			data.UUID,            //6
			data.Name,            //7
			data.AlternateName,   //8
			data.StartDate,       //9
			data.EndDate,         //10
			data.Description,     //11
			data.ActionParameter, //12
			data.CreatedBy,       //13
			data.CreatedAt,       //14
			data.UpdatedBy,       //15
			data.UpdatedAt,       //16
		)
	}
	return rows
}

func TestGetStepAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []StepAsset{targetData}
	stepAssetID := targetData.ID
	mockRows := AddStepAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WithArgs(*stepAssetID).WillReturnRows(mockRows)
	foundStepAsset, err := GetStepAsset(mock, stepAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStepAsset", err)
	}
	if cmp.Equal(*foundStepAsset, targetData) == false {
		t.Errorf("Expected StepAsset From Method GetStepAsset: %v is different from actual %v", foundStepAsset, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	stepAssetID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WithArgs(stepAssetID).WillReturnRows(noRows)
	foundStepAsset, err := GetStepAsset(mock, &stepAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStepAsset", err)
	}
	if foundStepAsset != nil {
		t.Errorf("Expected StepAsset From Method GetStepAsset: to be empty but got this: %v", foundStepAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	stepAssetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WithArgs(stepAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStepAsset, err := GetStepAsset(mock, &stepAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepAsset", err)
	}
	if foundStepAsset != nil {
		t.Errorf("Expected StepAsset From Method GetStepAsset: to be empty but got this: %v", foundStepAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	stepAssetID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WithArgs(stepAssetID).WillReturnRows(differentModelRows)
	foundStepAsset, err := GetStepAsset(mock, &stepAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepAsset", err)
	}
	if foundStepAsset != nil {
		t.Errorf("Expected StepAsset From Method GetStepAsset: to be empty but got this: %v", foundStepAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStepAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	stepAssetID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM step_assets").WithArgs(*stepAssetID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveStepAsset(mock, stepAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveStepAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStepAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	stepAssetID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveStepAsset(mock, &stepAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveStepAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	stepAssetID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM step_assets").WithArgs(stepAssetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveStepAsset(mock, &stepAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStepAssetToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WillReturnRows(mockRows)
	ids := []int{1}
	foundStepAssetList, err := GetStepAssets(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetStepAssets", err)
	}
	for i, sourceStepAsset := range dataList {
		if cmp.Equal(sourceStepAsset, foundStepAssetList[i]) == false {
			t.Errorf("Expected StepAsset From Method GetStepAssets: %v is different from actual %v", sourceStepAsset, foundStepAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1}
	foundStepAssetList, err := GetStepAssets(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepAssets", err)
	}
	if len(foundStepAssetList) != 0 {
		t.Errorf("Expected From Method GetStepAssets: to be empty but got this: %v", foundStepAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WillReturnRows(differentModelRows)
	ids := []int{1}
	foundStepAsset, err := GetStepAssets(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepAssets", err)
	}
	if foundStepAsset != nil {
		t.Errorf("Expected StepAsset From Method GetStepAssets: to be empty but got this: %v", foundStepAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStepAssetToMockRows(mock, dataList)
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WithArgs(pq.Array(UUIDList)).WillReturnRows(mockRows)
	foundStepAssetList, err := GetStepAssetsByUUIDs(mock, UUIDList)
	if err != nil {
		t.Fatalf("an error '%s' in GetStepAssetsByUUIDs", err)
	}
	for i, sourceStepAsset := range dataList {
		if cmp.Equal(sourceStepAsset, foundStepAssetList[i]) == false {
			t.Errorf("Expected StepAsset From Method GetStepAssetsByUUIDs: %v is different from actual %v", sourceStepAsset, foundStepAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WithArgs(pq.Array(UUIDList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStepAssetList, err := GetStepAssetsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepAssetsByUUIDs", err)
	}
	if len(foundStepAssetList) != 0 {
		t.Errorf("Expected From Method GetStepAssetsByUUIDs: to be empty but got this: %v", foundStepAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetsByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WithArgs(pq.Array(UUIDList)).WillReturnRows(differentModelRows)
	foundStepAsset, err := GetStepAssetsByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepAssetsByUUIDs", err)
	}
	if foundStepAsset != nil {
		t.Errorf("Expected StepAsset From Method GetStepAssetsByUUIDs: to be empty but got this: %v", foundStepAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStepAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddStepAssetToMockRows(mock, dataList)
	diffInDate := 2
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WithArgs(diffInDate).WillReturnRows(mockRows)
	foundStepAssetList, err := GetStartAndEndDateDiffStepAssets(mock, &diffInDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndDateDiffStepAssets", err)
	}
	for i, sourceStepAsset := range dataList {
		if cmp.Equal(sourceStepAsset, foundStepAssetList[i]) == false {
			t.Errorf("Expected StepAsset From Method GetStartAndEndDateDiffStepAssets: %v is different from actual %v", sourceStepAsset, foundStepAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStepAssetsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := -1
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WithArgs(diffInDate).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStepAssetList, err := GetStartAndEndDateDiffStepAssets(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffStepAssets", err)
	}
	if len(foundStepAssetList) != 0 {
		t.Errorf("Expected From Method GetStartAndEndDateDiffStepAssets: to be empty but got this: %v", foundStepAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffStepAssetsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := 1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WithArgs(diffInDate).WillReturnRows(differentModelRows)
	foundStepAsset, err := GetStartAndEndDateDiffStepAssets(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffStepAssets", err)
	}
	if foundStepAsset != nil {
		t.Errorf("Expected StepAsset From Method GetStartAndEndDateDiffStepAssets: to be empty but got this: %v", foundStepAsset)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateStepAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE step_assets").WithArgs(
		targetData.StepID,          //1
		targetData.AssetID,         //2
		targetData.SwapAssetID,     //3
		targetData.TargetPoolID,    //4
		targetData.UUID,            //5
		targetData.Name,            //6
		targetData.AlternateName,   //7
		targetData.StartDate,       //8
		targetData.EndDate,         //9
		targetData.Description,     //10
		targetData.ActionParameter, //11
		targetData.UpdatedBy,       //12
		targetData.ID,              //13
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateStepAsset(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateStepAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStepAssetOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateStepAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStepAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateStepAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateStepAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE step_assets").WithArgs(
		targetData.StepID,          //1
		targetData.AssetID,         //2
		targetData.SwapAssetID,     //3
		targetData.TargetPoolID,    //4
		targetData.UUID,            //5
		targetData.Name,            //6
		targetData.AlternateName,   //7
		targetData.StartDate,       //8
		targetData.EndDate,         //9
		targetData.Description,     //10
		targetData.ActionParameter, //11
		targetData.UpdatedBy,       //12
		targetData.ID,              //13
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateStepAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStepAsset(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO step_assets").WithArgs(
		targetData.StepID,          //1
		targetData.AssetID,         //2
		targetData.SwapAssetID,     //3
		targetData.TargetPoolID,    //4
		targetData.Name,            //5
		targetData.AlternateName,   //6
		targetData.StartDate,       //7
		targetData.EndDate,         //8
		targetData.Description,     //9
		targetData.ActionParameter, //10
		targetData.CreatedBy,       //11
	).WillReturnRows(pgxmock.NewRows([]string{"step_asset_id"}).AddRow(1))
	mock.ExpectCommit()
	stepID, err := InsertStepAsset(mock, &targetData)
	if stepID < 0 {
		t.Fatalf("stepID should not be negative ID: %d", stepID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertStepAsset", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStepAssetOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertStepAsset(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStepAssetOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO step_assets").WithArgs(
		targetData.StepID,          //1
		targetData.AssetID,         //2
		targetData.SwapAssetID,     //3
		targetData.TargetPoolID,    //4
		targetData.Name,            //5
		targetData.AlternateName,   //6
		targetData.StartDate,       //7
		targetData.EndDate,         //8
		targetData.Description,     //9
		targetData.ActionParameter, //10
		targetData.CreatedBy,       //11
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	stepID, err := InsertStepAsset(mock, &targetData)
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

func TestInsertStepAssetOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO step_assets").WithArgs(
		targetData.StepID,          //1
		targetData.AssetID,         //2
		targetData.SwapAssetID,     //3
		targetData.TargetPoolID,    //4
		targetData.Name,            //5
		targetData.AlternateName,   //6
		targetData.StartDate,       //7
		targetData.EndDate,         //8
		targetData.Description,     //9
		targetData.ActionParameter, //10
		targetData.CreatedBy,       //11
	).WillReturnRows(pgxmock.NewRows([]string{"step_asset_id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	stepID, err := InsertStepAsset(mock, &targetData)
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

func TestInsertStepAssets(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"step_assets"}, DBColumnsInsertStepAssets)
	targetData := TestAllData
	err = InsertStepAssets(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertStepAssets", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertStepAssetsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"step_assets"}, DBColumnsInsertStepAssets).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertStepAssets(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddStepAssetToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"pool_id = 1", "parent_step_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WillReturnRows(mockRows)
	foundStepAssetList, err := GetStepAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetStepAssetListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundStepAssetList[i]) == false {
			t.Errorf("Expected sourceData From Method GetStepAssetListByPagination: %v is different from actual %v", sourceData, foundStepAssetList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundStepAssetList, err := GetStepAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepAssetListByPagination", err)
	}
	if len(foundStepAssetList) != 0 {
		t.Errorf("Expected From Method GetStepAssetListByPagination: to be empty but got this: %v", foundStepAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStepAssetListByPaginationForCollectRowsErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM step_assets").WillReturnRows(differentModelRows)
	foundStepAssetList, err := GetStepAssetListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStepAssetListByPagination", err)
	}
	if len(foundStepAssetList) != 0 {
		t.Errorf("Expected From Method GetStepAssetListByPagination: to be empty but got this: %v", foundStepAssetList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalStepAssetsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM step_assets").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalStepAssetsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalStepAssetsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalStepAssetsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM step_assets").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalStepAssetsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalStepAssetsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalStepAssetsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
