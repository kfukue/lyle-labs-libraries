package exchange

import (
	"time"
)

type Exchange struct {
	ID            *int       `json:"id"`
	UUID          string     `json:"uuid"`
	Name          string     `json:"name"`
	AlternateName string     `json:"alternateName"`
	Url           string     `json:"url"`
	StartDate     *time.Time `json:"startDate"`
	EndDate       *time.Time `json:"endDate"`
	Description   string     `json:"description"`
	CreatedBy     string     `json:"createdBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedBy     string     `json:"updatedBy"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type ExchangeChain struct {
	UUID        string    `json:"uuid"`
	ExchangeID  *int      `json:"exchangeId"`
	ChainID     *int      `json:"chainId"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedBy   string    `json:"updatedBy"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
