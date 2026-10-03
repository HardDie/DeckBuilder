// Package errfmt turns errors from bound methods into what the window shows.
package errfmt

import (
	"unicode"
	"unicode/utf8"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

// Format is the Wails ErrorFormatter. Every error a bound method returns passes here.
// A known error shows its message; anything else shows a generic one,
// and its details go to the log.
func Format(err error) any {
	msg := apperr.Message(err)
	if msg == apperr.Unexpected {
		logger.Error.Println("unexpected error:", err.Error())
	}
	return capitalize(msg)
}

// capitalize upper-cases the first letter: errors are lowercase by Go convention,
// the window shows sentences.
func capitalize(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}
