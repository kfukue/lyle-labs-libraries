package strategymarketdataasset

import (
	"time"
)

type StrategyMarketDataAsset struct {
	ID            *int       `json:"id"`
	StrategyID    *int       `json:"strategyId"`
	BaseAssetID   *int       `json:"baseAssetId"`
	QuoteAssetID  *int       `json:"quoteAssetId"`
	Name          string     `json:"name"`
	UUID          string     `json:"uuid"`
	AlternateName string     `json:"alternateName"`
	StartDate     *time.Time `json:"startDate"`
	EndDate       *time.Time `json:"endDate"`
	Ticker        string     `json:"ticker"`
	Description   string     `json:"description"`
	SourceID      *int       `json:"sourceId"`
	FrequencyID   *int       `json:"frequencyId"`
	CreatedBy     string     `json:"createdBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedBy     string     `json:"updatedBy"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
