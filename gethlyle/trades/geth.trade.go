package gethlyletrades

import (
	"time"

	"github.com/kfukue/lyle-labs-libraries/asset"
	"github.com/shopspring/decimal"
)

type GethTrade struct {
	ID                     *int             `json:"id"`
	UUID                   string           `json:"uuid"`
	Name                   string           `json:"name"`
	AlternateName          string           `json:"alternateName"`
	AddressStr             string           `json:"addressStr"`
	AddressID              *int             `json:"addressId"`
	TradeDate              time.Time        `json:"tradeDate"`
	TxnHash                string           `json:"txnHash"`
	Token0Amount           *decimal.Decimal `json:"token0Amount"`
	Token0AmountDecimalAdj *decimal.Decimal `json:"token0AmountDecimalAdj"`
	Token1Amount           *decimal.Decimal `json:"token1Amount"`
	Token1AmountDecimalAdj *decimal.Decimal `json:"token1AmountDecimalAdj"`
	IsBuy                  *bool            `json:"isBuy"`
	Price                  *decimal.Decimal `json:"price"`
	PriceUSD               *decimal.Decimal `json:"priceUsd"`
	LPToken1PriceUSD       *decimal.Decimal `json:"lpToken1PriceUsd"`
	TotalAmountUSD         *decimal.Decimal `json:"totalAmountUsd"`
	Token0AssetId          *int             `json:"token0Id"`
	Token1AssetId          *int             `json:"token1Id"`
	GethProcessJobID       *int             `json:"gethProcessJobId"`
	StatusID               *int             `json:"statusId"`
	TradeTypeID            *int             `json:"tradeTypeId"`
	Description            string           `json:"description"`
	CreatedBy              string           `json:"createdBy"`
	CreatedAt              time.Time        `json:"createdAt"`
	UpdatedBy              string           `json:"updatedBy"`
	UpdatedAt              time.Time        `json:"updatedAt"`
	BaseAssetID            *int             `json:"baseAssetId"`
	OraclePriceUSD         *decimal.Decimal `json:"oraclePriceUsd"`
	OraclePriceAssetID     *int             `json:"oraclePriceAssetId"`
}

type NetTransferByAddress struct {
	TxnHash    string           `json:"txnHash"`
	AddressStr string           `json:"addressStr"`
	AssetID    *int             `json:"addressId"`
	NetAmount  *decimal.Decimal `json:"netAmount"`
	asset.Asset
}
