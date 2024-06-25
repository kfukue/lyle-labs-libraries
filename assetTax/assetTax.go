package assettax

import (
	"time"

	"github.com/shopspring/decimal"
)

// Asset
type AssetTax struct {
	TaxID           *int             `json:"taxId" db:"tax_id"`
	AssetID         *int             `json:"assetId" db:"asset_id"`
	UUID            string           `json:"uuid" db:"uuid"`
	Name            string           `json:"name" db:"name"`
	AlternateName   string           `json:"alternateName" db:"alternate_name"`
	TaxRateOverride *decimal.Decimal `json:"taxRateOverride" db:"tax_rate_override"`
	TaxRateTypeID   *int             `json:"taxRateTypeId" db:"tax_rate_type_id"`
	Description     string           `json:"description" db:"description"`
	CreatedBy       string           `json:"createdBy" db:"created_by"`
	CreatedAt       time.Time        `json:"createdAt" db:"created_at"`
	UpdatedBy       string           `json:"updatedBy" db:"updated_by"`
	UpdatedAt       time.Time        `json:"updatedAt" db:"updated_at"`
}
