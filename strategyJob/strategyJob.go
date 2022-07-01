package strategyjob

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Asset
type StrategyJob struct {
	StrategyID       *int      `json:"strategyId"`
	JobID            *int      `json:"jobId"`
	UUID             string    `json:"uuid"`
	Name             string    `json:"name"`
	AlternateName    string    `json:"alternateName"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
	Description      string    `json:"description"`
	StatusID         *int      `json:"statusId"`
	ResponseStatus   string    `json:"responseStatus"`
	RequestUrl       string    `json:"requestUrl"`
	RequestBody      string    `json:"requestBody"`
	RequestMethod    string    `json:"requestMethod"`
	ResponseData     string    `json:"responseData"`
	ResponseDataJson Attrs     `json:"responseDataJson"`
	CreatedBy        string    `json:"createdBy"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedBy        string    `json:"updatedBy"`
	UpdatedAt        time.Time `json:"updatedAt"`
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
