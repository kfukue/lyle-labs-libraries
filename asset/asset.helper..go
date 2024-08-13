package asset

import (
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"id",                    //1
	"uuid",                  //2
	"name",                  //3
	"alternate_name",        //4
	"cusip",                 //5
	"ticker",                //6
	"base_asset_id",         //7
	"quote_asset_id",        //8
	"description",           //9
	"asset_type_id",         //10
	"created_by",            //11
	"created_at",            //12
	"updated_by",            //13
	"updated_at",            //14
	"chain_id",              //15
	"category_id",           //16
	"sub_category_id",       //17
	"is_default_quote",      //18
	"ignore_market_data",    //19
	"decimals",              //20
	"contract_address",      //21
	"starting_block_number", //22
	"import_geth",           //23
	"import_geth_initial",   //24
}

var DBColumnsInsertAssets = []string{
	"uuid",                  //1
	"name",                  //2
	"alternate_name",        //3
	"cusip",                 //4
	"ticker",                //5
	"base_asset_id",         //6
	"quote_asset_id",        //7
	"description",           //8
	"asset_type_id",         //9
	"created_by",            //10
	"created_at",            //11
	"updated_by",            //12
	"updated_at",            //13
	"chain_id",              //14
	"category_id",           //15
	"sub_category_id",       //16
	"is_default_quote",      //17
	"ignore_market_data",    //18
	"decimals",              //19
	"contract_address",      //20
	"starting_block_number", //21
	"import_geth",           //22
	"import_geth_initial",   //23
}

var DBColumnsAssetWithSources = []string{
	"id",                             //1
	"uuid",                           //2
	"name",                           //3
	"alternate_name",                 //4
	"cusip",                          //5
	"ticker",                         //6
	"base_asset_id",                  //7
	"quote_asset_id",                 //8
	"description",                    //9
	"asset_type_id",                  //10
	"created_by",                     //11
	"created_at",                     //12
	"updated_by",                     //13
	"updated_at",                     //14
	"chain_id",                       //15
	"category_id",                    //16
	"sub_category_id",                //17
	"is_default_quote",               //18
	"ignore_market_data",             //19
	"decimals",                       //20
	"contract_address",               //21
	"starting_block_number",          //22
	"import_geth",                    //23
	"import_geth_initial",            //24
	"assetSources.source_id",         //25
	"assetSources.source_identifier", //26
}

var TestData1 = Asset{
	ID:                  utils.Ptr[int](1),                            //1
	UUID:                "880607ab-2833-4ad7-a231-b983a61c7b39",       //2
	Name:                "ETHER",                                      //3
	AlternateName:       "Ether",                                      //4
	Cusip:               "",                                           //5
	Ticker:              "ETH",                                        //6
	BaseAssetID:         utils.Ptr[int](1),                            //7
	QuoteAssetID:        utils.Ptr[int](2),                            //8
	Description:         "",                                           //9
	AssetTypeID:         utils.Ptr[int](1),                            //10
	CreatedBy:           "SYSTEM",                                     //11
	CreatedAt:           utils.SampleCreatedAtTime,                    //12
	UpdatedBy:           "SYSTEM",                                     //13
	UpdatedAt:           utils.SampleCreatedAtTime,                    //14
	ChainID:             utils.Ptr[int](2),                            //15
	CategoryID:          utils.Ptr[int](27),                           //16
	SubCategoryID:       utils.Ptr[int](10),                           //17
	IsDefaultQuote:      utils.Ptr[bool](true),                        //18
	IgnoreMarketData:    utils.Ptr[bool](false),                       //19
	Decimals:            utils.Ptr[int](1),                            //20
	ContractAddress:     "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", //21
	StartingBlockNumber: utils.Ptr[uint64](1),                         //22
	ImportGeth:          nil,                                          //23
	ImportGethInitial:   nil,                                          //24
}

var TestData2 = Asset{
	ID:                  utils.Ptr[int](2),                       //1
	UUID:                "880607ab-2833-4ad7-a231-b983a61c7b334", //2
	Name:                "BTC",                                   //3
	AlternateName:       "Bitcoin",                               //4
	Cusip:               "",                                      //5
	Ticker:              "BTC",                                   //6
	BaseAssetID:         utils.Ptr[int](2),                       //7
	QuoteAssetID:        utils.Ptr[int](1),                       //8
	Description:         "",                                      //9
	AssetTypeID:         utils.Ptr[int](1),                       //10
	CreatedBy:           "SYSTEM",                                //11
	CreatedAt:           utils.SampleCreatedAtTime,               //12
	UpdatedBy:           "SYSTEM",                                //13
	UpdatedAt:           utils.SampleCreatedAtTime,               //14
	ChainID:             utils.Ptr[int](1),                       //15
	CategoryID:          utils.Ptr[int](28),                      //16
	SubCategoryID:       utils.Ptr[int](20),                      //17
	IsDefaultQuote:      utils.Ptr[bool](true),                   //18
	IgnoreMarketData:    utils.Ptr[bool](false),                  //19
	Decimals:            utils.Ptr[int](1),                       //20
	ContractAddress:     "0x",                                    //21
	StartingBlockNumber: utils.Ptr[uint64](1),                    //22
	ImportGeth:          nil,                                     //23
	ImportGethInitial:   nil,                                     //24
}
var TestAllData = []Asset{TestData1, TestData2}

func AddAssetToMockRows(mock pgxmock.PgxPoolIface, dataList []Asset) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                  //1
			data.UUID,                //2
			data.Name,                //3
			data.AlternateName,       //4
			data.Cusip,               //5
			data.Ticker,              //6
			data.BaseAssetID,         //7
			data.QuoteAssetID,        //8
			data.Description,         //9
			data.AssetTypeID,         //10
			data.CreatedBy,           //11
			data.CreatedAt,           //12
			data.UpdatedBy,           //13
			data.UpdatedAt,           //14
			data.ChainID,             //15
			data.CategoryID,          //16
			data.SubCategoryID,       //17
			data.IsDefaultQuote,      //18
			data.IgnoreMarketData,    //19
			data.Decimals,            //20
			data.ContractAddress,     //21
			data.StartingBlockNumber, //22
			data.ImportGeth,          //23
			data.ImportGethInitial,   //24
		)
	}
	return rows
}

var dataAssetWithSources1 = AssetWithSources{
	Asset:            TestData1,
	SourceID:         utils.Ptr[int](1),
	SourceIdentifier: "ETH",
}

var dataAssetWithSources2 = AssetWithSources{
	Asset:            TestData2,
	SourceID:         utils.Ptr[int](1),
	SourceIdentifier: "XBT",
}

var TestAllDataAssetWithSources = []AssetWithSources{dataAssetWithSources1, dataAssetWithSources2}

func AddAssetWithSourcesToMockRows(mock pgxmock.PgxPoolIface, dataList []AssetWithSources) *pgxmock.Rows {
	rows := mock.NewRows(DBColumnsAssetWithSources)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                  //1
			data.UUID,                //2
			data.Name,                //3
			data.AlternateName,       //4
			data.Cusip,               //5
			data.Ticker,              //6
			data.BaseAssetID,         //7
			data.QuoteAssetID,        //8
			data.Description,         //9
			data.AssetTypeID,         //10
			data.CreatedBy,           //11
			data.CreatedAt,           //12
			data.UpdatedBy,           //13
			data.UpdatedAt,           //14
			data.ChainID,             //15
			data.CategoryID,          //16
			data.SubCategoryID,       //17
			data.IsDefaultQuote,      //18
			data.IgnoreMarketData,    //19
			data.Decimals,            //20
			data.ContractAddress,     //21
			data.StartingBlockNumber, //22
			data.ImportGeth,          //23
			data.ImportGethInitial,   //24
			data.SourceID,            //25
			data.SourceIdentifier,    //26
		)
	}
	return rows
}
