package account

import (
	"time"
)

// Asset
type Account struct {
	ID             *int      `json:"id" db:"id"`                           //1
	UUID           string    `json:"uuid" db:"uuid"`                       //2
	Name           string    `json:"name" db:"name"`                       //3
	AlternateName  string    `json:"alternateName" db:"alternate_name"`    //4
	Address        string    `json:"address" db:"address"`                 //5
	NameFromSource string    `json:"nameFromSource" db:"name_from_source"` //6
	PortfolioID    *int      `json:"portfolioId" db:"portfolio_id"`        //7
	SourceID       *int      `json:"sourceId" db:"source_id"`              //8
	AccountTypeID  *int      `json:"accountTypeId" db:"account_type_id"`   //9
	Description    string    `json:"description" db:"description"`         //10
	CreatedBy      string    `json:"createdBy" db:"created_by"`            //11
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`            //12
	UpdatedBy      string    `json:"updatedBy" db:"updated_by"`            //13
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`            //14
	ChainID        *int      `json:"chainId" db:"chain_id"`                //15
}
