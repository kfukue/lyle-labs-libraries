package utils

import (
	"strings"
)

func IndexOfInts(arr []int, val int) int {
	for pos, v := range arr {
		if v == val {
			return pos
		}
	}
	return -1
}
func IndexOfStrings(arr []string, val string) int {
	for pos, v := range arr {
		if v == val {
			return pos
		}
	}
	return -1
}

func CheckIfStringPrefixExistsInArray(arr []string, val string) bool {
	exists := false
	for _, v := range arr {
		exists = strings.HasPrefix(val, v)
		if exists == true {
			break
		}
	}
	return exists
}
