package exchange

import (
	"time"
)

type Exchange struct {
	ID             *int       `json:"id" db:"id"`
	UUID           string     `json:"uuid" db:"uuid"`
	Name           string     `json:"name" db:"name" db:"name"`
	AlternateName  string     `json:"alternateName" db:"alternate_name"`
	ExchangeTypeID *int       `json:"exchangeTypeId" db:"exchange_type_id"`
	Url            string     `json:"url" db:"url"`
	StartDate      *time.Time `json:"startDate" db:"start_date"`
	EndDate        *time.Time `json:"endDate" db:"end_date"`
	Description    string     `json:"description" db:"description"`
	CreatedBy      string     `json:"createdBy" db:"created_by"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedBy      string     `json:"updatedBy" db:"updated_by"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type ExchangeChain struct {
	UUID        string    `json:"uuid" db:"uuid"`
	ExchangeID  *int      `json:"exchangeId" db:"exchange_id"`
	ChainID     *int      `json:"chainId" db:"chain_id"`
	Description string    `json:"description" db:"description"`
	CreatedBy   string    `json:"createdBy" db:"created_by"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedBy   string    `json:"updatedBy" db:"updated_by"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
