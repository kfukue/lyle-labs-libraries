package exchange

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

var columns = []string{
	"id",               //1
	"uuid",             //2
	"name",             //3
	"alternate_name",   //4
	"exchange_type_id", //5
	"url",              //6
	"start_date",       //7
	"end_date",         //8
	"description",      //9
	"created_by",       //10
	"created_at",       //11
	"updated_by",       //12
	"updated_at",       //13
}
var columnsInsertExchanges = []string{
	"uuid",             //1
	"name",             //2
	"alternate_name",   //3
	"exchange_type_id", //4
	"url",              //5
	"start_date",       //6
	"end_date",         //7
	"description",      //8
	"created_by",       //9
	"created_at",       //10
	"updated_by",       //11
	"updated_at",       //12
}

var TestData1 = Exchange{
	ID:             utils.Ptr[int](1),
	UUID:           "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:           "Coinbase",
	AlternateName:  "CB",
	ExchangeTypeID: utils.Ptr[int](1),
	Url:            "https://coinbase.com",
	StartDate:      utils.Ptr[time.Time](utils.SampleCreatedAtTime),
	EndDate:        utils.Ptr[time.Time](utils.SampleCreatedAtTime),
	Description:    "",
	CreatedBy:      "SYSTEM",
	CreatedAt:      utils.SampleCreatedAtTime,
	UpdatedBy:      "SYSTEM",
	UpdatedAt:      utils.SampleCreatedAtTime,
}

var TestData2 = Exchange{
	ID:             utils.Ptr[int](2),
	UUID:           "880607ab-2833-4ad7-a231-b983a61cad34",
	Name:           "Binance",
	AlternateName:  "BNB",
	ExchangeTypeID: utils.Ptr[int](1),
	Url:            "https://binance.com",
	StartDate:      utils.Ptr[time.Time](utils.SampleCreatedAtTime),
	EndDate:        utils.Ptr[time.Time](utils.SampleCreatedAtTime),
	Description:    "",
	CreatedBy:      "SYSTEM",
	CreatedAt:      utils.SampleCreatedAtTime,
	UpdatedBy:      "SYSTEM",
	UpdatedAt:      utils.SampleCreatedAtTime,
}
var TestAllData = []Exchange{TestData1, TestData2}

func AddExchangeToMockRows(mock pgxmock.PgxPoolIface, dataList []Exchange) *pgxmock.Rows {
	rows := mock.NewRows(columns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,             //1
			data.UUID,           //2
			data.Name,           //3
			data.AlternateName,  //4
			data.ExchangeTypeID, //5
			data.Url,            //6
			data.StartDate,      //7
			data.EndDate,        //8
			data.Description,    //9
			data.CreatedBy,      //10
			data.CreatedAt,      //11
			data.UpdatedBy,      //12
			data.UpdatedAt,      //13

		)
	}
	return rows
}

var columnsExchangeChains = []string{
	"uuid",        //1
	"exchange_id", //2
	"chain_id",    //3
	"description", //4
	"created_by",  //5
	"created_at",  //6
	"updated_by",  //7
	"updated_at",  //8
}

var columnsInsertExchangeChains = []string{
	"uuid",        //1
	"exchange_id", //2
	"chain_id",    //3
	"description", //4
	"created_by",  //5
	"created_at",  //6
	"updated_by",  //7
	"updated_at",  //8
}

var data1ExchangeChain = ExchangeChain{
	UUID:        "880607ab-2833-4ad7-a231-b983a61c7b39",
	ExchangeID:  utils.Ptr[int](1),
	ChainID:     utils.Ptr[int](1),
	Description: "",
	CreatedBy:   "SYSTEM",
	CreatedAt:   utils.SampleCreatedAtTime,
	UpdatedBy:   "SYSTEM",
	UpdatedAt:   utils.SampleCreatedAtTime,
}

var data2ExchangeChain = ExchangeChain{
	UUID:        "880607ab-2833-4ad7-a231-b983a61c7b39",
	ExchangeID:  utils.Ptr[int](2),
	ChainID:     utils.Ptr[int](2),
	Description: "",
	CreatedBy:   "SYSTEM",
	CreatedAt:   utils.SampleCreatedAtTime,
	UpdatedBy:   "SYSTEM",
	UpdatedAt:   utils.SampleCreatedAtTime,
}
var allDataExchangeEchains = []ExchangeChain{data1ExchangeChain, data2ExchangeChain}

func AddExchangeChainToMockRows(mock pgxmock.PgxPoolIface, dataList []ExchangeChain) *pgxmock.Rows {
	rows := mock.NewRows(columnsExchangeChains)
	for _, data := range dataList {
		rows.AddRow(
			data.UUID,        //1
			data.ExchangeID,  //2
			data.ChainID,     //3
			data.Description, //4
			data.CreatedBy,   //5
			data.CreatedAt,   //6
			data.UpdatedBy,   //7
			data.UpdatedAt,   //8

		)
	}
	return rows
}

func TestGetExchange(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Exchange{targetData}
	exchangeID := targetData.ID
	mockRows := AddExchangeToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WithArgs(*exchangeID).WillReturnRows(mockRows)
	foundExchange, err := GetExchange(mock, exchangeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetExchange", err)
	}
	if cmp.Equal(*foundExchange, targetData) == false {
		t.Errorf("Expected Exchange From Method GetExchange: %v is different from actual %v", foundExchange, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangeForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	exchangeID := 999
	noRows := pgxmock.NewRows(columns)
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WithArgs(exchangeID).WillReturnRows(noRows)
	foundExchange, err := GetExchange(mock, &exchangeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetExchange", err)
	}
	if foundExchange != nil {
		t.Errorf("Expected Exchange From Method GetExchange: to be empty but got this: %v", foundExchange)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangeForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	exchangeID := -1
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WithArgs(exchangeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundExchange, err := GetExchange(mock, &exchangeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetExchange", err)
	}
	if foundExchange != nil {
		t.Errorf("Expected Exchange From Method GetExchange: to be empty but got this: %v", foundExchange)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangeForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	exchangeID := -1
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WithArgs(exchangeID).WillReturnRows(differentModelRows)
	foundExchange, err := GetExchange(mock, &exchangeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetExchange", err)
	}
	if foundExchange != nil {
		t.Errorf("Expected foundExchange From Method GetExchange: to be empty but got this: %v", foundExchange)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveExchange(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	exchangeID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM exchanges").WithArgs(*exchangeID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveExchange(mock, exchangeID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveExchange", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveExchangeOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveExchange(mock, &taxID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveExchangeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	exchangeID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM exchanges").WithArgs(exchangeID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveExchange(mock, &exchangeID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestGetExchangeList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Exchange{TestData1, TestData2}
	mockRows := AddExchangeToMockRows(mock, dataList)
	ids := []int{*TestData1.ID, *TestData2.ID}
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WillReturnRows(mockRows)
	foundExchanges, err := GetExchangeList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetExchangeList", err)
	}
	testExchanges := TestAllData
	for i, foundExchange := range foundExchanges {
		if cmp.Equal(foundExchange, testExchanges[i]) == false {
			t.Errorf("Expected Exchange From Method GetExchangeList: %v is different from actual %v", foundExchange, testExchanges[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangeListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	ids := []int{*TestData1.ID, *TestData2.ID}
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundExchanges, err := GetExchangeList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetExchangeList", err)
	}
	if len(foundExchanges) != 0 {
		t.Errorf("Expected From Method GetExchangeList: to be empty but got this: %v", foundExchanges)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangeListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	ids := []int{*TestData1.ID, *TestData2.ID}
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WillReturnRows(differentModelRows)
	foundExchanges, err := GetExchangeList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetExchangeList", err)
	}
	if foundExchanges != nil {
		t.Errorf("Expected foundExchanges From Method GetExchangeList: to be empty but got this: %v", foundExchanges)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangesByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Exchange{TestData1, TestData2}
	mockRows := AddExchangeToMockRows(mock, dataList)
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WithArgs(pq.Array(uuids)).WillReturnRows(mockRows)
	foundExchanges, err := GetExchangesByUUIDs(mock, uuids)
	if err != nil {
		t.Fatalf("an error '%s' in GetExchangesByUUIDs", err)
	}
	testExchanges := TestAllData
	for i, foundExchange := range foundExchanges {
		if cmp.Equal(foundExchange, testExchanges[i]) == false {
			t.Errorf("Expected Exchange From Method GetExchangesByUUIDs: %v is different from actual %v", foundExchange, testExchanges[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangesByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WithArgs(pq.Array(uuids)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundExchanges, err := GetExchangesByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetExchangesByUUIDs", err)
	}
	if len(foundExchanges) != 0 {
		t.Errorf("Expected From Method GetExchangesByUUIDs: to be empty but got this: %v", foundExchanges)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangesByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WithArgs(pq.Array(uuids)).WillReturnRows(differentModelRows)
	foundExchanges, err := GetExchangesByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetExchangesByUUIDs", err)
	}
	if foundExchanges != nil {
		t.Errorf("Expected foundExchanges From Method GetExchangesByUUIDs: to be empty but got this: %v", foundExchanges)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffExchanges(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []Exchange{TestData1}
	mockRows := AddExchangeToMockRows(mock, dataList)
	diffInDate := 10
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WithArgs(diffInDate).WillReturnRows(mockRows)
	foundExchanges, err := GetStartAndEndDateDiffExchanges(mock, &diffInDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndDateDiffExchanges", err)
	}
	testExchanges := TestAllData
	for i, foundExchange := range foundExchanges {
		if cmp.Equal(foundExchange, testExchanges[i]) == false {
			t.Errorf("Expected Exchange From Method GetStartAndEndDateDiffExchanges: %v is different from actual %v", foundExchange, testExchanges[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffExchangesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	diffInDate := -10
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WithArgs(diffInDate).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundExchanges, err := GetStartAndEndDateDiffExchanges(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffExchanges", err)
	}
	if len(foundExchanges) != 0 {
		t.Errorf("Expected From Method GetStartAndEndDateDiffExchanges: to be empty but got this: %v", foundExchanges)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndDateDiffExchangesForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	diffInDate := -10
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WithArgs(diffInDate).WillReturnRows(differentModelRows)
	foundExchanges, err := GetStartAndEndDateDiffExchanges(mock, &diffInDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndDateDiffExchanges", err)
	}
	if foundExchanges != nil {
		t.Errorf("Expected foundExchanges From Method GetStartAndEndDateDiffExchanges: to be empty but got this: %v", foundExchanges)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateExchange(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE exchanges").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.ExchangeTypeID, //3
		targetData.Url,            //4
		targetData.StartDate,      //5
		targetData.EndDate,        //6
		targetData.Description,    //7
		targetData.UpdatedBy,      //8
		targetData.ID,             //9

	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateExchange(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateExchange", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateExchangeOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateExchange(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateExchangeOnFailureAtBegin(t *testing.T) {
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
	err = UpdateExchange(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateExchangeOnFailure(t *testing.T) {
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
	mock.ExpectExec("^UPDATE exchanges").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.ExchangeTypeID, //3
		targetData.Url,            //4
		targetData.StartDate,      //5
		targetData.EndDate,        //6
		targetData.Description,    //7
		targetData.UpdatedBy,      //8
		targetData.ID,             //91
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateExchange(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertExchange(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO exchanges").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.ExchangeTypeID, //3
		targetData.Url,            //4
		targetData.StartDate,      //5
		targetData.EndDate,        //6
		targetData.Description,    //7
		targetData.CreatedBy,      //8
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	chainID, err := InsertExchange(mock, &targetData)
	if chainID < 0 {
		t.Fatalf("chainID should not be negative ID: %d", chainID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertExchange", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertExchangeOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertExchange(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertExchangeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO exchanges").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.ExchangeTypeID, //3
		targetData.Url,            //4
		targetData.StartDate,      //5
		targetData.EndDate,        //6
		targetData.Description,    //7
		targetData.CreatedBy,      //8
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	chainID, err := InsertExchange(mock, &targetData)
	if chainID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error chainID: %d", chainID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertExchangeOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO exchanges").WithArgs(
		targetData.Name,           //1
		targetData.AlternateName,  //2
		targetData.ExchangeTypeID, //3
		targetData.Url,            //4
		targetData.StartDate,      //5
		targetData.EndDate,        //6
		targetData.Description,    //7
		targetData.CreatedBy,      //8
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	chainID, err := InsertExchange(mock, &targetData)
	if chainID >= 0 {
		t.Fatalf("Expecting -1 for chainID because of error chainID: %d", chainID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertExchanges(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"exchanges"}, columnsInsertExchanges)
	targetData := TestAllData
	err = InsertExchanges(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertExchanges", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertExchangesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"exchanges"}, columnsInsertExchanges).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertExchanges(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateExchangeChainByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1ExchangeChain
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE exchange_chains").WithArgs(
		targetData.ExchangeID,  //1
		targetData.ChainID,     //2
		targetData.Description, //3
		targetData.UpdatedBy,   //4
		targetData.UUID,        //5
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateExchangeChainByUUID(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateExchangeChainByUUID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateExchangeChainByUUIDOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1ExchangeChain
	targetData.ExchangeID = nil
	err = UpdateExchangeChainByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateExchangeChainByUUIDOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1ExchangeChain
	targetData.ChainID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Error at begin"))
	err = UpdateExchangeChainByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateExchangeChainByUUIDOnFailureAtExec(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1ExchangeChain
	targetData.ChainID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE exchange_chains").WithArgs(
		targetData.ExchangeID,  //1
		targetData.ChainID,     //2
		targetData.Description, //3
		targetData.UpdatedBy,   //4
		targetData.UUID,        //5
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateExchangeChainByUUID(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertExchangeChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1ExchangeChain
	targetData.Description = "New Description"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO exchange_chains").WithArgs(
		targetData.UUID,        //1
		targetData.ExchangeID,  //2
		targetData.ChainID,     //3
		targetData.Description, //4
		targetData.CreatedBy,   //5
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	chainID, err := InsertExchangeChain(mock, &targetData)
	if chainID < 0 {
		t.Fatalf("chainID should not be negative ID: %d", chainID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertExchangeChain", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestInsertExchangeChainOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1ExchangeChain
	targetData.ChainID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Random SQL Error"))
	chainID, err := InsertExchangeChain(mock, &targetData)
	if chainID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error chainID: %d", chainID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestInsertExchangeChainOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1ExchangeChain
	targetData.ChainID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO exchange_chains").WithArgs(
		targetData.UUID,        //1
		targetData.ExchangeID,  //2
		targetData.ChainID,     //3
		targetData.Description, //4
		targetData.CreatedBy,   //5
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	chainID, err := InsertExchangeChain(mock, &targetData)
	if chainID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error chainID: %d", chainID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertExchangeChainOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1ExchangeChain
	targetData.ChainID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO exchange_chains").WithArgs(
		targetData.UUID,        //1
		targetData.ExchangeID,  //2
		targetData.ChainID,     //3
		targetData.Description, //4
		targetData.CreatedBy,   //5
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	chainID, err := InsertExchangeChain(mock, &targetData)
	if chainID >= 0 {
		t.Fatalf("Expecting -1 for chainID because of error chainID: %d", chainID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertExchangeChains(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"exchange_chains"}, columnsInsertExchangeChains)
	targetData := allDataExchangeEchains
	err = InsertExchangeChains(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertExchangeChains", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertExchangeChainsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"exchange_chains"}, columnsInsertExchangeChains).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := allDataExchangeEchains
	err = InsertExchangeChains(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangeListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddExchangeToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"exchange_type_id = 1", "id=1"}
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WillReturnRows(mockRows)
	foundExcchangeList, err := GetExchangeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetExchangeListByPagination", err)
	}
	for i, sourceExchange := range dataList {
		if cmp.Equal(sourceExchange, foundExcchangeList[i]) == false {
			t.Errorf("Expected foundExcchangeList From Method GetExchangeListByPagination: %v is different from actual %v", sourceExchange, foundExcchangeList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangeListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"exchange_type_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundExcchangeList, err := GetExchangeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetExchangeListByPagination", err)
	}
	if len(foundExcchangeList) != 0 {
		t.Errorf("Expected From Method GetExchangeListByPagination: to be empty but got this: %v", foundExcchangeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetExchangeListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"exchange_type_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM exchanges").WillReturnRows(differentModelRows)
	foundExcchangeList, err := GetExchangeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetExchangeListByPagination", err)
	}
	if foundExcchangeList != nil {
		t.Errorf("Expected From Method GetExchangeListByPagination: to be empty but got this: %v", foundExcchangeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalExchangeCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM exchanges").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalExchangeCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalExchangeCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalExchangeCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalExchangeCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM exchanges").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalExchangeCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalExchangeCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalExchangeCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
