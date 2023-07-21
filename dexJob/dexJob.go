package dexjob

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Asset
type DexTxnJob struct {
	ID                *int      `json:"Id"`
	JobID             *int      `json:"jobId"`
	UUID              string    `json:"uuid"`
	Name              string    `json:"name"`
	AlternateName     string    `json:"alternateName"`
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	Description       string    `json:"description"`
	StatusID          *int      `json:"statusId"`
	ChainID           *int      `json:"chainId"`
	ExchangeID        *int      `json:"exchangeId"`
	TransactionHashes []string  `json:"trasnactionHashes"`
	CreatedBy         string    `json:"createdBy"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedBy         string    `json:"updatedBy"`
	UpdatedAt         time.Time `json:"updatedAt"`
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
