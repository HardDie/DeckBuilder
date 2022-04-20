package games

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Requesting an existing game
//
// swagger:parameters RequestGame
type RequestGame struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// Game
//
// swagger:response ResponseGame
type ResponseGame struct {
	// In: body
	Body struct {
		games.GameInfo
	}
}

// swagger:route GET /games/{game} Games RequestGame
//
// Get game
//
// Get an existing game
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: ResponseGame
//       default: ResponseError
func ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameName := mux.Vars(r)["game"]
	item, e := games.ItemGame(gameName)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
	return
}
