package dexjob

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Asset
type DexTxnJob struct {
	ID                *int       `json:"id" db:"id"`
	JobID             *int       `json:"jobId" db:"job_id"`
	UUID              string     `json:"uuid" db:"uuid"`
	Name              string     `json:"name" db:"name" db:"name"`
	AlternateName     string     `json:"alternateName" db:"alternate_name"`
	StartDate         *time.Time `json:"startDate" db:"start_date"`
	EndDate           *time.Time `json:"endDate" db:"end_date"`
	Description       string     `json:"description" db:"description"`
	StatusID          *int       `json:"statusId" db:"status_id"`
	ChainID           *int       `json:"chainId" db:"chain_id"`
	ExchangeID        *int       `json:"exchangeId" db:"exchange_id"`
	TransactionHashes []string   `json:"trasnactionHashes" db:"transaction_hashes"`
	CreatedBy         string     `json:"createdBy" db:"created_by"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	UpdatedBy         string     `json:"updatedBy" db:"updated_by"`
	UpdatedAt         time.Time  `json:"updatedAt" db:"updated_at"`
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
