package structuredvaluetype

import (
	"time"
)

// StructuredValueType
type StructuredValueType struct {
	ID            *int      `json:"id" db:"id"`                        //1
	UUID          string    `json:"uuid" db:"uuid"`                    //2
	Name          string    `json:"name" db:"name"`                    //3
	AlternateName string    `json:"alternateName" db:"alternate_name"` //4
	CreatedBy     string    `json:"createdBy" db:"created_by"`         //5
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`         //6
	UpdatedBy     string    `json:"updatedBy" db:"updated_by"`         //7
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`         //8
}
