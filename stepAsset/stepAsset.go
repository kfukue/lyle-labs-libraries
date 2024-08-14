package stepasset

import (
	"time"
)

type StepAsset struct {
	ID              *int       `json:"id" db:"id"`                            //1
	StepID          *int       `json:"stepId" db:"step_id"`                   //2
	AssetID         *int       `json:"assetId" db:"asset_id"`                 //3
	SwapAssetID     *int       `json:"swapAssetId" db:"swap_asset_id"`        //4
	TargetPoolID    *int       `json:"targetPoolId" db:"target_pool_id"`      //5
	UUID            string     `json:"uuid" db:"uuid"`                        //6
	Name            string     `json:"name" db:"name"`                        //7
	AlternateName   string     `json:"alternateName" db:"alternate_name"`     //8
	StartDate       *time.Time `json:"startDate" db:"start_date"`             //9
	EndDate         *time.Time `json:"endDate" db:"end_date"`                 //10
	Description     string     `json:"description" db:"description"`          //11
	ActionParameter *float64   `json:"actionParameter" db:"action_parameter"` //12
	CreatedBy       string     `json:"createdBy" db:"created_by"`             //13
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`             //14
	UpdatedBy       string     `json:"updatedBy" db:"updated_by"`             //15
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`             //16
}
