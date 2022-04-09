package games

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	GamesRoute := route.PathPrefix("/games").Subrouter()
	GamesRoute.HandleFunc("", ListHandler).Methods(http.MethodGet)
	GamesRoute.HandleFunc("", CreateHandler).Methods(http.MethodPost)
	GamesRoute.HandleFunc("/{name}", DeleteHandler).Methods(http.MethodDelete)
	GamesRoute.HandleFunc("/{name}", ItemHandler).Methods(http.MethodGet)
	GamesRoute.HandleFunc("/{name}", UpdateHandler).Methods(http.MethodPatch)
}
