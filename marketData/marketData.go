package marketdata

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

// MarketData
type MarketDataQuote struct {
	MarketDataID         *int      `json:"marketDataId"`
	BaseAssetID          *int      `json:"baseAssetId"`
	QuoteAssetID         *int      `json:"quoteAssetId"`
	UUID                 string    `json:"uuid"`
	Name                 string    `json:"name"`
	AlternateName        string    `json:"alternateName"`
	Open                 *float64  `json:"open"`
	Close                *float64  `json:"close"`
	High24h              *float64  `json:"high24h"`
	Low24h               *float64  `json:"low24h"`
	Price                *float64  `json:"price"`
	Volume               *float64  `json:"volume"`
	MarketCap            *float64  `json:"marketCap"`
	Ticker               string    `json:"ticker"`
	Description          string    `json:"description"`
	SourceID             *int      `json:"sourceId"`
	FullyDilutedValution *float64  `json:"fullyDilutedValution"`
	Ath                  *float64  `json:"ath"`
	AthDate              time.Time `json:"athDate"`
	Atl                  *float64  `json:"atl"`
	AtlDate              time.Time `json:"atlDate"`
	PriceChange1h        *float64  `json:"priceChange1h"`
	PriceChange24h       *float64  `json:"priceChange24h"`
	PriceChange7d        *float64  `json:"priceChange7d"`
	PriceChange30d       *float64  `json:"priceChange30d"`
	PriceChange60d       *float64  `json:"priceChange60d"`
	PriceChange200d      *float64  `json:"priceChange200d"`
	PriceChange1y        *float64  `json:"priceChange1y"`
	CreatedBy            string    `json:"createdBy"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedBy            string    `json:"updatedBy"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type MarketDataQuoteResults struct {
	MarketDataQuote
	StartDate        time.Time
	EndDate          time.Time
	BaseAssetName    string
	BaseAssetTicker  string
	QuoteAssetName   string
	QuoteAssetTicker string
}
