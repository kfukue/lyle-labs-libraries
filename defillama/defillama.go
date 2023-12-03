package defillama

// Results DefiLlama Coin API (Using coingecko id);
type DefiLlamaCurrentPriceResults struct {
	Coins map[string]DefiLlamaCurrentPriceCoinResults `json:"coins"`
}
type DefiLlamaCurrentPriceCoinResults struct {
	Price      *float64 `json:"price"`
	Symbol     string   `json:"symbol"`
	TimeStamp  *float64 `json:"timestamp"`
	Confidence *float64 `json:"confidence"`
}

type DefiLlamaClosestBlockResults struct {
	Height    *float64 `json:"height"`
	TimeStamp *float64 `json:"timestamp"`
}
