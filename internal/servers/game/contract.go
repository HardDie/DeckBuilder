package game

import "net/http"

type Game interface {
	ImportHandler(w http.ResponseWriter, r *http.Request)
}
