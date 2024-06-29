package asset

import (
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"id",
	"uuid",
	"name",
	"alternate_name",
	"cusip",
	"ticker",
	"base_asset_id",
	"quote_asset_id",
	"description",
	"asset_type_id",
	"created_by",
	"created_at",
	"updated_by",
	"updated_at",
	"chain_id",
	"category_id",
	"sub_category_id",
	"is_default_quote",
	"ignore_market_data",
	"decimals",
	"contract_address",
	"starting_block_number",
	"import_geth",
	"import_geth_initial",
}

var DBColumnsAssetWithSources = []string{
	"id",
	"uuid",
	"name",
	"alternate_name",
	"cusip",
	"ticker",
	"base_asset_id",
	"quote_asset_id",
	"description",
	"asset_type_id",
	"created_by",
	"created_at",
	"updated_by",
	"updated_at",
	"chain_id",
	"category_id",
	"sub_category_id",
	"is_default_quote",
	"ignore_market_data",
	"decimals",
	"contract_address",
	"starting_block_number",
	"import_geth",
	"import_geth_initial",
	"assetSources.source_id",
	"assetSources.source_identifier",
}

var TestData1 = Asset{
	ID:                  utils.Ptr[int](1),
	UUID:                "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:                "ETHER",
	AlternateName:       "Ether",
	Cusip:               "",
	Ticker:              "ETH",
	BaseAssetID:         utils.Ptr[int](1),
	QuoteAssetID:        utils.Ptr[int](2),
	Description:         "",
	AssetTypeID:         utils.Ptr[int](1),
	CreatedBy:           "SYSTEM",
	CreatedAt:           utils.SampleCreatedAtTime,
	UpdatedBy:           "SYSTEM",
	UpdatedAt:           utils.SampleCreatedAtTime,
	ChainID:             utils.Ptr[int](2),
	CategoryID:          utils.Ptr[int](27),
	SubCategoryID:       utils.Ptr[int](10),
	IsDefaultQuote:      utils.Ptr[bool](true),
	IgnoreMarketData:    utils.Ptr[bool](false),
	Decimals:            utils.Ptr[int](1),
	ContractAddress:     "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
	StartingBlockNumber: utils.Ptr[uint64](1),
	ImportGeth:          nil,
	ImportGethInitial:   nil,
}

var TestData2 = Asset{
	ID:                  utils.Ptr[int](2),
	UUID:                "880607ab-2833-4ad7-a231-b983a61c7b334",
	Name:                "BTC",
	AlternateName:       "Bitcoin",
	Cusip:               "",
	Ticker:              "BTC",
	BaseAssetID:         utils.Ptr[int](2),
	QuoteAssetID:        utils.Ptr[int](1),
	Description:         "",
	AssetTypeID:         utils.Ptr[int](1),
	CreatedBy:           "SYSTEM",
	CreatedAt:           utils.SampleCreatedAtTime,
	UpdatedBy:           "SYSTEM",
	UpdatedAt:           utils.SampleCreatedAtTime,
	ChainID:             utils.Ptr[int](1),
	CategoryID:          utils.Ptr[int](28),
	SubCategoryID:       utils.Ptr[int](20),
	IsDefaultQuote:      utils.Ptr[bool](true),
	IgnoreMarketData:    utils.Ptr[bool](false),
	Decimals:            utils.Ptr[int](1),
	ContractAddress:     "0x",
	StartingBlockNumber: utils.Ptr[uint64](1),
	ImportGeth:          nil,
	ImportGethInitial:   nil,
}
var TestAllData = []Asset{TestData1, TestData2}

func AddAssetToMockRows(mock pgxmock.PgxPoolIface, dataList []Asset) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,
			data.UUID,
			data.Name,
			data.AlternateName,
			data.Cusip,
			data.Ticker,
			data.BaseAssetID,
			data.QuoteAssetID,
			data.Description,
			data.AssetTypeID,
			data.CreatedBy,
			data.CreatedAt,
			data.UpdatedBy,
			data.UpdatedAt,
			data.ChainID,
			data.CategoryID,
			data.SubCategoryID,
			data.IsDefaultQuote,
			data.IgnoreMarketData,
			data.Decimals,
			data.ContractAddress,
			data.StartingBlockNumber,
			data.ImportGeth,
			data.ImportGethInitial,
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
			data.ID,
			data.UUID,
			data.Name,
			data.AlternateName,
			data.Cusip,
			data.Ticker,
			data.BaseAssetID,
			data.QuoteAssetID,
			data.Description,
			data.AssetTypeID,
			data.CreatedBy,
			data.CreatedAt,
			data.UpdatedBy,
			data.UpdatedAt,
			data.ChainID,
			data.CategoryID,
			data.SubCategoryID,
			data.IsDefaultQuote,
			data.IgnoreMarketData,
			data.Decimals,
			data.ContractAddress,
			data.StartingBlockNumber,
			data.ImportGeth,
			data.ImportGethInitial,
			data.SourceID,
			data.SourceIdentifier,
		)
	}
	return rows
}
