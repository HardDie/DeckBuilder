package utils

import (
	"regexp"
	"strings"
)

const (
	MaxFilenameLength = 200
)

var (
	reg                = regexp.MustCompile(`[^\p{L}0-9_]+`)
	uniqForbiddenNames = map[string]struct{}{
		// windows
		"con":  {},
		"prn":  {},
		"aux":  {},
		"nul":  {},
		"com1": {},
		"com2": {},
		"com3": {},
		"com4": {},
		"com5": {},
		"com6": {},
		"com7": {},
		"com8": {},
		"com9": {},
		"com0": {},
		"lpt1": {},
		"lpt2": {},
		"lpt3": {},
		"lpt4": {},
		"lpt5": {},
		"lpt6": {},
		"lpt7": {},
		"lpt8": {},
		"lpt9": {},
		"lpt0": {},
	}
)

func NameToID(in string) string {
	// Convert all symbols to lowercase
	lower := strings.ToLower(in)
	// Replace all spaces to underscore symbol
	underscore := strings.ReplaceAll(lower, " ", "_")
	// Keep only letters, numbers and underscore symbols
	res := reg.ReplaceAllString(underscore, "")
	if len(res) > MaxFilenameLength {
		res = res[0:MaxFilenameLength]
	}
	if _, ok := uniqForbiddenNames[res]; ok {
		return ""
	}
	return res
}
func Allocate[T any](val T) *T {
	return &val
}
