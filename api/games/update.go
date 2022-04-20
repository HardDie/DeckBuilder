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
		games.UpdateGameRequest
	}
}

// Status of game update
//
// swagger:response ResponseUpdateGame
type ResponseUpdateGame struct {
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
	gameName := mux.Vars(r)["game"]
	req := &games.UpdateGameRequest{}
	e := network.RequestToObject(r.Body, &req)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	if e = games.UpdateGame(gameName, req); e != nil {
		network.ResponseError(w, e)
	}
	return
}
