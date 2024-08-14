package tax

import (
	"time"

	"github.com/shopspring/decimal"
)

// Tax
type Tax struct {
	ID                 *int             `json:"id" db:"id"`                                   //1
	UUID               string           `json:"uuid" db:"uuid"`                               //2
	Name               string           `json:"name" db:"name"`                               //3
	AlternateName      string           `json:"alternateName" db:"alternate_name"`            //4
	StartDate          *time.Time       `json:"startDate" db:"start_date"`                    //5
	EndDate            *time.Time       `json:"endDate" db:"end_date"`                        //6
	StartBlock         *int             `json:"startBlock" db:"start_block"`                  //7
	EndBlock           *int             `json:"endBlock" db:"end_block"`                      //8
	TaxRate            *decimal.Decimal `json:"taxRate" db:"tax_rate"`                        //9
	TaxRateTypeID      *int             `json:"taxRateTypeId" db:"tax_rate_type_id"`          //10
	ContractAddressStr string           `json:"contractAddressStr" db:"contract_address_str"` //11
	ContractAddressID  *int             `json:"contractAddressId" db:"contract_address_id"`   //12
	TaxTypeID          *int             `json:"taxTypeId" db:"tax_type_id"`                   //13
	Description        string           `json:"description" db:"description"`                 //14
	CreatedBy          string           `json:"createdBy" db:"created_by"`                    //15
	CreatedAt          time.Time        `json:"createdAt" db:"created_at"`                    //16
	UpdatedBy          string           `json:"updatedBy" db:"updated_by"`                    //17
	UpdatedAt          time.Time        `json:"updatedAt" db:"updated_at"`                    //18
}
