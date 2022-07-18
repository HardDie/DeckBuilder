package cards

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	CardsRoute := route.PathPrefix("/games/{game}/collections/{collection}/decks/{deck}/cards").Subrouter()
	CardsRoute.HandleFunc("", ListHandler).Methods(http.MethodGet)
	CardsRoute.HandleFunc("", CreateHandler).Methods(http.MethodPost)
	CardsRoute.HandleFunc("/{card}", DeleteHandler).Methods(http.MethodDelete)
	CardsRoute.HandleFunc("/{card}", ItemHandler).Methods(http.MethodGet)
	CardsRoute.HandleFunc("/{card}", UpdateHandler).Methods(http.MethodPatch)
}
