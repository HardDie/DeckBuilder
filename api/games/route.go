package games

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	GamesRoute := route.PathPrefix("/games").Subrouter()
	GamesRoute.HandleFunc("", ListHandler).Methods(http.MethodGet)
	GamesRoute.HandleFunc("", CreateHandler).Methods(http.MethodPost)
	GamesRoute.HandleFunc("/{game}", DeleteHandler).Methods(http.MethodDelete)
	GamesRoute.HandleFunc("/{game}", ItemHandler).Methods(http.MethodGet)
	GamesRoute.HandleFunc("/{game}", UpdateHandler).Methods(http.MethodPatch)
}
