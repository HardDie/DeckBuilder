package settings

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	SettingsRoute := route.PathPrefix("/settings").Subrouter()
	SettingsRoute.HandleFunc("", SettingsHandler).Methods(http.MethodGet)
	SettingsRoute.HandleFunc("", UpdateHandler).Methods(http.MethodPatch)
}
