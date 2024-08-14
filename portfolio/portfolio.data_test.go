package portfolio

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"id",                //1
	"uuid",              //2
	"name",              //3
	"alternate_name",    //4
	"start_date",        //5
	"end_date",          //6
	"user_email",        //7
	"description",       //8
	"base_asset_id",     //9
	"portfolio_type_id", //10
	"parent_id",         //11
	"created_by",        //12
	"created_at",        //13
	"updated_by",        //14
	"updated_at",        //15
}
var DBColumnsInsertPortfolios = []string{
	"uuid",              //1
	"name",              //2
	"alternate_name",    //3
	"start_date",        //4
	"end_date",          //5
	"user_email",        //6
	"description",       //7
	"base_asset_id",     //8
	"portfolio_type_id", //9
	"parent_id",         //10
	"created_by",        //11
	"created_at",        //12
	"updated_by",        //13
	"updated_at",        //14
}

var TestData1 = Portfolio{
	ID:              utils.Ptr[int](1),                               //1
	UUID:            "01ef85e8-2c26-441e-8c7f-71d79518ad72",          //2
	Name:            "Crypto ",                                       //3
	AlternateName:   "Crypto",                                        //4
	StartDate:       utils.Ptr[time.Time](utils.SampleCreatedAtTime), //5
	EndDate:         utils.Ptr[time.Time](utils.SampleCreatedAtTime), //6
	UserEmail:       "system@example.com",                            //7
	Description:     "",                                              //8
	BaseAssetID:     utils.Ptr[int](1),                               //9
	PortfolioTypeID: utils.Ptr[int](1),                               //10
	ParentID:        nil,                                             //11
	CreatedBy:       "SYSTEM",                                        //12
	CreatedAt:       utils.SampleCreatedAtTime,                       //13
	UpdatedBy:       "SYSTEM",                                        //14
	UpdatedAt:       utils.SampleCreatedAtTime,                       //15

}

var TestData2 = Portfolio{
	ID:              utils.Ptr[int](2),                               //1
	UUID:            "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",          //2
	Name:            "Defi",                                          //3
	AlternateName:   "Defi",                                          //4
	StartDate:       utils.Ptr[time.Time](utils.SampleCreatedAtTime), //5
	EndDate:         utils.Ptr[time.Time](utils.SampleCreatedAtTime), //6
	UserEmail:       "system@example.com",                            //7
	Description:     "",                                              //8
	BaseAssetID:     utils.Ptr[int](1),                               //9
	PortfolioTypeID: utils.Ptr[int](2),                               //10
	ParentID:        utils.Ptr[int](1),                               //11
	CreatedBy:       "SYSTEM",                                        //12
	CreatedAt:       utils.SampleCreatedAtTime,                       //13
	UpdatedBy:       "SYSTEM",                                        //14
	UpdatedAt:       utils.SampleCreatedAtTime,                       //15
}
var TestAllData = []Portfolio{TestData1, TestData2}

func AddPortfolioToMockRows(mock pgxmock.PgxPoolIface, dataList []Portfolio) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,              //1
			data.UUID,            //2
			data.Name,            //3
			data.AlternateName,   //4
			data.StartDate,       //5
			data.EndDate,         //6
			data.UserEmail,       //7
			data.Description,     //8
			data.BaseAssetID,     //9
			data.PortfolioTypeID, //10
			data.ParentID,        //11
			data.CreatedBy,       //12
			data.CreatedAt,       //13
			data.UpdatedBy,       //14
			data.UpdatedAt,       //15
		)
	}
	return rows
}

func TestGetPortfolio(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Portfolio{targetData}
	portfolioID := targetData.ID
	mockRows := AddPortfolioToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM portfolios").WithArgs(*portfolioID).WillReturnRows(mockRows)
	foundPortfolio, err := GetPortfolio(mock, portfolioID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPortfolio", err)
	}
	if cmp.Equal(*foundPortfolio, targetData) == false {
		t.Errorf("Expected Portfolio From Method GetPortfolio: %v is different from actual %v", foundPortfolio, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPortfolioForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	portfolioID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM portfolios").WithArgs(portfolioID).WillReturnRows(noRows)
	foundPortfolio, err := GetPortfolio(mock, &portfolioID)
	if err != nil {
		t.Fatalf("an error '%s' in GetPortfolio", err)
	}
	if foundPortfolio != nil {
		t.Errorf("Expected Portfolio From Method GetPortfolio: to be empty but got this: %v", foundPortfolio)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPortfolioForCollectRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	portfolioID := -1
	mock.ExpectQuery("^SELECT (.+) FROM portfolios").WithArgs(portfolioID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPortfolio, err := GetPortfolio(mock, &portfolioID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPortfolio", err)
	}
	if foundPortfolio != nil {
		t.Errorf("Expected Portfolio From Method GetPortfolio: to be empty but got this: %v", foundPortfolio)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPortfolioForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	portfolioID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM portfolios").WithArgs(portfolioID).WillReturnRows(differentModelRows)
	foundPortfolio, err := GetPortfolio(mock, &portfolioID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPortfolio", err)
	}
	if foundPortfolio != nil {
		t.Errorf("Expected Portfolio From Method GetPortfolio: to be empty but got this: %v", foundPortfolio)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePortfolio(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	portfolioID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM portfolios").WithArgs(*portfolioID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemovePortfolio(mock, portfolioID)
	if err != nil {
		t.Fatalf("an error '%s' in RemovePortfolio", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePortfolioOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	portfolioID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemovePortfolio(mock, &portfolioID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemovePortfolioOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	portfolioID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM portfolios").WithArgs(portfolioID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemovePortfolio(mock, &portfolioID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPortfolios(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddPortfolioToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM portfolios").WillReturnRows(mockRows)
	ids := []int{1}
	foundPortfolioList, err := GetPortfolios(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetPortfolios", err)
	}
	for i, sourcePortfolio := range dataList {
		if cmp.Equal(sourcePortfolio, foundPortfolioList[i]) == false {
			t.Errorf("Expected Portfolio From Method GetPortfolios: %v is different from actual %v", sourcePortfolio, foundPortfolioList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPortfoliosForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM portfolios").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1}
	foundPortfolioList, err := GetPortfolios(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPortfolios", err)
	}
	if len(foundPortfolioList) != 0 {
		t.Errorf("Expected From Method GetPortfolios: to be empty but got this: %v", foundPortfolioList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPortfoliosForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM portfolios").WillReturnRows(differentModelRows)
	ids := []int{1}
	foundPortfolio, err := GetPortfolios(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPortfolios", err)
	}
	if foundPortfolio != nil {
		t.Errorf("Expected Portfolio From Method GetPortfolios: to be empty but got this: %v", foundPortfolio)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdatePortfolio(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE portfolios").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.StartDate,       //3
		targetData.EndDate,         //4
		targetData.UserEmail,       //5
		targetData.Description,     //6
		targetData.BaseAssetID,     //7
		targetData.PortfolioTypeID, //8
		targetData.ParentID,        //9
		targetData.UpdatedBy,       //10
		targetData.ID,              //11
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdatePortfolio(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdatePortfolio", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePortfolioOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdatePortfolio(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePortfolioOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdatePortfolio(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdatePortfolioOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE portfolios").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.StartDate,       //3
		targetData.EndDate,         //4
		targetData.UserEmail,       //5
		targetData.Description,     //6
		targetData.BaseAssetID,     //7
		targetData.PortfolioTypeID, //8
		targetData.ParentID,        //9
		targetData.UpdatedBy,       //10
		targetData.ID,              //11
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdatePortfolio(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPortfolio(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO portfolios").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.StartDate,       //3
		targetData.EndDate,         //4
		targetData.UserEmail,       //5
		targetData.Description,     //6
		targetData.BaseAssetID,     //7
		targetData.PortfolioTypeID, //8
		targetData.ParentID,        //9
		targetData.CreatedBy,       //10
	).WillReturnRows(pgxmock.NewRows([]string{"portfolio_id", "job_id"}).AddRow(1, "return-uuid"))
	mock.ExpectCommit()
	portfolioID, uuid, err := InsertPortfolio(mock, &targetData)
	if portfolioID < 0 {
		t.Fatalf("portfolioID should not be negative ID: %d", portfolioID)
	}
	if uuid == "" {
		t.Fatalf("uuid should not be empty UUID: %s", uuid)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertPortfolio", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPortfolioOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertPortfolio(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPortfolioOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO portfolios").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.StartDate,       //3
		targetData.EndDate,         //4
		targetData.UserEmail,       //5
		targetData.Description,     //6
		targetData.BaseAssetID,     //7
		targetData.PortfolioTypeID, //8
		targetData.ParentID,        //9
		targetData.CreatedBy,       //10
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	portfolioID, uuid, err := InsertPortfolio(mock, &targetData)
	if portfolioID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error portfolioID: %d", portfolioID)
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

func TestInsertPortfolioOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO portfolios").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.StartDate,       //3
		targetData.EndDate,         //4
		targetData.UserEmail,       //5
		targetData.Description,     //6
		targetData.BaseAssetID,     //7
		targetData.PortfolioTypeID, //8
		targetData.ParentID,        //9
		targetData.CreatedBy,       //10
	).WillReturnRows(pgxmock.NewRows([]string{"portfolio_id", "job_id"}).AddRow(1, "return-uuid"))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	portfolioID, uuid, err := InsertPortfolio(mock, &targetData)
	if portfolioID >= 0 {
		t.Fatalf("Expecting -1 for portfolioID because of error portfolioID: %d", portfolioID)
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

func TestInsertPortfolios(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"portfolios"}, DBColumnsInsertPortfolios)
	targetData := TestAllData
	err = InsertPortfolios(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertPortfolios", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertPortfoliosOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"portfolios"}, DBColumnsInsertPortfolios).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertPortfolios(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPortfolioListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddPortfolioToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"target_asset_id = 1", "strategy_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM portfolios").WillReturnRows(mockRows)
	foundPortfolioList, err := GetPortfolioListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetPortfolioListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundPortfolioList[i]) == false {
			t.Errorf("Expected sourceData From Method GetPortfolioListByPagination: %v is different from actual %v", sourceData, foundPortfolioList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPortfolioListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM portfolios").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundPortfolioList, err := GetPortfolioListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPortfolioListByPagination", err)
	}
	if len(foundPortfolioList) != 0 {
		t.Errorf("Expected From Method GetPortfolioListByPagination: to be empty but got this: %v", foundPortfolioList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetPortfolioListByPaginationForCollectRowsErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM portfolios").WillReturnRows(differentModelRows)
	foundPortfolioList, err := GetPortfolioListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetPortfolioListByPagination", err)
	}
	if len(foundPortfolioList) != 0 {
		t.Errorf("Expected From Method GetPortfolioListByPagination: to be empty but got this: %v", foundPortfolioList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalPortfoliosCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM portfolios").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalPortfoliosCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalPortfoliosCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalPortfoliosCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM portfolios").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalPortfoliosCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalPortfoliosCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalPortfoliosCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
