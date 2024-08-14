package step

import (
	"time"
)

type Step struct {
	ID            *int      `json:"id" db:"id"`                        //1
	PoolID        *int      `json:"poolId" db:"pool_id"`               //2
	ParentStepId  *int      `json:"parentStepId" db:"parent_step_id"`  //3
	UUID          string    `json:"uuid" db:"uuid"`                    //4
	Name          string    `json:"name" db:"name"`                    //5
	AlternateName string    `json:"alternateName" db:"alternate_name"` //6
	StartDate     time.Time `json:"startDate" db:"start_date"`         //7
	EndDate       time.Time `json:"endDate" db:"end_date"`             //8
	Description   string    `json:"description" db:"description"`      //9
	ActionTypeID  *int      `json:"actionTypeId" db:"action_type_id"`  //10
	FunctionName  string    `json:"functionName" db:"function_name"`   //11
	StepOrder     *int      `json:"stepOrder"  db:"step_order"`        //12
	CreatedBy     string    `json:"createdBy" db:"created_by"`         //13
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`         //14
	UpdatedBy     string    `json:"updatedBy" db:"updated_by"`         //15
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`         //16
}
