package utils

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	expected := "production"
	os.Setenv("APP_ENV", expected)
	result := GetEnv()
	if result != expected {
		t.Errorf("GetEnv() = %s; want %s", result, expected)
	}
}
