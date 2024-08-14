package pool

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"id",              //1
	"target_asset_id", //2
	"strategy_id",     //3
	"account_id",      //4
	"uuid",            //5
	"name",            //6
	"alternate_name",  //7
	"start_date",      //8
	"end_date",        //9
	"description",     //10
	"chain_id",        //11
	"frequency_id",    //12
	"created_by",      //13
	"created_at",      //14
	"updated_by",      //15
	"updated_at",      //16
}
var DBColumnsInsertPools = []string{
	"target_asset_id", //1
	"strategy_id",     //2
	"account_id",      //3
	"uuid",            //5
	"name",            //6
	"alternate_name",  //6
	"start_date",      //7
	"end_date",        //8
	"description",     //9
	"chain_id",        //10
	"frequency_id",    //11
	"created_by",      //12
	"created_at",      //13
	"updated_by",      //14
	"updated_at",      //15
}

var TestData1 = Pool{
	ID:            utils.Ptr[int](1),                                     //1
	TargetAssetID: utils.Ptr[int](2),                                     //2
	StrategyID:    utils.Ptr[int](1),                                     //3
	AccountID:     utils.Ptr[int](1),                                     //4
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72",                //5
	Name:          "VTX Staking ",                                        //6
	AlternateName: "VTX",                                                 //7
	StartDate:     utils.Ptr[time.Time](utils.SampleCreatedAtTime),       //8
	EndDate:       utils.Ptr[time.Time](utils.SampleCreatedAtTime),       //9
	Description:   "Vector staking. 1. Stake VTX 2. Lock VTX (16 weeks)", //10
	ChainID:       utils.Ptr[int](65),                                    //11
	FrequencyID:   utils.Ptr[int](13),                                    //12
	CreatedBy:     "SYSTEM",                                              //13
	CreatedAt:     utils.SampleCreatedAtTime,                             //14
	UpdatedBy:     "SYSTEM",                                              //15
	UpdatedAt:     utils.SampleCreatedAtTime,                             //16

}

var TestData2 = Pool{
	ID:            utils.Ptr[int](2),                               //1
	TargetAssetID: utils.Ptr[int](3),                               //2
	StrategyID:    utils.Ptr[int](2),                               //3
	AccountID:     utils.Ptr[int](2),                               //4
	UUID:          "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",          //5
	Name:          "VELO/OP LP",                                    //6
	AlternateName: "VELO/OP LP",                                    //7
	StartDate:     utils.Ptr[time.Time](utils.SampleCreatedAtTime), //8
	EndDate:       utils.Ptr[time.Time](utils.SampleCreatedAtTime), //9
	Description:   "VELO/OP LP",                                    //10
	ChainID:       utils.Ptr[int](12),                              //11
	FrequencyID:   utils.Ptr[int](23),                              //12
	CreatedBy:     "SYSTEM",                                        //13
	CreatedAt:     utils.SampleCreatedAtTime,                       //14
	UpdatedBy:     "SYSTEM",                                        //15
	UpdatedAt:     utils.SampleCreatedAtTime,                       //16
}
var TestAllData = []Pool{TestData1, TestData2}

func AddPoolToMockRows(mock pgxmock.PgxPoolIface, dataList []Pool) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,            //1
			data.TargetAssetID, //2
			data.StrategyID,    //3
			data.AccountID,     //4
			data.UUID,          //5
			data.Name,          //6
			data.AlternateName, //7
			data.StartDate,     //8
			data.EndDate,       //9
			data.Description,   //10
			data.ChainID,       //11
			data.FrequencyID,   //12
			data.CreatedBy,     //13
			data.CreatedAt,     //14
			data.UpdatedBy,     //15
			data.UpdatedAt,     //16
		)
	}
	return rows
}

func TestGetPool(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Pool{targetData}
	poolID := targetData.ID
	mockRows := AddPoolToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(*poolID).WillReturnRows(mockRows)
	foundPool, err := GetPool(mock, poolID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPool", err)
	}
	if cmp.Equal(*foundPool, targetData) == false {
		t.Errorf("Expected Pool From Method GetPool: %v is different from actual %v", foundPool, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	poolID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(poolID).WillReturnRows(noRows)
	foundPool, err := GetPool(mock, &poolID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPool", err)
	}
	if foundPool != nil {
		t.Errorf("Expected Pool From Method GetPool: to be empty but got this: %v", foundPool)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	poolID := -1
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(poolID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPool, err := GetPool(mock, &poolID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPool", err)
	}
	if foundPool != nil {
		t.Errorf("Expected Pool From Method GetPool: to be empty but got this: %v", foundPool)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	poolID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(poolID).WillReturnRows(differentModelRows)
	foundPool, err := GetPool(mock, &poolID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPool", err)
	}
	if foundPool != nil {
		t.Errorf("Expected Pool From Method GetPool: to be empty but got this: %v", foundPool)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePool(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	poolID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM pools").WithArgs(*poolID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemovePool(mock, poolID)
	if err != nil {
		t.Fatalf("an error '%s' in RemovePool", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePoolOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	poolID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemovePool(mock, &poolID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePoolOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	poolID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM pools").WithArgs(poolID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemovePool(mock, &poolID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPools(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddPoolToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM pools").WillReturnRows(mockRows)
	ids := []int{1}
	foundPoolList, err := GetPools(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetPools", err)
	}
	for i, sourcePool := range dataList {
		if cmp.Equal(sourcePool, foundPoolList[i]) == false {
			t.Errorf("Expected Pool From Method GetPools: %v is different from actual %v", sourcePool, foundPoolList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM pools").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1}
	foundPoolList, err := GetPools(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPools", err)
	}
	if len(foundPoolList) != 0 {
		t.Errorf("Expected From Method GetPools: to be empty but got this: %v", foundPoolList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM pools").WillReturnRows(differentModelRows)
	ids := []int{1}
	foundPool, err := GetPools(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPools", err)
	}
	if foundPool != nil {
		t.Errorf("Expected Pool From Method GetPools: to be empty but got this: %v", foundPool)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolsByStrategyID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddPoolToMockRows(mock, dataList)
	strategyID := 1
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(strategyID).WillReturnRows(mockRows)
	foundPoolList, err := GetPoolsByStrategyID(mock, &strategyID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPoolsByStrategyID", err)
	}
	for i, sourcePool := range dataList {
		if cmp.Equal(sourcePool, foundPoolList[i]) == false {
			t.Errorf("Expected Pool From Method GetPoolsByStrategyID: %v is different from actual %v", sourcePool, foundPoolList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolsByStrategyIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	strategyID := -1
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(strategyID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPoolList, err := GetPoolsByStrategyID(mock, &strategyID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPoolsByStrategyID", err)
	}
	if len(foundPoolList) != 0 {
		t.Errorf("Expected From Method GetPoolsByStrategyID: to be empty but got this: %v", foundPoolList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolsByStrategyIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	strategyID := 1
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(strategyID).WillReturnRows(differentModelRows)
	foundPool, err := GetPoolsByStrategyID(mock, &strategyID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPoolsByStrategyID", err)
	}
	if foundPool != nil {
		t.Errorf("Expected Pool From Method GetPoolsByStrategyID: to be empty but got this: %v", foundPool)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestGetPoolsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddPoolToMockRows(mock, dataList)
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(pq.Array(uuids)).WillReturnRows(mockRows)
	foundPoolList, err := GetPoolsByUUIDs(mock, uuids)
	if err != nil {
		t.Fatalf("an error '%s' in GetPoolsByUUIDs", err)
	}
	for i, sourcePool := range dataList {
		if cmp.Equal(sourcePool, foundPoolList[i]) == false {
			t.Errorf("Expected Pool From Method GetPoolsByUUIDs: %v is different from actual %v", sourcePool, foundPoolList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(pq.Array(uuids)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPoolList, err := GetPoolsByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPoolsByUUIDs", err)
	}
	if len(foundPoolList) != 0 {
		t.Errorf("Expected From Method GetPoolsByUUIDs: to be empty but got this: %v", foundPoolList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolsByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(pq.Array(uuids)).WillReturnRows(differentModelRows)
	foundPool, err := GetPoolsByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPoolsByUUIDs", err)
	}
	if foundPool != nil {
		t.Errorf("Expected Pool From Method GetPoolsByUUIDs: to be empty but got this: %v", foundPool)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestGetStartAndEndDateDiffPools(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddPoolToMockRows(mock, dataList)
	diffInDate := 1
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(diffInDate).WillReturnRows(mockRows)
	foundPoolList, err := GetStartAndEndDateDiffPools(mock, &diffInDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndDateDiffPools", err)
	}
	for i, sourcePool := range dataList {
		if cmp.Equal(sourcePool, foundPoolList[i]) == false {
			t.Errorf("Expected Pool From Method GetStartAndEndDateDiffPools: %v is different from actual %v", sourcePool, foundPoolList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffPoolsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := -1
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(diffInDate).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPoolList, err := GetStartAndEndDateDiffPools(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffPools", err)
	}
	if len(foundPoolList) != 0 {
		t.Errorf("Expected From Method GetStartAndEndDateDiffPools: to be empty but got this: %v", foundPoolList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffPoolsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	diffInDate := 1
	mock.ExpectQuery("^SELECT (.+) FROM pools").WithArgs(diffInDate).WillReturnRows(differentModelRows)
	foundPool, err := GetStartAndEndDateDiffPools(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffPools", err)
	}
	if foundPool != nil {
		t.Errorf("Expected Pool From Method GetStartAndEndDateDiffPools: to be empty but got this: %v", foundPool)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdatePool(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE pools").WithArgs(
		targetData.TargetAssetID, //1
		targetData.StrategyID,    //2
		targetData.AccountID,     //3
		targetData.UUID,          //4
		targetData.Name,          //5
		targetData.AlternateName, //6
		targetData.StartDate,     //7
		targetData.EndDate,       //8
		targetData.Description,   //9
		targetData.ChainID,       //10
		targetData.FrequencyID,   //11
		targetData.UpdatedBy,     //12
		targetData.ID,            //13
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdatePool(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdatePool", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePoolOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdatePool(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePoolOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdatePool(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePoolOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE pools").WithArgs(
		targetData.TargetAssetID, //1
		targetData.StrategyID,    //2
		targetData.AccountID,     //3
		targetData.UUID,          //4
		targetData.Name,          //5
		targetData.AlternateName, //6
		targetData.StartDate,     //7
		targetData.EndDate,       //8
		targetData.Description,   //9
		targetData.ChainID,       //10
		targetData.FrequencyID,   //11
		targetData.UpdatedBy,     //12
		targetData.ID,            //13
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdatePool(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPool(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO pools").WithArgs(
		targetData.TargetAssetID, //1
		targetData.StrategyID,    //2
		targetData.AccountID,     //3
		// &targetData.UUID, //
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.StartDate,     //6
		targetData.EndDate,       //7
		targetData.Description,   //8
		targetData.ChainID,       //9
		targetData.FrequencyID,   //10
		targetData.CreatedBy,     //11
	).WillReturnRows(pgxmock.NewRows([]string{"pool_id", "pool_uuid"}).AddRow(1, "return-uuid"))
	mock.ExpectCommit()
	poolID, uuid, err := InsertPool(mock, &targetData)
	if poolID < 0 {
		t.Fatalf("poolID should not be negative ID: %d", poolID)
	}
	if uuid == "" {
		t.Fatalf("uuid should not be empty UUID: %s", uuid)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertPool", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPoolOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertPool(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPoolOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO pools").WithArgs(
		targetData.TargetAssetID, //1
		targetData.StrategyID,    //2
		targetData.AccountID,     //3
		// &targetData.UUID, //
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.StartDate,     //6
		targetData.EndDate,       //7
		targetData.Description,   //8
		targetData.ChainID,       //9
		targetData.FrequencyID,   //10
		targetData.CreatedBy,     //11
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	poolID, uuid, err := InsertPool(mock, &targetData)
	if poolID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error poolID: %d", poolID)
	}
	if uuid != "" {
		t.Fatalf("Expecting empty for uuid because of error uuid: %s", uuid)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPoolOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO pools").WithArgs(
		targetData.TargetAssetID, //1
		targetData.StrategyID,    //2
		targetData.AccountID,     //3
		// &targetData.UUID, //
		targetData.Name,          //4
		targetData.AlternateName, //5
		targetData.StartDate,     //6
		targetData.EndDate,       //7
		targetData.Description,   //8
		targetData.ChainID,       //9
		targetData.FrequencyID,   //10
		targetData.CreatedBy,     //11
	).WillReturnRows(pgxmock.NewRows([]string{"pool_id", "pool_uuid"}).AddRow(1, "return-uuid"))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	poolID, uuid, err := InsertPool(mock, &targetData)
	if poolID >= 0 {
		t.Fatalf("Expecting -1 for poolID because of error poolID: %d", poolID)
	}
	if uuid != "" {
		t.Fatalf("Expecting empty for uuid because of error uuid: %s", uuid)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPools(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"pools"}, DBColumnsInsertPools)
	targetData := TestAllData
	err = InsertPools(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertPools", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPoolsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"pools"}, DBColumnsInsertPools).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertPools(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddPoolToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"target_asset_id = 1", "strategy_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM pools").WillReturnRows(mockRows)
	foundPoolList, err := GetPoolListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetPoolListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundPoolList[i]) == false {
			t.Errorf("Expected sourceData From Method GetPoolListByPagination: %v is different from actual %v", sourceData, foundPoolList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"strategy_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM pools").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPoolList, err := GetPoolListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPoolListByPagination", err)
	}
	if len(foundPoolList) != 0 {
		t.Errorf("Expected From Method GetPoolListByPagination: to be empty but got this: %v", foundPoolList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPoolListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"strategy_id = -1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM pools").WillReturnRows(differentModelRows)
	foundPoolList, err := GetPoolListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPoolListByPagination", err)
	}
	if len(foundPoolList) != 0 {
		t.Errorf("Expected From Method GetPoolListByPagination: to be empty but got this: %v", foundPoolList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalPoolsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM pools").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalPoolsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalPoolsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalPoolsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM pools").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalPoolsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalPoolsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalPoolsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
