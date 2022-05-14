package utils

import (
	"log"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/exp/constraints"
	"tts_deck_build/internal/fs"
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
	return reg.ReplaceAllString(underscore, "")
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

// TODO: Remove panic
func GetFilenameFromURL(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, filename := path.Split(u.Path)
	extension := filepath.Ext(filename)
	nameOnly := fs.GetFilenameWithoutExt(filename)
	return NameToID(nameOnly) + extension
}
