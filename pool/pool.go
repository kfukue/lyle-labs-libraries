package pool

import (
	"time"
)

type Pool struct {
	ID            *int      `json:"id" db:"id"`                         //1
	TargetAssetID *int      `json:"targetAssetId" db:"target_asset_id"` //2
	StrategyID    *int      `json:"strategyId" db:"strategy_id"`        //3
	AccountID     *int      `json:"accountId" db:"account_id"`          //4
	UUID          string    `json:"uuid" db:"uuid"`                     //5
	Name          string    `json:"name" db:"name"`                     //6
	AlternateName string    `json:"alternateName" db:"alternate_name"`  //7
	StartDate     time.Time `json:"startDate" db:"start_date"`          //8
	EndDate       time.Time `json:"endDate" db:"end_date"`              //9
	Description   string    `json:"description" db:"description"`       //10
	ChainID       *int      `json:"chainId" db:"chain_id"`              //11
	FrequencyID   *int      `json:"frequencyId" db:"frequency_id"`      //12
	CreatedBy     string    `json:"createdBy" db:"created_by"`          //13
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`          //14
	UpdatedBy     string    `json:"updatedBy" db:"updated_by"`          //15
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`          //16
}
