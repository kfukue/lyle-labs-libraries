package trade

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
	"id",                       //1
	"parent_trade_id",          //2
	"from_account_id",          //3
	"to_account_id",            //4
	"asset_id",                 //5
	"source_id",                //6
	"uuid",                     //7
	"transaction_id",           //8
	"order_id",                 //9
	"trade_id",                 //10
	"name",                     //11
	"alternate_name",           //12
	"trade_type_id",            //13
	"trade_date",               //14
	"settle_date",              //15
	"transfer_date",            //16
	"from_quantity",            //17
	"to_quantity",              //18
	"price",                    //19
	"total_amount",             //20
	"fees_amount",              //21
	"fees_asset_id",            //22
	"realized_return_amount",   //23
	"realized_return_asset_id", //24
	"cost_basis_amount",        //25
	"cost_basis_trade_id",      //26
	"description",              //27
	"is_active",                //28
	"source_data",              //29
	"created_by",               //30
	"created_at",               //31
	"updated_by",               //32
	"updated_at",               //33
}
var DBColumnsInsertTrades = []string{
	"parent_trade_id",          //1
	"from_account_id",          //2
	"to_account_id",            //3
	"asset_id",                 //4
	"source_id",                //5
	"uuid",                     //6
	"transaction_id",           //7
	"order_id",                 //8
	"trade_id",                 //9
	"name",                     //10
	"alternate_name",           //11
	"trade_type_id",            //12
	"trade_date",               //13
	"settle_date",              //14
	"transfer_date",            //15
	"from_quantity",            //16
	"to_quantity",              //17
	"price",                    //18
	"total_amount",             //19
	"fees_amount",              //20
	"fees_asset_id",            //21
	"realized_return_amount",   //22
	"realized_return_asset_id", //23
	"cost_basis_amount",        //24
	"cost_basis_trade_id",      //25
	"description",              //26
	"is_active",                //27
	"source_data",              //28
	"created_by",               //29
	"created_at",               //30
	"updated_by",               //31
	"updated_at",               //32
}

var TestData1 = Trade{
	ID:                      utils.Ptr[int](1),                                                //1
	ParentTradeID:           nil,                                                              //2
	FromAccountID:           utils.Ptr[int](1),                                                //3
	ToAccountID:             utils.Ptr[int](2),                                                //4
	AssetID:                 utils.Ptr[int](1),                                                //5
	SourceID:                utils.Ptr[int](1),                                                //6
	UUID:                    "01ef85e8-2c26-441e-8c7f-71d79518ad72",                           //7
	TransactionIDFromSource: "7577343b2dd6f990303d41cb8c09b2178217b000bddef3d4abe3a0a3d12342", //8
	OrderIDFromSource:       "7577343b2dd6f990303d41cb8c09b2178217b000bddef3d4a343a0a3d12333", //9
	TradeIDFromSource:       "7577343b2dd6f990303d41cb8c09b2178217b000bddef3d4abe23423d12321", //10
	Name:                    "",                                                               //11
	AlternateName:           "",                                                               //12
	TradeTypeID:             utils.Ptr[int](13),                                               //13
	TradeDate:               utils.Ptr[time.Time](utils.SampleCreatedAtTime),                  //14
	SettleDate:              utils.Ptr[time.Time](utils.SampleCreatedAtTime),                  //15
	TransferDate:            utils.Ptr[time.Time](utils.SampleCreatedAtTime),                  //16
	FromQuantity:            utils.Ptr[float64](0),                                            //17
	ToQuantity:              utils.Ptr[float64](535.2),                                        //18
	Price:                   utils.Ptr[float64](643.36),                                       //19
	TotalAmount:             utils.Ptr[float64](10000),                                        //20
	FeesAmount:              utils.Ptr[float64](12.2),                                         //21
	FeesAssetID:             utils.Ptr[int](2),                                                //22
	RealizedReturnAmount:    utils.Ptr[float64](100),                                          //23
	RealizedReturnAssetID:   utils.Ptr[int](3),                                                //24
	CostBasisAmount:         utils.Ptr[float64](121),                                          //25
	CostBasisTradeID:        utils.Ptr[int](3),                                                //26
	Description:             "",                                                               //27
	IsActive:                true,                                                             //28
	SourceData:              nil,                                                              //29
	CreatedBy:               "SYSTEM",                                                         //30
	CreatedAt:               utils.SampleCreatedAtTime,                                        //31
	UpdatedBy:               "SYSTEM",                                                         //32
	UpdatedAt:               utils.SampleCreatedAtTime,                                        //33
}

var TestData2 = Trade{
	ID:                      utils.Ptr[int](2),                                                //1
	ParentTradeID:           utils.Ptr[int](1),                                                //2
	FromAccountID:           utils.Ptr[int](3),                                                //3
	ToAccountID:             utils.Ptr[int](4),                                                //4
	AssetID:                 utils.Ptr[int](3),                                                //5
	SourceID:                utils.Ptr[int](5),                                                //6
	UUID:                    "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                           //7
	TransactionIDFromSource: "7577343b2dd6f990303d41cb8c09b2178217b000bddef3d4abe3a0a3d34322", //8
	OrderIDFromSource:       "7577343b2dd6f990303d41cb8c09b2178217b000bddef3343223a0a3d12342", //9
	TradeIDFromSource:       "7577343b2dd6f990303d41cb8c09b2178217b000bdabe3455be3a0a3d12342", //10
	Name:                    "",                                                               //11
	AlternateName:           "",                                                               //12
	TradeTypeID:             utils.Ptr[int](2),                                                //13
	TradeDate:               utils.Ptr[time.Time](utils.SampleCreatedAtTime),                  //14
	SettleDate:              utils.Ptr[time.Time](utils.SampleCreatedAtTime),                  //15
	TransferDate:            utils.Ptr[time.Time](utils.SampleCreatedAtTime),                  //16
	FromQuantity:            utils.Ptr[float64](3425.2),                                       //17
	ToQuantity:              utils.Ptr[float64](0),                                            //18
	Price:                   utils.Ptr[float64](11111.2),                                      //19
	TotalAmount:             utils.Ptr[float64](34335.2),                                      //20
	FeesAmount:              utils.Ptr[float64](5.2),                                          //21
	FeesAssetID:             utils.Ptr[int](2),                                                //22
	RealizedReturnAmount:    utils.Ptr[float64](110.2),                                        //23
	RealizedReturnAssetID:   utils.Ptr[int](3),                                                //24
	CostBasisAmount:         utils.Ptr[float64](1100.2),                                       //25
	CostBasisTradeID:        utils.Ptr[int](1),                                                //26
	Description:             "",                                                               //27
	IsActive:                true,                                                             //28
	SourceData:              nil,                                                              //29
	CreatedBy:               "SYSTEM",                                                         //30
	CreatedAt:               utils.SampleCreatedAtTime,                                        //31
	UpdatedBy:               "SYSTEM",                                                         //32
	UpdatedAt:               utils.SampleCreatedAtTime,                                        //33
}
var TestAllData = []Trade{TestData1, TestData2}

func AddTradeToMockRows(mock pgxmock.PgxPoolIface, dataList []Trade) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                      //1
			data.ParentTradeID,           //2
			data.FromAccountID,           //3
			data.ToAccountID,             //4
			data.AssetID,                 //5
			data.SourceID,                //6
			data.UUID,                    //7
			data.TransactionIDFromSource, //8
			data.OrderIDFromSource,       //9
			data.TradeIDFromSource,       //10
			data.Name,                    //11
			data.AlternateName,           //12
			data.TradeTypeID,             //13
			data.TradeDate,               //14
			data.SettleDate,              //15
			data.TransferDate,            //16
			data.FromQuantity,            //17
			data.ToQuantity,              //18
			data.Price,                   //19
			data.TotalAmount,             //20
			data.FeesAmount,              //21
			data.FeesAssetID,             //22
			data.RealizedReturnAmount,    //23
			data.RealizedReturnAssetID,   //24
			data.CostBasisAmount,         //25
			data.CostBasisTradeID,        //26
			data.Description,             //27
			data.IsActive,                //28
			data.SourceData,              //29
			data.CreatedBy,               //30
			data.CreatedAt,               //31
			data.UpdatedBy,               //32
			data.UpdatedAt,               //33
		)
	}
	return rows
}

var DBColumnsMinMaxTradeDates = []string{
	"min_date", //1
	"max_date", //2

}

var TestData1MinMaxTradeDates = MinMaxTradeDates{
	MinTradeDate: &utils.SampleCreatedAtTime,
	MaxTradeDate: &utils.SampleCreatedAtTime,
}

var TestData2MinMaxTradeDates = MinMaxTradeDates{
	MinTradeDate: &utils.SampleCreatedAtTime,
	MaxTradeDate: &utils.SampleCreatedAtTime,
}
var TestAllDataMinMaxTradeDates = []MinMaxTradeDates{TestData1MinMaxTradeDates, TestData2MinMaxTradeDates}

func AddMinMaxTradeDatesToMockRows(mock pgxmock.PgxPoolIface, dataList []MinMaxTradeDates) *pgxmock.Rows {
	rows := mock.NewRows(DBColumnsMinMaxTradeDates)
	for _, data := range dataList {
		rows.AddRow(
			data.MinTradeDate, //1
			data.MaxTradeDate, //2
		)
	}
	return rows
}

func TestGetTrade(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Trade{targetData}
	tradeID := targetData.ID
	mockRows := AddTradeToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(*tradeID).WillReturnRows(mockRows)
	foundTrade, err := GetTrade(mock, tradeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTrade", err)
	}
	if cmp.Equal(*foundTrade, targetData) == false {
		t.Errorf("Expected Trade From Method GetTrade: %v is different from actual %v", foundTrade, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	tradeID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(tradeID).WillReturnRows(noRows)
	foundTrade, err := GetTrade(mock, &tradeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTrade", err)
	}
	if foundTrade != nil {
		t.Errorf("Expected Trade From Method GetTrade: to be empty but got this: %v", foundTrade)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	tradeID := -1
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(tradeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTrade, err := GetTrade(mock, &tradeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTrade", err)
	}
	if foundTrade != nil {
		t.Errorf("Expected Trade From Method GetTrade: to be empty but got this: %v", foundTrade)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	tradeID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(tradeID).WillReturnRows(differentModelRows)
	foundTrade, err := GetTrade(mock, &tradeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTrade", err)
	}
	if foundTrade != nil {
		t.Errorf("Expected foundTrade From Method GetTrade: to be empty but got this: %v", foundTrade)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByTransactionIDFromSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Trade{targetData}
	transactionIDFromSource := targetData.TransactionIDFromSource
	fromAccountID := targetData.FromAccountID
	toAccountID := targetData.ToAccountID

	mockRows := AddTradeToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(transactionIDFromSource, *fromAccountID, *toAccountID).WillReturnRows(mockRows)
	foundTrade, err := GetTradeByTransactionIDFromSource(mock, transactionIDFromSource, fromAccountID, toAccountID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTradeByTransactionIDFromSource", err)
	}
	if cmp.Equal(*foundTrade, targetData) == false {
		t.Errorf("Expected Trade From Method GetTradeByTransactionIDFromSource: %v is different from actual %v", foundTrade, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByTransactionIDFromSourceForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionIDFromSource := "test"
	fromAccountID := 999
	toAccountID := 777
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(transactionIDFromSource, fromAccountID, toAccountID).WillReturnRows(noRows)
	foundTrade, err := GetTradeByTransactionIDFromSource(mock, transactionIDFromSource, &fromAccountID, &toAccountID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTradeByTransactionIDFromSource", err)
	}
	if foundTrade != nil {
		t.Errorf("Expected Trade From Method GetTradeByTransactionIDFromSource: to be empty but got this: %v", foundTrade)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByTransactionIDFromSourceForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionIDFromSource := "test"
	fromAccountID := -1
	toAccountID := -1
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(transactionIDFromSource, fromAccountID, toAccountID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTrade, err := GetTradeByTransactionIDFromSource(mock, transactionIDFromSource, &fromAccountID, &toAccountID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeByTransactionIDFromSource", err)
	}
	if foundTrade != nil {
		t.Errorf("Expected Trade From Method GetTradeByTransactionIDFromSource: to be empty but got this: %v", foundTrade)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByTransactionIDFromSourceForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	transactionIDFromSource := "test"
	fromAccountID := -1
	toAccountID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(transactionIDFromSource, fromAccountID, toAccountID).WillReturnRows(differentModelRows)
	foundTrade, err := GetTradeByTransactionIDFromSource(mock, transactionIDFromSource, &fromAccountID, &toAccountID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeByTransactionIDFromSource", err)
	}
	if foundTrade != nil {
		t.Errorf("Expected foundTrade From Method GetTradeByTransactionIDFromSource: to be empty but got this: %v", foundTrade)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByStartAndEndDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTradeToMockRows(mock, dataList)
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnRows(mockRows)
	foundTradeList, err := GetTradeByStartAndEndDates(mock, startDate, endDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetTradeByStartAndEndDates", err)
	}
	testMarketDataList := TestAllData
	for i, foundTrade := range foundTradeList {
		if cmp.Equal(foundTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected Trade From Method GetTradeByStartAndEndDates: %v is different from actual %v", foundTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByStartAndEndDatesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTradeList, err := GetTradeByStartAndEndDates(mock, startDate, endDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeByStartAndEndDates", err)
	}
	if len(foundTradeList) != 0 {
		t.Errorf("Expected From Method GetTradeByStartAndEndDates: to be empty but got this: %v", foundTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByStartAndEndDatesForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnRows(differentModelRows)
	foundTradeList, err := GetTradeByStartAndEndDates(mock, startDate, endDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeByStartAndEndDates", err)
	}
	if foundTradeList != nil {
		t.Errorf("Expected foundTradeList From Method GetTradeByStartAndEndDates: to be empty but got this: %v", foundTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByFromAccountStartAndEndDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTradeToMockRows(mock, dataList)
	fromAccountID := 1
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(fromAccountID, startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnRows(mockRows)
	foundTradeList, err := GetTradeByFromAccountStartAndEndDates(mock, &fromAccountID, startDate, endDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetTradeByFromAccountStartAndEndDates", err)
	}
	testMarketDataList := TestAllData
	for i, foundTrade := range foundTradeList {
		if cmp.Equal(foundTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected Trade From Method GetTradeByFromAccountStartAndEndDates: %v is different from actual %v", foundTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByFromAccountStartAndEndDatesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	fromAccountID := -999
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(fromAccountID, startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTradeList, err := GetTradeByFromAccountStartAndEndDates(mock, &fromAccountID, startDate, endDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeByFromAccountStartAndEndDates", err)
	}
	if len(foundTradeList) != 0 {
		t.Errorf("Expected From Method GetTradeByFromAccountStartAndEndDates: to be empty but got this: %v", foundTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByFromAccountStartAndEndDatesForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	fromAccountID := -1
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(fromAccountID, startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnRows(differentModelRows)
	foundTradeList, err := GetTradeByFromAccountStartAndEndDates(mock, &fromAccountID, startDate, endDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeByFromAccountStartAndEndDates", err)
	}
	if foundTradeList != nil {
		t.Errorf("Expected foundTradeList From Method GetTradeByFromAccountStartAndEndDates: to be empty but got this: %v", foundTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByToAccountStartAndEndDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTradeToMockRows(mock, dataList)
	toAccountID := 1
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(toAccountID, startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnRows(mockRows)
	foundTradeList, err := GetTradeByToAccountStartAndEndDates(mock, &toAccountID, startDate, endDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetTradeByToAccountStartAndEndDates", err)
	}
	testMarketDataList := TestAllData
	for i, foundTrade := range foundTradeList {
		if cmp.Equal(foundTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected Trade From Method GetTradeByToAccountStartAndEndDates: %v is different from actual %v", foundTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByToAccountStartAndEndDatesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	toAccountID := -999
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(toAccountID, startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTradeList, err := GetTradeByToAccountStartAndEndDates(mock, &toAccountID, startDate, endDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeByToAccountStartAndEndDates", err)
	}
	if len(foundTradeList) != 0 {
		t.Errorf("Expected From Method GetTradeByToAccountStartAndEndDates: to be empty but got this: %v", foundTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeByToAccountStartAndEndDatesForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	toAccountID := -1
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WithArgs(toAccountID, startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnRows(differentModelRows)
	foundTradeList, err := GetTradeByToAccountStartAndEndDates(mock, &toAccountID, startDate, endDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeByToAccountStartAndEndDates", err)
	}
	if foundTradeList != nil {
		t.Errorf("Expected foundTradeList From Method GetTradeByToAccountStartAndEndDates: to be empty but got this: %v", foundTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTrade(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	tradeID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM trades").WithArgs(*tradeID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveTrade(mock, tradeID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveTrade", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTradeOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	tradeID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveTrade(mock, &tradeID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTradeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	tradeID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM trades").WithArgs(tradeID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveTrade(mock, &tradeID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTradeToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WillReturnRows(mockRows)
	foundTradeList, err := GetTradeList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTradeList", err)
	}
	testMarketDataList := TestAllData
	for i, foundTrade := range foundTradeList {
		if cmp.Equal(foundTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected Trade From Method GetTradeList: %v is different from actual %v", foundTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM trades").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTradeList, err := GetTradeList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeList", err)
	}
	if len(foundTradeList) != 0 {
		t.Errorf("Expected From Method GetTradeList: to be empty but got this: %v", foundTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WillReturnRows(differentModelRows)
	foundTradeList, err := GetTradeList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeList", err)
	}
	if foundTradeList != nil {
		t.Errorf("Expected foundTradeList From Method GetTradeList: to be empty but got this: %v", foundTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTrade(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE trades").WithArgs(
		targetData.ParentTradeID,           //1
		targetData.FromAccountID,           //2
		targetData.ToAccountID,             //3
		targetData.AssetID,                 //4
		targetData.SourceID,                //5
		targetData.TransactionIDFromSource, //6
		targetData.OrderIDFromSource,       //7
		targetData.TradeIDFromSource,       //8
		targetData.Name,                    //9
		targetData.AlternateName,           //10
		targetData.TradeTypeID,             //11
		targetData.TradeDate,               //12
		targetData.SettleDate,              //13
		targetData.TransferDate,            //14
		targetData.FromQuantity,            //15
		targetData.ToQuantity,              //16
		targetData.Price,                   //17
		targetData.TotalAmount,             //18
		targetData.FeesAmount,              //19
		targetData.FeesAssetID,             //20
		targetData.RealizedReturnAmount,    //21
		targetData.RealizedReturnAssetID,   //22
		targetData.CostBasisAmount,         //23
		targetData.CostBasisTradeID,        //24
		targetData.Description,             //25
		targetData.IsActive,                //26
		targetData.SourceData,              //27
		targetData.UpdatedBy,               //28
		targetData.ID,                      //29
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateTrade(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateTrade", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTradeOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateTrade(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTradeOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateTrade(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTradeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE trades").WithArgs(
		targetData.ParentTradeID,           //1
		targetData.FromAccountID,           //2
		targetData.ToAccountID,             //3
		targetData.AssetID,                 //4
		targetData.SourceID,                //5
		targetData.TransactionIDFromSource, //6
		targetData.OrderIDFromSource,       //7
		targetData.TradeIDFromSource,       //8
		targetData.Name,                    //9
		targetData.AlternateName,           //10
		targetData.TradeTypeID,             //11
		targetData.TradeDate,               //12
		targetData.SettleDate,              //13
		targetData.TransferDate,            //14
		targetData.FromQuantity,            //15
		targetData.ToQuantity,              //16
		targetData.Price,                   //17
		targetData.TotalAmount,             //18
		targetData.FeesAmount,              //19
		targetData.FeesAssetID,             //20
		targetData.RealizedReturnAmount,    //21
		targetData.RealizedReturnAssetID,   //22
		targetData.CostBasisAmount,         //23
		targetData.CostBasisTradeID,        //24
		targetData.Description,             //25
		targetData.IsActive,                //26
		targetData.SourceData,              //27
		targetData.UpdatedBy,               //28
		targetData.ID,                      //29
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateTrade(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTrade(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO trades").WithArgs(
		targetData.ParentTradeID,           //1
		targetData.FromAccountID,           //2
		targetData.ToAccountID,             //3
		targetData.AssetID,                 //4
		targetData.SourceID,                //5
		targetData.TransactionIDFromSource, //6
		targetData.OrderIDFromSource,       //7
		targetData.TradeIDFromSource,       //8
		targetData.Name,                    //9
		targetData.AlternateName,           //10
		targetData.TradeTypeID,             //11
		targetData.TradeDate,               //12
		targetData.SettleDate,              //13
		targetData.TransferDate,            //14
		targetData.FromQuantity,            //15
		targetData.ToQuantity,              //16
		targetData.Price,                   //17
		targetData.TotalAmount,             //18
		targetData.FeesAmount,              //19
		targetData.FeesAssetID,             //20
		targetData.RealizedReturnAmount,    //21
		targetData.RealizedReturnAssetID,   //22
		targetData.CostBasisAmount,         //23
		targetData.CostBasisTradeID,        //24
		targetData.Description,             //25
		targetData.IsActive,                //26
		targetData.SourceData,              //27
		targetData.CreatedBy,               //28
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	tradeID, err := InsertTrade(mock, &targetData)
	if tradeID < 0 {
		t.Fatalf("tradeID should not be negative ID: %d", tradeID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertTrade", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTradesOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertTrade(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTradeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO trades").WithArgs(
		targetData.ParentTradeID,           //1
		targetData.FromAccountID,           //2
		targetData.ToAccountID,             //3
		targetData.AssetID,                 //4
		targetData.SourceID,                //5
		targetData.TransactionIDFromSource, //6
		targetData.OrderIDFromSource,       //7
		targetData.TradeIDFromSource,       //8
		targetData.Name,                    //9
		targetData.AlternateName,           //10
		targetData.TradeTypeID,             //11
		targetData.TradeDate,               //12
		targetData.SettleDate,              //13
		targetData.TransferDate,            //14
		targetData.FromQuantity,            //15
		targetData.ToQuantity,              //16
		targetData.Price,                   //17
		targetData.TotalAmount,             //18
		targetData.FeesAmount,              //19
		targetData.FeesAssetID,             //20
		targetData.RealizedReturnAmount,    //21
		targetData.RealizedReturnAssetID,   //22
		targetData.CostBasisAmount,         //23
		targetData.CostBasisTradeID,        //24
		targetData.Description,             //25
		targetData.IsActive,                //26
		targetData.SourceData,              //27
		targetData.CreatedBy,               //28
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	tradeID, err := InsertTrade(mock, &targetData)
	if tradeID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error tradeID: %d", tradeID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTradeOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO trades").WithArgs(
		targetData.ParentTradeID,           //1
		targetData.FromAccountID,           //2
		targetData.ToAccountID,             //3
		targetData.AssetID,                 //4
		targetData.SourceID,                //5
		targetData.TransactionIDFromSource, //6
		targetData.OrderIDFromSource,       //7
		targetData.TradeIDFromSource,       //8
		targetData.Name,                    //9
		targetData.AlternateName,           //10
		targetData.TradeTypeID,             //11
		targetData.TradeDate,               //12
		targetData.SettleDate,              //13
		targetData.TransferDate,            //14
		targetData.FromQuantity,            //15
		targetData.ToQuantity,              //16
		targetData.Price,                   //17
		targetData.TotalAmount,             //18
		targetData.FeesAmount,              //19
		targetData.FeesAssetID,             //20
		targetData.RealizedReturnAmount,    //21
		targetData.RealizedReturnAssetID,   //22
		targetData.CostBasisAmount,         //23
		targetData.CostBasisTradeID,        //24
		targetData.Description,             //25
		targetData.IsActive,                //26
		targetData.SourceData,              //27
		targetData.CreatedBy,               //28
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	tradeID, err := InsertTrade(mock, &targetData)
	if tradeID >= 0 {
		t.Fatalf("Expecting -1 for tradeID because of error tradeID: %d", tradeID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTrades(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"trades"}, DBColumnsInsertTrades)
	targetData := TestAllData
	err = InsertTrades(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertTrades", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTradesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"trades"}, DBColumnsInsertTrades).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertTrades(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinAndMaxTradeDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1MinMaxTradeDates
	dataList := []MinMaxTradeDates{targetData}
	mockRows := AddMinMaxTradeDatesToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WillReturnRows(mockRows)
	foundDates, err := GetMinAndMaxTradeDates(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetMinAndMaxTradeDates", err)
	}
	if cmp.Equal(*foundDates, targetData) == false {
		t.Errorf("Expected foundDates From Method GetMinAndMaxTradeDates: %v is different from actual %v", foundDates, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinAndMaxTradeDatesForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	noRows := pgxmock.NewRows(DBColumnsMinMaxTradeDates)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WillReturnRows(noRows)
	foundDates, err := GetMinAndMaxTradeDates(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetMinAndMaxTradeDates", err)
	}
	if foundDates != nil {
		t.Errorf("Expected foundDates From Method GetMinAndMaxTradeDates: to be empty but got this: %v", foundDates)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetMinAndMaxTradeDatesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM trades").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundDates, err := GetMinAndMaxTradeDates(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTrade", err)
	}
	if foundDates != nil {
		t.Errorf("Expected foundDates From Method GetMinAndMaxTradeDates: to be empty but got this: %v", foundDates)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTradeToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"from_account_id = 1", "asset_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM trades").WillReturnRows(mockRows)
	foundTradeList, err := GetTradeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetTradeListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundTradeList[i]) == false {
			t.Errorf("Expected sourceData From Method GetTradeListByPagination: %v is different from actual %v", sourceData, foundTradeList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"from_account_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM trades").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTradeList, err := GetTradeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeListByPagination", err)
	}
	if len(foundTradeList) != 0 {
		t.Errorf("Expected From Method GetTradeListByPagination: to be empty but got this: %v", foundTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTradeListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"from_account_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM trades").WillReturnRows(differentModelRows)
	foundTaxList, err := GetTradeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTradeListByPagination", err)
	}
	if len(foundTaxList) != 0 {
		t.Errorf("Expected From Method GetTradeListByPagination: to be empty but got this: %v", foundTaxList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalTradesCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM trades").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalTradesCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalTradesCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalTradesCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM trades").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalTradesCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalTradesCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalTradesCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
