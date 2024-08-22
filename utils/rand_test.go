package utils

import (
	"testing"
)

func TestRandomStringWithCharset(t *testing.T) {
	length := 10
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := RandomStringWithCharset(length, charset)
	if len(result) != length {
		t.Errorf("RandomStringWithCharset(%d, %s) returned a string of length %d; want %d", length, charset, len(result), length)
	}
}

func TestRandomString(t *testing.T) {
	length := 10
	result := RandomString(length)
	if len(result) != length {
		t.Errorf("RandomString(%d) returned a string of length %d; want %d", length, len(result), length)
	}
}
