package catalog

import (
	"errors"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/apperr"
)

func TestImageWarning(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "no_error", want: ""},
		{
			name: "known_error",
			err:  apperr.Withf(apperr.ErrDownloadTimeout, "the image download timed out after %d s", 120),
			want: "Image was not saved: the image download timed out after 120 s",
		},
		{
			name: "unexpected_error",
			err:  errors.New("disk full"),
			want: "Image was not saved: " + apperr.Unexpected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ImageWarning(tt.err); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSaveWarning(t *testing.T) {
	failed := apperr.With(apperr.ErrDownloadFailed, "the image could not be downloaded: the server answered 404")
	tests := []struct {
		name     string
		err      error
		hasImage bool
		want     string
	}{
		{name: "image_saved", hasImage: true, want: ""},
		{name: "no_image", want: CardNoImage},
		{name: "failed_old_image_kept", err: failed, hasImage: true, want: "Image was not saved: the image could not be downloaded: the server answered 404"},
		{name: "failed_and_no_image", err: failed, want: "Image was not saved: the image could not be downloaded: the server answered 404 " + CardNoImage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SaveWarning(tt.err, tt.hasImage, CardNoImage); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
