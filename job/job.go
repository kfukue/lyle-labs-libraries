package job

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Job
type Job struct {
	ID               *int      `json:"id" db:"id"`                               //1
	UUID             string    `json:"uuid" db:"uuid"`                           //2
	Name             string    `json:"name" db:"name"`                           //3
	AlternateName    string    `json:"alternateName" db:"alternate_name"`        //4
	StartDate        time.Time `json:"startDate" db:"start_date"`                //5
	EndDate          time.Time `json:"endDate" db:"end_date"`                    //6
	Description      string    `json:"description" db:"description"`             //7
	StatusID         *int      `json:"statusId" db:"status_id"`                  //8
	ResponseStatus   string    `json:"responseStatus" db:"response_status"`      //9
	RequestUrl       string    `json:"requestUrl" db:"request_url"`              //10
	RequestBody      string    `json:"requestBody" db:"request_body"`            //11
	RequestMethod    string    `json:"requestMethod" db:"request_method"`        //12
	ResponseData     string    `json:"responseData" db:"response_data"`          //13
	ResponseDataJson Attrs     `json:"responseDataJson" db:"response_data_json"` //14
	JobCategoryID    *int      `json:"jobCategoryId" db:"job_category_id"`       //15
	CreatedBy        string    `json:"createdBy" db:"created_by"`                //16
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`                //17
	UpdatedBy        string    `json:"updatedBy" db:"updated_by"`                //18
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`                //19
}

type Attrs map[string]interface{}

func (a Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
