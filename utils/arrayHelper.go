package utils

import (
	"strconv"
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

func SplitToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}

func SliceAtoi(sa []string) ([]int, error) {
	si := make([]int, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.Atoi(a)
		if err != nil {
			return si, err
		}
		si = append(si, i)
	}
	return si, nil
}
