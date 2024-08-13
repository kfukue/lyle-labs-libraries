package portfolio

import (
	"time"
)

// Portfolio
type Portfolio struct {
	ID              *int       `json:"id" db:"id"`                             //1
	UUID            string     `json:"uuid" db:"uuid"`                         //2
	Name            string     `json:"name" db:"name"`                         //3
	AlternateName   string     `json:"alternateName" db:"alternate_name"`      //4
	StartDate       *time.Time `json:"startDate" db:"start_date"`              //5
	EndDate         *time.Time `json:"endDate" db:"end_date"`                  //6
	UserEmail       string     `json:"userEmail"  db:"user_email"`             //7
	Description     string     `json:"description" db:"description"`           //8
	BaseAssetID     *int       `json:"baseAssetId" db:"base_asset_id"`         //9
	PortfolioTypeID *int       `json:"portfolioTypeId" db:"portfolio_type_id"` //10
	ParentID        *int       `json:"parentId" db:"parent_id"`                //11
	CreatedBy       string     `json:"createdBy" db:"created_by"`              //12
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`              //13
	UpdatedBy       string     `json:"updatedBy" db:"updated_by"`              //14
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`              //15
}
