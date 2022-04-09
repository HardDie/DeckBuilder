package games

import (
	"net/http"

	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/utils"
)

// Request to create a game
//
// swagger:parameters RequestCreateGame
type RequestCreateGame struct {
	// In: body
	// Required: true
	Body struct {
		games.CreateGameRequest
	}
}

// Status of game creation
//
// swagger:response ResponseCreateGame
type ResponseCreateGame struct {
}

// swagger:route POST /games Games RequestCreateGame
//
// Create game
//
// Allows you to create a new game
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
//       200: ResponseCreateGame
//       default: ResponseError
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	req := &games.CreateGameRequest{}
	err := utils.RequestToObject(r.Body, &req)
	if err != nil {
		errors.ResponseError(w, errors.InternalError.AddMessage(err.Error()))
		return
	}

	e := games.CreateGame(req)
	if e != nil {
		errors.ResponseError(w, e)
	}
	return
}
