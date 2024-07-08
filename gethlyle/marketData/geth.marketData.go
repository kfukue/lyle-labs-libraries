package gethlylemarketdata

import (
	"time"

	decimal "github.com/shopspring/decimal"
)

// MarketData
type GethMarketData struct {
	ID                *int              `json:"id" db:"id"`                                //1
	UUID              string            `json:"uuid" db:"uuid"`                            //2
	Name              string            `json:"name" db:"name"`                            //3
	AlternateName     string            `json:"alternateName" db:"alternate_name"`         //4
	StartDate         time.Time         `json:"startDate" db:"start_date"`                 //5
	EndDate           time.Time         `json:"endDate" db:"end_date"`                     //6
	AssetID           *int              `json:"assetId" db:"asset_id"`                     //7
	OpenUSD           *decimal.Decimal  `json:"openUsd" db:"open_usd"`                     //8
	CloseUSD          *decimal.Decimal  `json:"closeUsd" db:"close_usd"`                   //9
	HighUSD           *decimal.Decimal  `json:"highUsd" db:"high_usd"`                     //10
	LowUSD            *decimal.Decimal  `json:"lowUsd" db:"low_usd"`                       //11
	PriceUSD          *decimal.Decimal  `json:"priceUsd" db:"price_usd"`                   //12
	VolumeUSD         *decimal.Decimal  `json:"volumeUsd" db:"volume_usd"`                 //13
	MarketCapUSD      *decimal.Decimal  `json:"marketCapUsd" db:"market_cap_usd"`          //14
	Ticker            string            `json:"ticker" db:"ticker"`                        //15
	Description       string            `json:"description" db:"description"`              //16
	IntervalID        *int              `json:"intervalId" db:"interval_id"`               //17
	MarketDataTypeID  *int              `json:"marketDataTypeId" db:"market_data_type_id"` //18
	SourceID          *int              `json:"sourceId" db:"source_id"`                   //19
	TotalSupply       *decimal.Decimal  `json:"totalSupply" db:"total_supply"`             //20
	MaxSupply         *decimal.Decimal  `json:"maxSupply" db:"max_supply"`                 //21
	CirculatingSupply *decimal.Decimal  `json:"circulatingSupply" db:"circulating_supply"` //22
	Sparkline7d       []decimal.Decimal `json:"sparkline7d" db:"sparkline_7d"`             //23
	CreatedBy         string            `json:"createdBy" db:"created_by"`                 //24
	CreatedAt         time.Time         `json:"createdAt" db:"created_at"`                 //25
	UpdatedBy         string            `json:"updatedBy" db:"updated_by"`                 //26
	UpdatedAt         time.Time         `json:"updatedAt" db:"updated_at"`                 //27
	GethProcessJobID  *int              `json:"gethProcessJobId" db:"geth_process_job_id"` //28
}
