package utils

import (
	"testing"

	"github.com/jackc/pgtype"
	shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"
	decimal "github.com/shopspring/decimal"
)

func TestConvertFloatToPgTypeNumeric(t *testing.T) {
	value := 3.14
	expected := &pgtype.DataType{
		Value: &shopspring.Numeric{
			Decimal: decimal.NewFromFloat32(float32(value)),
		},
		Name: "numeric",
		OID:  pgtype.NumericOID,
	}

	result := ConvertFloatToPgTypeNumeric(&value)

	if result == nil {
		t.Errorf("ConvertFloatToPgTypeNumeric(%v) = nil; want %v", value, expected)
	} else if result.Name != expected.Name || result.OID != expected.OID || !result.Value.(*shopspring.Numeric).Decimal.Equal(expected.Value.(*shopspring.Numeric).Decimal) {
		t.Errorf("ConvertFloatToPgTypeNumeric(%v) = %v; want %v", value, result, expected)
	}

	expected = nil

	result = ConvertFloatToPgTypeNumeric(nil)

	if result != expected {
		t.Errorf("ConvertFloatToPgTypeNumeric(%v) = %v; want %v", value, result, expected)
	}
}
