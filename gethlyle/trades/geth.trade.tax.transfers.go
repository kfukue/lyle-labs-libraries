package gethlyletrades

import (
	"time"
)

// Asset
type GethTradeTaxTransfer struct {
	GethTradeID    *int      `json:"gethTradeId" db:"geth_trade_id"`       //1
	GethTransferID *int      `json:"gethTransferId" db:"geth_transfer_id"` //2
	TaxID          *int      `json:"taxId" db:"tax_id"`                    //3
	UUID           string    `json:"uuid" db:"uuid"`                       //4
	Name           string    `json:"name" db:"name"`                       //5
	AlternateName  string    `json:"alternateName" db:"alternate_name"`    //6
	Description    string    `json:"description" db:"description"`         //7
	CreatedBy      string    `json:"createdBy" db:"created_by"`            //8
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`            //9
	UpdatedBy      string    `json:"updatedBy" db:"updated_by"`            //10
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`            //11
}
