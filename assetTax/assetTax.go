package assettax

import (
	"time"

	"github.com/shopspring/decimal"
)

// Asset
type AssetTax struct {
	TaxID           *int             `json:"taxId" db:"tax_id"`                      //1
	AssetID         *int             `json:"assetId" db:"asset_id"`                  //2
	UUID            string           `json:"uuid" db:"uuid"`                         //3
	Name            string           `json:"name" db:"name"`                         //4
	AlternateName   string           `json:"alternateName" db:"alternate_name"`      //5
	TaxRateOverride *decimal.Decimal `json:"taxRateOverride" db:"tax_rate_override"` //6
	TaxRateTypeID   *int             `json:"taxRateTypeId" db:"tax_rate_type_id"`    //7
	Description     string           `json:"description" db:"description"`           //8
	CreatedBy       string           `json:"createdBy" db:"created_by"`              //9
	CreatedAt       time.Time        `json:"createdAt" db:"created_at"`              //10
	UpdatedBy       string           `json:"updatedBy" db:"updated_by"`              //11
	UpdatedAt       time.Time        `json:"updatedAt" db:"updated_at"`              //12
}
