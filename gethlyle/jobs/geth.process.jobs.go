package gethlylejobs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type GethProcessJob struct {
	ID               *int      `json:"id"`
	UUID             string    `json:"uuid"`
	Name             string    `json:"name"`
	AlternateName    string    `json:"alternateName"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
	Description      string    `json:"description"`
	StatusID         *int      `json:"statusId"`
	JobCategoryID    *int      `json:"jobCategoryId"`
	ImportTypeID     *int      `json:"importTypeId"`
	ChainID          *int      `json:"chainId"`
	StartBlockNumber *uint64   `json:"startBlockNumber"`
	EndBlockNumber   *uint64   `json:"endBlockNumber"`
	CreatedBy        string    `json:"createdBy"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedBy        string    `json:"updatedBy"`
	UpdatedAt        time.Time `json:"updatedAt"`
	AssetID          *int      `json:"assetId"`
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
