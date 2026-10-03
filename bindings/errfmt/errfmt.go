// Package errfmt turns errors from bound methods into what the window shows.
package errfmt

import (
	"log/slog"
	"unicode"
	"unicode/utf8"

	"github.com/HardDie/DeckBuilder/internal/apperr"
)

// Format is the Wails ErrorFormatter. Every error a bound method returns passes here.
// It returns Text, and logs the details of an unexpected error.
func Format(err error) any {
	if apperr.Message(err) == apperr.Unexpected {
		slog.Error("unexpected error", "err", err)
	}
	return Text(err)
}

// Text is what the window shows for err: a known error's message as a sentence,
// or the generic message. It does not log.
func Text(err error) string {
	return capitalize(apperr.Message(err))
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
