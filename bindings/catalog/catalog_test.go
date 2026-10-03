package catalog

import (
	"errors"
	"testing"

	er "github.com/HardDie/DeckBuilder/internal/errors"
)

func TestImageWarning(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "no_error", want: ""},
		{
			name: "catalog_error_without_http_prefix",
			err:  er.NetworkTimeout.AddMessage("download timed out after 120 s"),
			want: "Image was not saved: download timed out after 120 s",
		},
		{
			name: "plain_error",
			err:  errors.New("disk full"),
			want: "Image was not saved: disk full",
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
	failed := er.NetworkBadResponse.AddMessage("server answered 404")
	tests := []struct {
		name     string
		err      error
		hasImage bool
		want     string
	}{
		{name: "image_saved", hasImage: true, want: ""},
		{name: "no_image", want: CardNoImage},
		{name: "failed_old_image_kept", err: failed, hasImage: true, want: "Image was not saved: server answered 404"},
		{name: "failed_and_no_image", err: failed, want: "Image was not saved: server answered 404 " + CardNoImage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SaveWarning(tt.err, tt.hasImage, CardNoImage); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
