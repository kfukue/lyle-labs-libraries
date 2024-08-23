package asset

import (
	"strings"
	"time"
)

// Asset
type Asset struct {
	ID                  *int      `json:"id" db:"id"`                                      //1
	UUID                string    `json:"uuid" db:"uuid"`                                  //2
	Name                string    `json:"name" db:"name"`                                  //3
	AlternateName       string    `json:"alternateName" db:"alternate_name"`               //4
	Cusip               string    `json:"cusip" db:"cusip"`                                //5
	Ticker              string    `json:"ticker" db:"ticker"`                              //6
	BaseAssetID         *int      `json:"baseAssetId" db:"base_asset_id"`                  //7
	QuoteAssetID        *int      `json:"quoteAssetId" db:"quote_asset_id"`                //8
	Description         string    `json:"description" db:"description"`                    //9
	AssetTypeID         *int      `json:"assetTypeId" db:"asset_type_id"`                  //10
	CreatedBy           string    `json:"createdBy" db:"created_by"`                       //11
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`                       //12
	UpdatedBy           string    `json:"updatedBy" db:"updated_by"`                       //13
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`                       //14
	ChainID             *int      `json:"chainId" db:"chain_id"`                           //15
	CategoryID          *int      `json:"categoryId" db:"category_id"`                     //16
	SubCategoryID       *int      `json:"subCategoryId" db:"sub_category_id"`              //17
	IsDefaultQuote      *bool     `json:"isDefaultQuote" db:"is_default_quote"`            //18
	IgnoreMarketData    *bool     `json:"ignoreMarketData" db:"ignore_market_data"`        //19
	Decimals            *int      `json:"decimals" db:"decimals"`                          //20
	ContractAddress     string    `json:"contractAddress" db:"contract_address"`           //21
	StartingBlockNumber *uint64   `json:"startingBlockNumber" db:"starting_block_number"`  //22
	ImportGeth          *bool     `json:"importGeth" db:"import_geth"`                     //23
	ImportGethInitial   *bool     `json:"importGethInitial" db:"import_geth_initial"`      //24
	ChainlinkUSDAddress *string   `json:"chainlinkUSDAddress" db:"chainlink_usd_address"`  //25
	ChainlinkUSDChainID *int      `json:"chainlinkUSDChainId" db:"chainlink_usd_chain_id"` //26
}

// Asset
type AssetWithSources struct {
	Asset
	SourceID         *int   `json:"sourceId" db:"assetSources.source_id"`
	SourceIdentifier string `json:"sourceIdentifier" db:"assetSources.source_identifier"`
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
