package pool

import (
	"time"
)

type Pool struct {
	ID            *int       `json:"id"`
	TargetAssetID *int       `json:"targetAssetId"`
	StrategyID    *int       `json:"strategyId"`
	AccountID     *int       `json:"accountId"`
	Name          string     `json:"name"`
	UUID          string     `json:"uuid"`
	AlternateName string     `json:"alternateName"`
	StartDate     *time.Time `json:"startDate"`
	EndDate       *time.Time `json:"endDate"`
	Description   string     `json:"description"`
	ChainID       *int       `json:"chainId"`
	FrequencyID   *int       `json:"frequencyId"`
	CreatedBy     string     `json:"createdBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedBy     string     `json:"updatedBy"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
