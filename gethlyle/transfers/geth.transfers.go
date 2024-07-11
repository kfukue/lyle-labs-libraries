package gethlyletransfers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type GethTransfer struct {
	ID               *int             `json:"id" db:"id"`                                //1
	UUID             string           `json:"uuid" db:"uuid"`                            //2
	ChainID          *int             `json:"chainId" db:"chain_id"`                     //3
	TokenAddress     string           `json:"tokenAddress" db:"token_address"`           //4
	TokenAddressID   *int             `json:"tokenAddressId" db:"token_address_id"`      //5
	AssetID          *int             `json:"assetId" db:"asset_id"`                     //6
	BlockNumber      *uint64          `json:"blockNumber" db:"block_number"`             //7
	IndexNumber      *uint            `json:"indexNumber" db:"index_number"`             //8
	TransferDate     time.Time        `json:"transferDate" db:"transfer_date"`           //9
	TxnHash          string           `json:"txnHash" db:"txn_hash"`                     //10
	SenderAddress    string           `json:"senderAddress" db:"sender_address"`         //11
	SenderAddressID  *int             `json:"senderAddressID" db:"sender_address_id"`    //12
	ToAddress        string           `json:"toAddress" db:"to_address"`                 //13
	ToAddressID      *int             `json:"toAddressID" db:"to_address_id"`            //14
	Amount           *decimal.Decimal `json:"amount" db:"amount"`                        //15
	Description      string           `json:"description" db:"description"`              //16
	CreatedBy        string           `json:"createdBy" db:"created_by"`                 //17
	CreatedAt        time.Time        `json:"createdAt" db:"created_at"`                 //18
	UpdatedBy        string           `json:"updatedBy" db:"updated_by"`                 //19
	UpdatedAt        time.Time        `json:"updatedAt" db:"updated_at"`                 //20
	GethProcessJobID *int             `json:"gethProcessJobId" db:"geth_process_job_id"` //21
	TopicsStr        []string         `json:"topicsStr" db:"topics_str"`                 //22
	StatusID         *int             `json:"statusId" db:"status_id"`                   //23
	BaseAssetID      *int             `json:"baseAssetId" db:"base_asset_id"`            //24
	TransferTypeID   *int             `json:"transferTypeId" db:"transfer_type_id"`      //25
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
