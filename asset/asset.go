package asset

import (
	"time"
)

// Asset
type Asset struct {
	ID                  *int      `json:"id  db:"id"`
	UUID                string    `json:"uuid db:"uuid"`
	Name                string    `json:"name db:"name"`
	AlternateName       string    `json:"alternateName db:"alternate_name"`
	Cusip               string    `json:"cusip db:"cusip"`
	Ticker              string    `json:"ticker db:"ticker"`
	BaseAssetID         *int      `json:"baseAssetId db:"base_asset_id"`
	QuoteAssetID        *int      `json:"quoteAssetId db:"quote_asset_id"`
	Description         string    `json:"description db:"description"`
	AssetTypeID         *int      `json:"assetTypeId db:vasset_type_id"`
	CreatedBy           string    `json:"createdBy db:"created_by"`
	CreatedAt           time.Time `json:"createdAt db:"created_at"`
	UpdatedBy           string    `json:"updatedBy db:"udpated_by"`
	UpdatedAt           time.Time `json:"updatedAt db:"updated_at"`
	ChainID             *int      `json:"chainId" db:"chain_id`
	CategoryID          *int      `json:"categoryId" db:"category_id`
	SubCategoryID       *int      `json:"subCategoryId" db:"sub_category_id`
	IsDefaultQuote      *bool     `json:"isDefaultQuote" db:"is_default_quote`
	IgnoreMarketData    *bool     `json:"ignoreMarketData" db:"ignore_market_data`
	Decimals            *int      `json:"decimals" db:"decimals`
	ContractAddress     string    `json:"contractAddress" db:"contract_address`
	StartingBlockNumber *uint64   `json:"startingBlock" db:"starting_block_number`
}

// Asset
type AssetWithSources struct {
	Asset
	SourceID         *int   `json:"sourceId"`
	SourceIdentifier string `json:"sourceIdentifier"`
}
