package search

import "net/http"

type Search interface {
	RootHandler(w http.ResponseWriter, r *http.Request)
	GameHandler(w http.ResponseWriter, r *http.Request)
	CollectionHandler(w http.ResponseWriter, r *http.Request)
}
