package transactionstep

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Asset
type TransactionStep struct {
	TransactionID *int      `json:"transactionId"` //1
	StepID        *int      `json:"stepId"`        //2
	UUID          string    `json:"uuid"`          //3
	Name          string    `json:"name"`          //4
	AlternateName string    `json:"alternateName"` //5
	Description   string    `json:"description"`   //6
	CreatedBy     string    `json:"createdBy"`     //7
	CreatedAt     time.Time `json:"createdAt"`     //8
	UpdatedBy     string    `json:"updatedBy"`     //9
	UpdatedAt     time.Time `json:"updatedAt"`     //10
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
