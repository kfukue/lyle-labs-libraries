package liquiditypool

import (
	"time"
)

type LiquidityPool struct {
	ID                *int      `json:"id"`
	UUID              string    `json:"uuid"`
	Name              string    `json:"name"`
	AlternateName     string    `json:"alternateName"`
	PairAddress       string    `json:"pairAddress"`
	ChainID           *int      `json:"chainId"`
	ExchangeID        *int      `json:"exchangeId"`
	Token0ID          *int      `json:token0`
	Token1ID          *int      `json:token1`
	Url               string    `json:"url"`
	StartBlock        *int      `json:"startBlock"`
	LatestBlockSynced *int      `json:"latestBlockSynced"`
	CreatedTxnHash    string    `json:"CreatedTxnHash"`
	IsActive          bool      `json:isActive`
	Description       string    `json:"description"`
	CreatedBy         string    `json:"createdBy"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedBy         string    `json:"updatedBy"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type LiquidityPoolAsset struct {
	UUID            string    `json:"uuid"`
	LiquidityPoolID *int      `json:"liquidityPoolId"`
	AssetID         *int      `json:"assetId"`
	TokenNumber     *int      `json:TokenNumber`
	Name            string    `json:name`
	AlternateName   string    `json:"alternateName"`
	Description     string    `json:"description"`
	CreatedBy       string    `json:"createdBy"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedBy       string    `json:"updatedBy"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
