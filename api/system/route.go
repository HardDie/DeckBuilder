package system

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	SettingsRoute := route.PathPrefix("/system").Subrouter()
	SettingsRoute.HandleFunc("/quit", QuitHandler).Methods(http.MethodDelete)
}
