package assetchain

import (
	"time"
)

type AssetChain struct {
	AssetID                          *int      `json:"assetId" db:"asset_id"`
	ChainID                          *int      `json:"chainId" db:"chain_id"`
	ChainlinkDataFeedContractAddress string    `json:"chainlinkDataFeedContractAddress" db:"chainlink_data_feed_contract_address"`
	CreatedBy                        string    `json:"createdBy" db:"created_by"`
	CreatedAt                        time.Time `json:"createdAt" db:"created_at"`
	UpdatedBy                        string    `json:"updatedBy" db:"updated_by"`
	UpdatedAt                        time.Time `json:"updatedAt" db:"updated_at"`
}

type AssetChainIface interface {
	GetAssetChain(assetID, chainID *int) (*AssetChain, error)
	GetAssetChainList(assetIDs, chainIDs []int) ([]AssetChain, error)
	InsertAssetChain(feed *AssetChain) error
	UpdateAssetChain(feed *AssetChain) error
	RemoveAssetChain(assetID, chainID *int) error
}
