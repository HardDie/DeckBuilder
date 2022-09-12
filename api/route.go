package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/api/cards"
	"tts_deck_build/api/collections"
	"tts_deck_build/api/decks"
	"tts_deck_build/api/games"
	"tts_deck_build/api/generator"
	"tts_deck_build/api/images"
	"tts_deck_build/api/system"
	"tts_deck_build/api/web"
)

func GetRoutes() *mux.Router {
	routes := mux.NewRouter().StrictSlash(false)
	ApiRoute := routes.PathPrefix("/api").Subrouter()
	web.Init(routes)
	games.Init(ApiRoute)
	collections.Init(ApiRoute)
	decks.Init(ApiRoute)
	cards.Init(ApiRoute)
	images.Init(ApiRoute)
	system.Init(ApiRoute)
	generator.Init(ApiRoute)
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
