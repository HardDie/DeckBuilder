package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/api/web"
	"tts_deck_build/internal/api"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/server"
)

func GetRoutes() *mux.Router {
	routes := mux.NewRouter().StrictSlash(false)

	web.Init(routes)

	cfg := config.GetConfig()

	api.RegisterGameServer(routes, server.NewGameServer(cfg))
	api.RegisterCollectionServer(routes, server.NewCollectionServer(cfg))
	api.RegisterDeckServer(routes, server.NewDeckServer(cfg))
	api.RegisterCardServer(routes, server.NewCardServer(cfg))

	api.RegisterImageServer(routes, server.NewImageServer(cfg))
	api.RegisterSystemServer(routes, server.NewSystemServer(cfg))
	api.RegisterGeneratorServer(routes, server.NewGeneratorServer(cfg))

	routes.Use(corsMiddleware)
	return routes
}

// CORS headers
func corsSetupHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// CORS Headers middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsSetupHeaders(w)
		next.ServeHTTP(w, r)
	})
}
