package sourcejob

import (
	"time"
)

// SourceJob
type SourceJob struct {
	SourceID    *int      `json:"sourceId" db:"source_id"`      //1
	JobID       *int      `json:"jobId" db:"job_id"`            //2
	UUID        string    `json:"uuid" db:"uuid"`               //3
	Description string    `json:"description" db:"description"` //4
	CreatedBy   string    `json:"createdBy" db:"created_by"`    //5
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`    //6
	UpdatedBy   string    `json:"updatedBy" db:"updated_by"`    //7
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`    //8
}
