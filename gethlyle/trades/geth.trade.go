package gethlyletrades

import (
	"time"

	"github.com/kfukue/lyle-labs-libraries/asset"
	"github.com/shopspring/decimal"
)

type GethTrade struct {
	ID                     *int             `json:"id" db:"id"`                                            //1
	UUID                   string           `json:"uuid" db:"uuid"`                                        //2
	Name                   string           `json:"name" db:"name"`                                        //3
	AlternateName          string           `json:"alternateName" db:"alternate_name"`                     //4
	AddressStr             string           `json:"addressStr" db:"address_str"`                           //5
	AddressID              *int             `json:"addressId" db:"address_id"`                             //6
	TradeDate              time.Time        `json:"tradeDate" db:"trade_date"`                             //7
	TxnHash                string           `json:"txnHash" db:"txn_hash"`                                 //8
	Token0Amount           *decimal.Decimal `json:"token0Amount" db:"token0_amount"`                       //9
	Token0AmountDecimalAdj *decimal.Decimal `json:"token0AmountDecimalAdj" db:"token0_amount_decimal_adj"` //10
	Token1Amount           *decimal.Decimal `json:"token1Amount" db:"token1_amount"`                       //11
	Token1AmountDecimalAdj *decimal.Decimal `json:"token1AmountDecimalAdj" db:"token1_amount_decimal_adj"` //12
	IsBuy                  *bool            `json:"isBuy" db:"is_buy"`                                     //13
	Price                  *decimal.Decimal `json:"price" db:"price"`                                      //14
	PriceUSD               *decimal.Decimal `json:"priceUsd" db:"price_usd"`                               //15
	LPToken1PriceUSD       *decimal.Decimal `json:"lpToken1PriceUsd" db:"lp_token1_price_usd"`             //16
	TotalAmountUSD         *decimal.Decimal `json:"totalAmountUsd" db:"total_amount_usd"`                  //17
	Token0AssetId          *int             `json:"token0Id" db:"token0_asset_id"`                         //18
	Token1AssetId          *int             `json:"token1Id" db:"token1_asset_id"`                         //19
	GethProcessJobID       *int             `json:"gethProcessJobId" db:"geth_process_job_id"`             //20
	StatusID               *int             `json:"statusId" db:"status_id"`                               //21
	TradeTypeID            *int             `json:"tradeTypeId" db:"trade_type_id"`                        //22
	Description            string           `json:"description" db:"description"`                          //23
	CreatedBy              string           `json:"createdBy" db:"created_by"`                             //24
	CreatedAt              time.Time        `json:"createdAt" db:"created_at"`                             //25
	UpdatedBy              string           `json:"updatedBy" db:"updated_by"`                             //26
	UpdatedAt              time.Time        `json:"updatedAt" db:"updated_at"`                             //27
	BaseAssetID            *int             `json:"baseAssetId" db:"base_asset_id"`                        //28
	OraclePriceUSD         *decimal.Decimal `json:"oraclePriceUsd" db:"oracle_price_usd"`                  //29
	OraclePriceAssetID     *int             `json:"oraclePriceAssetId" db:"oracle_price_asset_id"`         //30
}

type NetTransferByAddress struct {
	TxnHash    string           `json:"txnHash" db:"txn_hash"`       //1
	AddressStr string           `json:"addressStr" db:"address_str"` //2
	AssetID    *int             `json:"addressId" db:"asset_id"`     //3
	NetAmount  *decimal.Decimal `json:"netAmount" db:"net_amount"`   //4
	asset.Asset
}
