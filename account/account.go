package account

import (
	"time"
)

// Asset
type Account struct {
	ID             *int      `json:"id"`
	UUID           string    `json:"uuid"`
	Name           string    `json:"name"`
	AlternateName  string    `json:"alternateName"`
	Address        string    `json:"address"`
	NameFromSource string    `json:"nameFromSource"`
	PortfolioID    *int      `json:"portfolioId"`
	SourceID       *int      `json:"sourceId"`
	AccountTypeID  *int      `json:"accountTypeId"`
	Description    string    `json:"description"`
	CreatedBy      string    `json:"createdBy"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedBy      string    `json:"updatedBy"`
	UpdatedAt      time.Time `json:"updatedAt"`
	ChainID        *int      `json:"chainId"`
}
