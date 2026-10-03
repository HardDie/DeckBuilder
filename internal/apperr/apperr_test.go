package apperr

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestKindsAndMessages(t *testing.T) {
	detailed := Withf(ErrImageTooLarge, "image is too large: %dx%d, max 128 megapixels", 20000, 20000)
	wrapped := fmt.Errorf("load game munchkin: %w", ErrGameNotFound)
	unexpected := fmt.Errorf("read settings: %w", io.ErrUnexpectedEOF)

	tests := []struct {
		name    string
		err     error
		is      []error // errors.Is must hold for each
		isNot   []error
		message string
		text    string // err.Error(), what the log shows
	}{
		{
			name:    "plain",
			err:     ErrDeckNotFound,
			is:      []error{ErrDeckNotFound, ErrNotFound},
			isNot:   []error{ErrGameNotFound, ErrInvalid},
			message: "deck not found",
			text:    "deck not found",
		},
		{
			name:    "with_details",
			err:     detailed,
			is:      []error{ErrImageTooLarge, ErrInvalid},
			isNot:   []error{ErrUnsupportedImage, ErrNotFound},
			message: "image is too large: 20000x20000, max 128 megapixels",
			text:    "image is too large: 20000x20000, max 128 megapixels",
		},
		{
			name:    "wrapped_with_context",
			err:     wrapped,
			is:      []error{ErrGameNotFound, ErrNotFound},
			message: "game not found",
			text:    "load game munchkin: game not found",
		},
		{
			name:    "unexpected",
			err:     unexpected,
			is:      []error{io.ErrUnexpectedEOF},
			isNot:   []error{ErrNotFound, ErrAlreadyExists, ErrInvalid, ErrBusy},
			message: Unexpected,
			text:    "read settings: unexpected EOF",
		},
		{
			name:    "busy",
			err:     ErrRenderInProgress,
			is:      []error{ErrRenderInProgress, ErrBusy},
			message: "a render is already running",
			text:    "a render is already running",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, target := range tt.is {
				if !errors.Is(tt.err, target) {
					t.Errorf("errors.Is(%v, %v) = false", tt.err, target)
				}
			}
			for _, target := range tt.isNot {
				if errors.Is(tt.err, target) {
					t.Errorf("errors.Is(%v, %v) = true", tt.err, target)
				}
			}
			if got := Message(tt.err); got != tt.message {
				t.Errorf("Message = %q, want %q", got, tt.message)
			}
			if got := tt.err.Error(); got != tt.text {
				t.Errorf("Error = %q, want %q", got, tt.text)
			}
		})
	}
}

// Every error a user can meet belongs to exactly one kind.
func TestEveryErrorHasOneKind(t *testing.T) {
	kinds := []error{ErrNotFound, ErrAlreadyExists, ErrInvalid, ErrBusy}
	all := []*Error{
		ErrGameNotFound, ErrCollectionNotFound, ErrDeckNotFound, ErrCardNotFound,
		ErrGameImageNotFound, ErrCollectionImageNotFound, ErrDeckImageNotFound, ErrCardImageNotFound,
		ErrGameExists, ErrCollectionExists, ErrDeckExists,
		ErrGameImageExists, ErrCollectionImageExists, ErrDeckImageExists, ErrCardImageExists,
		ErrBadName, ErrBadCardID, ErrBadArchive, ErrUnsupportedImage, ErrImageTooLarge,
		ErrDownloadBadURL, ErrDownloadFailed, ErrDownloadTimeout, ErrBadRenderFile, ErrBadMappingFile,
		ErrRenderInProgress, ErrMissingImages, ErrNothingForTTS,
	}
	for _, e := range all {
		n := 0
		for _, k := range kinds {
			if errors.Is(e, k) {
				n++
			}
		}
		if n != 1 {
			t.Errorf("%q belongs to %d kinds, want 1", e.Error(), n)
		}
	}
}
