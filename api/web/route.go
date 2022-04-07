package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init(route *mux.Router) {
	staticDir := "/web"
	route.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
}
