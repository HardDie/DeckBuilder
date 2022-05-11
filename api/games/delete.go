package games

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Request to delete a game
//
// swagger:parameters RequestDeleteGame
type RequestDeleteGame struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// Game deletion status
//
// swagger:response ResponseDeleteGame
type ResponseDeleteGame struct {
}

// swagger:route DELETE /games/{game} Games RequestDeleteGame
//
// Delete game
//
// Allows you to delete an existing game
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
//       200: ResponseDeleteGame
//       default: ResponseError
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	e := games.NewService().Delete(gameID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
