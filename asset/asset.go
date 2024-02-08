package asset

import (
	"strings"
	"time"
)

// Asset
type Asset struct {
	ID                  *int      `json:"id"`
	UUID                string    `json:"uuid"`
	Name                string    `json:"name"`
	AlternateName       string    `json:"alternateName"`
	Cusip               string    `json:"cusip"`
	Ticker              string    `json:"ticker"`
	BaseAssetID         *int      `json:"baseAssetId"`
	QuoteAssetID        *int      `json:"quoteAssetId"`
	Description         string    `json:"description"`
	AssetTypeID         *int      `json:"assetTypeId"`
	CreatedBy           string    `json:"createdBy"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedBy           string    `json:"updatedBy"`
	UpdatedAt           time.Time `json:"updatedAt"`
	ChainID             *int      `json:"chainId"`
	CategoryID          *int      `json:"categoryId"`
	SubCategoryID       *int      `json:"subCategoryId"`
	IsDefaultQuote      *bool     `json:"isDefaultQuote"`
	IgnoreMarketData    *bool     `json:"ignoreMarketData"`
	Decimals            *int      `json:"decimals"`
	ContractAddress     string    `json:"contractAddress"`
	StartingBlockNumber *uint64   `json:"startingBlockNumber"`
	ImportGeth          *bool     `json:"importGeth"`
	ImportGethInitial   *bool     `json:"importGethInitial"`
}

// Asset
type AssetWithSources struct {
	Asset
	SourceID         *int   `json:"sourceId"`
	SourceIdentifier string `json:"sourceIdentifier"`
}

// lowercase for the contract address
func CreateLookupByContractAddressFromAssetList(assetList []Asset) map[string]*Asset {
	assetLookup := make(map[string]*Asset)
	for _, asset := range assetList {
		if asset.ContractAddress != "" {
			loweredAddress := strings.ToLower(asset.ContractAddress)
			_, foundAsset := assetLookup[loweredAddress]
			if !foundAsset {
				assetLookup[loweredAddress] = &asset
			}
		}
	}
	return assetLookup
}
