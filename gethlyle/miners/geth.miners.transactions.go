package gethlyleminers

import (
	"time"
)

type GethMinerTransaction struct {
	MinerID       *int      `json:"minerId"`
	TransactionID *int      `json:"transactionId"`
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	AlternateName string    `json:"alternateName"`
	Description   string    `json:"description"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedBy     string    `json:"updatedBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
