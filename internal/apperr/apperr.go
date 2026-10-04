// Package apperr holds the errors the app can explain to the user.
//
// Every error has a clear message, shown in the window as-is, and belongs to one kind.
// Callers check the exact error; the loopback HTTP server checks only the kind
// to choose a status code. Errors carry no transport details.
//
//	errors.Is(err, apperr.ErrDeckNotFound) // exactly which problem
//	errors.Is(err, apperr.ErrNotFound)     // what sort of problem
//
// Anything else is an unexpected error: wrap it with context (fmt.Errorf with %w)
// for the log. Message turns it into "Something went wrong" for the window.
package apperr

import (
	"errors"
	"fmt"
)

// Kinds: what sort of problem an error is.
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalid       = errors.New("invalid")
	ErrBusy          = errors.New("busy")
)

// Error is an error with a message meant for the user.
type Error struct {
	msg    string
	parent error // the kind, or the error this one details
}

func (e *Error) Error() string { return e.msg }
func (e *Error) Unwrap() error { return e.parent }

func newErr(msg string, kind error) *Error {
	return &Error{msg: msg, parent: kind}
}

// With returns err with a more precise message for the user.
// errors.Is still matches err and its kind.
func With(err *Error, msg string) error {
	return &Error{msg: msg, parent: err}
}

// Withf is With with a format string.
func Withf(err *Error, format string, args ...any) error {
	return With(err, fmt.Sprintf(format, args...))
}

// Unexpected is the message shown for errors the app cannot explain.
const Unexpected = "Something went wrong. Details are in the log."

// Message returns the text to show the user: the message of the first
// apperr error in err's chain, or Unexpected when there is none.
func Message(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.msg
	}
	return Unexpected
}

// Catalog: not found.
var (
	ErrGameNotFound       = newErr("game not found", ErrNotFound)
	ErrCollectionNotFound = newErr("collection not found", ErrNotFound)
	ErrDeckNotFound       = newErr("deck not found", ErrNotFound)
	ErrCardNotFound       = newErr("card not found", ErrNotFound)

	ErrGameImageNotFound       = newErr("game image not found", ErrNotFound)
	ErrCollectionImageNotFound = newErr("collection image not found", ErrNotFound)
	ErrDeckImageNotFound       = newErr("deck image not found", ErrNotFound)
	ErrCardImageNotFound       = newErr("card image not found", ErrNotFound)
)

// Catalog: already exists.
var (
	ErrGameExists       = newErr("a game with this name already exists", ErrAlreadyExists)
	ErrCollectionExists = newErr("a collection with this name already exists", ErrAlreadyExists)
	ErrDeckExists       = newErr("a deck with this name already exists", ErrAlreadyExists)

	ErrGameImageExists       = newErr("game image already exists", ErrAlreadyExists)
	ErrCollectionImageExists = newErr("collection image already exists", ErrAlreadyExists)
	ErrDeckImageExists       = newErr("deck image already exists", ErrAlreadyExists)
	ErrCardImageExists       = newErr("card image already exists", ErrAlreadyExists)
)

// Input the app cannot use.
var (
	ErrBadName    = newErr("this name cannot be used", ErrInvalid)
	ErrBadCardID  = newErr("invalid card id", ErrInvalid)
	ErrBadArchive = newErr("not a DeckBuilder game archive", ErrInvalid)

	ErrUnsupportedImage = newErr("not a supported image (png, jpeg, gif)", ErrInvalid)
	ErrImageTooLarge    = newErr("image is too large", ErrInvalid)

	ErrDownloadBadURL  = newErr("the image link is not a valid URL", ErrInvalid)
	ErrDownloadFailed  = newErr("the image could not be downloaded", ErrInvalid)
	ErrDownloadTimeout = newErr("the image download timed out", ErrInvalid)

	ErrBadRenderFile  = newErr("not a DeckBuilder render file", ErrInvalid)
	ErrBadMappingFile = newErr("the mapping file is not valid", ErrInvalid)

	ErrBadCardScale = newErr("card scale must be between 0.1 and 10", ErrInvalid)
)

// Rendering.
var (
	ErrRenderInProgress = newErr("a render is already running", ErrBusy)
	ErrMissingImages    = newErr("render needs an image for every deck and card", ErrInvalid)
)

// Tabletop Simulator.
var (
	ErrNothingForTTS = newErr("nothing to send to Tabletop Simulator", ErrNotFound)
)
