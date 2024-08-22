package utils

import (
	"testing"
)

func TestPtr(t *testing.T) {
	// Test case 1: Test with an integer
	num := 42
	ptr := Ptr(num)
	if *ptr != num {
		t.Errorf("Ptr(%d) = %v; want %v", num, *ptr, num)
	}

	// Test case 2: Test with a string
	str := "Hello, World!"
	strPtr := Ptr[string](str)
	if *strPtr != str {
		t.Errorf("Ptr(%s) = %v; want %v", str, *strPtr, str)
	}

	// Test case 3: Test with a struct
	type Person struct {
		Name string
		Age  int
	}
	person := Person{Name: "John Doe", Age: 30}
	personPtr := Ptr[Person](person)
	if *personPtr != person {
		t.Errorf("Ptr(%v) = %v; want %v", person, *personPtr, person)
	}
}
