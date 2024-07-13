package gethlyletrades

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/asset"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/shopspring/decimal"
)

var DBColumns = []string{
	"id",                        //1
	"uuid",                      //2
	"name",                      //3
	"alternate_name",            //4
	"address_str",               //5
	"address_id",                //6
	"trade_date",                //7
	"txn_hash",                  //8
	"token0_amount",             //9
	"token0_amount_decimal_adj", //10
	"token1_amount",             //11
	"token1_amount_decimal_adj", //12
	"is_buy",                    //13
	"price",                     //14
	"price_usd",                 //15
	"lp_token1_price_usd",       //16
	"total_amount_usd",          //17
	"token0_asset_id",           //18
	"token1_asset_id",           //19
	"geth_process_job_id",       //20
	"status_id",                 //21
	"trade_type_id",             //22
	"description",               //23
	"created_by",                //24
	"created_at",                //25
	"updated_by",                //26
	"updated_at",                //27
	"base_asset_id",             //28
	"oracle_price_usd",          //29
	"oracle_price_asset_id",     //30
}
var DBColumnsInsertGethTrades = []string{
	"uuid",                      //1
	"name",                      //2
	"alternate_name",            //3
	"address_str",               //4
	"address_id",                //5
	"trade_date",                //6
	"txn_hash",                  //7
	"token0_amount",             //8
	"token0_amount_decimal_adj", //9
	"token1_amount",             //10
	"token1_amount_decimal_adj", //11
	"is_buy",                    //12
	"price",                     //13
	"price_usd",                 //14
	"lp_token1_price_usd",       //15
	"total_amount_usd",          //16
	"token0_asset_id",           //17
	"token1_asset_id",           //18
	"geth_process_job_id",       //19
	"status_id",                 //20
	"trade_type_id",             //21
	"description",               //22
	"created_by",                //23
	"created_at",                //24
	"updated_by",                //25
	"updated_at",                //26
	"base_asset_id",             //27
	"oracle_price_usd",          //28
	"oracle_price_asset_id",     //29
}

var TestData1 = GethTrade{
	ID:                     utils.Ptr[int](1),                                                                   //1
	UUID:                   "01ef85e8-2c26-441e-8c7f-71d79518ad72",                                              //2
	Name:                   "PEPE/WETH",                                                                         //3
	AlternateName:          "PEPE/WETH",                                                                         //4
	AddressStr:             "0xd2203a02d4b1D070e9F194A1A88956209e7791B7",                                        //5
	AddressID:              utils.Ptr[int](798584),                                                              //6
	TradeDate:              utils.SampleCreatedAtTime,                                                           //7
	TxnHash:                "0xf5f20f10458168136a02a06534969c232da05e5cbe7b562fe807e74c0ae8c670",                //8
	Token0Amount:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(6365181906890837627808795)),         //9
	Token0AmountDecimalAdj: utils.Ptr[decimal.Decimal](decimal.NewFromFloat(6365181.9068908376278088)),          //10
	Token1Amount:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(-20000000000000000)),                //11
	Token1AmountDecimalAdj: utils.Ptr[decimal.Decimal](decimal.NewFromFloat(-0.02)),                             //12
	IsBuy:                  utils.Ptr[bool](true),                                                               //12                        //13
	Price:                  utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.000000003142094)),                 //13
	PriceUSD:               utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.00000928857452)),                  //14
	LPToken1PriceUSD:       utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.00000306905761197174)),            //16
	TotalAmountUSD:         utils.Ptr[decimal.Decimal](decimal.NewFromFloat(59.123466475511246811122063111776)), //17
	Token0AssetId:          utils.Ptr[int](535),                                                                 //18
	Token1AssetId:          utils.Ptr[int](35),                                                                  //19
	GethProcessJobID:       utils.Ptr[int](4588207),                                                             //20
	StatusID:               utils.Ptr[int](53),                                                                  //21
	TradeTypeID:            utils.Ptr[int](2),                                                                   //22
	Description:            "Imported by Geth Dex Analyzer",                                                     //23
	CreatedBy:              "SYSTEM",                                                                            //24
	CreatedAt:              utils.SampleCreatedAtTime,                                                           //25
	UpdatedBy:              "SYSTEM",                                                                            //26
	UpdatedAt:              utils.SampleCreatedAtTime,                                                           //27
	BaseAssetID:            utils.Ptr[int](1),                                                                   //28
	OraclePriceUSD:         utils.Ptr[decimal.Decimal](decimal.NewFromFloat(2986.03364013)),                     //29
	OraclePriceAssetID:     utils.Ptr[int](530),                                                                 //30
}

var TestData2 = GethTrade{
	ID:                     utils.Ptr[int](2),
	UUID:                   "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",
	Name:                   "Moon Tropica/USDT",                                                       //3
	AlternateName:          "Moon Tropica/USDT",                                                       //4
	AddressStr:             "0x859bFc051c93dDD08163C1AAe645269F142c1841",                              //5
	AddressID:              utils.Ptr[int](524975),                                                    //6
	TradeDate:              utils.SampleCreatedAtTime,                                                 //7
	TxnHash:                "0x8dad48e40a54b154d524e6b649787bcba5d1f57c3796a666787803acc1b28a6a",      //8
	Token0Amount:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(-100741838000000000000)),  //9
	Token0AmountDecimalAdj: utils.Ptr[decimal.Decimal](decimal.NewFromFloat(-100.741838)),             //10
	Token1Amount:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(124791127470024425)),      //11
	Token1AmountDecimalAdj: utils.Ptr[decimal.Decimal](decimal.NewFromFloat(124791127.470024425)),     //12
	IsBuy:                  utils.Ptr[bool](false),                                                    //12                        //13
	Price:                  utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.0012387219644536)),      //13
	PriceUSD:               utils.Ptr[decimal.Decimal](decimal.NewFromFloat(8.281857094022)),          //14
	LPToken1PriceUSD:       utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.00000306905761197174)),  //16
	TotalAmountUSD:         utils.Ptr[decimal.Decimal](decimal.NewFromFloat(-834.329505705115092436)), //17
	Token0AssetId:          utils.Ptr[int](546),                                                       //18
	Token1AssetId:          utils.Ptr[int](23896),                                                     //19
	GethProcessJobID:       utils.Ptr[int](4617997),                                                   //20
	StatusID:               utils.Ptr[int](53),                                                        //21
	TradeTypeID:            utils.Ptr[int](2),                                                         //22
	Description:            "Imported by Geth Dex Analyzer",                                           //23
	CreatedBy:              "SYSTEM",                                                                  //24
	CreatedAt:              utils.SampleCreatedAtTime,                                                 //25
	UpdatedBy:              "SYSTEM",                                                                  //26
	UpdatedAt:              utils.SampleCreatedAtTime,                                                 //27
	BaseAssetID:            utils.Ptr[int](1),                                                         //28
	OraclePriceUSD:         utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.99857206)),              //29
	OraclePriceAssetID:     utils.Ptr[int](546),                                                       //30
}
var TestAllData = []GethTrade{TestData1, TestData2}

func AddGethTradeToMockRows(mock pgxmock.PgxPoolIface, dataList []GethTrade) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                     //1
			data.UUID,                   //2
			data.Name,                   //3
			data.AlternateName,          //4
			data.AddressStr,             //5
			data.AddressID,              //6
			data.TradeDate,              //7
			data.TxnHash,                //8
			data.Token0Amount,           //9
			data.Token0AmountDecimalAdj, //10
			data.Token1Amount,           //11
			data.Token1AmountDecimalAdj, //12
			data.IsBuy,                  //13
			data.Price,                  //14
			data.PriceUSD,               //15
			data.LPToken1PriceUSD,       //16
			data.TotalAmountUSD,         //17
			data.Token0AssetId,          //18
			data.Token1AssetId,          //19
			data.GethProcessJobID,       //20
			data.StatusID,               //21
			data.TradeTypeID,            //22
			data.Description,            //23
			data.CreatedBy,              //24
			data.CreatedAt,              //25
			data.UpdatedBy,              //26
			data.UpdatedAt,              //27
			data.BaseAssetID,            //28
			data.OraclePriceUSD,         //29
			data.OraclePriceAssetID,     //30
		)
	}
	return rows
}

var DBColumnsNetTransferByAddress = []string{
	"txn_hash",              //1
	"address_str",           //2
	"asset_id",              //3
	"net_amount",            //4
	"id",                    //5
	"uuid",                  //6
	"name",                  //7
	"alternate_name",        //8
	"cusip",                 //9
	"ticker",                //10
	"base_asset_id",         //11
	"quote_asset_id",        //12
	"description",           //13
	"asset_type_id",         //14
	"created_by",            //15
	"created_at",            //16
	"updated_by",            //17
	"updated_at",            //18
	"chain_id",              //19
	"category_id",           //20
	"sub_category_id",       //21
	"is_default_quote",      //22
	"ignore_market_data",    //23
	"decimals",              //24
	"contract_address",      //25
	"starting_block_number", //26
	"import_geth",           //27
	"import_geth_initial",   //28
}

var TestData1NetTransferByAddress = NetTransferByAddress{
	TxnHash:    "0xf5f20f10458168136a02a06534969c232da05e5cbe7b562fe807e74c0ae8c670",
	AddressStr: "0xd2203a02d4b1D070e9F194A1A88956209e7791B7",
	AssetID:    utils.Ptr[int](1),
	NetAmount:  utils.Ptr[decimal.Decimal](decimal.NewFromFloat(59.123466475511246811122063111776)),
	Asset:      asset.TestData1,
}

var TestData2NetTransferByAddress = NetTransferByAddress{
	TxnHash:    "0x8dad48e40a54b154d524e6b649787bcba5d1f57c3796a666787803acc1b28a6a",
	AddressStr: "0x859bFc051c93dDD08163C1AAe645269F142c1841",
	NetAmount:  utils.Ptr[decimal.Decimal](decimal.NewFromFloat(-834.329505705115092436)),
	Asset:      asset.TestData2,
}
var TestAllDataNetTransferByAddress = []NetTransferByAddress{TestData1NetTransferByAddress, TestData2NetTransferByAddress}

func AddNetTransferByAddressToMockRows(mock pgxmock.PgxPoolIface, dataList []NetTransferByAddress) *pgxmock.Rows {
	rows := mock.NewRows(DBColumnsNetTransferByAddress)
	for _, data := range dataList {
		rows.AddRow(
			data.TxnHash,             //1
			data.AddressStr,          //2
			data.AssetID,             //3
			data.NetAmount,           //4
			data.ID,                  //5
			data.UUID,                //6
			data.Name,                //7
			data.AlternateName,       //8
			data.Cusip,               //9
			data.Ticker,              //10
			data.BaseAssetID,         //11
			data.QuoteAssetID,        //12
			data.Description,         //13
			data.AssetTypeID,         //14
			data.CreatedBy,           //15
			data.CreatedAt,           //16
			data.UpdatedBy,           //17
			data.UpdatedAt,           //18
			data.ChainID,             //19
			data.CategoryID,          //20
			data.SubCategoryID,       //21
			data.IsDefaultQuote,      //22
			data.IgnoreMarketData,    //23
			data.Decimals,            //24
			data.ContractAddress,     //25
			data.StartingBlockNumber, //26
			data.ImportGeth,          //27
			data.ImportGethInitial,   //28
		)
	}
	return rows
}

func TestGetGethTrade(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethTrade{targetData}
	gethTradeID := targetData.ID
	mockRows := AddGethTradeToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(*gethTradeID).WillReturnRows(mockRows)
	foundGethTrade, err := GetGethTrade(mock, gethTradeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTrade", err)
	}
	if cmp.Equal(*foundGethTrade, targetData) == false {
		t.Errorf("Expected GethTrade From Method GetGethTrade: %v is different from actual %v", foundGethTrade, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(gethTradeID).WillReturnRows(noRows)
	foundGethTrade, err := GetGethTrade(mock, &gethTradeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTrade", err)
	}
	if foundGethTrade != nil {
		t.Errorf("Expected GethTrade From Method GetGethTrade: to be empty but got this: %v", foundGethTrade)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(gethTradeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTrade, err := GetGethTrade(mock, &gethTradeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTrade", err)
	}
	if foundGethTrade != nil {
		t.Errorf("Expected GethTrade From Method GetGethTrade: to be empty but got this: %v", foundGethTrade)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeByStartAndEndDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTradeToMockRows(mock, dataList)
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnRows(mockRows)
	foundGethTradeList, err := GetGethTradeByStartAndEndDates(mock, startDate, endDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeByStartAndEndDates", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTrade := range foundGethTradeList {
		if cmp.Equal(foundGethTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTrade From Method GetGethTradeByStartAndEndDates: %v is different from actual %v", foundGethTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeByStartAndEndDatesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startDate := utils.SampleCreatedAtTime
	endDate := startDate.Add(time.Hour * 24)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeList, err := GetGethTradeByStartAndEndDates(mock, startDate, endDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeByStartAndEndDates", err)
	}
	if len(foundGethTradeList) != 0 {
		t.Errorf("Expected From Method GetGethTradeByStartAndEndDates: to be empty but got this: %v", foundGethTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeByFromAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTrade{TestData1}
	mockRows := AddGethTradeToMockRows(mock, dataList)
	addressStr := TestData1.AddressStr
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(addressStr).WillReturnRows(mockRows)
	foundGethTradeList, err := GetGethTradeByFromAddress(mock, addressStr)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeByFromAddress", err)
	}
	testMarketDataList := dataList
	for i, foundGethTrade := range foundGethTradeList {
		if cmp.Equal(foundGethTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTrade From Method GetGethTradeByFromAddress: %v is different from actual %v", foundGethTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeByFromAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	addressStr := TestData1.AddressStr
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(addressStr).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeList, err := GetGethTradeByFromAddress(mock, addressStr)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeByFromAddress", err)
	}
	if len(foundGethTradeList) != 0 {
		t.Errorf("Expected From Method GetGethTradeByFromAddress: to be empty but got this: %v", foundGethTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeByFromAddressId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTrade{TestData1}
	mockRows := AddGethTradeToMockRows(mock, dataList)
	addressID := TestData1.AddressID
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(*addressID).WillReturnRows(mockRows)
	foundGethTradeList, err := GetGethTradeByFromAddressId(mock, addressID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeByFromAddressId", err)
	}
	testMarketDataList := dataList
	for i, foundGethTrade := range foundGethTradeList {
		if cmp.Equal(foundGethTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTrade From Method GetGethTradeByFromAddressId: %v is different from actual %v", foundGethTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeByFromAddressIdForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	addressID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(addressID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeList, err := GetGethTradeByFromAddressId(mock, &addressID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeByFromAddressId", err)
	}
	if len(foundGethTradeList) != 0 {
		t.Errorf("Expected From Method GetGethTradeByFromAddressId: to be empty but got this: %v", foundGethTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTrade{TestData1, TestData2}
	mockRows := AddGethTradeToMockRows(mock, dataList)
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(pq.Array(uuids)).WillReturnRows(mockRows)
	foundGethTradeList, err := GetGethTradeByUUIDs(mock, uuids)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeByUUIDs", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTrade := range foundGethTradeList {
		if cmp.Equal(foundGethTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTrade From Method GetGethTradeByUUIDs: %v is different from actual %v", foundGethTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()

	uuids := []string{"uuid-invalid-1", "uuid-invalid-2"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(pq.Array(uuids)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeList, err := GetGethTradeByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeByUUIDs", err)
	}
	if len(foundGethTradeList) != 0 {
		t.Errorf("Expected From Method GetGethTradeByUUIDs: to be empty but got this: %v", foundGethTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetNetTransfersByTxnHashAndAddressStrs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []NetTransferByAddress{TestData1NetTransferByAddress, TestData2NetTransferByAddress}
	mockRows := AddNetTransferByAddressToMockRows(mock, dataList)
	txnHash := TestData1.TxnHash
	addressStr := TestData1.AddressStr
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectQuery("^WITH to_address as").WithArgs(txnHash, addressStr, *baseAssetID).WillReturnRows(mockRows)
	foundGethTradeList, err := GetNetTransfersByTxnHashAndAddressStrs(mock, txnHash, addressStr, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetNetTransfersByTxnHashAndAddressStrs", err)
	}
	testMarketDataList := dataList
	for i, foundGethTrade := range foundGethTradeList {
		if cmp.Equal(foundGethTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTrade From Method GetNetTransfersByTxnHashAndAddressStrs: %v is different from actual %v", foundGethTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetNetTransfersByTxnHashAndAddressStrsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	txnHash := TestData1.TxnHash
	addressStr := TestData1.AddressStr
	baseAssetID := -1
	mock.ExpectQuery("^WITH to_address as").WithArgs(txnHash, addressStr, baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeList, err := GetNetTransfersByTxnHashAndAddressStrs(mock, txnHash, addressStr, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetNetTransfersByTxnHashAndAddressStrs", err)
	}
	if len(foundGethTradeList) != 0 {
		t.Errorf("Expected From Method GetNetTransfersByTxnHashAndAddressStrs: to be empty but got this: %v", foundGethTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetFromNetTransfersByTxnHashesAndAddressStrs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []NetTransferByAddress{TestData1NetTransferByAddress, TestData2NetTransferByAddress}
	mockRows := AddNetTransferByAddressToMockRows(mock, dataList)
	txnHashes := []string{TestData1.TxnHash, TestData2.TxnHash}
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectQuery("^WITH to_address as").WithArgs(pq.Array(txnHashes), *baseAssetID).WillReturnRows(mockRows)
	foundGethTradeList, err := GetFromNetTransfersByTxnHashesAndAddressStrs(mock, txnHashes, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetFromNetTransfersByTxnHashesAndAddressStrs", err)
	}
	testMarketDataList := dataList
	for i, foundGethTrade := range foundGethTradeList {
		if cmp.Equal(foundGethTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTrade From Method GetFromNetTransfersByTxnHashesAndAddressStrs: %v is different from actual %v", foundGethTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetFromNetTransfersByTxnHashesAndAddressStrsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	txnHashes := []string{TestData1.TxnHash, TestData2.TxnHash}
	baseAssetID := -1
	mock.ExpectQuery("^WITH to_address as").WithArgs(pq.Array(txnHashes), baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeList, err := GetFromNetTransfersByTxnHashesAndAddressStrs(mock, txnHashes, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetFromNetTransfersByTxnHashesAndAddressStrs", err)
	}
	if len(foundGethTradeList) != 0 {
		t.Errorf("Expected From Method GetFromNetTransfersByTxnHashesAndAddressStrs: to be empty but got this: %v", foundGethTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndBlockForNewTradesByBaseAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startBlockNumber := uint64(10000)
	endBlockNumber := uint64(20000)
	baseAssetID := TestData1.BaseAssetID
	mockRows := mock.NewRows([]string{"start_block_number", "end_block_number"}).AddRow(startBlockNumber, endBlockNumber)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WithArgs(*baseAssetID).WillReturnRows(mockRows)
	startBlockNumberResult, endBlockNumberResult, err := GetStartAndEndBlockForNewTradesByBaseAssetID(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndBlockForNewTradesByBaseAssetID", err)
	}
	if cmp.Equal(*startBlockNumberResult, startBlockNumber) == false {
		t.Errorf("Expected startBlockNumber From Method GetStartAndEndBlockForNewTradesByBaseAssetID: %v is different from actual %v", startBlockNumberResult, startBlockNumber)
	}
	if cmp.Equal(*endBlockNumberResult, endBlockNumber) == false {
		t.Errorf("Expected endBlockNumber From Method GetStartAndEndBlockForNewTradesByBaseAssetID: %v is different from actual %v", endBlockNumberResult, endBlockNumber)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndBlockForNewTradesByBaseAssetIDForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := 999
	noRows := mock.NewRows([]string{"start_block_number", "end_block_number"})
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WithArgs(baseAssetID).WillReturnRows(noRows)
	startBlockNumberResult, endBlockNumberResult, err := GetStartAndEndBlockForNewTradesByBaseAssetID(mock, &baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetStartAndEndBlockForNewTradesByBaseAssetID", err)
	}
	if startBlockNumberResult != nil {
		t.Errorf("Expected startBlockNumberResult From Method GetStartAndEndBlockForNewTradesByBaseAssetID: to be empty but got this: %v", startBlockNumberResult)
	}
	if endBlockNumberResult != nil {
		t.Errorf("Expected endBlockNumberResult From Method GetStartAndEndBlockForNewTradesByBaseAssetID: to be empty but got this: %v", endBlockNumberResult)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetStartAndEndBlockForNewTradesByBaseAssetIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_trade_swaps").WithArgs(baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	startBlockNumberResult, endBlockNumberResult, err := GetStartAndEndBlockForNewTradesByBaseAssetID(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetStartAndEndBlockForNewTradesByBaseAssetID", err)
	}
	if startBlockNumberResult != nil {
		t.Errorf("Expected startBlockNumberResult From Method GetStartAndEndBlockForNewTradesByBaseAssetID: to be empty but got this: %v", startBlockNumberResult)
	}
	if endBlockNumberResult != nil {
		t.Errorf("Expected endBlockNumberResult From Method GetStartAndEndBlockForNewTradesByBaseAssetID: to be empty but got this: %v", endBlockNumberResult)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTrade(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	gethTradeID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_trades").WithArgs(*gethTradeID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethTrade(mock, gethTradeID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethTrade", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTradeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTradeID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_trades").WithArgs(gethTradeID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethTrade(mock, &gethTradeID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestDeleteGethTradesByBaseAssetId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	baseAssetID := targetData.BaseAssetID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_trades").WithArgs(*baseAssetID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = DeleteGethTradesByBaseAssetId(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in DeleteGethTradesByBaseAssetId", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestDeleteGethTradesByBaseAssetIdOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_trades").WithArgs(baseAssetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = DeleteGethTradesByBaseAssetId(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTrade{TestData1, TestData2}

	mockRows := AddGethTradeToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WillReturnRows(mockRows)
	foundGethTradeList, err := GetGethTradeList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeList", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTrade := range foundGethTradeList {
		if cmp.Equal(foundGethTrade, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTrade From Method GetGethTradeList: %v is different from actual %v", foundGethTrade, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeList, err := GetGethTradeList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeList", err)
	}
	if len(foundGethTradeList) != 0 {
		t.Errorf("Expected From Method GetGethTradeList: to be empty but got this: %v", foundGethTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTrade(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_trades").WithArgs(
		targetData.Name,                   //1
		targetData.AlternateName,          //2
		targetData.AddressStr,             //3
		targetData.AddressID,              //4
		targetData.TradeDate,              //5
		targetData.TxnHash,                //6
		targetData.Token0Amount,           //7
		targetData.Token0AmountDecimalAdj, //8
		targetData.Token1Amount,           //9
		targetData.Token1AmountDecimalAdj, //10
		targetData.IsBuy,                  //11
		targetData.Price,                  //12
		targetData.PriceUSD,               //13
		targetData.LPToken1PriceUSD,       //14
		targetData.TotalAmountUSD,         //15
		targetData.Token0AssetId,          //16
		targetData.Token1AssetId,          //17
		targetData.GethProcessJobID,       //18
		targetData.StatusID,               //19
		targetData.TradeTypeID,            //20
		targetData.Description,            //21
		targetData.UpdatedBy,              //22
		targetData.BaseAssetID,            //23
		targetData.OraclePriceUSD,         //24
		targetData.OraclePriceAssetID,     //25
		targetData.ID,                     //26
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethTrade(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethTrade", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTradeOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethTrade(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTradeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_trades").WithArgs(
		targetData.Name,                   //1
		targetData.AlternateName,          //2
		targetData.AddressStr,             //3
		targetData.AddressID,              //4
		targetData.TradeDate,              //5
		targetData.TxnHash,                //6
		targetData.Token0Amount,           //7
		targetData.Token0AmountDecimalAdj, //8
		targetData.Token1Amount,           //9
		targetData.Token1AmountDecimalAdj, //10
		targetData.IsBuy,                  //11
		targetData.Price,                  //12
		targetData.PriceUSD,               //13
		targetData.LPToken1PriceUSD,       //14
		targetData.TotalAmountUSD,         //15
		targetData.Token0AssetId,          //16
		targetData.Token1AssetId,          //17
		targetData.GethProcessJobID,       //18
		targetData.StatusID,               //19
		targetData.TradeTypeID,            //20
		targetData.Description,            //21
		targetData.UpdatedBy,              //22
		targetData.BaseAssetID,            //23
		targetData.OraclePriceUSD,         //24
		targetData.OraclePriceAssetID,     //25
		targetData.ID,                     //26
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethTrade(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTrade(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_trades").WithArgs(
		targetData.Name,                   //1
		targetData.AlternateName,          //2
		targetData.AddressStr,             //3
		targetData.AddressID,              //4
		targetData.TradeDate,              //5
		targetData.TxnHash,                //6
		targetData.Token0Amount,           //7
		targetData.Token0AmountDecimalAdj, //8
		targetData.Token1Amount,           //9
		targetData.Token1AmountDecimalAdj, //10
		targetData.IsBuy,                  //11
		targetData.Price,                  //12
		targetData.PriceUSD,               //13
		targetData.LPToken1PriceUSD,       //14
		targetData.TotalAmountUSD,         //15
		targetData.Token0AssetId,          //16
		targetData.Token1AssetId,          //17
		targetData.GethProcessJobID,       //18
		targetData.StatusID,               //19
		targetData.TradeTypeID,            //20
		targetData.Description,            //21
		targetData.CreatedBy,              //22
		targetData.BaseAssetID,            //23
		targetData.OraclePriceUSD,         //24
		targetData.OraclePriceAssetID,     //25
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit()
	gethTradeID, newUUID, err := InsertGethTrade(mock, &targetData)
	if gethTradeID < 0 {
		t.Fatalf("gethTradeID should not be negative ID: %d", gethTradeID)
	}
	if newUUID == "" {
		t.Fatalf("newUUID should not be empty string: %s", newUUID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTrade", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_trades").WithArgs(
		targetData.Name,                   //1
		targetData.AlternateName,          //2
		targetData.AddressStr,             //3
		targetData.AddressID,              //4
		targetData.TradeDate,              //5
		targetData.TxnHash,                //6
		targetData.Token0Amount,           //7
		targetData.Token0AmountDecimalAdj, //8
		targetData.Token1Amount,           //9
		targetData.Token1AmountDecimalAdj, //10
		targetData.IsBuy,                  //11
		targetData.Price,                  //12
		targetData.PriceUSD,               //13
		targetData.LPToken1PriceUSD,       //14
		targetData.TotalAmountUSD,         //15
		targetData.Token0AssetId,          //16
		targetData.Token1AssetId,          //17
		targetData.GethProcessJobID,       //18
		targetData.StatusID,               //19
		targetData.TradeTypeID,            //20
		targetData.Description,            //21
		targetData.CreatedBy,              //22
		targetData.BaseAssetID,            //23
		targetData.OraclePriceUSD,         //24
		targetData.OraclePriceAssetID,     //25
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethTradeID, newUUID, err := InsertGethTrade(mock, &targetData)
	if gethTradeID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethTradeID: %d", gethTradeID)
	}
	if newUUID != "" {
		t.Fatalf("on failure newUUID should be empty string: %s", newUUID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradeOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_trades").WithArgs(
		targetData.Name,                   //1
		targetData.AlternateName,          //2
		targetData.AddressStr,             //3
		targetData.AddressID,              //4
		targetData.TradeDate,              //5
		targetData.TxnHash,                //6
		targetData.Token0Amount,           //7
		targetData.Token0AmountDecimalAdj, //8
		targetData.Token1Amount,           //9
		targetData.Token1AmountDecimalAdj, //10
		targetData.IsBuy,                  //11
		targetData.Price,                  //12
		targetData.PriceUSD,               //13
		targetData.LPToken1PriceUSD,       //14
		targetData.TotalAmountUSD,         //15
		targetData.Token0AssetId,          //16
		targetData.Token1AssetId,          //17
		targetData.GethProcessJobID,       //18
		targetData.StatusID,               //19
		targetData.TradeTypeID,            //20
		targetData.Description,            //21
		targetData.CreatedBy,              //22
		targetData.BaseAssetID,            //23
		targetData.OraclePriceUSD,         //24
		targetData.OraclePriceAssetID,     //25
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethTradeID, newUUID, err := InsertGethTrade(mock, &targetData)
	if gethTradeID >= 0 {
		t.Fatalf("Expecting -1 for gethTradeID because of error gethTradeID: %d", gethTradeID)
	}
	if newUUID != "" {
		t.Fatalf("on failure newUUID should be empty string: %s", newUUID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTrades(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_trades"}, DBColumnsInsertGethTrades)
	targetData := TestAllData
	err = InsertGethTrades(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTrades", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTradesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_trades"}, DBColumnsInsertGethTrades).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertGethTrades(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLatestGethTradeFromAssetIDAnDate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethTrade{targetData}
	assetID := targetData.BaseAssetID
	asOfDate := utils.SampleCreatedAtTime
	isBefore := true
	mockRows := AddGethTradeToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(*assetID, asOfDate.Format(utils.LayoutPostgres)).WillReturnRows(mockRows)
	foundGethTrade, err := GetLatestGethTradeFromAssetIDAnDate(mock, assetID, asOfDate, &isBefore)
	if err != nil {
		t.Fatalf("an error '%s' in GetLatestGethTradeFromAssetIDAnDate", err)
	}
	if cmp.Equal(*foundGethTrade, targetData) == false {
		t.Errorf("Expected GethTrade From Method GetLatestGethTradeFromAssetIDAnDate: %v is different from actual %v", foundGethTrade, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLatestGethTradeFromAssetIDAnDateForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 111111
	asOfDate := utils.SampleCreatedAtTime
	isBefore := true
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(assetID, asOfDate.Format(utils.LayoutPostgres)).WillReturnRows(noRows)
	foundGethTrade, err := GetLatestGethTradeFromAssetIDAnDate(mock, &assetID, asOfDate, &isBefore)
	if err != nil {
		t.Fatalf("an error '%s' in GetLatestGethTradeFromAssetIDAnDate", err)
	}
	if foundGethTrade != nil {
		t.Errorf("Expected GethTrade From Method GetLatestGethTradeFromAssetIDAnDate: to be empty but got this: %v", foundGethTrade)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLatestGethTradeFromAssetIDAnDateForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := -1
	asOfDate := utils.SampleCreatedAtTime
	isBefore := true
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WithArgs(assetID, asOfDate.Format(utils.LayoutPostgres)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTrade, err := GetLatestGethTradeFromAssetIDAnDate(mock, &assetID, asOfDate, &isBefore)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetLatestGethTradeFromAssetIDAnDate", err)
	}
	if foundGethTrade != nil {
		t.Errorf("Expected GethTrade From Method GetLatestGethTradeFromAssetIDAnDate: to be empty but got this: %v", foundGethTrade)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTradeToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WillReturnRows(mockRows)
	foundGethTradeList, err := GetGethTradeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTradeListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundGethTradeList[i]) == false {
			t.Errorf("Expected sourceData From Method GetGethTradeListByPagination: %v is different from actual %v", sourceData, foundGethTradeList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTradeListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_trades").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTradeList, err := GetGethTradeListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTradeListByPagination", err)
	}
	if len(foundGethTradeList) != 0 {
		t.Errorf("Expected From Method GetGethTradeListByPagination: to be empty but got this: %v", foundGethTradeList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalTransactionsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_trades").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalTransactionsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalTransactionsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalTransactionsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_trades").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalTransactionsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalTransactionsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalTransactionsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
