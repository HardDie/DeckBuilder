package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type ITTSServer interface {
	DataHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterTTSServer(route *mux.Router, srv ITTSServer) {
	route.HandleFunc("/api/tts/data", srv.DataHandler).Methods(http.MethodGet)
}
