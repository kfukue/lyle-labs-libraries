package gethlyleminers

import (
	"time"
)

type GethMinerTransactionInput struct {
	MinerID            *int      `json:"minerId" db:"miner_id"`                        //1
	TransactionInputID *int      `json:"transactionInputId" db:"transaction_input_id"` //2
	UUID               string    `json:"uuid" db:"uuid"`                               //3
	Name               string    `json:"name" db:"name"`                               //4
	AlternateName      string    `json:"alternateName" db:"alternate_name"`            //5
	Description        string    `json:"description" db:"description"`                 //6
	CreatedBy          string    `json:"createdBy" db:"created_by"`                    //7
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`                    //8
	UpdatedBy          string    `json:"updatedBy" db:"updated_by"`                    //9
	UpdatedAt          time.Time `json:"updatedAt" db:"updated_at"`                    //10
}
