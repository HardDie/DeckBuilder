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

	DecksRoute := CollectionsRoute.PathPrefix("/{collection}/decks").Subrouter()
	DecksRoute.HandleFunc("/{deck}/image", DeckHandler).Methods(http.MethodGet)

	CardsRoute := DecksRoute.PathPrefix("/{deck}/cards").Subrouter()
	CardsRoute.HandleFunc("/{card}/image", CardHandler).Methods(http.MethodGet)
}
