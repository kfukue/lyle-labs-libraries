package biginthelper

import (
	"database/sql/driver"
	"fmt"
	"math/big"

	"github.com/jackc/pgtype"
)

type BigInt big.Int

func ConvertToBigInt(b *big.Int) (*BigInt, error) {
	var convertedBigInt BigInt
	if b != nil {
		err := convertedBigInt.Scan(b.String())
		return &convertedBigInt, err
	}
	return nil, nil
}

func ConvertToBigIntFromStr(b string) (*BigInt, error) {
	var convertedBigInt BigInt
	if b != "" {
		err := convertedBigInt.Scan(b)
		return &convertedBigInt, err
	}
	return nil, nil
}
func (b *BigInt) Value() (driver.Value, error) {
	if b != nil {
		return (*big.Int)(b).String(), nil
	}
	return nil, nil
}

func (b *BigInt) RawValue() (*big.Int, error) {
	if b != nil {
		return (*big.Int)(b), nil
	}
	return nil, nil
}

func (b *BigInt) CreatePgNumeric() (*pgtype.Numeric, error) {
	if b != nil {
		pgTypeNumeric := &pgtype.Numeric{}
		bDriveValue, _ := b.Value()
		pgTypeNumeric.Set(bDriveValue)
		return pgTypeNumeric, nil
	} else {
		pgTypeNumeric := &pgtype.Numeric{Status: pgtype.Null}
		return pgTypeNumeric, nil
	}
}

func (b *BigInt) Scan(value interface{}) error {
	if value == nil {
		b = nil
	}

	switch t := value.(type) {
	case string:
		_, ok := (*big.Int)(b).SetString(value.(string), 10)
		if !ok {
			return fmt.Errorf("failed to load value to []uint8: %v", value)
		}
	case []uint8:
		_, ok := (*big.Int)(b).SetString(string(value.([]uint8)), 10)
		if !ok {
			return fmt.Errorf("failed to load value to []uint8: %v", value)
		}
	default:
		return fmt.Errorf("could not scan type %T into BigInt", t)
	}

	return nil
}

type BigFloat big.Float

func ConvertToBigFloat(b *big.Float) (*BigFloat, error) {
	var convertedBigFloat BigFloat
	if b != nil {
		err := convertedBigFloat.Scan(b.String())
		return &convertedBigFloat, err
	}
	return nil, nil
}
func (b *BigFloat) Value() (driver.Value, error) {
	if b != nil {
		return fmt.Sprintf("%f", (*big.Float)(b)), nil
	}
	return nil, nil
}
func (b *BigFloat) RawValue() *big.Float {
	if b != nil {
		return (*big.Float)(b)
	}
	return nil
}
func (b *BigFloat) CreatePgNumeric() (*pgtype.Numeric, error) {
	if b != nil {
		pgTypeNumeric := &pgtype.Numeric{}
		bDriveValue, _ := b.Value()
		pgTypeNumeric.Set(bDriveValue)
		return pgTypeNumeric, nil
	} else {
		pgTypeNumeric := &pgtype.Numeric{Status: pgtype.Null}
		return pgTypeNumeric, nil
	}
}

func (b *BigFloat) Scan(value interface{}) error {
	// fmt.Printf("Scan value: %v\n", value)
	if value == nil {
		b = nil
	}

	switch t := value.(type) {
	case string:
		// _, ok := (*big.Int)(b).SetString(value.(string), 10)
		_, ok := (*big.Float)(b).SetString(value.(string))
		if !ok {
			return fmt.Errorf("failed to load value to []uint8: %v", value)
		}
	case float64:
		(*big.Float)(b).SetFloat64(value.(float64))
	case []uint8:
		_, ok := (*big.Float)(b).SetString(string(value.([]uint8)))
		if !ok {
			return fmt.Errorf("failed to load value to []uint8: %v", value)
		}
	default:
		return fmt.Errorf("could not scan type %T into BigFloat", t)
	}

	return nil
}
