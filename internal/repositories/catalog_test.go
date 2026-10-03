package repositories

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveImage(t *testing.T) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 4, 4))); err != nil {
		t.Fatal(err)
	}
	pngBytes := buf.Bytes()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ok.png" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write(pngBytes)
	}))
	defer srv.Close()
	okURL, missingURL := srv.URL+"/ok.png", srv.URL+"/missing.png"

	tests := []struct {
		name    string
		oldURL  string
		newURL  string
		newFile []byte
		want    ImageChange
		wantErr bool
	}{
		{name: "create_without_image"},
		{name: "same_url_keeps_file", oldURL: okURL, newURL: okURL, want: ImageChange{URL: okURL}},
		{name: "clear", oldURL: okURL, want: ImageChange{Clear: true}},
		{name: "new_url", newURL: okURL, want: ImageChange{URL: okURL, Data: pngBytes}},
		{name: "new_file", newFile: pngBytes, want: ImageChange{Data: pngBytes}},
		{name: "new_file_replaces_url", oldURL: okURL, newFile: pngBytes, want: ImageChange{Data: pngBytes}},
		{name: "file_wins_over_url", newURL: missingURL, newFile: pngBytes, want: ImageChange{Data: pngBytes}},
		{name: "bad_url_on_create", newURL: missingURL, want: ImageChange{}, wantErr: true},
		{name: "bad_url_keeps_old", oldURL: okURL, newURL: missingURL, want: ImageChange{URL: okURL}, wantErr: true},
		{name: "bad_file_keeps_old", oldURL: okURL, newURL: okURL, newFile: []byte("text"), want: ImageChange{URL: okURL}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveImage(tt.oldURL, tt.newURL, tt.newFile)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
