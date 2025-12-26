package aiModels

import (
	"time"
)

// AIModel
type AIModel struct {
	ID            *int      `json:"id" db:"id"`                        //1
	UUID          string    `json:"uuid" db:"uuid"`                    //2
	Name          string    `json:"name" db:"name"`                    //3
	AlternateName string    `json:"alternateName" db:"alternate_name"` //4
	URL           string    `json:"url" db:"url"`                      //5
	Ticker        string    `json:"ticker" db:"ticker"`                //6
	Description   string    `json:"description" db:"description"`      //7
	OllamaName    string    `json:"ollamaName" db:"ollama_name"`       //8
	ParamsSize    int64     `json:"paramsSize" db:"params_size"`       //9
	QuantizedSize string    `json:"quantizedSize" db:"quantiized_size"`//10
	BaseModelID   *int      `json:"baseModelId" db:"base_model_id"`    //11
	CreatedBy     string    `json:"createdBy" db:"created_by"`         //12
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`         //13
	UpdatedBy     string    `json:"updatedBy" db:"updated_by"`         //14
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`         //15
}
