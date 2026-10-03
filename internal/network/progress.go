package network

import (
	"io"
	"sync"
)

// DownloadState is one reading of the current image download.
// Total is 0 when the server sent no Content-Length.
type DownloadState struct {
	Active bool
	Done   int64
	Total  int64
}

// The GUI saves one entity at a time, so at most one download runs.
// The window polls this state while a save is pending.
var (
	downloadMu    sync.Mutex
	downloadState DownloadState
)

// DownloadProgress returns the current reading.
func DownloadProgress() DownloadState {
	downloadMu.Lock()
	defer downloadMu.Unlock()
	return downloadState
}

func beginDownload(total int64) {
	downloadMu.Lock()
	defer downloadMu.Unlock()
	downloadState = DownloadState{Active: true, Total: max(total, 0)}
}

func addDownloaded(n int) {
	downloadMu.Lock()
	defer downloadMu.Unlock()
	downloadState.Done += int64(n)
}

func endDownload() {
	downloadMu.Lock()
	defer downloadMu.Unlock()
	downloadState = DownloadState{}
}

// progressReader counts the bytes read into the download state.
type progressReader struct {
	r io.Reader
}

func (p progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	addDownloaded(n)
	return n, err
}
