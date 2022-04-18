package collections

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	CollectionsRoute := route.PathPrefix("/games/{gameName}/collections").Subrouter()
	CollectionsRoute.HandleFunc("", ListHandler).Methods(http.MethodGet)
	CollectionsRoute.HandleFunc("", CreateHandler).Methods(http.MethodPost)
	// CollectionsRoute.HandleFunc("/{name}", DeleteHandler).Methods(http.MethodDelete)
	// CollectionsRoute.HandleFunc("/{name}", ItemHandler).Methods(http.MethodGet)
	// CollectionsRoute.HandleFunc("/{name}", UpdateHandler).Methods(http.MethodPatch)
}
