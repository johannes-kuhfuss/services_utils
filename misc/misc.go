package misc

import (
	"slices"
	"strings"
)

func SliceContainsString(slice []string, searchStr string) bool {
	return slices.Contains(slice, searchStr)
}

func SliceContainsStringCI(slice []string, searchStr string) bool {
	var ciSlice []string
	for _, val := range slice {
		ciSlice = append(ciSlice, strings.ToLower(val))
	}
	return slices.Contains(ciSlice, searchStr)
}
