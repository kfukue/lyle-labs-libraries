package gethlylemarketdata

import (
	"time"
)

// MarketData
type MarketData struct {
	ID                *int      `json:"id"`
	UUID              string    `json:"uuid"`
	Name              string    `json:"name"`
	AlternateName     string    `json:"alternateName"`
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	AssetID           *int      `json:"assetId"`
	OpenUSD           *float64  `json:"openUsd"`
	CloseUSD          *float64  `json:"closeUsd"`
	HighUSD           *float64  `json:"highUsd"`
	LowUSD            *float64  `json:"lowUsd"`
	PriceUSD          *float64  `json:"priceUsd"`
	VolumeUSD         *float64  `json:"volumeUsd"`
	MarketCapUSD      *float64  `json:"marketCapUsd"`
	Ticker            string    `json:"ticker"`
	Description       string    `json:"description"`
	IntervalID        *int      `json:"intervalId"`
	MarketDataTypeID  *int      `json:"marketDataTypeId"`
	SourceID          *int      `json:"sourceId"`
	TotalSupply       *float64  `json:"totalSupply"`
	MaxSupply         *float64  `json:"maxSupply"`
	CirculatingSupply *float64  `json:"circulatingSupply"`
	Sparkline7d       []float64 `json:"sparkline7d"`
	CreatedBy         string    `json:"createdBy"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedBy         string    `json:"updatedBy"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
