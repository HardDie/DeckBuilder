package decks

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	DecksRoute := route.PathPrefix("/games/{game}/collections/{collection}/decks").Subrouter()
	DecksRoute.HandleFunc("", ListHandler).Methods(http.MethodGet)
	DecksRoute.HandleFunc("/{deck}", ItemHandler).Methods(http.MethodGet)
}
