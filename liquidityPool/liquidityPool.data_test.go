package liquiditypool

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/asset"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"id",                                //1
	"uuid",                              //2
	"name",                              //3
	"alternate_name",                    //4
	"pair_address",                      //5
	"chain_id",                          //6
	"exchange_id",                       //7
	"liquidity_pool_type_id",            //8
	"token0_id",                         //9
	"token1_id",                         //10
	"url",                               //11
	"start_block",                       //12
	"latest_block_synced",               //13
	"created_txn_hash",                  //14
	"IsActive",                          //15
	"description",                       //16
	"created_by",                        //17
	"created_at",                        //18
	"updated_by",                        //19
	"updated_at",                        //20
	"base_asset_id",                     //21
	"quote_asset_id",                    //22
	"quote_asset_chainlink_address_usd", //23
}
var DBColumnsInsertLiquidityPools = []string{
	"uuid",                              //1
	"name",                              //2
	"alternate_name",                    //3
	"pair_address",                      //4
	"chain_id",                          //5
	"exchange_id",                       //6
	"liquidity_pool_type_id",            //7
	"token0_id",                         //8
	"token1_id",                         //9
	"url",                               //10
	"start_block",                       //11
	"latest_block_synced",               //12
	"created_txn_hash",                  //13
	"IsActive",                          //14
	"description",                       //15
	"created_by",                        //16
	"created_at",                        //17
	"updated_by",                        //18
	"updated_at",                        //19
	"base_asset_id",                     //20
	"quote_asset_id",                    //21
	"quote_asset_chainlink_address_usd", //22
}

var TestData1 = LiquidityPool{
	ID:                         utils.Ptr[int](1),                                                    //1
	UUID:                       "01ef85e8-2c26-441e-8c7f-71d79518ad72",                               //2
	Name:                       "PEPE/WETH Uniswap V2",                                               //3
	AlternateName:              "PEPE/WETH Uniswap V2",                                               //4
	PairAddress:                "0xA43fe16908251ee70EF74718545e4FE6C5cCEc9f",                         //5
	ChainID:                    utils.Ptr[int](1),                                                    //6
	ExchangeID:                 utils.Ptr[int](1),                                                    //7
	LiquidityPoolTypeID:        utils.Ptr[int](1),                                                    //8
	Token0ID:                   utils.Ptr[int](1),                                                    //9
	Token1ID:                   utils.Ptr[int](2),                                                    //10
	Url:                        "",                                                                   //11
	StartBlock:                 utils.Ptr[int](17046833),                                             //12
	LatestBlockSynced:          utils.Ptr[int](1),                                                    //13
	CreatedTxnHash:             "0x273894b35d8c30d32e1ffa22ee6aa320cc9f55f2adbba0583594ed47c031f6f6", //14
	IsActive:                   true,                                                                 //15
	Description:                "",                                                                   //16
	CreatedBy:                  "SYSTEM",                                                             //17
	CreatedAt:                  utils.SampleCreatedAtTime,                                            //18
	UpdatedBy:                  "SYSTEM",                                                             //19
	UpdatedAt:                  utils.SampleCreatedAtTime,                                            //20
	BaseAssetID:                utils.Ptr[int](1),                                                    //21
	QuoteAssetID:               utils.Ptr[int](2),                                                    //22
	QuoteAssetChainlinkAddress: "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419",                         //23
}

var TestData2 = LiquidityPool{
	ID:                         utils.Ptr[int](2),                                                             //1
	UUID:                       "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                                        //2
	Name:                       "WOJAK/WETH UNI V2",                                                           //3
	AlternateName:              "WOJAK/WETH UNI V2",                                                           //4
	PairAddress:                "0x0F23d49bC92Ec52FF591D091b3e16c937034496E",                                  //5
	ChainID:                    utils.Ptr[int](1),                                                             //6
	ExchangeID:                 utils.Ptr[int](1),                                                             //7
	LiquidityPoolTypeID:        utils.Ptr[int](1),                                                             //8
	Token0ID:                   utils.Ptr[int](1),                                                             //9
	Token1ID:                   utils.Ptr[int](3),                                                             //10
	Url:                        "https://dexscreener.com/ethereum/0x0f23d49bc92ec52ff591d091b3e16c937034496e", //11
	StartBlock:                 utils.Ptr[int](17069842),                                                      //12
	LatestBlockSynced:          utils.Ptr[int](1),                                                             //13
	CreatedTxnHash:             "0x5bf984bf37a135428f21738421743f4375b7b6512ff03c14fab3a99d8156a3b4",          //14
	IsActive:                   true,                                                                          //15
	Description:                "",                                                                            //16
	CreatedBy:                  "SYSTEM",                                                                      //17
	CreatedAt:                  utils.SampleCreatedAtTime,                                                     //18
	UpdatedBy:                  "SYSTEM",                                                                      //19
	UpdatedAt:                  utils.SampleCreatedAtTime,                                                     //20
	BaseAssetID:                utils.Ptr[int](1),                                                             //21
	QuoteAssetID:               utils.Ptr[int](2),                                                             //22
	QuoteAssetChainlinkAddress: "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419",                                  //23
}
var TestAllData = []LiquidityPool{TestData1, TestData2}

func AddLiquidityPoolToMockRows(mock pgxmock.PgxPoolIface, dataList []LiquidityPool) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                         //1
			data.UUID,                       //2
			data.Name,                       //3
			data.AlternateName,              //4
			data.PairAddress,                //5
			data.ChainID,                    //6
			data.ExchangeID,                 //7
			data.LiquidityPoolTypeID,        //8
			data.Token0ID,                   //9
			data.Token1ID,                   //10
			data.Url,                        //11
			data.StartBlock,                 //12
			data.LatestBlockSynced,          //13
			data.CreatedTxnHash,             //14
			data.IsActive,                   //15
			data.Description,                //16
			data.CreatedBy,                  //17
			data.CreatedAt,                  //18
			data.UpdatedBy,                  //19
			data.UpdatedAt,                  //20
			data.BaseAssetID,                //21
			data.QuoteAssetID,               //22
			data.QuoteAssetChainlinkAddress, //23
		)
	}
	return rows
}

var DBColumnsLiquidityPoolWithTokens = []string{
	"id",                                //1
	"uuid",                              //2
	"name",                              //3
	"alternate_name",                    //4
	"pair_address",                      //5
	"chain_id",                          //6
	"exchange_id",                       //7
	"liquidity_pool_type_id",            //8
	"token0_id",                         //9
	"token1_id",                         //10
	"url",                               //11
	"start_block",                       //12
	"latest_block_synced",               //13
	"created_txn_hash",                  //14
	"IsActive",                          //15
	"description",                       //16
	"created_by",                        //17
	"created_at",                        //18
	"updated_by",                        //19
	"updated_at",                        //20
	"base_asset_id",                     //21
	"quote_asset_id",                    //22
	"quote_asset_chainlink_address_usd", //23
	"id",                                //24 //Token0 Asset
	"uuid",                              //25
	"name",                              //26
	"alternate_name",                    //27
	"cusip",                             //28
	"ticker",                            //29
	"base_asset_id",                     //30
	"quote_asset_id",                    //31
	"description",                       //32
	"asset_type_id",                     //33
	"created_by",                        //34
	"created_at",                        //35
	"updated_by",                        //36
	"updated_at",                        //37
	"chain_id",                          //38
	"category_id",                       //39
	"sub_category_id",                   //40
	"is_default_quote",                  //41
	"ignore_market_data",                //42
	"decimals",                          //43
	"contract_address",                  //44
	"starting_block_number",             //45
	"import_geth",                       //46
	"import_geth_initial",               //47
	"id",                                //48 //Token1 Asset
	"uuid",                              //49
	"name",                              //50
	"alternate_name",                    //51
	"cusip",                             //52
	"ticker",                            //53
	"base_asset_id",                     //54
	"quote_asset_id",                    //55
	"description",                       //56
	"asset_type_id",                     //57
	"created_by",                        //58
	"created_at",                        //59
	"updated_by",                        //60
	"updated_at",                        //61
	"chain_id",                          //62
	"category_id",                       //63
	"sub_category_id",                   //64
	"is_default_quote",                  //65
	"ignore_market_data",                //66
	"decimals",                          //67
	"contract_address",                  //68
	"starting_block_number",             //69
	"import_geth",                       //70
	"import_geth_initial",               //71
}

func AddLiquidityPoolWithTokensToMockRows(mock pgxmock.PgxPoolIface, dataList []LiquidityPoolWithTokens) *pgxmock.Rows {
	rows := mock.NewRows(DBColumnsLiquidityPoolWithTokens)
	for _, data := range dataList {
		rows.AddRow(
			data.LiquidityPool.ID,                         //1
			data.LiquidityPool.UUID,                       //2
			data.LiquidityPool.Name,                       //3
			data.LiquidityPool.AlternateName,              //4
			data.LiquidityPool.PairAddress,                //5
			data.LiquidityPool.ChainID,                    //6
			data.LiquidityPool.ExchangeID,                 //7
			data.LiquidityPool.LiquidityPoolTypeID,        //8
			data.LiquidityPool.Token0ID,                   //9
			data.LiquidityPool.Token1ID,                   //10
			data.LiquidityPool.Url,                        //11
			data.LiquidityPool.StartBlock,                 //12
			data.LiquidityPool.LatestBlockSynced,          //13
			data.LiquidityPool.CreatedTxnHash,             //14
			data.LiquidityPool.IsActive,                   //15
			data.LiquidityPool.Description,                //16
			data.LiquidityPool.CreatedBy,                  //17
			data.LiquidityPool.CreatedAt,                  //18
			data.LiquidityPool.UpdatedBy,                  //19
			data.LiquidityPool.UpdatedAt,                  //20
			data.LiquidityPool.BaseAssetID,                //21
			data.LiquidityPool.QuoteAssetID,               //22
			data.LiquidityPool.QuoteAssetChainlinkAddress, //23
			data.Token0.ID,                                //token0 24
			data.Token0.UUID,                              //25
			data.Token0.Name,                              //26
			data.Token0.AlternateName,                     //27
			data.Token0.Cusip,                             //28
			data.Token0.Ticker,                            //29
			data.Token0.BaseAssetID,                       //30
			data.Token0.QuoteAssetID,                      //31
			data.Token0.Description,                       //32
			data.Token0.AssetTypeID,                       //33
			data.Token0.CreatedBy,                         //34
			data.Token0.CreatedAt,                         //35
			data.Token0.UpdatedBy,                         //36
			data.Token0.UpdatedAt,                         //37
			data.Token0.ChainID,                           //38
			data.Token0.CategoryID,                        //39
			data.Token0.SubCategoryID,                     //40
			data.Token0.IsDefaultQuote,                    //41
			data.Token0.IgnoreMarketData,                  //42
			data.Token0.Decimals,                          //43
			data.Token0.ContractAddress,                   //44
			data.Token0.StartingBlockNumber,               //45
			data.Token0.ImportGeth,                        //46
			data.Token0.ImportGethInitial,                 //457
			data.Token1.ID,                                //token1 48
			data.Token1.UUID,                              //49
			data.Token1.Name,                              //50
			data.Token1.AlternateName,                     //51
			data.Token1.Cusip,                             //52
			data.Token1.Ticker,                            //53
			data.Token1.BaseAssetID,                       //54
			data.Token1.QuoteAssetID,                      //55
			data.Token1.Description,                       //56
			data.Token1.AssetTypeID,                       //57
			data.Token1.CreatedBy,                         //58
			data.Token1.CreatedAt,                         //59
			data.Token1.UpdatedBy,                         //60
			data.Token1.UpdatedAt,                         //61
			data.Token1.ChainID,                           //62
			data.Token1.CategoryID,                        //63
			data.Token1.SubCategoryID,                     //64
			data.Token1.IsDefaultQuote,                    //65
			data.Token1.IgnoreMarketData,                  //66
			data.Token1.Decimals,                          //67
			data.Token1.ContractAddress,                   //68
			data.Token1.StartingBlockNumber,               //69
			data.Token1.ImportGeth,                        //70
			data.Token1.ImportGethInitial,                 //71
		)
	}
	return rows
}

func TestGetLiquidityPool(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []LiquidityPool{targetData}
	liquidityPoolID := targetData.ID
	mockRows := AddLiquidityPoolToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WithArgs(*liquidityPoolID).WillReturnRows(mockRows)
	foundLiquidityPool, err := GetLiquidityPool(mock, liquidityPoolID)
	if err != nil {
		t.Fatalf("an error '%s' in GetLiquidityPool", err)
	}
	if cmp.Equal(*foundLiquidityPool, targetData) == false {
		t.Errorf("Expected LiquidityPool From Method GetLiquidityPool: %v is different from actual %v", foundLiquidityPool, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	liquidityPoolID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WithArgs(liquidityPoolID).WillReturnRows(noRows)
	foundLiquidityPool, err := GetLiquidityPool(mock, &liquidityPoolID)
	if err != nil {
		t.Fatalf("an error '%s' in GetLiquidityPool", err)
	}
	if foundLiquidityPool != nil {
		t.Errorf("Expected LiquidityPool From Method GetLiquidityPool: to be empty but got this: %v", foundLiquidityPool)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	liquidityPoolID := -1
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WithArgs(liquidityPoolID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundLiquidityPool, err := GetLiquidityPool(mock, &liquidityPoolID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetLiquidityPool", err)
	}
	if foundLiquidityPool != nil {
		t.Errorf("Expected LiquidityPool From Method GetLiquidityPool: to be empty but got this: %v", foundLiquidityPool)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveLiquidityPool(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	liquidityPoolID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM liquidity_pools").WithArgs(*liquidityPoolID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveLiquidityPool(mock, liquidityPoolID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveLiquidityPool", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveLiquidityPoolOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	liquidityPoolID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM liquidity_pools").WithArgs(liquidityPoolID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveLiquidityPool(mock, &liquidityPoolID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPools(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddLiquidityPoolToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WillReturnRows(mockRows)
	foundLiquidityPoolList, err := GetLiquidityPools(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetLiquidityPools", err)
	}
	for i, sourceLiquidityPool := range dataList {
		if cmp.Equal(sourceLiquidityPool, foundLiquidityPoolList[i]) == false {
			t.Errorf("Expected LiquidityPool From Method GetLiquidityPools: %v is different from actual %v", sourceLiquidityPool, foundLiquidityPoolList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundLiquidityPoolList, err := GetLiquidityPools(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetLiquidityPools", err)
	}
	if len(foundLiquidityPoolList) != 0 {
		t.Errorf("Expected From Method GetLiquidityPools: to be empty but got this: %v", foundLiquidityPoolList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddLiquidityPoolToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WillReturnRows(mockRows)
	ids := []int{}
	foundLiquidityPoolList, err := GetLiquidityPoolList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetLiquidityPoolList", err)
	}
	for i, sourceLiquidityPool := range dataList {
		if cmp.Equal(sourceLiquidityPool, foundLiquidityPoolList[i]) == false {
			t.Errorf("Expected LiquidityPool From Method GetLiquidityPoolList: %v is different from actual %v", sourceLiquidityPool, foundLiquidityPoolList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{}
	foundLiquidityPoolList, err := GetLiquidityPoolList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetLiquidityPoolList", err)
	}
	if len(foundLiquidityPoolList) != 0 {
		t.Errorf("Expected From Method GetLiquidityPoolList: to be empty but got this: %v", foundLiquidityPoolList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolListByBaseAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	testData1LiquidityPoolWithToken := LiquidityPoolWithTokens{}
	testData1LiquidityPoolWithToken.Token0 = asset.TestData1
	testData1LiquidityPoolWithToken.LiquidityPool = TestData1
	testData1LiquidityPoolWithToken.Token1 = asset.TestData2
	dataList := []LiquidityPoolWithTokens{testData1LiquidityPoolWithToken}
	baseAssetID := TestData1.BaseAssetID
	mockRows := AddLiquidityPoolWithTokensToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WithArgs(*baseAssetID).WillReturnRows(mockRows)
	foundLiquidityPoolList, err := GetLiquidityPoolListByBaseAssetID(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetLiquidityPoolListByBaseAssetID", err)
	}
	for i, sourceLiquidityPool := range dataList {
		if cmp.Equal(sourceLiquidityPool, foundLiquidityPoolList[i]) == false {
			t.Errorf("Expected LiquidityPool From Method GetLiquidityPoolListByBaseAssetID: %v is different from actual %v", sourceLiquidityPool, foundLiquidityPoolList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolListByBaseAssetIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WithArgs(*baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundLiquidityPoolList, err := GetLiquidityPoolListByBaseAssetID(mock, baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetLiquidityPoolListByBaseAssetID", err)
	}
	if len(foundLiquidityPoolList) != 0 {
		t.Errorf("Expected From Method GetLiquidityPoolListByBaseAssetID: to be empty but got this: %v", foundLiquidityPoolList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddLiquidityPoolToMockRows(mock, dataList)
	uuidList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WithArgs(pq.Array(uuidList)).WillReturnRows(mockRows)
	foundLiquidityPoolList, err := GetLiquidityPoolsByUUIDs(mock, uuidList)
	if err != nil {
		t.Fatalf("an error '%s' in GetLiquidityPoolsByUUIDs", err)
	}
	for i, sourceLiquidityPool := range dataList {
		if cmp.Equal(sourceLiquidityPool, foundLiquidityPoolList[i]) == false {
			t.Errorf("Expected LiquidityPool From Method GetLiquidityPoolsByUUIDs: %v is different from actual %v", sourceLiquidityPool, foundLiquidityPoolList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuidList := []string{"test", "test2"}
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WithArgs(pq.Array(uuidList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundLiquidityPoolList, err := GetLiquidityPoolsByUUIDs(mock, uuidList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetLiquidityPoolsByUUIDs", err)
	}
	if len(foundLiquidityPoolList) != 0 {
		t.Errorf("Expected From Method GetLiquidityPoolsByUUIDs: to be empty but got this: %v", foundLiquidityPoolList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateLiquidityPool(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE liquidity_pools").WithArgs(
		targetData.Name,                       //1
		targetData.AlternateName,              //2
		targetData.PairAddress,                //3
		targetData.ChainID,                    //4
		targetData.ExchangeID,                 //5
		targetData.LiquidityPoolTypeID,        //6
		targetData.Token0ID,                   //7
		targetData.Token1ID,                   //8
		targetData.Url,                        //9
		targetData.StartBlock,                 //10
		targetData.LatestBlockSynced,          //11
		targetData.CreatedTxnHash,             //12
		targetData.IsActive,                   //13
		targetData.Description,                //14
		targetData.UpdatedBy,                  //15
		targetData.BaseAssetID,                //16
		targetData.QuoteAssetID,               //17
		targetData.QuoteAssetChainlinkAddress, //18
		targetData.ID,                         //19
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateLiquidityPool(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateLiquidityPool", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateLiquidityPoolOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateLiquidityPool(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateLiquidityPoolOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE liquidity_pools").WithArgs(
		targetData.Name,                       //1
		targetData.AlternateName,              //2
		targetData.PairAddress,                //3
		targetData.ChainID,                    //4
		targetData.ExchangeID,                 //5
		targetData.LiquidityPoolTypeID,        //6
		targetData.Token0ID,                   //7
		targetData.Token1ID,                   //8
		targetData.Url,                        //9
		targetData.StartBlock,                 //10
		targetData.LatestBlockSynced,          //11
		targetData.CreatedTxnHash,             //12
		targetData.IsActive,                   //13
		targetData.Description,                //14
		targetData.UpdatedBy,                  //15
		targetData.BaseAssetID,                //16
		targetData.QuoteAssetID,               //17
		targetData.QuoteAssetChainlinkAddress, //18
		targetData.ID,                         //19
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateLiquidityPool(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertLiquidityPool(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO liquidity_pools").WithArgs(
		targetData.Name,                       //1
		targetData.AlternateName,              //2
		targetData.PairAddress,                //3
		targetData.ChainID,                    //4
		targetData.ExchangeID,                 //5
		targetData.LiquidityPoolTypeID,        //6
		targetData.Token0ID,                   //7
		targetData.Token1ID,                   //8
		targetData.Url,                        //9
		targetData.StartBlock,                 //10
		targetData.LatestBlockSynced,          //11
		targetData.CreatedTxnHash,             //12
		targetData.IsActive,                   //13
		targetData.Description,                //14
		targetData.CreatedBy,                  //15
		targetData.BaseAssetID,                //16
		targetData.QuoteAssetID,               //17
		targetData.QuoteAssetChainlinkAddress, //18
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit()
	liquidityPoolID, newUUID, err := InsertLiquidityPool(mock, &targetData)
	if liquidityPoolID < 0 {
		t.Fatalf("liquidityPoolID should not be negative ID: %d", liquidityPoolID)
	}
	if newUUID == "" {
		t.Fatalf("newUUID should not be empty string: %s", newUUID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertLiquidityPool", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertLiquidityPoolOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Token0ID = nil
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO liquidity_pools").WithArgs(
		targetData.Name,                       //1
		targetData.AlternateName,              //2
		targetData.PairAddress,                //3
		targetData.ChainID,                    //4
		targetData.ExchangeID,                 //5
		targetData.LiquidityPoolTypeID,        //6
		targetData.Token0ID,                   //7
		targetData.Token1ID,                   //8
		targetData.Url,                        //9
		targetData.StartBlock,                 //10
		targetData.LatestBlockSynced,          //11
		targetData.CreatedTxnHash,             //12
		targetData.IsActive,                   //13
		targetData.Description,                //14
		targetData.CreatedBy,                  //15
		targetData.BaseAssetID,                //16
		targetData.QuoteAssetID,               //17
		targetData.QuoteAssetChainlinkAddress, //18
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	liquidityPoolID, newUUID, err := InsertLiquidityPool(mock, &targetData)
	if liquidityPoolID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error liquidityPoolID: %d", liquidityPoolID)
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

func TestInsertLiquidityPoolOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO liquidity_pools").WithArgs(
		targetData.Name,                       //1
		targetData.AlternateName,              //2
		targetData.PairAddress,                //3
		targetData.ChainID,                    //4
		targetData.ExchangeID,                 //5
		targetData.LiquidityPoolTypeID,        //6
		targetData.Token0ID,                   //7
		targetData.Token1ID,                   //8
		targetData.Url,                        //9
		targetData.StartBlock,                 //10
		targetData.LatestBlockSynced,          //11
		targetData.CreatedTxnHash,             //12
		targetData.IsActive,                   //13
		targetData.Description,                //14
		targetData.CreatedBy,                  //15
		targetData.BaseAssetID,                //16
		targetData.QuoteAssetID,               //17
		targetData.QuoteAssetChainlinkAddress, //18
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	liquidityPoolID, newUUID, err := InsertLiquidityPool(mock, &targetData)
	if liquidityPoolID >= 0 {
		t.Fatalf("Expecting -1 for liquidityPoolID because of error liquidityPoolID: %d", liquidityPoolID)
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

func TestInsertLiquidityPools(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"liquidity_pools"}, DBColumnsInsertLiquidityPools)
	targetData := TestAllData
	err = InsertLiquidityPools(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertLiquidityPools", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertLiquidityPoolsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"liquidity_pools"}, DBColumnsInsertLiquidityPools).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertLiquidityPools(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddLiquidityPoolToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"chain_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WillReturnRows(mockRows)
	foundChains, err := GetLiquidityPoolListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetLiquidityPoolListByPagination", err)
	}
	testChains := dataList
	for i, foundChain := range foundChains {
		if cmp.Equal(foundChain, testChains[i]) == false {
			t.Errorf("Expected Chain From Method GetLiquidityPoolListByPagination: %v is different from actual %v", foundChain, testChains[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetLiquidityPoolListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"chain_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM liquidity_pools").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChains, err := GetLiquidityPoolListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetLiquidityPoolListByPagination", err)
	}
	if len(foundChains) != 0 {
		t.Errorf("Expected From Method GetLiquidityPoolListByPagination: to be empty but got this: %v", foundChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalLiquidityPoolCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM liquidity_pools").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalLiquidityPoolCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalLiquidityPoolCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalLiquidityPoolCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM liquidity_pools").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalLiquidityPoolCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalLiquidityPoolCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalLiquidityPoolCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
