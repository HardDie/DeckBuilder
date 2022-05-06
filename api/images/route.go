package images

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	GamesRoute := route.PathPrefix("/games").Subrouter()
	GamesRoute.HandleFunc("/{game}/image", GameHandler).Methods(http.MethodGet)

	CollectionsRoute := GamesRoute.PathPrefix("/{game}/collections").Subrouter()
	CollectionsRoute.HandleFunc("/{collection}/image", CollectionHandler).Methods(http.MethodGet)
}
