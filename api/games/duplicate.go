package games

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Request to duplicate a game
//
// swagger:parameters RequestDuplicateGame
type RequestDuplicateGame struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: body
	// Required: true
	Body struct {
		// Required: true
		games.DuplicateGameDTO
	}
}

// Status of game duplicate
//
// swagger:response ResponseDuplicateGame
type ResponseDuplicateGame struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data games.GameInfo `json:"data"`
	}
}

// swagger:route POST /games/{game}/duplicate Games RequestDuplicateGame
//
// Duplicate game
//
// Allows you to create a copy of an existing game
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
//       200: ResponseDuplicateGame
//       default: ResponseError
func DuplicateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	dto := &games.DuplicateGameDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := games.NewService().Duplicate(gameID, dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
