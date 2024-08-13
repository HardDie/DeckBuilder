package replace

import "net/http"

type Replace interface {
	PrepareHandler(w http.ResponseWriter, r *http.Request)
	ReplaceHandler(w http.ResponseWriter, r *http.Request)
}
