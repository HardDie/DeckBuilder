package api

import (
	"github.com/gorilla/mux"
	"tts_deck_build/api/games"
	"tts_deck_build/api/web"
)

func GetRoutes() *mux.Router {
	routes := mux.NewRouter().StrictSlash(false)
	web.Init(routes)
	games.Init(routes)
	return routes
}
