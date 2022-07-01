package trade

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

// Asset
type Trade struct {
	ID                      *int           `json:"id"`
	ParentTradeID           *int           `json:"parentTradeId"`
	FromAccountID           *int           `json:"fromAccountId"`
	ToAccountID             *int           `json:"toAccountId"`
	AssetID                 *int           `json:"assetId"`
	SourceID                *int           `json:"sourceId"`
	UUID                    string         `json:"uuid"`
	TransactionIDFromSource string         `json:"transactionIdFromSource"`
	OrderIDFromSource       string         `json:"orderIdFromSource"`
	TradeIDFromSource       string         `json:"tradeIdFromSource"`
	Name                    string         `json:"name"`
	AlternateName           string         `json:"alternateName"`
	TradeTypeID             *int           `json:"tradeTypeId"`
	TradeDate               time.Time      `json:"tradeDate"`
	SettleDate              time.Time      `json:"settleDate"`
	TransferDate            time.Time      `json:"transferDate"`
	FromQuantity            *float64       `json:"fromQuantity"`
	ToQuantity              *float64       `json:"toQuantity"`
	Price                   *float64       `json:"price"`
	TotalAmount             *float64       `json:"totalAmount"`
	FeesAmount              *float64       `json:"feesAmount"`
	FeesAssetID             *int           `json:"feesAssetId"`
	RealizedReturnAmount    *float64       `json:"realizedReturnAmount"`
	RealizedReturnAssetID   *int           `json:"realizedReturnAssetId"`
	CostBasisAmount         *float64       `json:"costBasisAmount"`
	CostBasisTradeID        *int           `json:"costBasisTradeId"`
	Description             string         `json:"description"`
	SourceData              types.JSONText `json:"sourceData"`
	IsActive                bool           `json:"isActive"`
	CreatedBy               string         `json:"createdBy"`
	CreatedAt               time.Time      `json:"createdAt"`
	UpdatedBy               string         `json:"updatedBy"`
	UpdatedAt               time.Time      `json:"updatedAt"`
}

type MinMaxTradeDates struct {
	MinTradeDate *time.Time
	MaxTradeDate *time.Time
}
