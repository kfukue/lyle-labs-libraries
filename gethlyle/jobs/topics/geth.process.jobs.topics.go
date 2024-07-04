package gethlylejobstopics

import (
	"time"
)

type GethProcessJobTopic struct {
	ID               *int      `json:"id" db:"id"`                                //1
	GethProcessJobID *int      `json:"gethProcessJobId" db:"geth_process_job_id"` //2
	UUID             string    `json:"uuid" db:"uuid"`                            //3
	Name             string    `json:"name" db:"name"`                            //4
	AlternateName    string    `json:"alternateName" db:"alternate_name"`         //5
	Description      string    `json:"description" db:"description"`              //6
	StatusID         *int      `json:"statusId" db:"status_id"`                   //7
	TopicStr         string    `json:"topicStr" db:"topic_str"`                   //8
	CreatedBy        string    `json:"createdBy" db:"created_by"`                 //9
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`                 //10
	UpdatedBy        string    `json:"updatedBy" db:"updated_by"`                 //11
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`                 //12

}
