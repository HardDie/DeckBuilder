package generator

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	GamesRoute := route.PathPrefix("/games").Subrouter()
	GamesRoute.HandleFunc("/{game}/generate", GameHandler).Methods(http.MethodPost)
}
