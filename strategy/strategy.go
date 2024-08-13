package strategy

import (
	"time"
)

// Strategy
type Strategy struct {
	ID             *int      `json:"id" db:"id"`                           //1
	UUID           string    `json:"uuid" db:"uuid"`                       //2
	Name           string    `json:"name" db:"name"`                       //3
	AlternateName  string    `json:"alternateName" db:"alternate_name"`    //4
	StartDate      time.Time `json:"startDate" db:"start_date"`            //5
	EndDate        time.Time `json:"endDate" db:"end_date"`                //6
	Description    string    `json:"description" db:"description"`         //7
	StrategyTypeID *int      `json:"strategyTypeId" db:"strategy_type_id"` //8
	CreatedBy      string    `json:"createdBy" db:"created_by"`            //9
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`            //10
	UpdatedBy      string    `json:"updatedBy" db:"updated_by"`            //11
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`            //12
}
