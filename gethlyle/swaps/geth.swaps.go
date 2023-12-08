package gethlyleswaps

import (
	"time"

	"github.com/shopspring/decimal"
)

type GethSwap struct {
	ID                 *int             `json:"id"`
	UUID               string           `json:"uuid"`
	ChainID            *int             `json:"chainId"`
	ExchangeID         *int             `json:"exchangeId"`
	BlockNumber        *uint64          `json:"blockNumber"`
	IndexNumber        *uint            `json:"indexNumber"`
	SwapDate           time.Time        `json:"swapDate"`
	TradeTypeID        *int             `json:"tradeTypeId"`
	TxnHash            string           `json:"txnHash"`
	MakerAddress       string           `json:"makerAddress"`
	MakerAddressID     *int             `json:"makerAddressId"`
	IsBuy              *bool            `json:"isBuy"`
	Price              *decimal.Decimal `json:"price"`
	PriceUSD           *decimal.Decimal `json:"priceUsd"`
	Token1PriceUSD     *decimal.Decimal `json:"token1PriceUsd"`
	TotalAmountUSD     *decimal.Decimal `json:"totalAmountUsd"`
	PairAddress        string           `json:"pairAddress"`
	LiquidityPoolID    *int             `json:"liquidityPoolId"`
	Token0AssetId      *int             `json:"token0Id"`
	Token1AssetId      *int             `json:"token1Id"`
	Token0Amount       *decimal.Decimal `json:"token0Amount"`
	Token1Amount       *decimal.Decimal `json:"token1Amount"`
	Description        string           `json:"description"`
	CreatedBy          string           `json:"createdBy"`
	CreatedAt          time.Time        `json:"createdAt"`
	UpdatedBy          string           `json:"updatedBy"`
	UpdatedAt          time.Time        `json:"updatedAt"`
	GethProcessJobID   *int             `json:"gethProcessJobId"`
	TopicsStr          []string         `json:"topicsStr"`
	StatusID           *int             `json:"statusId"`
	BaseAssetID        *int             `json:"baseAssetId"`
	OraclePriceUSD     *decimal.Decimal `json:"oraclePriceUsd"`
	OraclePriceAssetID *int             `json:"oraclePriceAssetId"`
}

type GethSwapAudit struct {
	GethSwap
	GethSwapAuditId      *int `json:"gethSwapAuditId"`
	GethProcessVlogJobID *int `json:"gethProcessVlogJobId"`
	InsertTypeID         *int `json:"insertTypeId"`
}
