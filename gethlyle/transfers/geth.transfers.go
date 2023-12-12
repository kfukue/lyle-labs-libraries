package gethlyletransfers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type GethTransfer struct {
	ID               *int             `json:"id"`
	UUID             string           `json:"uuid"`
	ChainID          *int             `json:"chainId"`
	TokenAddress     string           `json:"tokenAddress"`
	TokenAddressID   *int             `json:"tokenAddressId"`
	AssetID          *int             `json:"assetId"`
	BlockNumber      *uint64          `json:"blockNumber"`
	IndexNumber      *uint            `json:"indexNumber"`
	TransferDate     time.Time        `json:"transferDate"`
	TxnHash          string           `json:"txnHash"`
	SenderAddress    string           `json:"senderAddress"`
	SenderAddressID  *int             `json:"senderAddressID"`
	ToAddress        string           `json:"toAddress"`
	ToAddressID      *int             `json:"toAddressID"`
	Amount           *decimal.Decimal `json:"amount"`
	Description      string           `json:"description"`
	CreatedBy        string           `json:"createdBy"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedBy        string           `json:"updatedBy"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	GethProcessJobID *int             `json:"gethProcessJobId"`
	TopicsStr        []string         `json:"topicsStr"`
	StatusID         *int             `json:"statusId"`
	BaseAssetID      *int             `json:"baseAssetId"`
	TransferTypeID   *int             `json:"transferTypeId"`
}

type GethTransferAudit struct {
	GethTransfer
	GethTransferAuditId  *int `json:"gethTransferAuditId" db:"geth_transfer_audit_id"`
	GethProcessVlogJobID *int `json:"gethProcessVlogJobId" db:"geth_process_vlog_job_id"`
	InsertTypeID         *int `json:"insertTypeId" db:"insert_type_id"`
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
