package tax

import (
	"time"

	"github.com/shopspring/decimal"
)

// Tax
type Tax struct {
	ID                 *int             `json:"id"`
	UUID               string           `json:"uuid"`
	Name               string           `json:"name"`
	AlternateName      string           `json:"alternateName"`
	StartDate          *time.Time       `json:"startDate"`
	EndDate            *time.Time       `json:"endDate"`
	StartBlock         *int             `json:"startBlock"`
	EndBlock           *int             `json:"endBlock"`
	TaxRate            *decimal.Decimal `json:"taxRate"`
	TaxRateTypeID      *int             `json:"taxRateTypeId"`
	ContractAddressStr string           `json:"contractAddressStr"`
	ContractAddressID  *int             `json:"contractAddressId"`
	TaxTypeID          *int             `json:"taxTypeId"`
	Description        string           `json:"description"`
	CreatedBy          string           `json:"createdBy"`
	CreatedAt          time.Time        `json:"createdAt"`
	UpdatedBy          string           `json:"updatedBy"`
	UpdatedAt          time.Time        `json:"updatedAt"`
}
