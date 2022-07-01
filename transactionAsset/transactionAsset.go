package transactionasset

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Asset
type TransactionAsset struct {
	TransactionID         *int      `json:"transactionId"`
	AssetID               *int      `json:"assetId"`
	UUID                  string    `json:"uuid"`
	Name                  string    `json:"name"`
	AlternateName         string    `json:"alternateName"`
	Description           string    `json:"description"`
	Quantity              *float64  `json:"quantity"`
	QuantityUSD           *float64  `json:"quantityUsd"`
	MarketDataID          *int      `json:"marketDataId"`
	ManualExchangeRateUSD *float64  `json:"manualExchangeRateUsd"`
	CreatedBy             string    `json:"createdBy"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedBy             string    `json:"updatedBy"`
	UpdatedAt             time.Time `json:"updatedAt"`
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
