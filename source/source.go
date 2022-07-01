package source

import (
	"time"
)

// Source
type Source struct {
	ID            *int      `json:"id"`
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	AlternateName string    `json:"alternateName"`
	URL           string    `json:"url"`
	Ticker        string    `json:"ticker"`
	Description   string    `json:"description"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedBy     string    `json:"updatedBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
