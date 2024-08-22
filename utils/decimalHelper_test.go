package utils

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestConvertFloatToDecimal(t *testing.T) {
	// Test case 1: Valid float64 value
	val := 123.456
	expected := decimal.NewFromFloat32(float32(val))
	result := ConvertFloatToDecimal(&val)
	if result == nil || !result.Equal(expected) {
		t.Errorf("ConvertFloatToDecimal(%v) = %v; want %v", val, result, expected)
	}

	// Test case 2: Nil value
	var nilVal *float64
	result = ConvertFloatToDecimal(nilVal)
	if result != nil {
		t.Errorf("ConvertFloatToDecimal(nil) = %v; want nil", result)
	}
}
