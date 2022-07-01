package step

import (
	"time"
)

type Step struct {
	ID            *int       `json:"id"`
	PoolID        *int       `json:"poolId"`
	ParentStepId  *int       `json:"parentStepId"`
	Name          string     `json:"name"`
	UUID          string     `json:"uuid"`
	AlternateName string     `json:"alternateName"`
	StartDate     *time.Time `json:"startDate"`
	EndDate       *time.Time `json:"endDate"`
	Description   string     `json:"description"`
	ActionTypeID  *int       `json:"actionTypeId"`
	FunctionName  string     `json:"functionName"`
	StepOrder     *int       `json:"stepOrder"`
	CreatedBy     string     `json:"createdBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedBy     string     `json:"updatedBy"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
