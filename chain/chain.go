package chain

import (
	"time"
)

// Asset
type Chain struct {
	ID               *int      `json:"id" db:"id"`                               //1
	UUID             string    `json:"uuid" db:"uuid"`                           //2
	BaseAssetID      *int      `json:"baseAssetId" db:"base_asset_id"`           //3
	Name             string    `json:"name" db:"name" db:"name"`                 //4
	AlternateName    string    `json:"alternateName" db:"alternate_name"`        //5
	Address          string    `json:"address" db:"address"`                     //6
	ChainTypeID      *int      `json:"chainTypeId" db:"chain_type_id"`           //7
	Description      string    `json:"description" db:"description"`             //8
	CreatedBy        string    `json:"createdBy" db:"created_by"`                //9
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`                //10
	UpdatedBy        string    `json:"updatedBy" db:"updated_by"`                //11
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`                //12
	RpcURL           string    `json:"rpcUrl" db:"rpc_url"`                      //13
	ChainID          *int      `json:"chainId" db:"chain_id"`                    //14
	BlockExplorerURL string    `json:"blockExplorerUrl" db:"block_explorer_url"` //15
	RpcURLDev        string    `json:"rpcUrlDev" db:"rpc_url_dev"`               //16
	RpcURLProd       string    `json:"rpcUrlProd" db:"rpc_url_prod"`             //17
	RpcURLArchive    string    `json:"rpcUrlArchive" db:"rpc_url_archive"`       //18
}
