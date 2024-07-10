package gethlyletransactions

import (
	"time"

	"github.com/shopspring/decimal"
)

type GethTransaction struct {
	ID                          *int             `json:"id" db:"id"`                                                      //1
	UUID                        string           `json:"uuid" db:"uuid"`                                                  //2
	ChainID                     *int             `json:"chainId" db:"chain_id"`                                           //3
	ExchangeID                  *int             `json:"exchangeId" db:"exchange_id"`                                     //4
	BlockNumber                 *uint64          `json:"blockNumber" db:"block_number"`                                   //5
	IndexNumber                 *uint            `json:"indexNumber" db:"index_number"`                                   //6
	TxnDate                     time.Time        `json:"txnDate" db:"txn_date"`                                           //7
	TxnHash                     string           `json:"txnHash" db:"txn_hash"`                                           //8
	FromAddress                 string           `json:"senderAddress" db:"from_address"`                                 //9
	FromAddressID               *int             `json:"senderAddressID" db:"from_address_id"`                            //10
	ToAddress                   string           `json:"toAddress" db:"to_address"`                                       //11
	ToAddressID                 *int             `json:"toAddressID" db:"to_address_id"`                                  //12
	InteractedContractAddress   string           `json:"interactedContractAddress" db:"interacted_contract_address"`      //13
	InteractedContractAddressID *int             `json:"interactedContractAddressId" db:"interacted_contract_address_id"` //14
	NativeAssetID               *int             `json:"nativeAssetId" db:"native_asset_id"`                              //15
	GethProcessJobID            *int             `json:"gethProcessJobId" db:"geth_process_job_id"`                       //16
	Value                       *decimal.Decimal `json:"value" db:"value"`                                                //17
	GethTransctionInputId       *int             `json:"gethTransctionInputId" db:"geth_transction_input_id"`             //18
	StatusID                    *int             `json:"statusId" db:"status_id"`                                         //19
	Description                 string           `json:"description" db:"description"`                                    //20
	CreatedBy                   string           `json:"createdBy" db:"created_by"`                                       //21
	CreatedAt                   time.Time        `json:"createdAt" db:"created_at"`                                       //22
	UpdatedBy                   string           `json:"updatedBy" db:"updated_by"`                                       //23
	UpdatedAt                   time.Time        `json:"updatedAt" db:"updated_at"`                                       //24
}
