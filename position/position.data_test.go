package position

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
	"id",             //1
	"uuid",           //2
	"name",           //3
	"alternate_name", //4
	"account_id",     //5
	"portfolio_id",   //6
	"frequency_id",   //7
	"start_date",     //8
	"end_date",       //9
	"base_asset_id",  //10
	"quote_asset_id", //11
	"quantity",       //12
	"cost_basis",     //13
	"profit",         //14
	"total_amount",   //15
	"description",    //16
	"created_by",     //17
	"created_at",     //18
	"updated_by",     //19
	"updated_at",     //20
}
var DBColumnsInsertPositions = []string{
	"uuid",           //1
	"name",           //2
	"alternate_name", //3
	"account_id",     //4
	"portfolio_id",   //5
	"frequency_id",   //6
	"start_date",     //7
	"end_date",       //8
	"base_asset_id",  //9
	"quote_asset_id", //10
	"quantity",       //11
	"cost_basis",     //12
	"profit",         //13
	"total_amount",   //14
	"description",    //15
	"created_by",     //16
	"created_at",     //17
	"updated_by",     //18
	"updated_at",     //19
}

var TestData1 = Position{
	ID:            utils.Ptr[int](1),                                     //1
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72",                //2
	Name:          "November 22, 2013 Base : Kyber Network, Quote : USD", //3
	AlternateName: "November 22, 2013 Base : Kyber Network, Quote : USD", //4
	AccountID:     utils.Ptr[int](1),                                     //5
	PortfolioID:   utils.Ptr[int](1),                                     //6
	FrequnecyID:   utils.Ptr[int](1),                                     //7
	StartDate:     utils.SampleCreatedAtTime,                             //8
	EndDate:       utils.SampleCreatedAtTime,                             //9
	BaseAssetID:   utils.Ptr[int](2),                                     //10
	QuoteAssetID:  utils.Ptr[int](1),                                     //11
	Quantity:      utils.Ptr[float64](100.1),                             //12
	CostBasis:     utils.Ptr[float64](90),                                //13
	Profit:        utils.Ptr[float64](10.5),                              //14
	TotalAmount:   utils.Ptr[float64](1001),                              //15
	Description:   "",                                                    //16
	CreatedBy:     "SYSTEM",                                              //17
	CreatedAt:     utils.SampleCreatedAtTime,                             //18
	UpdatedBy:     "SYSTEM",                                              //19
	UpdatedAt:     utils.SampleCreatedAtTime,                             //20

}

var TestData2 = Position{
	ID:            utils.Ptr[int](2),                              //1
	UUID:          "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",         //2
	Name:          "November 22, 2013 Base : Orchid, Quote : USD", //3
	AlternateName: "November 22, 2013 Base : Orchid, Quote : USD", //4
	AccountID:     utils.Ptr[int](2),                              //5
	PortfolioID:   utils.Ptr[int](2),                              //6
	FrequnecyID:   utils.Ptr[int](1),                              //7
	StartDate:     utils.SampleCreatedAtTime,                      //8
	EndDate:       utils.SampleCreatedAtTime,                      //9
	BaseAssetID:   utils.Ptr[int](4),                              //10
	QuoteAssetID:  utils.Ptr[int](1),                              //11
	Quantity:      utils.Ptr[float64](1001),                       //12
	CostBasis:     utils.Ptr[float64](10.1),                       //13
	Profit:        utils.Ptr[float64](1.5),                        //14
	TotalAmount:   utils.Ptr[float64](1002),                       //15
	Description:   "",                                             //16
	CreatedBy:     "SYSTEM",                                       //17
	CreatedAt:     utils.SampleCreatedAtTime,                      //18
	UpdatedBy:     "SYSTEM",                                       //19
	UpdatedAt:     utils.SampleCreatedAtTime,                      //20
}
var TestAllData = []Position{TestData1, TestData2}

func AddPositionToMockRows(mock pgxmock.PgxPoolIface, dataList []Position) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,            //1
			data.UUID,          //2
			data.Name,          //3
			data.AlternateName, //4
			data.AccountID,     //5
			data.PortfolioID,   //6
			data.FrequnecyID,   //7
			data.StartDate,     //8
			data.EndDate,       //9
			data.BaseAssetID,   //10
			data.QuoteAssetID,  //11
			data.Quantity,      //12
			data.CostBasis,     //13
			data.Profit,        //14
			data.TotalAmount,   //15
			data.Description,   //16
			data.CreatedBy,     //17
			data.CreatedAt,     //18
			data.UpdatedBy,     //19
			data.UpdatedAt,     //20
		)
	}
	return rows
}

func TestGetPosition(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Position{targetData}
	positionID := targetData.ID
	mockRows := AddPositionToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM positions").WithArgs(*positionID).WillReturnRows(mockRows)
	foundPosition, err := GetPosition(mock, positionID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPosition", err)
	}
	if cmp.Equal(*foundPosition, targetData) == false {
		t.Errorf("Expected Position From Method GetPosition: %v is different from actual %v", foundPosition, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM positions").WithArgs(positionID).WillReturnRows(noRows)
	foundPosition, err := GetPosition(mock, &positionID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPosition", err)
	}
	if foundPosition != nil {
		t.Errorf("Expected Position From Method GetPosition: to be empty but got this: %v", foundPosition)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionID := -1
	mock.ExpectQuery("^SELECT (.+) FROM positions").WithArgs(positionID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPosition, err := GetPosition(mock, &positionID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPosition", err)
	}
	if foundPosition != nil {
		t.Errorf("Expected Position From Method GetPosition: to be empty but got this: %v", foundPosition)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM positions").WithArgs(positionID).WillReturnRows(differentModelRows)
	foundPosition, err := GetPosition(mock, &positionID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPosition", err)
	}
	if foundPosition != nil {
		t.Errorf("Expected Position From Method GetPosition: to be empty but got this: %v", foundPosition)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePosition(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	positionID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM positions").WithArgs(*positionID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemovePosition(mock, positionID)
	if err != nil {
		t.Fatalf("an error '%s' in RemovePosition", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePositionOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	positionID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemovePosition(mock, &positionID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePositionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	positionID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM positions").WithArgs(positionID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemovePosition(mock, &positionID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositions(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddPositionToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM positions").WillReturnRows(mockRows)
	ids := []int{1}
	foundPositionList, err := GetPositions(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetPositions", err)
	}
	for i, sourcePosition := range dataList {
		if cmp.Equal(sourcePosition, foundPositionList[i]) == false {
			t.Errorf("Expected Position From Method GetPositions: %v is different from actual %v", sourcePosition, foundPositionList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM positions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1}
	foundPositionList, err := GetPositions(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositions", err)
	}
	if len(foundPositionList) != 0 {
		t.Errorf("Expected From Method GetPositions: to be empty but got this: %v", foundPositionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM positions").WillReturnRows(differentModelRows)
	ids := []int{1}
	foundPosition, err := GetPositions(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositions", err)
	}
	if foundPosition != nil {
		t.Errorf("Expected Position From Method GetPositions: to be empty but got this: %v", foundPosition)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdatePosition(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE positions").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.UserEmail,      //5
		targetData.Description,    //6
		targetData.BaseAssetID,    //7
		targetData.PositionTypeID, //8
		targetData.ParentID,       //9
		targetData.UpdatedBy,      //10
		targetData.ID,             //11
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdatePosition(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdatePosition", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePositionOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdatePosition(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePositionOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdatePosition(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePositionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE positions").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.UserEmail,      //5
		targetData.Description,    //6
		targetData.BaseAssetID,    //7
		targetData.PositionTypeID, //8
		targetData.ParentID,       //9
		targetData.UpdatedBy,      //10
		targetData.ID,             //11
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdatePosition(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPosition(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO positions").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.UserEmail,      //5
		targetData.Description,    //6
		targetData.BaseAssetID,    //7
		targetData.PositionTypeID, //8
		targetData.ParentID,       //9
		targetData.CreatedBy,      //10
	).WillReturnRows(pgxmock.NewRows([]string{"market_data_id", "job_id"}).AddRow(1, "return-uuid"))
	mock.ExpectCommit()
	positionID, uuid, err := InsertPosition(mock, &targetData)
	if positionID < 0 {
		t.Fatalf("positionID should not be negative ID: %d", positionID)
	}
	if uuid == "" {
		t.Fatalf("uuid should not be empty UUID: %s", uuid)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertPosition", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPositionOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertPosition(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPositionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO positions").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.UserEmail,      //5
		targetData.Description,    //6
		targetData.BaseAssetID,    //7
		targetData.PositionTypeID, //8
		targetData.ParentID,       //9
		targetData.CreatedBy,      //10
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	positionID, uuid, err := InsertPosition(mock, &targetData)
	if positionID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error positionID: %d", positionID)
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

func TestInsertPositionOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO positions").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.StartDate,      //3
		targetData.EndDate,        //4
		targetData.UserEmail,      //5
		targetData.Description,    //6
		targetData.BaseAssetID,    //7
		targetData.PositionTypeID, //8
		targetData.ParentID,       //9
		targetData.CreatedBy,      //10
	).WillReturnRows(pgxmock.NewRows([]string{"market_data_id", "job_id"}).AddRow(1, "return-uuid"))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	positionID, uuid, err := InsertPosition(mock, &targetData)
	if positionID >= 0 {
		t.Fatalf("Expecting -1 for positionID because of error positionID: %d", positionID)
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

func TestInsertPositions(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"positions"}, DBColumnsInsertPositions)
	targetData := TestAllData
	err = InsertPositions(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertPositions", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPositionsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"positions"}, DBColumnsInsertPositions).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertPositions(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddPositionToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"target_asset_id = 1", "strategy_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM positions").WillReturnRows(mockRows)
	foundPositionList, err := GetPositionListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetPositionListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundPositionList[i]) == false {
			t.Errorf("Expected sourceData From Method GetPositionListByPagination: %v is different from actual %v", sourceData, foundPositionList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM positions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPositionList, err := GetPositionListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionListByPagination", err)
	}
	if len(foundPositionList) != 0 {
		t.Errorf("Expected From Method GetPositionListByPagination: to be empty but got this: %v", foundPositionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPositionListByPaginationForCollectRowsErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM positions").WillReturnRows(differentModelRows)
	foundPositionList, err := GetPositionListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPositionListByPagination", err)
	}
	if len(foundPositionList) != 0 {
		t.Errorf("Expected From Method GetPositionListByPagination: to be empty but got this: %v", foundPositionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalPositionsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM positions").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalPositionsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalPositionsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalPositionsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM positions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalPositionsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalPositionsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalPositionsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
