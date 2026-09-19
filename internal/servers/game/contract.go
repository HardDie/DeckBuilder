package game

import "net/http"

type Game interface {
	ExportHandler(w http.ResponseWriter, r *http.Request)
	ImportHandler(w http.ResponseWriter, r *http.Request)
}
