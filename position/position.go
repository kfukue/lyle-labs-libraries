package position

import (
	"time"
)

// Position
type Position struct {
	ID            *int       `json:"id" db:"id"`                        //1
	UUID          string     `json:"uuid" db:"uuid"`                    //2
	Name          string     `json:"name" db:"name"`                    //3
	AlternateName string     `json:"alternateName" db:"alternate_name"` //4
	AccountID     *int       `json:"accountId"  db:"account_id"`        //5
	PortfolioID   *int       `json:"portfolioId" db:"portfolio_id"`     //6
	FrequnecyID   *int       `json:"frequencyId" db:"frequency_id"`     //7
	StartDate     *time.Time `json:"startDate" db:"start_date"`         //8
	EndDate       *time.Time `json:"endDate" db:"end_date"`             //9
	BaseAssetID   *int       `json:"baseAssetId" db:"base_asset_id"`    //10
	QuoteAssetID  *int       `json:"quoteAssetId" db:"quote_asset_id"`  //11
	Quantity      *float64   `json:"quantity" db:"quantity"`            //12
	CostBasis     *float64   `json:"costBasis" db:"cost_basis"`         //13
	Profit        *float64   `json:"profit" db:"profit"`                //14
	TotalAmount   *float64   `json:"totalAmount" db:"total_amount"`     //15
	Description   string     `json:"description" db:"description"`      //16
	CreatedBy     string     `json:"createdBy" db:"created_by"`         //17
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`         //18
	UpdatedBy     string     `json:"updatedBy" db:"updated_by"`         //19
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`         //20
}

type PositionQuery struct {
	AccountID   *int       `schema:"accountId" db:"account_id"`      //1
	PortfolioID *int       `schema:"portfolioId" db:"portfolio_id"`  //2
	FrequencyID *int       `schema:"frequencyId" db:"frequency_id"`  //3
	StartDate   *time.Time `schema:"startDate" db:"start_date"`      //4
	EndDate     *time.Time `schema:"endDate" db:"end_date"`          //5
	BaseAssetID *int       `schema:"baseAssetId" db:"base_asset_id"` //6
}
