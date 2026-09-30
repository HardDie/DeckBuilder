package api

import (
	"net/http"

	"github.com/gorilla/mux"

	serversImage "github.com/HardDie/DeckBuilder/internal/servers/image"
)

func RegisterImageServer(route *mux.Router, srv serversImage.Image) {
	GamesRoute := route.PathPrefix("/api/games").Subrouter()
	GamesRoute.HandleFunc("/{game}/image", srv.GameHandler).Methods(http.MethodGet)

	CollectionsRoute := GamesRoute.PathPrefix("/{game}/collections").Subrouter()
	CollectionsRoute.HandleFunc("/{collection}/image", srv.CollectionHandler).Methods(http.MethodGet)

	DecksRoute := CollectionsRoute.PathPrefix("/{collection}/decks").Subrouter()
	DecksRoute.HandleFunc("/{deck}/image", srv.DeckHandler).Methods(http.MethodGet)

	CardsRoute := DecksRoute.PathPrefix("/{deck}/cards").Subrouter()
	CardsRoute.HandleFunc("/{card}/image", srv.CardHandler).Methods(http.MethodGet)
}
