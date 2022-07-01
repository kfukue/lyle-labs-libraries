package sourcejob

import (
	"time"
)

// Asset
type SourceJob struct {
	SourceID    *int      `json:"sourceId"`
	JobID       *int      `json:"jobId"`
	UUID        string    `json:"uuid"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedBy   string    `json:"updatedBy"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
