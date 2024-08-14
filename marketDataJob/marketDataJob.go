package marketdatajob

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Asset
type MarketDataJob struct {
	MarketDataID     *int      `json:"marketDataId" db:"market_data_id"`         //1
	JobID            *int      `json:"jobId" db:"job_id"`                        //2
	UUID             string    `json:"uuid" db:"uuid"`                           //3
	Name             string    `json:"name" db:"name"`                           //4
	AlternateName    string    `json:"alternateName" db:"alternate_name"`        //5
	StartDate        time.Time `json:"startDate" db:"start_date"`                //6
	EndDate          time.Time `json:"endDate" db:"end_date"`                    //7
	Description      string    `json:"description" db:"description"`             //8
	StatusID         *int      `json:"statusId" db:"status_id"`                  //9
	ResponseStatus   string    `json:"responseStatus" db:"response_status"`      //10
	RequestUrl       string    `json:"requestUrl" db:"request_url"`              //11
	RequestBody      string    `json:"requestBody" db:"request_body"`            //12
	RequestMethod    string    `json:"requestMethod" db:"request_method"`        //13
	ResponseData     string    `json:"responseData" db:"response_data"`          //14
	ResponseDataJson Attrs     `json:"responseDataJson" db:"response_data_json"` //15
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
