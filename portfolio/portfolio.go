package portfolio

import (
	"time"
)

// Portfolio
type Portfolio struct {
	ID              *int       `json:"id"`
	UUID            string     `json:"uuid"`
	Name            string     `json:"name"`
	AlternateName   string     `json:"alternateName"`
	StartDate       *time.Time `json:"startDate"`
	EndDate         *time.Time `json:"endDate"`
	UserEmail       string     `json:"userEmail"`
	Description     string     `json:"description"`
	BaseAssetID     *int       `json:"baseAssetId"`
	PortfolioTypeID *int       `json:"portfolioTypeId"`
	ParentID        *int       `json:"parentId"`
	CreatedBy       string     `json:"createdBy"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedBy       string     `json:"updatedBy"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}
