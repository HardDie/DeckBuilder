package catalog

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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

// The window sends uploaded files as base64 strings (S13). encoding/json, which Wails
// uses for binding arguments, decodes them into []byte; the old number array still works.
func TestWriteRequestImageFileFormats(t *testing.T) {
	want := []byte{0x89, 'P', 'N', 'G', 0, 255}
	tests := []struct {
		name string
		body string
		want []byte
	}{
		{"base64", `{"imageFile":"` + base64.StdEncoding.EncodeToString(want) + `"}`, want},
		{"number_array", `{"imageFile":[137,80,78,71,0,255]}`, want},
		{"empty_string", `{"imageFile":""}`, nil},
		{"missing", `{}`, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req WriteRequest
			if err := json.Unmarshal([]byte(tt.body), &req); err != nil {
				t.Fatal(err)
			}
			if got := req.ImageBytes(); !bytes.Equal(got, tt.want) || (tt.want == nil && got != nil) {
				t.Fatalf("ImageBytes %v, want %v", got, tt.want)
			}
		})
	}

	// A plain []byte argument, as in game.Import and replace.Prepare.
	var arg []byte
	if err := json.Unmarshal([]byte(`"`+base64.StdEncoding.EncodeToString(want)+`"`), &arg); err != nil || !bytes.Equal(arg, want) {
		t.Fatalf("[]byte argument %v, err %v", arg, err)
	}
}
