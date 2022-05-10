package utils

import (
	"regexp"
	"strings"
)

var (
	reg = regexp.MustCompile("[^a-zA-Z0-9_]+")
)

func NameToId(in string) string {
	// Convert all symbols to lowercase
	lower := strings.ToLower(in)
	// Replace all spaces to underscore symbol
	underscore := strings.ReplaceAll(lower, " ", "_")
	// Keep only letters, numbers and underscore symbols
	return reg.ReplaceAllString(underscore, "")
}

func Allocate[T any](val T) *T {
	return &val
}
