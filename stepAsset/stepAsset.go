package stepasset

import (
	"time"
)

type StepAsset struct {
	ID              *int       `json:"id"`
	StepID          *int       `json:"stepId"`
	AssetID         *int       `json:"assetId"`
	SwapAssetID     *int       `json:"swapAssetId"`
	TargetPoolID    *int       `json:"targetPoolId"`
	Name            string     `json:"name"`
	UUID            string     `json:"uuid"`
	AlternateName   string     `json:"alternateName"`
	StartDate       *time.Time `json:"startDate"`
	EndDate         *time.Time `json:"endDate"`
	Description     string     `json:"description"`
	ActionParameter *float64   `json:"actionParameter"`
	CreatedBy       string     `json:"createdBy"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedBy       string     `json:"updatedBy"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}
