package biginthelper

import (
	"math/big"
	"reflect"
	"testing"
)

func TestConvertToBigInt(t *testing.T) {
	b := big.NewInt(123)
	expected := &BigInt{}
	expected.Scan("123")

	result, err := ConvertToBigInt(b)
	if err != nil {
		t.Errorf("ConvertToBigInt(%v) returned an error: %v", b, err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ConvertToBigInt(%v) = %v; want %v", b, result, expected)
	}

	b = nil
	expected = nil

	result, err = ConvertToBigInt(b)
	if err != nil {
		t.Errorf("ConvertToBigInt(%v) returned an error: %v", b, err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ConvertToBigInt(%v) = %v; want %v", b, result, expected)
	}
}

func TestBigInt_Value(t *testing.T) {
	b := &BigInt{}
	b.Scan("123")
	expected := "123"

	result, err := b.Value()
	if err != nil {
		t.Errorf("BigInt.Value() returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("BigInt.Value() = %v; want %v", result, expected)
	}

	b = nil

	result, err = b.Value()
	if err != nil {
		t.Errorf("BigInt.Value() returned an error: %v", err)
	}

	if result != nil {
		t.Errorf("BigInt.Value() = %v; want %v", result, nil)
	}
}

// Add more tests for the remaining functions...
