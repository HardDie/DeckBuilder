package game

import "net/http"

type Game interface {
	CreateHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
	DuplicateHandler(w http.ResponseWriter, r *http.Request)
	ExportHandler(w http.ResponseWriter, r *http.Request)
	ImportHandler(w http.ResponseWriter, r *http.Request)
	ItemHandler(w http.ResponseWriter, r *http.Request)
	ListHandler(w http.ResponseWriter, r *http.Request)
	UpdateHandler(w http.ResponseWriter, r *http.Request)
}
