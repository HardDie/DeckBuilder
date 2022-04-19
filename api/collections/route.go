package collections

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	CollectionsRoute := route.PathPrefix("/games/{gameName}/collections").Subrouter()
	CollectionsRoute.HandleFunc("", ListHandler).Methods(http.MethodGet)
	CollectionsRoute.HandleFunc("", CreateHandler).Methods(http.MethodPost)
	CollectionsRoute.HandleFunc("/{collectionName}", DeleteHandler).Methods(http.MethodDelete)
	CollectionsRoute.HandleFunc("/{collectionName}", ItemHandler).Methods(http.MethodGet)
	CollectionsRoute.HandleFunc("/{collectionName}", UpdateHandler).Methods(http.MethodPatch)
}
