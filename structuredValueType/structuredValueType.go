package structuredvaluetype

import (
	"time"
)

// StructuredValueType
type StructuredValueType struct {
	ID            *int      `json:"id"`
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	AlternateName string    `json:"alternateName"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedBy     string    `json:"updatedBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
