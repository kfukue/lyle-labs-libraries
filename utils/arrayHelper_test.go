package utils

import (
	"testing"
)

func TestIndexOfInts(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	val := 3
	expected := 2
	result := IndexOfInts(arr, val)
	if result != expected {
		t.Errorf("IndexOfInts(%v, %d) = %d; want %d", arr, val, result, expected)
	}

	val = 6
	expected = -1
	result = IndexOfInts(arr, val)
	if result != expected {
		t.Errorf("IndexOfInts(%v, %d) = %d; want %d", arr, val, result, expected)
	}
}

func TestIndexOfStrings(t *testing.T) {
	arr := []string{"a", "b", "c", "d"}
	val := "c"
	expected := 2
	result := IndexOfStrings(arr, val)
	if result != expected {
		t.Errorf("IndexOfStrings(%v, %s) = %d; want %d", arr, val, result, expected)
	}

	val = "e"
	expected = -1
	result = IndexOfStrings(arr, val)
	if result != expected {
		t.Errorf("IndexOfStrings(%v, %s) = %d; want %d", arr, val, result, expected)
	}
}

func TestCheckIfStringPrefixExistsInArray(t *testing.T) {
	arr := []string{"pre", "fix", "test"}
	val := "prefix"
	expected := true
	result := CheckIfStringPrefixExistsInArray(arr, val)
	if result != expected {
		t.Errorf("CheckIfStringPrefixExistsInArray(%v, %s) = %v; want %v", arr, val, result, expected)
	}

	val = "noprefix"
	expected = false
	result = CheckIfStringPrefixExistsInArray(arr, val)
	if result != expected {
		t.Errorf("CheckIfStringPrefixExistsInArray(%v, %s) = %v; want %v", arr, val, result, expected)
	}
}

func TestSplitToString(t *testing.T) {
	arr := []int{1, 2, 3}
	sep := ","
	expected := "1,2,3"
	result := SplitToString(arr, sep)
	if result != expected {
		t.Errorf("SplitToString(%v, %s) = %s; want %s", arr, sep, result, expected)
	}

	arr = []int{}
	expected = ""
	result = SplitToString(arr, sep)
	if result != expected {
		t.Errorf("SplitToString(%v, %s) = %s; want %s", arr, sep, result, expected)
	}
}

func TestSliceAtoi(t *testing.T) {
	arr := []string{"1", "2", "3"}
	expected := []int{1, 2, 3}
	result, err := SliceAtoi(arr)
	if err != nil || !equal(result, expected) {
		t.Errorf("SliceAtoi(%v) = %v, %v; want %v, nil", arr, result, err, expected)
	}

	arr = []string{"1", "a", "3"}
	expected = []int{1}
	result, err = SliceAtoi(arr)
	if err == nil {
		t.Errorf("SliceAtoi(%v) = %v, %v; want error", arr, result, err)
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
