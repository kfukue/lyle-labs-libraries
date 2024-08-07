package source

import (
	"time"
)

// Source
type Source struct {
	ID            *int      `json:"id" db:"id"`                        //1
	UUID          string    `json:"uuid" db:"uuid"`                    //2
	Name          string    `json:"name" db:"name"`                    //3
	AlternateName string    `json:"alternateName" db:"alternate_name"` //4
	URL           string    `json:"url" db:"url"`                      //5
	Ticker        string    `json:"ticker" db:"ticker"`                //6
	Description   string    `json:"description" db:"description"`      //7
	CreatedBy     string    `json:"createdBy" db:"created_by"`         //8
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`         //9
	UpdatedBy     string    `json:"updatedBy" db:"updated_by"`         //10
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`         //11
}
