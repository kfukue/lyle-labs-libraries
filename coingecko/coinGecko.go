//go:build !skip
// +build !skip

package coingecko

// assets response struct from coinGecko pro api
type CoinGeckoAssetsResponse struct {
	Error  []interface{}                     `json:"error"`
	Result map[string]CoinGeckoAssetsResults `json:"result"`
}

// Results portion of assets
type CoinGeckoAssetsResults struct {
	AClass          string   `json:"aclass"`
	AltName         string   `json:"altname"`
	Decimals        *float64 `json:"decimals"`
	DisplayDecimals *float64 `json:"display_decimals"`
}

type CoinGeckoAssetPairsResponse struct {
	Error  []interface{}                         `json:"error"`
	Result map[string]CoinGeckoAssetPairsResults `json:"result"`
}

// Results portion of assets
type CoinGeckoAssetPairsResults struct {
	AltName           string      `json:"altname"`
	WSName            string      `json:"wsname"`
	AClassBase        *float64    `json:"aclass_base"`
	Base              string      `json:"base"`
	AClassQuote       string      `json:"aclass_quote"`
	Quote             string      `json:"quote"`
	Lot               string      `json:"lot"`
	PairDecimals      *float64    `json:"pair_decimals"`
	LotDecimals       *float64    `json:"lot_decimals"`
	LotMultiplier     *float64    `json:"lot_multiplier"`
	LeverageBuy       []float64   `json:"leverage_buy"`
	LeverageSell      []float64   `json:"leverage_sell"`
	Fees              [][]float64 `json:"fees"`
	FeesMaker         [][]float64 `json:"fees_maker"`
	FeeVolumeCurrency string      `json:"fee_volume_currency"`
	MarginCall        *float64    `json:"margin_call"`
	MarginStop        *float64    `json:"margin_stop"`
	OrderMin          string      `json:"ordermin"`
}
