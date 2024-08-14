package transaction

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Transaction
type Transaction struct {
	ID            *int       `json:"id" db:"id"`                         //1
	UUID          string     `json:"uuid" db:"uuid"`                     //2
	Name          string     `json:"name" db:"name"`                     //3
	AlternateName string     `json:"alternateName" db:"alternate_name"`  //4
	StartDate     *time.Time `json:"startDate" db:"start_date"`          //5
	EndDate       *time.Time `json:"endDate" db:"end_date"`              //6
	Description   string     `json:"description" db:"description"`       //7
	TxHash        string     `json:"txHash"  db:"tx_hash"`               //8
	StatusID      *int       `json:"statusId" db:"status_id"`            //9
	FromAccountID *int       `json:"fromAccountId" db:"from_account_id"` //10
	ToAccountID   *int       `json:"toAccountId" db:"to_account_id"`     //11
	ChainID       *int       `json:"chainId" db:"chain_id"`              //12
	CreatedBy     string     `json:"createdBy" db:"created_by"`          //13
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`          //14
	UpdatedBy     string     `json:"updatedBy" db:"updated_by"`          //15
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`          //16
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
