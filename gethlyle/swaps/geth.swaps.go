package gethlyleswaps

import (
	"time"

	"github.com/shopspring/decimal"
)

type GethSwap struct {
	ID                 *int             `json:"id" db:"id"`                                    //1
	UUID               string           `json:"uuid" db:"uuid"`                                //2
	ChainID            *int             `json:"chainId" db:"chain_id"`                         //3
	ExchangeID         *int             `json:"exchangeId" db:"exchange_id"`                   //4
	BlockNumber        *uint64          `json:"blockNumber" db:"block_number"`                 //5
	IndexNumber        *uint            `json:"indexNumber" db:"index_number"`                 //6
	SwapDate           *time.Time       `json:"swapDate" db:"swap_date"`                       //7
	TradeTypeID        *int             `json:"tradeTypeId" db:"trade_type_id"`                //8
	TxnHash            string           `json:"txnHash" db:"txn_hash"`                         //9
	MakerAddress       string           `json:"makerAddress" db:"maker_address"`               //10
	MakerAddressID     *int             `json:"makerAddressId" db:"maker_address_id"`          //11
	IsBuy              *bool            `json:"isBuy" db:"is_buy"`                             //12
	Price              *decimal.Decimal `json:"price" db:"price"`                              //13
	PriceUSD           *decimal.Decimal `json:"priceUsd" db:"price_usd"`                       //14
	Token1PriceUSD     *decimal.Decimal `json:"token1PriceUsd" db:"token1_price_usd"`          //15
	TotalAmountUSD     *decimal.Decimal `json:"totalAmountUsd" db:"total_amount_usd"`          //16
	PairAddress        string           `json:"pairAddress" db:"pair_address"`                 //17
	LiquidityPoolID    *int             `json:"liquidityPoolId" db:"liquidity_pool_id"`        //18
	Token0AssetId      *int             `json:"token0Id" db:"token0_asset_id"`                 //19
	Token1AssetId      *int             `json:"token1Id" db:"token1_asset_id"`                 //20
	Token0Amount       *decimal.Decimal `json:"token0Amount" db:"token0_amount"`               //21
	Token1Amount       *decimal.Decimal `json:"token1Amount" db:"token1_amount"`               //22
	Description        string           `json:"description" db:"description"`                  //23
	CreatedBy          string           `json:"createdBy" db:"created_by"`                     //24
	CreatedAt          time.Time        `json:"createdAt" db:"created_at"`                     //25
	UpdatedBy          string           `json:"updatedBy" db:"updated_by"`                     //26
	UpdatedAt          time.Time        `json:"updatedAt" db:"updated_at"`                     //27
	GethProcessJobID   *int             `json:"gethProcessJobId" db:"geth_process_job_id"`     //28
	TopicsStr          []string         `json:"topicsStr" db:"topics_str"`                     //29
	StatusID           *int             `json:"statusId" db:"status_id"`                       //30
	BaseAssetID        *int             `json:"baseAssetId" db:"base_asset_id"`                //31
	OraclePriceUSD     *decimal.Decimal `json:"oraclePriceUsd" db:"oracle_price_usd"`          //32
	OraclePriceAssetID *int             `json:"oraclePriceAssetId" db:"oracle_price_asset_id"` //33
}
