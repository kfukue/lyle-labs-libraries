package structuredvalue

import (
	"time"
)

// StructuredValue
type StructuredValue struct {
	ID                    *int      `json:"id" db:"id"`                                           //1
	UUID                  string    `json:"uuid" db:"uuid"`                                       //2
	Name                  string    `json:"name" db:"name"`                                       //3
	AlternateName         string    `json:"alternateName" db:"alternate_name"`                    //4
	StructuredValueTypeID *int      `json:"structuredValueTypeId"  db:"structured_value_type_id"` //5
	CreatedBy             string    `json:"createdBy" db:"created_by"`                            //6
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`                            //7
	UpdatedBy             string    `json:"updatedBy" db:"updated_by"`                            //8
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`                            //9
}
