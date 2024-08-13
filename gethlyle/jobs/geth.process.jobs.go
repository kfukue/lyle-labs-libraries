package gethlylejobs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type GethProcessJob struct {
	ID               *int       `json:"id" db:"id"`                               //1
	UUID             string     `json:"uuid" db:"uuid"`                           //2
	Name             string     `json:"name" db:"name"`                           //3
	AlternateName    string     `json:"alternateName" db:"alternate_name"`        //4
	StartDate        time.Time  `json:"startDate" db:"start_date"`                //5
	EndDate          *time.Time `json:"endDate" db:"end_date"`                    //6
	Description      string     `json:"description" db:"description"`             //7
	StatusID         *int       `json:"statusId" db:"status_id"`                  //8
	JobCategoryID    *int       `json:"jobCategoryId" db:"job_category_id"`       //9
	ImportTypeID     *int       `json:"importTypeId" db:"import_type_id"`         //10
	ChainID          *int       `json:"chainId" db:"chain_id"`                    //11
	StartBlockNumber *uint64    `json:"startBlockNumber" db:"start_block_number"` //12
	EndBlockNumber   *uint64    `json:"endBlockNumber" db:"end_block_number"`     //13
	CreatedBy        string     `json:"createdBy" db:"created_by"`                //14
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`                //15
	UpdatedBy        string     `json:"updatedBy" db:"updated_by"`                //16
	UpdatedAt        time.Time  `json:"updatedAt" db:"updated_at"`                //17
	AssetID          *int       `json:"assetId"  db:"asset_id"`                   //18
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
