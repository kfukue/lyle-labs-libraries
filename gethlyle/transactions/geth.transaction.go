package gethlyletransactions

import (
	"time"

	"github.com/shopspring/decimal"
)

type GethTransaction struct {
	ID                          *int             `json:"id"`
	UUID                        string           `json:"uuid"`
	ChainID                     *int             `json:"chainId"`
	ExchangeID                  *int             `json:"exchangeId"`
	BlockNumber                 *uint64          `json:"blockNumber"`
	IndexNumber                 *uint            `json:"indexNumber"`
	TxnDate                     time.Time        `json:"txnDate"`
	TxnHash                     string           `json:"txnHash"`
	FromAddress                 string           `json:"senderAddress"`
	FromAddressID               *int             `json:"senderAddressID"`
	ToAddress                   string           `json:"toAddress"`
	ToAddressID                 *int             `json:"toAddressID"`
	InteractedContractAddress   string           `json:"interactedContractAddress"`
	InteractedContractAddressID *int             `json:"interactedContractAddressID"`
	NativeAssetID               *int             `json:"nativeAssetId"`
	GethProcessJobID            *int             `json:"gethProcessJobId"`
	Value                       *decimal.Decimal `json:"value"`
	GethTransctionInputId       *int             `json:"gethTransctionInputId"`
	StatusID                    *int             `json:"statusId"`
	Description                 string           `json:"description"`
	CreatedBy                   string           `json:"createdBy"`
	CreatedAt                   time.Time        `json:"createdAt"`
	UpdatedBy                   string           `json:"updatedBy"`
	UpdatedAt                   time.Time        `json:"updatedAt"`
}
