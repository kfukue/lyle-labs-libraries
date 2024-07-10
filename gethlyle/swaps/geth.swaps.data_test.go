package gethlyleswaps

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/shopspring/decimal"
)

var DBColumns = []string{
	"id",                    //1
	"uuid",                  //2
	"chain_id",              //3
	"exchange_id",           //4
	"block_number",          //5
	"index_number",          //6
	"swap_date",             //7
	"trade_type_id",         //8
	"txn_hash",              //9
	"maker_address",         //10
	"maker_address_id",      //11
	"is_buy",                //12
	"price",                 //13
	"price_usd",             //14
	"token1_price_usd",      //15
	"total_amount_usd",      //16
	"pair_address",          //17
	"liquidity_pool_id",     //18
	"token0_asset_id",       //19
	"token1_asset_id",       //20
	"token0_amount",         //21
	"token1_amount",         //22
	"description",           //23
	"created_by",            //24
	"created_at",            //25
	"updated_by",            //26
	"updated_at",            //27
	"geth_process_job_id",   //28
	"topics_str",            //29
	"status_id",             //30
	"base_asset_id",         //31
	"oracle_price_usd",      //32
	"oracle_price_asset_id", //33
}
var DBColumnsInsertGethSwaps = []string{
	"uuid",                  //1
	"chain_id",              //2
	"exchange_id",           //3
	"block_number",          //4
	"index_number",          //5
	"swap_date",             //6
	"trade_type_id",         //7
	"txn_hash",              //8
	"maker_address",         //9
	"maker_address_id",      //10
	"is_buy",                //11
	"price",                 //12
	"price_usd",             //13
	"token1_price_usd",      //14
	"total_amount_usd",      //15
	"pair_address",          //16
	"liquidity_pool_id",     //17
	"token0_asset_id",       //18
	"token1_asset_id",       //19
	"token0_amount",         //20
	"token1_amount",         //21
	"description",           //22
	"created_by",            //23
	"created_at",            //24
	"updated_by",            //25
	"updated_at",            //26
	"geth_process_job_id",   //27
	"topics_str",            //28
	"status_id",             //29
	"base_asset_id",         //30
	"oracle_price_usd",      //31
	"oracle_price_asset_id", //32
}

var TestData1 = GethSwap{
	ID:                 utils.Ptr[int](1),                                                                                                              //1
	UUID:               "01ef85e8-2c26-441e-8c7f-71d79518ad72",                                                                                         //2
	ChainID:            utils.Ptr[int](1),                                                                                                              //3
	ExchangeID:         utils.Ptr[int](2),                                                                                                              //4
	BlockNumber:        utils.Ptr[uint64](17387265),                                                                                                    //5
	IndexNumber:        utils.Ptr[uint](76),                                                                                                            //6
	SwapDate:           utils.SampleCreatedAtTime,                                                                                                      //7
	TradeTypeID:        utils.Ptr[int](2),                                                                                                              //8
	TxnHash:            "0x67775b7b31ff14d7a52c883e5ffe1a10cbdacb28c59728c5a78948863aa31b3b",                                                           //9
	MakerAddress:       "0x00000000000124d994209fbB955E0217B5C2ECA1",                                                                                   //10
	MakerAddressID:     utils.Ptr[int](12404),                                                                                                          //11
	IsBuy:              utils.Ptr[bool](false),                                                                                                         //12
	Price:              utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.0000000016482182)),                                                           //13
	PriceUSD:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.00000306905761197174)),                                                       //14
	Token1PriceUSD:     utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.01)),                                                                         //15
	TotalAmountUSD:     utils.Ptr[decimal.Decimal](decimal.NewFromFloat(-2093.4883905833930747)),                                                       //16
	PairAddress:        "0xd101821c56B4405Af4A376cBe81FA0dC90207dC2",                                                                                   //17
	LiquidityPoolID:    utils.Ptr[int](1),                                                                                                              //18
	Token0AssetId:      utils.Ptr[int](1),                                                                                                              //19
	Token1AssetId:      utils.Ptr[int](2),                                                                                                              //20
	Token0Amount:       utils.Ptr[decimal.Decimal](decimal.NewFromFloat(-682127433000000000000000000)),                                                 //21
	Token1Amount:       utils.Ptr[decimal.Decimal](decimal.NewFromFloat(1124294835794079411)),                                                          //22
	Description:        "Imported by Geth Dex Analyzer",                                                                                                //23
	CreatedBy:          "SYSTEM",                                                                                                                       //24
	CreatedAt:          utils.SampleCreatedAtTime,                                                                                                      //25
	UpdatedBy:          "SYSTEM",                                                                                                                       //26
	UpdatedAt:          utils.SampleCreatedAtTime,                                                                                                      //27
	GethProcessJobID:   utils.Ptr[int](10),                                                                                                             //28
	TopicsStr:          []string{"Swap(address,uint256,uint256,uint256,uint256,address)", "Swap(address,address,int256,int256,uint160,uint128,int24)"}, //29
	StatusID:           utils.Ptr[int](52),                                                                                                             //30
	BaseAssetID:        utils.Ptr[int](1),                                                                                                              //31
	OraclePriceUSD:     utils.Ptr[decimal.Decimal](decimal.NewFromFloat(1862.0457)),                                                                    //32
	OraclePriceAssetID: utils.Ptr[int](1),                                                                                                              //33
}

var TestData2 = GethSwap{
	ID:                 utils.Ptr[int](2),                                                                                                              //1
	UUID:               "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                                                                                         //2
	ChainID:            utils.Ptr[int](1),                                                                                                              //3
	ExchangeID:         utils.Ptr[int](2),                                                                                                              //4
	BlockNumber:        utils.Ptr[uint64](20267173),                                                                                                    //5
	IndexNumber:        utils.Ptr[uint](19),                                                                                                            //6
	SwapDate:           utils.SampleCreatedAtTime,                                                                                                      //7
	TradeTypeID:        utils.Ptr[int](2),                                                                                                              //8
	TxnHash:            "0xcdae57cf75ad8f3b8051ba0d8a9bcfc247e1416910f4850238d0e956ab9b82d5",                                                           //9
	MakerAddress:       "0xe75eD6F453c602Bd696cE27AF11565eDc9b46B0D",                                                                                   //10
	MakerAddressID:     utils.Ptr[int](9730),                                                                                                           //11
	IsBuy:              utils.Ptr[bool](true),                                                                                                          //12
	Price:              utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.0196828489853681)),                                                           //13
	PriceUSD:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(60.637139547661794576378192)),                                                  //14
	Token1PriceUSD:     nil,                                                                                                                            //15
	TotalAmountUSD:     utils.Ptr[decimal.Decimal](decimal.NewFromFloat(1684.2238633602744855)),                                                        //16
	PairAddress:        "0x53E79ef1Cf6aC0cDF6f1743C3BE3ad48fA3c5657",                                                                                   //17
	LiquidityPoolID:    utils.Ptr[int](2),                                                                                                              //18
	Token0AssetId:      utils.Ptr[int](3),                                                                                                              //19
	Token1AssetId:      utils.Ptr[int](4),                                                                                                              //20
	Token0Amount:       utils.Ptr[decimal.Decimal](decimal.NewFromFloat(-546699996196503552)),                                                          //21
	Token1Amount:       utils.Ptr[decimal.Decimal](decimal.NewFromFloat(27775450424016896000)),                                                         //22
	Description:        "Imported by Geth Dex Analyzer",                                                                                                //23
	CreatedBy:          "SYSTEM",                                                                                                                       //24
	CreatedAt:          utils.SampleCreatedAtTime,                                                                                                      //25
	UpdatedBy:          "SYSTEM",                                                                                                                       //26
	UpdatedAt:          utils.SampleCreatedAtTime,                                                                                                      //27
	GethProcessJobID:   utils.Ptr[int](10),                                                                                                             //28
	TopicsStr:          []string{"Swap(address,uint256,uint256,uint256,uint256,address)", "Swap(address,address,int256,int256,uint160,uint128,int24)"}, //29
	StatusID:           utils.Ptr[int](52),                                                                                                             //30
	BaseAssetID:        utils.Ptr[int](1),                                                                                                              //31
	OraclePriceUSD:     utils.Ptr[decimal.Decimal](decimal.NewFromFloat(3080.70948432)),                                                                //32
	OraclePriceAssetID: utils.Ptr[int](1),
}
var TestAllData = []GethSwap{TestData1, TestData2}

func AddGethSwapToMockRows(mock pgxmock.PgxPoolIface, dataList []GethSwap) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                 //1
			data.UUID,               //2
			data.ChainID,            //3
			data.ExchangeID,         //4
			data.BlockNumber,        //5
			data.IndexNumber,        //6
			data.SwapDate,           //7
			data.TradeTypeID,        //8
			data.TxnHash,            //9
			data.MakerAddress,       //10
			data.MakerAddressID,     //11
			data.IsBuy,              //12
			data.Price,              //13
			data.PriceUSD,           //14
			data.Token1PriceUSD,     //15
			data.TotalAmountUSD,     //16
			data.PairAddress,        //17
			data.LiquidityPoolID,    //18
			data.Token0AssetId,      //19
			data.Token1AssetId,      //20
			data.Token0Amount,       //21
			data.Token1Amount,       //22
			data.Description,        //23
			data.CreatedBy,          //24
			data.CreatedAt,          //25
			data.UpdatedBy,          //26
			data.UpdatedAt,          //27
			data.GethProcessJobID,   //28
			data.TopicsStr,          //29
			data.StatusID,           //30
			data.BaseAssetID,        //31
			data.OraclePriceUSD,     //32
			data.OraclePriceAssetID, //33
		)
	}
	return rows
}

func TestGetGethSwapByBlockChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	dataList := []GethSwap{targetData}
	txnHash := targetData.TxnHash
	blockNumber := targetData.BlockNumber
	indexNumber := targetData.IndexNumber
	makerAddressID := targetData.MakerAddressID
	liquidityPoolID := targetData.LiquidityPoolID
	mockRows := AddGethSwapToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(txnHash, *blockNumber, *indexNumber, *makerAddressID, *liquidityPoolID).WillReturnRows(mockRows)
	foundGethSwap, err := GetGethSwapByBlockChain(mock, txnHash, blockNumber, indexNumber, makerAddressID, liquidityPoolID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapByBlockChain", err)
	}
	if cmp.Equal(*foundGethSwap, targetData) == false {
		t.Errorf("Expected GethSwap From Method GetGethSwapByBlockChain: %v is different from actual %v", foundGethSwap, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByBlockChainForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	txnHash := targetData.TxnHash
	blockNumber := targetData.BlockNumber
	indexNumber := targetData.IndexNumber
	makerAddressID := 99999
	liquidityPoolID := targetData.LiquidityPoolID
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(txnHash, *blockNumber, *indexNumber, makerAddressID, *liquidityPoolID).WillReturnRows(noRows)
	foundGethSwap, err := GetGethSwapByBlockChain(mock, txnHash, blockNumber, indexNumber, &makerAddressID, liquidityPoolID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapByBlockChain", err)
	}
	if foundGethSwap != nil {
		t.Errorf("Expected GethSwap From Method GetGethSwapByBlockChain: to be empty but got this: %v", foundGethSwap)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByBlockChainForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	txnHash := targetData.TxnHash
	blockNumber := targetData.BlockNumber
	indexNumber := targetData.IndexNumber
	makerAddressID := -1
	liquidityPoolID := targetData.LiquidityPoolID
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(txnHash, *blockNumber, *indexNumber, makerAddressID, *liquidityPoolID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwap, err := GetGethSwapByBlockChain(mock, txnHash, blockNumber, indexNumber, &makerAddressID, liquidityPoolID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwapByBlockChain", err)
	}
	if foundGethSwap != nil {
		t.Errorf("Expected GethSwap From Method GetGethSwapByBlockChain: to be empty but got this: %v", foundGethSwap)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestGetGethSwap(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethSwap{targetData}
	gethSwapID := targetData.ID
	mockRows := AddGethSwapToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*gethSwapID).WillReturnRows(mockRows)
	foundGethSwap, err := GetGethSwap(mock, gethSwapID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwap", err)
	}
	if cmp.Equal(*foundGethSwap, targetData) == false {
		t.Errorf("Expected GethSwap From Method GetGethSwap: %v is different from actual %v", foundGethSwap, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethSwapID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(gethSwapID).WillReturnRows(noRows)
	foundGethSwap, err := GetGethSwap(mock, &gethSwapID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwap", err)
	}
	if foundGethSwap != nil {
		t.Errorf("Expected GethSwap From Method GetGethSwap: to be empty but got this: %v", foundGethSwap)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethSwapID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(gethSwapID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwap, err := GetGethSwap(mock, &gethSwapID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwap", err)
	}
	if foundGethSwap != nil {
		t.Errorf("Expected GethSwap From Method GetGethSwap: to be empty but got this: %v", foundGethSwap)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByStartAndEndDates(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethSwap{TestData1, TestData2}
	startDate := utils.SampleCreatedAtTime
	endDate := utils.SampleCreatedAtTime
	mockRows := AddGethSwapToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnRows(mockRows)
	foundGethSwapList, err := GetGethSwapByStartAndEndDates(mock, startDate, endDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapByStartAndEndDates", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethSwap := range foundGethSwapList {
		if cmp.Equal(foundGethSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethSwap From Method GetGethSwapByStartAndEndDates: %v is different from actual %v", foundGethSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByStartAndEndDatesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	startDate := time.Now()
	endDate := utils.SampleCreatedAtTime.Add(time.Hour * 24)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwapList, err := GetGethSwapByStartAndEndDates(mock, startDate, endDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwapByStartAndEndDates", err)
	}
	if len(foundGethSwapList) != 0 {
		t.Errorf("Expected From Method GetGethSwapByStartAndEndDates: to be empty but got this: %v", foundGethSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByFromMakerAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethSwap{TestData1, TestData2}
	mockRows := AddGethSwapToMockRows(mock, dataList)
	makerAddress := TestData1.MakerAddress
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(makerAddress).WillReturnRows(mockRows)
	foundGethSwapList, err := GetGethSwapByFromMakerAddress(mock, makerAddress)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapByFromMakerAddress", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethSwap := range foundGethSwapList {
		if cmp.Equal(foundGethSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethSwap From Method GetGethSwapByFromMakerAddress: %v is different from actual %v", foundGethSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByFromMakerAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	makerAddress := "0x-invalid-1"
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(makerAddress).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwapList, err := GetGethSwapByFromMakerAddress(mock, makerAddress)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwapByFromMakerAddress", err)
	}
	if len(foundGethSwapList) != 0 {
		t.Errorf("Expected From Method GetGethSwapByFromMakerAddress: to be empty but got this: %v", foundGethSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByFromMakerAddressId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethSwap{TestData1, TestData2}
	mockRows := AddGethSwapToMockRows(mock, dataList)
	makerAddressID := TestData1.MakerAddressID
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*makerAddressID).WillReturnRows(mockRows)
	foundGethSwapList, err := GetGethSwapByFromMakerAddressId(mock, makerAddressID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapByFromMakerAddressId", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethSwap := range foundGethSwapList {
		if cmp.Equal(foundGethSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethSwap From Method GetGethSwapByFromMakerAddressId: %v is different from actual %v", foundGethSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByFromMakerAddressIdForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	makerAddressID := TestData1.MakerAddressID
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*makerAddressID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwapList, err := GetGethSwapByFromMakerAddressId(mock, makerAddressID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwapByFromMakerAddressId", err)
	}
	if len(foundGethSwapList) != 0 {
		t.Errorf("Expected From Method GetGethSwapByFromMakerAddressId: to be empty but got this: %v", foundGethSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByFromMakerAddressIdAndBeforeBlockNumber(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethSwap{TestData1, TestData2}
	mockRows := AddGethSwapToMockRows(mock, dataList)
	baseAssetID := TestData1.BaseAssetID
	makerAddressID := TestData1.MakerAddressID
	blockNumber := TestData1.BlockNumber
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*baseAssetID, *makerAddressID, *blockNumber).WillReturnRows(mockRows)
	foundGethSwapList, err := GetGethSwapByFromMakerAddressIdAndBeforeBlockNumber(mock, baseAssetID, makerAddressID, blockNumber)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapByFromMakerAddressIdAndBeforeBlockNumber", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethSwap := range foundGethSwapList {
		if cmp.Equal(foundGethSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethSwap From Method GetGethSwapByFromMakerAddressIdAndBeforeBlockNumber: %v is different from actual %v", foundGethSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByFromMakerAddressIdAndBeforeBlockNumberForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := TestData1.BaseAssetID
	makerAddressID := -1
	blockNumber := TestData1.BlockNumber
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*baseAssetID, makerAddressID, *blockNumber).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwapList, err := GetGethSwapByFromMakerAddressIdAndBeforeBlockNumber(mock, baseAssetID, &makerAddressID, blockNumber)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwapByFromMakerAddressIdAndBeforeBlockNumber", err)
	}
	if len(foundGethSwapList) != 0 {
		t.Errorf("Expected From Method GetGethSwapByFromMakerAddressIdAndBeforeBlockNumber: to be empty but got this: %v", foundGethSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByFromBaseAssetAndBeforeBlockNumber(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethSwap{TestData1, TestData2}
	mockRows := AddGethSwapToMockRows(mock, dataList)
	baseAssetID := TestData1.BaseAssetID
	blockNumber := TestData1.BlockNumber
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*baseAssetID, *blockNumber).WillReturnRows(mockRows)
	foundGethSwapList, err := GetGethSwapByFromBaseAssetAndBeforeBlockNumber(mock, baseAssetID, blockNumber)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapByFromBaseAssetAndBeforeBlockNumber", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethSwap := range foundGethSwapList {
		if cmp.Equal(foundGethSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethSwap From Method GetGethSwapByFromBaseAssetAndBeforeBlockNumber: %v is different from actual %v", foundGethSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByFromBaseAssetAndBeforeBlockNumberForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	blockNumber := TestData1.BlockNumber
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID, *blockNumber).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwapList, err := GetGethSwapByFromBaseAssetAndBeforeBlockNumber(mock, &baseAssetID, blockNumber)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwapByFromBaseAssetAndBeforeBlockNumber", err)
	}
	if len(foundGethSwapList) != 0 {
		t.Errorf("Expected From Method GetGethSwapByFromBaseAssetAndBeforeBlockNumber: to be empty but got this: %v", foundGethSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByTxnHash(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethSwap{TestData1, TestData2}
	mockRows := AddGethSwapToMockRows(mock, dataList)
	baseAssetID := TestData1.BaseAssetID
	txnHash := TestData1.TxnHash
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(txnHash, utils.EOA_ADDRESS_TYPE_STRUCTURED_VALUE_ID, *baseAssetID).WillReturnRows(mockRows)
	foundGethSwapList, err := GetGethSwapByTxnHash(mock, txnHash, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapByTxnHash", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethSwap := range foundGethSwapList {
		if cmp.Equal(foundGethSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethSwap From Method GetGethSwapByTxnHash: %v is different from actual %v", foundGethSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapByTxnHashForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	txnHash := TestData1.TxnHash
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(txnHash, utils.EOA_ADDRESS_TYPE_STRUCTURED_VALUE_ID, baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwapList, err := GetGethSwapByTxnHash(mock, txnHash, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwapByTxnHash", err)
	}
	if len(foundGethSwapList) != 0 {
		t.Errorf("Expected From Method GetGethSwapByTxnHash: to be empty but got this: %v", foundGethSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapsByTxnHashes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethSwap{TestData1, TestData2}

	mockRows := AddGethSwapToMockRows(mock, dataList)
	baseAssetID := TestData1.BaseAssetID
	txnHashes := []string{TestData1.TxnHash, TestData2.TxnHash}
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(pq.Array(txnHashes), utils.EOA_ADDRESS_TYPE_STRUCTURED_VALUE_ID, *baseAssetID).WillReturnRows(mockRows)
	foundGethSwapList, err := GetGethSwapsByTxnHashes(mock, txnHashes, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapsByTxnHashes", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethSwap := range foundGethSwapList {
		if cmp.Equal(foundGethSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethSwap From Method GetGethSwapsByTxnHashes: %v is different from actual %v", foundGethSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapsByTxnHashesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	txnHashes := []string{TestData1.TxnHash, TestData2.TxnHash}
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(pq.Array(txnHashes), utils.EOA_ADDRESS_TYPE_STRUCTURED_VALUE_ID, baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwapList, err := GetGethSwapsByTxnHashes(mock, txnHashes, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwapsByTxnHashes", err)
	}
	if len(foundGethSwapList) != 0 {
		t.Errorf("Expected From Method GetGethSwapsByTxnHashes: to be empty but got this: %v", foundGethSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctTransactionHashesFromAssetIdAndStartingBlock(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := TestData1.BaseAssetID
	startingBlock := utils.Ptr[uint64](1)
	txnHashResults := []string{TestData1.TxnHash, TestData2.TxnHash}
	mockRows := mock.NewRows([]string{"txn_hash"}).AddRow(TestData1.TxnHash).AddRow(TestData2.TxnHash)

	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*startingBlock, *assetID).WillReturnRows(mockRows)
	foundTxnHashList, err := GetDistinctTransactionHashesFromAssetIdAndStartingBlock(mock, assetID, startingBlock)
	if err != nil {
		t.Fatalf("an error '%s' in GetDistinctTransactionHashesFromAssetIdAndStartingBlock", err)
	}
	for i, foundTxnHash := range foundTxnHashList {
		if cmp.Equal(foundTxnHash, txnHashResults[i]) == false {
			t.Errorf("Expected foundTxnHash From Method GetDistinctTransactionHashesFromAssetIdAndStartingBlock: %v is different from actual %v", foundTxnHash, txnHashResults[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctTransactionHashesFromAssetIdAndStartingBlockForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := TestData1.BaseAssetID
	startingBlock := utils.Ptr[uint64](1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*startingBlock, *assetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTxnHashList, err := GetDistinctTransactionHashesFromAssetIdAndStartingBlock(mock, assetID, startingBlock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDistinctTransactionHashesFromAssetIdAndStartingBlock", err)
	}
	if len(foundTxnHashList) != 0 {
		t.Errorf("Expected From Method GetDistinctTransactionHashesFromAssetIdAndStartingBlock: to be empty but got this: %v", foundTxnHashList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetHighestBlockFromBaseAssetId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()

	baseAssetID := TestData1.BaseAssetID
	highestBlockNumberResult := TestData1.BlockNumber
	targetData := *highestBlockNumberResult
	mockRows := mock.NewRows([]string{"block_number"}).AddRow(*highestBlockNumberResult)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*baseAssetID).WillReturnRows(mockRows)
	highestBlockNumber, err := GetHighestBlockFromBaseAssetId(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetHighestBlockFromBaseAssetId", err)
	}
	if cmp.Equal(*highestBlockNumber, targetData) == false {
		t.Errorf("Expected highestBlockNumber From Method GetHighestBlockFromBaseAssetId: %v is different from actual %v", *highestBlockNumber, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetHighestBlockFromBaseAssetIdForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := 99999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID).WillReturnRows(noRows)
	highestBlockNumber, err := GetHighestBlockFromBaseAssetId(mock, &baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetHighestBlockFromBaseAssetId", err)
	}
	if *highestBlockNumber > 0 {
		t.Errorf("Expected highestBlockNumber From Method GetHighestBlockFromBaseAssetId: to be 0 but got this: %v", *highestBlockNumber)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetHighestBlockFromBaseAssetIdForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	highestBlockNumber, err := GetHighestBlockFromBaseAssetId(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetHighestBlockFromBaseAssetId", err)
	}
	if highestBlockNumber != nil {
		t.Errorf("Expected highestBlockNumber From Method GetHighestBlockFromBaseAssetId: to be empty but got this: %v", highestBlockNumber)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctMakerAddressesFromBaseTokenAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := TestData1.BaseAssetID
	makerAddressIDResults := []int{*TestData1.MakerAddressID, *TestData2.MakerAddressID}
	mockRows := mock.NewRows([]string{"maker_address_id"}).AddRow(*TestData1.MakerAddressID).AddRow(*TestData2.MakerAddressID)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*baseAssetID).WillReturnRows(mockRows)
	foundMakerAddressIDs, err := GetDistinctMakerAddressesFromBaseTokenAssetID(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetDistinctMakerAddressesFromBaseTokenAssetID", err)
	}
	for i, foundMakerAddressID := range foundMakerAddressIDs {
		if cmp.Equal(foundMakerAddressID, makerAddressIDResults[i]) == false {
			t.Errorf("Expected foundMakerAddressID From Method GetDistinctMakerAddressesFromBaseTokenAssetID: %v is different from actual %v", foundMakerAddressID, makerAddressIDResults[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctMakerAddressesFromBaseTokenAssetIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1

	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMakerAddressIDs, err := GetDistinctMakerAddressesFromBaseTokenAssetID(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDistinctMakerAddressesFromBaseTokenAssetID", err)
	}
	if len(foundMakerAddressIDs) != 0 {
		t.Errorf("Expected From Method GetDistinctMakerAddressesFromBaseTokenAssetID: to be empty but got this: %v", foundMakerAddressIDs)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethSwap(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	gethSwapID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_swaps").WithArgs(*gethSwapID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethSwap(mock, gethSwapID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethSwap", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethSwapOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethSwapID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_swaps").WithArgs(gethSwapID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethSwap(mock, &gethSwapID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethSwapsFromAssetIDAndStartBlockNumber(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	baseAssetID := targetData.BaseAssetID
	startBlockNumber := targetData.BlockNumber
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_swaps").WithArgs(*baseAssetID, *startBlockNumber).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethSwapsFromAssetIDAndStartBlockNumber(mock, baseAssetID, startBlockNumber)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethSwapsFromAssetIDAndStartBlockNumber", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethSwapsFromAssetIDAndStartBlockNumberOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainID := -1
	startBlockNumber := uint64(10000)
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_swaps").WithArgs(chainID, startBlockNumber).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethSwapsFromAssetIDAndStartBlockNumber(mock, &chainID, &startBlockNumber)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestDeleteGethSwapsByBaseAssetId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	baseAssetID := targetData.BaseAssetID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_swaps").WithArgs(*baseAssetID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = DeleteGethSwapsByBaseAssetId(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in DeleteGethSwapsByBaseAssetId", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestDeleteGethSwapsByBaseAssetIdOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_swaps").WithArgs(baseAssetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = DeleteGethSwapsByBaseAssetId(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethSwap{TestData1, TestData2}

	mockRows := AddGethSwapToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WillReturnRows(mockRows)
	foundGethSwapList, err := GetGethSwapList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapList", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethSwap := range foundGethSwapList {
		if cmp.Equal(foundGethSwap, testMarketDataList[i]) == false {
			t.Errorf("Expected GethSwap From Method GetGethSwapList: %v is different from actual %v", foundGethSwap, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwapList, err := GetGethSwapList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwapList", err)
	}
	if len(foundGethSwapList) != 0 {
		t.Errorf("Expected From Method GetGethSwapList: to be empty but got this: %v", foundGethSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethSwap(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_swaps").WithArgs(
		targetData.ChainID,             //1
		targetData.ExchangeID,          //2
		targetData.BlockNumber,         //3
		targetData.IndexNumber,         //4
		targetData.SwapDate,            //5
		targetData.TradeTypeID,         //6
		targetData.TxnHash,             //7
		targetData.MakerAddress,        //8
		targetData.MakerAddressID,      //9
		targetData.IsBuy,               //10
		targetData.Price,               //11
		targetData.PriceUSD,            //12
		targetData.Token1PriceUSD,      //13
		targetData.TotalAmountUSD,      //14
		targetData.PairAddress,         //15
		targetData.LiquidityPoolID,     //16
		targetData.Token0AssetId,       //17
		targetData.Token1AssetId,       //18
		targetData.Token0Amount,        //19
		targetData.Token1Amount,        //20
		targetData.Description,         //21
		targetData.UpdatedBy,           //22
		targetData.GethProcessJobID,    //23
		pq.Array(targetData.TopicsStr), //24
		targetData.StatusID,            //25
		targetData.BaseAssetID,         //26
		targetData.OraclePriceUSD,      //27
		targetData.OraclePriceAssetID,  //28
		targetData.ID,                  //29
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethSwap(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethSwap", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethSwapOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethSwap(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethSwapOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_swaps").WithArgs(
		targetData.ChainID,             //1
		targetData.ExchangeID,          //2
		targetData.BlockNumber,         //3
		targetData.IndexNumber,         //4
		targetData.SwapDate,            //5
		targetData.TradeTypeID,         //6
		targetData.TxnHash,             //7
		targetData.MakerAddress,        //8
		targetData.MakerAddressID,      //9
		targetData.IsBuy,               //10
		targetData.Price,               //11
		targetData.PriceUSD,            //12
		targetData.Token1PriceUSD,      //13
		targetData.TotalAmountUSD,      //14
		targetData.PairAddress,         //15
		targetData.LiquidityPoolID,     //16
		targetData.Token0AssetId,       //17
		targetData.Token1AssetId,       //18
		targetData.Token0Amount,        //19
		targetData.Token1Amount,        //20
		targetData.Description,         //21
		targetData.UpdatedBy,           //22
		targetData.GethProcessJobID,    //23
		pq.Array(targetData.TopicsStr), //24
		targetData.StatusID,            //25
		targetData.BaseAssetID,         //26
		targetData.OraclePriceUSD,      //27
		targetData.OraclePriceAssetID,  //28
		targetData.ID,                  //29
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethSwap(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethSwap(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_swaps").WithArgs(
		targetData.ChainID,             //1
		targetData.ExchangeID,          //2
		targetData.BlockNumber,         //3
		targetData.IndexNumber,         //4
		targetData.SwapDate,            //5
		targetData.TradeTypeID,         //6
		targetData.TxnHash,             //7
		targetData.MakerAddress,        //8
		targetData.MakerAddressID,      //9
		targetData.IsBuy,               //10
		targetData.Price,               //11
		targetData.PriceUSD,            //12
		targetData.Token1PriceUSD,      //13
		targetData.TotalAmountUSD,      //14
		targetData.PairAddress,         //15
		targetData.LiquidityPoolID,     //16
		targetData.Token0AssetId,       //17
		targetData.Token1AssetId,       //18
		targetData.Token0Amount,        //19
		targetData.Token1Amount,        //20
		targetData.Description,         //21
		targetData.CreatedBy,           //22
		targetData.GethProcessJobID,    //23
		pq.Array(targetData.TopicsStr), //24
		targetData.StatusID,            //25
		targetData.BaseAssetID,         //26
		targetData.OraclePriceUSD,      //27
		targetData.OraclePriceAssetID,  //28
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit()
	gethSwapID, newUUID, err := InsertGethSwap(mock, &targetData)
	if gethSwapID < 0 {
		t.Fatalf("gethSwapID should not be negative ID: %d", gethSwapID)
	}
	if newUUID == "" {
		t.Fatalf("newUUID should not be empty string: %s", newUUID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethSwap", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethSwapOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.BlockNumber = nil
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_swaps").WithArgs(
		targetData.ChainID,             //1
		targetData.ExchangeID,          //2
		targetData.BlockNumber,         //3
		targetData.IndexNumber,         //4
		targetData.SwapDate,            //5
		targetData.TradeTypeID,         //6
		targetData.TxnHash,             //7
		targetData.MakerAddress,        //8
		targetData.MakerAddressID,      //9
		targetData.IsBuy,               //10
		targetData.Price,               //11
		targetData.PriceUSD,            //12
		targetData.Token1PriceUSD,      //13
		targetData.TotalAmountUSD,      //14
		targetData.PairAddress,         //15
		targetData.LiquidityPoolID,     //16
		targetData.Token0AssetId,       //17
		targetData.Token1AssetId,       //18
		targetData.Token0Amount,        //19
		targetData.Token1Amount,        //20
		targetData.Description,         //21
		targetData.CreatedBy,           //22
		targetData.GethProcessJobID,    //23
		pq.Array(targetData.TopicsStr), //24
		targetData.StatusID,            //25
		targetData.BaseAssetID,         //26
		targetData.OraclePriceUSD,      //27
		targetData.OraclePriceAssetID,  //28
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethSwapID, newUUID, err := InsertGethSwap(mock, &targetData)
	if gethSwapID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethSwapID: %d", gethSwapID)
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

func TestInsertGethSwapOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_swaps").WithArgs(
		targetData.ChainID,             //1
		targetData.ExchangeID,          //2
		targetData.BlockNumber,         //3
		targetData.IndexNumber,         //4
		targetData.SwapDate,            //5
		targetData.TradeTypeID,         //6
		targetData.TxnHash,             //7
		targetData.MakerAddress,        //8
		targetData.MakerAddressID,      //9
		targetData.IsBuy,               //10
		targetData.Price,               //11
		targetData.PriceUSD,            //12
		targetData.Token1PriceUSD,      //13
		targetData.TotalAmountUSD,      //14
		targetData.PairAddress,         //15
		targetData.LiquidityPoolID,     //16
		targetData.Token0AssetId,       //17
		targetData.Token1AssetId,       //18
		targetData.Token0Amount,        //19
		targetData.Token1Amount,        //20
		targetData.Description,         //21
		targetData.CreatedBy,           //22
		targetData.GethProcessJobID,    //23
		pq.Array(targetData.TopicsStr), //24
		targetData.StatusID,            //25
		targetData.BaseAssetID,         //26
		targetData.OraclePriceUSD,      //27
		targetData.OraclePriceAssetID,  //28
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethSwapID, newUUID, err := InsertGethSwap(mock, &targetData)
	if gethSwapID >= 0 {
		t.Fatalf("Expecting -1 for gethSwapID because of error gethSwapID: %d", gethSwapID)
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

func TestInsertGethSwaps(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_swaps"}, DBColumnsInsertGethSwaps)
	targetData := TestAllData
	err = InsertGethSwaps(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethSwaps", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethSwapsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_swaps"}, DBColumnsInsertGethSwaps).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertGethSwaps(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetNullAddressStrsFromSwaps(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []string{TestData1.MakerAddress, TestData2.MakerAddress}
	assetID := TestData1.BaseAssetID
	mockRows := mock.NewRows([]string{"address"}).AddRow(TestData1.MakerAddress).AddRow(TestData2.MakerAddress)
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(*assetID).WillReturnRows(mockRows)
	foundNullAddresses, err := GetNullAddressStrsFromSwaps(mock, assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetNullAddressStrsFromSwaps", err)
	}
	for i, nullAddress := range foundNullAddresses {
		if cmp.Equal(nullAddress, dataList[i]) == false {
			t.Errorf("Expected GethSwap From Method GetNullAddressStrsFromSwaps: %v is different from actual %v", nullAddress, dataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetNullAddressStrsFromSwapsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WithArgs(assetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethSwapList, err := GetNullAddressStrsFromSwaps(mock, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetNullAddressStrsFromSwaps", err)
	}
	if len(foundGethSwapList) != 0 {
		t.Errorf("Expected From Method GetNullAddressStrsFromSwaps: to be empty but got this: %v", foundGethSwapList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethSwapAddresses(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_swaps").WithArgs(*baseAssetID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethSwapAddresses(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethSwapAddresses", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethSwapAddressesOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethSwapAddresses(mock, baseAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethSwapAddressesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_swaps").WithArgs(*baseAssetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethSwapAddresses(mock, baseAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethSwapToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WillReturnRows(mockRows)
	foundChains, err := GetGethSwapListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethSwapListByPagination", err)
	}
	testChains := dataList
	for i, foundChain := range foundChains {
		if cmp.Equal(foundChain, testChains[i]) == false {
			t.Errorf("Expected Chain From Method GetGethSwapListByPagination: %v is different from actual %v", foundChain, testChains[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethSwapListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM geth_swaps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChains, err := GetGethSwapListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethSwapListByPagination", err)
	}
	if len(foundChains) != 0 {
		t.Errorf("Expected From Method GetGethSwapListByPagination: to be empty but got this: %v", foundChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethSwapsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_swaps").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalGethSwapsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalGethSwapsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalGethSwapsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_swaps").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalGethSwapsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalGethSwapsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalGethSwapsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
