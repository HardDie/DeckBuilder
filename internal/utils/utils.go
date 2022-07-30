package utils

import (
	"regexp"
	"strings"

	"golang.org/x/exp/constraints"

	"tts_deck_build/internal/config"
)

var (
	reg = regexp.MustCompile("[^a-zA-Z0-9_]+")
)

func NameToID(in string) string {
	// Convert all symbols to lowercase
	lower := strings.ToLower(in)
	// Replace all spaces to underscore symbol
	underscore := strings.ReplaceAll(lower, " ", "_")
	// Keep only letters, numbers and underscore symbols
	res := reg.ReplaceAllString(underscore, "")
	if len(res) > config.MaxFilenameLength {
		res = res[0:config.MaxFilenameLength]
	}
	return res
}

func Allocate[T any](val T) *T {
	return &val
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// TODO: Remove
func GetFilenameFromURL(link string) string {
	return ""
}
