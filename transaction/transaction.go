package transaction

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Transaction
type Transaction struct {
	ID            *int      `json:"id"`
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	AlternateName string    `json:"alternateName"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	Description   string    `json:"description"`
	TxHash        string    `json:"txHash"`
	StatusID      *int      `json:"statusId"`
	FromAccountID *int      `json:"fromAccountId"`
	ToAccountID   *int      `json:"toAccountId"`
	ChainID       *int      `json:"chainId"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedBy     string    `json:"updatedBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
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
