package network

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadProgress(t *testing.T) {
	const half = 10
	tests := []struct {
		name      string
		sendSize  bool
		wantTotal int64
	}{
		{name: "known_size", sendSize: true, wantTotal: 2 * half},
		{name: "unknown_size", sendSize: false, wantTotal: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			release := make(chan struct{})
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if tt.sendSize {
					w.Header().Set("Content-Length", strconv.Itoa(2*half))
				}
				_, _ = w.Write([]byte(strings.Repeat("x", half)))
				w.(http.Flusher).Flush()
				<-release
				_, _ = w.Write([]byte(strings.Repeat("x", half)))
			}))
			defer srv.Close()

			assert.Equal(t, DownloadState{}, DownloadProgress(), "idle before")

			result := make(chan error, 1)
			go func() {
				_, err := download(srv.Client(), srv.URL, 1<<20)
				result <- err
			}()

			// Mid-download: the first half is counted.
			got := waitProgress(t, func(s DownloadState) bool { return s.Done == half })
			assert.True(t, got.Active)
			assert.Equal(t, tt.wantTotal, got.Total)

			close(release)
			require.NoError(t, <-result)
			assert.Equal(t, DownloadState{}, DownloadProgress(), "idle after")
		})
	}
}

func TestDownloadProgressClearedOnError(t *testing.T) {
	// No Content-Length, so the limit trips while reading, after the state was set.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("x"))
		w.(http.Flusher).Flush()
		_, _ = w.Write([]byte(strings.Repeat("x", 32)))
	}))
	defer srv.Close()

	_, err := download(srv.Client(), srv.URL, 16)
	assert.Error(t, err)
	assert.Equal(t, DownloadState{}, DownloadProgress())
}

func waitProgress(t *testing.T, ok func(DownloadState) bool) DownloadState {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if s := DownloadProgress(); ok(s) {
			return s
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("progress never matched, last %+v", DownloadProgress())
	return DownloadState{}
}
