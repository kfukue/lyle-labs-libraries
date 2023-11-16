package gethlyletrades

import (
	"time"
)

// Asset
type GethTradeSwap struct {
	GethTradeID   *int      `json:"gethTradeId"`
	GethSwapID    *int      `json:"gethSwapId"`
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	AlternateName string    `json:"alternateName"`
	Description   string    `json:"description"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedBy     string    `json:"updatedBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
