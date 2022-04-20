package decks

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	DecksRoute := route.PathPrefix("/games/{gameName}/collections/{collectionName}/decks").Subrouter()
	DecksRoute.HandleFunc("", ListHandler).Methods(http.MethodGet)
	DecksRoute.HandleFunc("/{deckName}", ItemHandler).Methods(http.MethodGet)
}
