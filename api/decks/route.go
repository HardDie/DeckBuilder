package decks

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	DecksRoute := route.PathPrefix("/games/{game}/collections/{collection}/decks").Subrouter()
	DecksRoute.HandleFunc("", ListHandler).Methods(http.MethodGet)
	DecksRoute.HandleFunc("", CreateHandler).Methods(http.MethodPost)
	DecksRoute.HandleFunc("/{deck}", DeleteHandler).Methods(http.MethodDelete)
	DecksRoute.HandleFunc("/{deck}", ItemHandler).Methods(http.MethodGet)
	DecksRoute.HandleFunc("/{deck}", UpdateHandler).Methods(http.MethodPatch)
	GameRoute := route.PathPrefix("/games/{game}").Subrouter()
	GameRoute.HandleFunc("/decks", AllDecksHandler).Methods(http.MethodGet)
}
