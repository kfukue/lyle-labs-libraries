package gethlyleminers

import (
	"time"
)

type GethMiner struct {
	ID                  *int      `json:"id"`
	UUID                string    `json:"uuid"`
	Name                string    `json:"name"`
	AlternateName       string    `json:"alternateName"`
	ChainID             *int      `json:"chainId"`
	ExchangeID          *int      `json:"exchangeId"`
	StartingBlockNumber *int      `json:"startingBlockNumber"`
	CreatedTxnHash      string    `json:"createdTxnHash"`
	LastBlockNumber     *uint64   `json:"lastBlockNumber"`
	ContractAddress     string    `json:"contractAddress"`
	ContractAddressID   *int      `json:"contractAddressId"`
	DeveloperAddress    string    `json:"developerAddress"`
	DeveloperAddressID  *int      `json:"developerAddressId"`
	MiningAssetID       *int      `json:"miningAssetId"`
	Description         string    `json:"description"`
	CreatedBy           string    `json:"createdBy"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedBy           string    `json:"updatedBy"`
	UpdatedAt           time.Time `json:"updatedAt"`
}
