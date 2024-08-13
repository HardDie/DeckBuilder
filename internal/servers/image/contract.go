package image

import "net/http"

type Image interface {
	CardHandler(w http.ResponseWriter, r *http.Request)
	CollectionHandler(w http.ResponseWriter, r *http.Request)
	DeckHandler(w http.ResponseWriter, r *http.Request)
	GameHandler(w http.ResponseWriter, r *http.Request)
}
