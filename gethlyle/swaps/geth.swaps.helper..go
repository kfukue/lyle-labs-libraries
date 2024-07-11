package gethlyleswaps

import (
	"github.com/kfukue/lyle-labs-libraries/utils"
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
