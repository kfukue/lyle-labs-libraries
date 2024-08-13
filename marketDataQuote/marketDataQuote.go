package marketdataquote

import (
	"time"
)

type MarketDataQuote struct {
	MarketDataID         *int      `json:"marketDataId" db:"market_data_id"`                 //1
	BaseAssetID          *int      `json:"baseAssetId" db:"base_asset_id"`                   //2
	QuoteAssetID         *int      `json:"quoteAssetId" db:"quote_asset_id"`                 //3
	UUID                 string    `json:"uuid" db:"uuid"`                                   //4
	Name                 string    `json:"name" db:"name"`                                   //5
	AlternateName        string    `json:"alternateName" db:"alternate_name"`                //6
	Open                 *float64  `json:"open" db:"open"`                                   //7
	Close                *float64  `json:"close" db:"close"`                                 //8
	High24h              *float64  `json:"high24h" db:"high_24h"`                            //9
	Low24h               *float64  `json:"low24h" db:"low_24h"`                              //10
	Price                *float64  `json:"price" db:"price"`                                 //11
	Volume               *float64  `json:"volume" db:"volume"`                               //12
	MarketCap            *float64  `json:"marketCap" db:"market_cap"`                        //13
	Ticker               string    `json:"ticker" db:"ticker"`                               //14
	Description          string    `json:"description" db:"description"`                     //15
	SourceID             *int      `json:"sourceId" db:"source_id"`                          //16
	FullyDilutedValution *float64  `json:"fullyDilutedValution" db:"fully_diluted_valution"` //17
	Ath                  *float64  `json:"ath" db:"ath"`                                     //18
	AthDate              time.Time `json:"athDate" db:"ath_date"`                            //19
	Atl                  *float64  `json:"atl" db:"atl"`                                     //20
	AtlDate              time.Time `json:"atlDate" db:"atl_date"`                            //21
	PriceChange1h        *float64  `json:"priceChange1h" db:"price_change_1h"`               //22
	PriceChange24h       *float64  `json:"priceChange24h" db:"price_change_24h"`             //23
	PriceChange7d        *float64  `json:"priceChange7d" db:"price_change_7d"`               //24
	PriceChange30d       *float64  `json:"priceChange30d" db:"price_change_30d"`             //25
	PriceChange60d       *float64  `json:"priceChange60d" db:"price_change_60d"`             //26
	PriceChange200d      *float64  `json:"priceChange200d" db:"price_change_200d"`           //27
	PriceChange1y        *float64  `json:"priceChange1y" db:"price_change_1y"`               //28
	CreatedBy            string    `json:"createdBy" db:"created_by"`                        //29
	CreatedAt            time.Time `json:"createdAt" db:"created_at"`                        //30
	UpdatedBy            string    `json:"updatedBy" db:"updated_by"`                        //31
	UpdatedAt            time.Time `json:"updatedAt" db:"updated_at"`                        //32
}

type MarketDataQuoteResults struct {
	StartDate        time.Time `json:"startDate" db:"start_date"`                //1
	EndDate          time.Time `json:"endDate" db:"end_date"`                    //2
	BaseAssetName    string    `json:"baseAssetName" db:"base_asset_name"`       //3
	BaseAssetTicker  string    `json:"baseAssetTicker" db:"base_asset_ticker"`   //4
	QuoteAssetName   string    `json:"quoteAssetName" db:"quote_asset_name"`     //5
	QuoteAssetTicker string    `json:"quoteAssetTicker" db:"quote_asset_ticker"` //6
	MarketDataQuote
}
