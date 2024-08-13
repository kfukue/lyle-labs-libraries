package transactionasset

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Asset
type TransactionAsset struct {
	TransactionID         *int      `json:"transactionId" db:"transaction_id"`                   //1
	AssetID               *int      `json:"assetId" db:"asset_id"`                               //2
	UUID                  string    `json:"uuid" db:"uuid"`                                      //3
	Name                  string    `json:"name" db:"name"`                                      //4
	AlternateName         string    `json:"alternateName" db:"alternate_name"`                   //5
	Description           string    `json:"description" db:"description"`                        //6
	Quantity              *float64  `json:"quantity" db:"quantity"`                              //7
	QuantityUSD           *float64  `json:"quantityUsd" db:"quantity_usd"`                       //8
	MarketDataID          *int      `json:"marketDataId" db:"market_data_id"`                    //9
	ManualExchangeRateUSD *float64  `json:"manualExchangeRateUsd" db:"manual_exchange_rate_usd"` //10
	CreatedBy             string    `json:"createdBy" db:"created_by"`                           //11
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`                           //12
	UpdatedBy             string    `json:"updatedBy" db:"updated_by"`                           //13
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`                           //14
}

type Attrs map[string]interface{}

func (a Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
