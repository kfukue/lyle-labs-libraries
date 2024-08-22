package utils

import (
	"reflect"
	"testing"
)

func TestGetFiltersFromQueryStrings(t *testing.T) {
	query := map[string][]string{
		"id":   {"123"},
		"name": {"John"},
		"age":  {"30"},
	}

	filterKeys := []string{"id", "name"}

	expected := []string{
		"id = 123",
		"name LIKE '%John%'",
	}

	result := GetFiltersFromQueryStrings(query, filterKeys)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetFiltersFromQueryStrings(%v, %v) = %v; want %v", query, filterKeys, result, expected)
	}
}
