package gethlyletransactions

import (
	"time"
)

type GethTransactionInput struct {
	ID              *int      `json:"id"`
	UUID            string    `json:"uuid"`
	Name            string    `json:"name"`
	AlternateName   string    `json:"alternateName"`
	FunctionName    string    `json:"functionName"`
	MethodIDStr     string    `json:"methodIdStr"`
	NumOfParameters *int      `json:"numOfParameters"`
	Description     string    `json:"description"`
	CreatedBy       string    `json:"createdBy"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedBy       string    `json:"updatedBy"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
