package collections

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	CollectionsRoute := route.PathPrefix("/games/{game}/collections").Subrouter()
	CollectionsRoute.HandleFunc("", ListHandler).Methods(http.MethodGet)
	CollectionsRoute.HandleFunc("", CreateHandler).Methods(http.MethodPost)
	CollectionsRoute.HandleFunc("/{collection}", DeleteHandler).Methods(http.MethodDelete)
	CollectionsRoute.HandleFunc("/{collection}", ItemHandler).Methods(http.MethodGet)
	CollectionsRoute.HandleFunc("/{collection}", UpdateHandler).Methods(http.MethodPatch)
}
