package liquiditypool

import (
	"time"

	"github.com/kfukue/lyle-labs-libraries/asset"
)

type LiquidityPool struct {
	ID                         *int      `json:"id"`
	UUID                       string    `json:"uuid"`
	Name                       string    `json:"name"`
	AlternateName              string    `json:"alternateName"`
	PairAddress                string    `json:"pairAddress"`
	ChainID                    *int      `json:"chainId"`
	ExchangeID                 *int      `json:"exchangeId"`
	LiquidityPoolTypeID        *int      `json:"liquidityPoolTypeId"`
	Token0ID                   *int      `json:"token0Id"`
	Token1ID                   *int      `json:"token1Id"`
	Url                        string    `json:"url"`
	StartBlock                 *int      `json:"startBlock"`
	LatestBlockSynced          *int      `json:"latestBlockSynced"`
	CreatedTxnHash             string    `json:"createdTxnHash"`
	IsActive                   bool      `json:"isActive"`
	Description                string    `json:"description"`
	CreatedBy                  string    `json:"createdBy"`
	CreatedAt                  time.Time `json:"createdAt"`
	UpdatedBy                  string    `json:"updatedBy"`
	UpdatedAt                  time.Time `json:"updatedAt"`
	BaseAssetID                *int      `json:"baseAssetId"`
	QuoteAssetID               *int      `json:"quoteAssetId"`
	QuoteAssetChainlinkAddress string    `json:"quoteAssetChainlinkAddress"`
}

type LiquidityPoolWithTokens struct {
	LiquidityPool
	Token0 asset.Asset `json:"token0"`
	Token1 asset.Asset `json:"token1"`
}

type LiquidityPoolAsset struct {
	UUID            string    `json:"uuid"`
	LiquidityPoolID *int      `json:"liquidityPoolId"`
	AssetID         *int      `json:"assetId"`
	TokenNumber     *int      `json:"tokenNumber"`
	Name            string    `json:"name"`
	AlternateName   string    `json:"alternateName"`
	Description     string    `json:"description"`
	CreatedBy       string    `json:"createdBy"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedBy       string    `json:"updatedBy"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
