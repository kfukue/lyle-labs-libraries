package assetchainlinkfeed

import (
	"time"
)

type AssetChainLinkFeed struct {
	AssetID                          *int      `json:"asset_id" db:"asset_id"`
	ChainID                          *int      `json:"chain_id" db:"chain_id"`
	ChainlinkDataFeedContractAddress string    `json:"chainlink_data_feed_contract_address" db:"chainlink_data_feed_contract_address"`
	CreatedBy                        string    `json:"created_by" db:"created_by"`
	CreatedAt                        time.Time `json:"created_at" db:"created_at"`
	UpdatedBy                        string    `json:"updated_by" db:"updated_by"`
	UpdatedAt                        time.Time `json:"updated_at" db:"updated_at"`
}

type AssetChainLinkFeedIface interface {
	GetAssetChainLinkFeed(assetID, chainID *int) (*AssetChainLinkFeed, error)
	GetAssetChainLinkFeedList(assetIDs, chainIDs []int) ([]AssetChainLinkFeed, error)
	InsertAssetChainLinkFeed(feed *AssetChainLinkFeed) error
	UpdateAssetChainLinkFeed(feed *AssetChainLinkFeed) error
	RemoveAssetChainLinkFeed(assetID, chainID *int) error
}
