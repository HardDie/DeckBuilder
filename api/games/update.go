package games

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Request to update a game
//
// swagger:parameters RequestUpdateGame
type RequestUpdateGame struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: body
	// Required: true
	Body struct {
		// Required: true
		games.UpdateGameDTO
	}
}

// Status of game update
//
// swagger:response ResponseUpdateGame
type ResponseUpdateGame struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data games.GameInfo `json:"data"`
	}
}

// swagger:route PATCH /games/{game} Games RequestUpdateGame
//
// Update game
//
// Allows you to update an existing game
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
//       200: ResponseUpdateGame
//       default: ResponseError
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameId := mux.Vars(r)["game"]
	dto := &games.UpdateGameDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := games.NewService().Update(gameId, dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
