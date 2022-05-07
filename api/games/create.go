package games

import (
	"net/http"

	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Request to create a game
//
// swagger:parameters RequestCreateGame
type RequestCreateGame struct {
	// In: body
	// Required: true
	Body struct {
		games.CreateGameDTO
	}
}

// Status of game creation
//
// swagger:response ResponseCreateGame
type ResponseCreateGame struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data games.GameInfo `json:"data"`
	}
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
	dto := &games.CreateGameDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := games.NewService().Create(dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
