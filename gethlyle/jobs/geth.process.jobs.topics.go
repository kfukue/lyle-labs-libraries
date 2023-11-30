package gethlylejobs

import (
	"time"
)

type GethProcessJobTopic struct {
	ID               *int      `json:"id"`
	GethProcessJobID *int      `json:"gethProcessJobId"`
	UUID             string    `json:"uuid"`
	Name             string    `json:"name"`
	AlternateName    string    `json:"alternateName"`
	Description      string    `json:"description"`
	StatusID         *int      `json:"statusId"`
	TopicStr         string    `json:"topicStr"`
	CreatedBy        string    `json:"createdBy"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedBy        string    `json:"updatedBy"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
