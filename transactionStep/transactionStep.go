package transactionstep

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Asset
type TransactionStep struct {
	TransactionID *int      `json:"transactionId" db:"transaction_id"` //1
	StepID        *int      `json:"stepId" db:"step_id"`               //2
	UUID          string    `json:"uuid" db:"uuid"`                    //3
	Name          string    `json:"name" db:"name"`                    //4
	AlternateName string    `json:"alternateName" db:"alternate_name"` //5
	Description   string    `json:"description" db:"description"`      //6
	CreatedBy     string    `json:"createdBy" db:"created_by"`         //7
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`         //8
	UpdatedBy     string    `json:"updatedBy" db:"updated_by"`         //9
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`         //10
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
