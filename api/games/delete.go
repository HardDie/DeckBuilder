package games

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
)

// Request to delete a game
//
// swagger:parameters RequestDeleteGame
type RequestDeleteGame struct {
	// In: path
	// Required: true
	Name string `json:"name"`
}

// Game deletion status
//
// swagger:response ResponseDeleteGame
type ResponseDeleteGame struct {
}

// swagger:route DELETE /games/{name} Games RequestDeleteGame
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
	name := mux.Vars(r)["name"]
	e := games.DeleteGame(name)
	if e != nil {
		errors.ResponseError(w, e)
	}
	return
}
