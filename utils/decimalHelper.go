package utils

import (
	decimal "github.com/shopspring/decimal"
)

func ConvertFloatToDecimal(value *float64) *decimal.Decimal {
	if value != nil {
		results := decimal.NewFromFloat32(float32(*value))
		return &results
	} else {
		return nil
	}
}
