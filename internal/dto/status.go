package dto

type Status struct {
	Type     string  `json:"type"`
	Message  string  `json:"message"`
	Progress float32 `json:"progress"`
	Status   string  `json:"status"`
}

// DownloadStatus is the current image download for the save overlay.
// Total and Percent are 0 when the size is unknown.
type DownloadStatus struct {
	Active  bool    `json:"active"`
	Done    int64   `json:"done"`
	Total   int64   `json:"total"`
	Percent float32 `json:"percent"`
}
