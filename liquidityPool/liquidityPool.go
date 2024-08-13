package liquiditypool

import (
	"time"

	"github.com/kfukue/lyle-labs-libraries/v2/asset"
)

type LiquidityPool struct {
	ID                         *int      `json:"id" db:"id"`                                                        //1
	UUID                       string    `json:"uuid" db:"uuid"`                                                    //2
	Name                       string    `json:"name" db:"name"`                                                    //3
	AlternateName              string    `json:"alternateName" db:"alternate_name"`                                 //4
	PairAddress                string    `json:"pairAddress" db:"pair_address"`                                     //5
	ChainID                    *int      `json:"chainId" db:"chain_id"`                                             //6
	ExchangeID                 *int      `json:"exchangeId" db:"exchange_id"`                                       //7
	LiquidityPoolTypeID        *int      `json:"liquidityPoolTypeId" db:"liquidity_pool_type_id"`                   //8
	Token0ID                   *int      `json:"token0Id" db:"token0_id"`                                           //9
	Token1ID                   *int      `json:"token1Id" db:"token1_id"`                                           //10
	Url                        string    `json:"url" db:"url"`                                                      //11
	StartBlock                 *int      `json:"startBlock" db:"start_block"`                                       //12
	LatestBlockSynced          *int      `json:"latestBlockSynced" db:"latest_block_synced"`                        //13
	CreatedTxnHash             string    `json:"createdTxnHash" db:"created_txn_hash"`                              //14
	IsActive                   bool      `json:"isActive" db:"is_active"`                                           //15
	Description                string    `json:"description" db:"description"`                                      //16
	CreatedBy                  string    `json:"createdBy" db:"created_by"`                                         //17
	CreatedAt                  time.Time `json:"createdAt" db:"created_at"`                                         //18
	UpdatedBy                  string    `json:"updatedBy" db:"updated_by"`                                         //19
	UpdatedAt                  time.Time `json:"updatedAt" db:"updated_at"`                                         //20
	BaseAssetID                *int      `json:"baseAssetId" db:"base_asset_id"`                                    //21
	QuoteAssetID               *int      `json:"quoteAssetId" db:"quote_asset_id"`                                  //22
	QuoteAssetChainlinkAddress string    `json:"quoteAssetChainlinkAddress" db:"quote_asset_chainlink_address_usd"` //23
}

type LiquidityPoolWithTokens struct {
	LiquidityPool
	Token0 asset.Asset `json:"token0"`
	Token1 asset.Asset `json:"token1"`
}
