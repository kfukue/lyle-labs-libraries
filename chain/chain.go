package chain

import (
	"time"
)

// Asset
type Chain struct {
	ID               *int      `json:"id" db:"id"`
	UUID             string    `json:"uuid" db:"uuid"`
	BaseAssetID      *int      `json:"baseAssetId" db:"base_asset_id"`
	Name             string    `json:"name" db:"name" db:"name"`
	AlternateName    string    `json:"alternateName" db:"alternate_name"`
	Address          string    `json:"address" db:"address"`
	ChainTypeID      *int      `json:"chainTypeId" db:"chain_type_id"`
	Description      string    `json:"description" db:"description"`
	CreatedBy        string    `json:"createdBy" db:"created_by"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
	UpdatedBy        string    `json:"updatedBy" db:"updated_by"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`
	RpcURL           string    `json:"rpcUrl" db:"rpc_url"`
	ChainID          *int      `json:"chainId" db:"chain_id"`
	BlockExplorerURL string    `json:"blockExplorerUrl" db:"block_explorer_url"`
	RpcURLDev        string    `json:"rpcUrlDev" db:"rpc_url_dev"`
	RpcURLProd       string    `json:"rpcUrlProd" db:"rpc_url_prod"`
	RpcURLArchive    string    `json:"rpcUrlArchive" db:"rpc_url_archive"`
}
