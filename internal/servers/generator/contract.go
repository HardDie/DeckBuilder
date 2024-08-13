package generator

import "net/http"

type Generator interface {
	GameHandler(w http.ResponseWriter, r *http.Request)
}
