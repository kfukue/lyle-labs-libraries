package transactionjob

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Asset
type TransactionJob struct {
	TransactionID    *int      `json:"transactionId"`    //1
	JobID            *int      `json:"jobId"`            //2
	UUID             string    `json:"uuid"`             //3
	Name             string    `json:"name"`             //4
	AlternateName    string    `json:"alternateName"`    //5
	StartDate        time.Time `json:"startDate"`        //6
	EndDate          time.Time `json:"endDate"`          //7
	Description      string    `json:"description"`      //8
	StatusID         *int      `json:"statusId"`         //9
	ResponseStatus   string    `json:"responseStatus"`   //10
	RequestUrl       string    `json:"requestUrl"`       //11
	RequestBody      string    `json:"requestBody"`      //12
	RequestMethod    string    `json:"requestMethod"`    //13
	ResponseData     string    `json:"responseData"`     //14
	ResponseDataJson Attrs     `json:"responseDataJson"` //15
	CreatedBy        string    `json:"createdBy"`        //16
	CreatedAt        time.Time `json:"createdAt"`        //17
	UpdatedBy        string    `json:"updatedBy"`        //18
	UpdatedAt        time.Time `json:"updatedAt"`        //19
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
