package utils

import (
	"regexp"
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
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
