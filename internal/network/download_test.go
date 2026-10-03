package network

import (
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HardDie/DeckBuilder/internal/apperr"
)

func TestDownload(t *testing.T) {
	const limit = 16

	tests := []struct {
		name    string
		handler http.HandlerFunc
		want    string
		wantErr *apperr.Error
		errText string
	}{
		{
			name: "ok",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte("image"))
			},
			want: "image",
		},
		{
			name: "exactly_the_limit",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(strings.Repeat("x", limit)))
			},
			want: strings.Repeat("x", limit),
		},
		{
			name: "not_found",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr: apperr.ErrDownloadFailed,
			errText: "server answered 404",
		},
		{
			name: "content_length_over_limit",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(strings.Repeat("x", limit+1)))
			},
			wantErr: apperr.ErrImageTooLarge,
			errText: "larger than",
		},
		{
			name: "chunked_over_limit",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				// Flushing before the end drops Content-Length.
				_, _ = w.Write([]byte("x"))
				w.(http.Flusher).Flush()
				_, _ = w.Write([]byte(strings.Repeat("x", limit)))
			},
			wantErr: apperr.ErrImageTooLarge,
			errText: "larger than",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			got, err := download(srv.Client(), srv.URL, limit)
			if tt.wantErr == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.want, string(got))
				return
			}
			assert.True(t, stderrors.Is(err, tt.wantErr), "err %v", err)
			assert.Contains(t, err.Error(), tt.errText)
		})
	}
}

func TestDownloadTimeout(t *testing.T) {
	tests := []struct {
		name        string
		sendHeaders bool
	}{
		{name: "before_headers"},
		{name: "while_reading_body", sendHeaders: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			release := make(chan struct{})
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.sendHeaders {
					w.WriteHeader(http.StatusOK)
					w.(http.Flusher).Flush()
				}
				select {
				case <-release:
				case <-r.Context().Done():
				}
			}))
			defer srv.Close()
			defer close(release)

			client := srv.Client()
			client.Timeout = 100 * time.Millisecond

			start := time.Now()
			_, err := download(client, srv.URL, 1<<20)
			assert.True(t, stderrors.Is(err, apperr.ErrDownloadTimeout), "err %v", err)
			assert.Less(t, time.Since(start), 5*time.Second, "download did not honor the timeout")
		})
	}
}
