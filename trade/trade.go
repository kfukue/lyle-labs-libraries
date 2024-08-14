package trade

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

// Asset
type Trade struct {
	ID                      *int           `json:"id" db:"id"`                                           //1
	ParentTradeID           *int           `json:"parentTradeId" db:"parent_trade_id"`                   //2
	FromAccountID           *int           `json:"fromAccountId" db:"from_account_id"`                   //3
	ToAccountID             *int           `json:"toAccountId" db:"to_account_id"`                       //4
	AssetID                 *int           `json:"assetId" db:"asset_id"`                                //5
	SourceID                *int           `json:"sourceId" db:"source_id"`                              //6
	UUID                    string         `json:"uuid" db:"uuid"`                                       //7
	TransactionIDFromSource string         `json:"transactionIdFromSource" db:"transaction_id"`          //8
	OrderIDFromSource       string         `json:"orderIdFromSource" db:"order_id"`                      //9
	TradeIDFromSource       string         `json:"tradeIdFromSource" db:"trade_id"`                      //10
	Name                    string         `json:"name" db:"name"`                                       //11
	AlternateName           string         `json:"alternateName" db:"alternate_name"`                    //12
	TradeTypeID             *int           `json:"tradeTypeId" db:"trade_type_id"`                       //13
	TradeDate               *time.Time     `json:"tradeDate" db:"trade_date"`                            //14
	SettleDate              *time.Time     `json:"settleDate" db:"settle_date"`                          //15
	TransferDate            *time.Time     `json:"transferDate" db:"transfer_date"`                      //16
	FromQuantity            *float64       `json:"fromQuantity" db:"from_quantity"`                      //17
	ToQuantity              *float64       `json:"toQuantity" db:"to_quantity"`                          //18
	Price                   *float64       `json:"price"  db:"price"`                                    //19
	TotalAmount             *float64       `json:"totalAmount"  db:"total_amount"`                       //20
	FeesAmount              *float64       `json:"feesAmount"  db:"fees_amount"`                         //21
	FeesAssetID             *int           `json:"feesAssetId"  db:"fees_asset_id"`                      //22
	RealizedReturnAmount    *float64       `json:"realizedReturnAmount"  db:"realized_return_amount"`    //23
	RealizedReturnAssetID   *int           `json:"realizedReturnAssetId"  db:"realized_return_asset_id"` //24
	CostBasisAmount         *float64       `json:"costBasisAmount"  db:"cost_basis_amount"`              //25
	CostBasisTradeID        *int           `json:"costBasisTradeId"  db:"cost_basis_trade_id"`           //26
	Description             string         `json:"description" db:"description"`                         //27
	IsActive                bool           `json:"isActive"  db:"is_active"`                             //28
	SourceData              types.JSONText `json:"sourceData"  db:"source_data"`                         //29
	CreatedBy               string         `json:"createdBy" db:"created_by"`                            //30
	CreatedAt               time.Time      `json:"createdAt" db:"created_at"`                            //31
	UpdatedBy               string         `json:"updatedBy" db:"updated_by"`                            //32
	UpdatedAt               time.Time      `json:"updatedAt" db:"updated_at"`                            //33
}

type MinMaxTradeDates struct {
	MinTradeDate *time.Time
	MaxTradeDate *time.Time
}
