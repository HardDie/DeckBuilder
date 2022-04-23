package utils

import (
	"log"
	"net/url"
	"path"
	"regexp"
	"strings"

	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
func CleanTitle(in string) string {
	res := strings.ReplaceAll(in, " / ", "_")
	res = strings.ReplaceAll(res, "/", "_")
	res = strings.ReplaceAll(res, "!", "")
	res = strings.ReplaceAll(res, "'", "")
	res = strings.ReplaceAll(res, ".", "")
	return strings.ReplaceAll(res, " ", "_")
}
func GetFilenameFromUrl(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, filename := path.Split(u.Path)
	return filename
}

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
