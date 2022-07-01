package position

import (
	"time"
)

// Position
type Position struct {
	ID            *int      `json:"id"`
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	AlternateName string    `json:"alternateName"`
	AccountID     *int      `json:"accountId"`
	PortfolioID   *int      `json:"portfolioId"`
	FrequnecyID   *int      `json:"frequencyId"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	BaseAssetID   *int      `json:"baseAssetId"`
	QuoteAssetID  *int      `json:"quoteAssetId"`
	Quantity      *float64  `json:"quantity"`
	CostBasis     *float64  `json:"costBasis"`
	Profit        *float64  `json:"profit"`
	TotalAmount   *float64  `json:"totalAmount"`
	Description   string    `json:"description"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedBy     string    `json:"updatedBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type PositionQuery struct {
	AccountID   *int       `schema:"accountId"`
	PortfolioID *int       `schema:"portfolioId"`
	FrequencyID *int       `schema:"frequencyId"`
	StartDate   *time.Time `schema:"startDate"`
	EndDate     *time.Time `schema:"endDate"`
	BaseAssetID *int       `schema:"baseAssetId"`
}
