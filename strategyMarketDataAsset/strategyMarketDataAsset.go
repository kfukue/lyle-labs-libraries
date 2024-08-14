package strategymarketdataasset

import (
	"time"
)

type StrategyMarketDataAsset struct {
	ID            *int       `json:"id" db:"id"`                        //1
	StrategyID    *int       `json:"strategyId"  db:"strategy_id"`      //2
	BaseAssetID   *int       `json:"baseAssetId" db:"base_asset_id"`    //3
	QuoteAssetID  *int       `json:"quoteAssetId" db:"quote_asset_id"`  //4
	UUID          string     `json:"uuid" db:"uuid"`                    //5
	Name          string     `json:"name" db:"name"`                    //6
	AlternateName string     `json:"alternateName" db:"alternate_name"` //7
	StartDate     *time.Time `json:"startDate" db:"start_date"`         //8
	EndDate       *time.Time `json:"endDate" db:"end_date"`             //9
	Ticker        string     `json:"ticker" db:"ticker"`                //10
	Description   string     `json:"description" db:"description"`      //11
	SourceID      *int       `json:"sourceId" db:"source_id"`           //12
	FrequencyID   *int       `json:"frequencyId"  db:"frequency_id"`    //13
	CreatedBy     string     `json:"createdBy" db:"created_by"`         //14
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`         //15
	UpdatedBy     string     `json:"updatedBy" db:"updated_by"`         //16
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`         //17
}
