package repositories

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"

	er "github.com/HardDie/DeckBuilder/internal/errors"
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

func TestMapFsentry(t *testing.T) {
	missing := func(path ...string) error { return &fsentry.BadPathError{Path: path} }
	deck := FsentrySentinels{Exist: er.DeckExist, NotExist: er.DeckNotExists}

	tests := []struct {
		name      string
		err       error
		sentinels FsentrySentinels
		want      error
	}{
		{name: "nil", err: nil, want: nil},
		{name: "missing_game", err: missing("games", "g"), sentinels: deck, want: er.GameNotExists},
		{name: "missing_collection", err: missing("games", "g", "c"), sentinels: deck, want: er.CollectionNotExists},
		{name: "missing_deck", err: missing("games", "g", "c", "d"), sentinels: deck, want: er.DeckNotExists},
		{name: "exist", err: fmt.Errorf("x: %w", fsentry.ErrExist), sentinels: deck, want: er.DeckExist},
		{name: "not_exist", err: fmt.Errorf("x: %w", fsentry.ErrNotExist), sentinels: deck, want: er.DeckNotExists},
		{name: "not_exist_without_sentinel", err: fsentry.ErrNotExist, want: er.InternalError},
		{name: "bad_name", err: fsentry.ErrBadName, sentinels: deck, want: er.BadName},
		{name: "bare_bad_path", err: fsentry.ErrBadPath, sentinels: deck, want: er.InternalError},
		{name: "other", err: errors.New("disk"), sentinels: deck, want: er.InternalError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapFsentry(tt.err, tt.sentinels)
			if tt.want == nil {
				assert.NoError(t, got)
				return
			}
			assert.ErrorIs(t, got, tt.want)
		})
	}
}
