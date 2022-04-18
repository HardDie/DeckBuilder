package games

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Requesting an existing game
//
// swagger:parameters RequestGame
type RequestGame struct {
	// In: path
	// Required: true
	Name string `json:"name"`
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

// swagger:route GET /games/{name} Games RequestGame
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
	name := mux.Vars(r)["name"]
	item, e := games.ItemGame(name)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	_, err := w.Write(network.ToJson(item))
	errors.IfErrorLog(err)
	return
}
