package gethlyleminers

import (
	"time"
)

type GethMinerTransactionInput struct {
	MinerID            *int      `json:"minerId"`
	TransactionInputID *int      `json:"transactionInputId"`
	UUID               string    `json:"uuid"`
	Name               string    `json:"name"`
	AlternateName      string    `json:"alternateName"`
	Description        string    `json:"description"`
	CreatedBy          string    `json:"createdBy"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedBy          string    `json:"updatedBy"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type GethMinerWithTransactionInput struct {
	GethMinerTransactionInput
	GethMiner
}
