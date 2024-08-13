package gethlyletransactions

import (
	"time"
)

type GethTransactionInput struct {
	ID              *int      `json:"id" db:"id"`                             //1
	UUID            string    `json:"uuid" db:"uuid"`                         //2
	Name            string    `json:"name" db:"name"`                         //3
	AlternateName   string    `json:"alternateName" db:"alternate_name"`      //4
	FunctionName    string    `json:"functionName" db:"function_name"`        //5
	MethodIDStr     string    `json:"methodIdStr" db:"method_id_str"`         //6
	NumOfParameters *int      `json:"numOfParameters" db:"num_of_parameters"` //7
	Description     string    `json:"description" db:"description"`           //8
	CreatedBy       string    `json:"createdBy" db:"created_by"`              //9
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`              //10
	UpdatedBy       string    `json:"updatedBy" db:"updated_by"`              //11
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`              //12
}
