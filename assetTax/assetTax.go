package assettax

import (
	"time"

	"github.com/shopspring/decimal"
)

// Asset
type AssetTax struct {
	TaxID           *int             `json:"taxId"`
	AssetID         *int             `json:"assetId"`
	UUID            string           `json:"uuid"`
	Name            string           `json:"name"`
	AlternateName   string           `json:"alternateName"`
	TaxRateOverride *decimal.Decimal `json:"taxRateOverride"`
	Description     string           `json:"description"`
	CreatedBy       string           `json:"createdBy"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedBy       string           `json:"updatedBy"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}
