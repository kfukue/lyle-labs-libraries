package account

import (
	"time"
)

// Asset
type Account struct {
	ID             *int      `json:"id" db:"id"`
	UUID           string    `json:"uuid" db:"uuid"`
	Name           string    `json:"name" db:"name"`
	AlternateName  string    `json:"alternateName" db:"alternate_name"`
	Address        string    `json:"address" db:"address"`
	NameFromSource string    `json:"nameFromSource" db:"name_from_source"`
	PortfolioID    *int      `json:"portfolioId" db:"portfolio_id"`
	SourceID       *int      `json:"sourceId" db:"source_id"`
	AccountTypeID  *int      `json:"accountTypeId" db:"account_type_id"`
	Description    string    `json:"description" db:"description"`
	CreatedBy      string    `json:"createdBy" db:"created_by"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedBy      string    `json:"updatedBy" db:"updated_by"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
	ChainID        *int      `json:"chainId" db:"chain_id"`
}
