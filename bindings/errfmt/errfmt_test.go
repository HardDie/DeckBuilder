package errfmt

import (
	"errors"
	"fmt"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/apperr"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"known", apperr.ErrGameNotFound, "Game not found"},
		{"with_details", apperr.Withf(apperr.ErrDownloadTimeout, "the image download timed out after %d s", 120), "The image download timed out after 120 s"},
		{"wrapped_known", fmt.Errorf("open game: %w", apperr.ErrGameNotFound), "Game not found"},
		{"already_capitalized", apperr.With(apperr.ErrMissingImages, "Render needs an image for every deck and card."), "Render needs an image for every deck and card."},
		{"unexpected", errors.New("disk full"), apperr.Unexpected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.err); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
