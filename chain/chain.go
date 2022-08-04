package chain

import (
	"time"
)

// Asset
type Chain struct {
	ID               *int      `json:"id"`
	UUID             string    `json:"uuid"`
	BaseAssetID      *int      `json:"baseAssetId"`
	Name             string    `json:"name"`
	AlternateName    string    `json:"alternateName"`
	Address          string    `json:"address"`
	ChainTypeID      *int      `json:"chainTypeId"`
	Description      string    `json:"description"`
	CreatedBy        string    `json:"createdBy"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedBy        string    `json:"updatedBy"`
	UpdatedAt        time.Time `json:"updatedAt"`
	RpcURL           string    `json:"rpcUrl"`
	ChainID          *int      `json:"chainId"`
	BlockExplorerURL string    `json:"blockExplorerUrl"`
}
