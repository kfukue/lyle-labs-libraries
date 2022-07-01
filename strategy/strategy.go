package strategy

import (
	"time"
)

// Strategy
type Strategy struct {
	ID             *int       `json:"id"`
	UUID           string     `json:"uuid"`
	Name           string     `json:"name"`
	AlternateName  string     `json:"alternateName"`
	StartDate      *time.Time `json:"startDate"`
	EndDate        *time.Time `json:"endDate"`
	Description    string     `json:"description"`
	StrategyTypeID *int       `json:"strategyTypeId"`
	CreatedBy      string     `json:"createdBy"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedBy      string     `json:"updatedBy"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
