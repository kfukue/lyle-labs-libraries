package gethlylemarketdata

import (
	"time"

	decimal "github.com/shopspring/decimal"
)

// MarketData
type MarketData struct {
	ID                *int              `json:"id"`
	UUID              string            `json:"uuid"`
	Name              string            `json:"name"`
	AlternateName     string            `json:"alternateName"`
	StartDate         time.Time         `json:"startDate"`
	EndDate           time.Time         `json:"endDate"`
	AssetID           *int              `json:"assetId"`
	OpenUSD           *decimal.Decimal  `json:"openUsd"`
	CloseUSD          *decimal.Decimal  `json:"closeUsd"`
	HighUSD           *decimal.Decimal  `json:"highUsd"`
	LowUSD            *decimal.Decimal  `json:"lowUsd"`
	PriceUSD          *decimal.Decimal  `json:"priceUsd"`
	VolumeUSD         *decimal.Decimal  `json:"volumeUsd"`
	MarketCapUSD      *decimal.Decimal  `json:"marketCapUsd"`
	Ticker            string            `json:"ticker"`
	Description       string            `json:"description"`
	IntervalID        *int              `json:"intervalId"`
	MarketDataTypeID  *int              `json:"marketDataTypeId"`
	SourceID          *int              `json:"sourceId"`
	TotalSupply       *decimal.Decimal  `json:"totalSupply"`
	MaxSupply         *decimal.Decimal  `json:"maxSupply"`
	CirculatingSupply *decimal.Decimal  `json:"circulatingSupply"`
	Sparkline7d       []decimal.Decimal `json:"sparkline7d"`
	CreatedBy         string            `json:"createdBy"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedBy         string            `json:"updatedBy"`
	UpdatedAt         time.Time         `json:"updatedAt"`
	GethProcessJobID  *int              `json:"gethProcessJobId"`
}
