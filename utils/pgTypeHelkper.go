package utils

import (
	"github.com/jackc/pgtype"
	shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"
	decimal "github.com/shopspring/decimal"
)

func ConvertFloatToPgTypeNumeric(value *float64) *pgtype.DataType {
	numericObj := &shopspring.Numeric{}
	if value != nil {
		numericObj.Decimal = decimal.NewFromFloat32(float32(*value))
		results := &pgtype.DataType{
			Value: numericObj,
			Name:  "numeric",
			OID:   pgtype.NumericOID,
		}
		return results
	} else {
		return nil
	}
}
