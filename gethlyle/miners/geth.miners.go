package gethlyleminers

import (
	"time"
)

type GethMiner struct {
	ID                  *int      `json:"id" db:"id"`                                     //1
	UUID                string    `json:"uuid" db:"uuid"`                                 //2
	Name                string    `json:"name" db:"name"`                                 //3
	AlternateName       string    `json:"alternateName" db:"alternate_name"`              //4
	ChainID             *int      `json:"chainId" db:"chain_id"`                          //5
	ExchangeID          *int      `json:"exchangeId" db:"exchange_id"`                    //6
	StartingBlockNumber *int      `json:"startingBlockNumber" db:"starting_block_number"` //7
	CreatedTxnHash      string    `json:"createdTxnHash" db:"created_txn_hash"`           //8
	LastBlockNumber     *uint64   `json:"lastBlockNumber" db:"last_block_number"`         //9
	ContractAddress     string    `json:"contractAddress" db:"contract_address"`          //10
	ContractAddressID   *int      `json:"contractAddressId" db:"contract_address_id"`     //11
	DeveloperAddress    string    `json:"developerAddress" db:"developer_address"`        //12
	DeveloperAddressID  *int      `json:"developerAddressId" db:"developer_address_id"`   //13
	MiningAssetID       *int      `json:"miningAssetId" db:"mining_asset_id"`             //14
	Description         string    `json:"description" db:"description"`                   //15
	CreatedBy           string    `json:"createdBy" db:"created_by"`                      //16
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`                      //17
	UpdatedBy           string    `json:"updatedBy" db:"updated_by"`                      //18
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`                      //19
}
