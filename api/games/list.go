package games

import (
	"net/http"

	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Requesting a list of existing games
//
// swagger:parameters RequestListOfGames
type RequestListOfGames struct {
	// In: query
	// Required: false
	Sort string `json:"sort"`
}

// List of games
//
// swagger:response ResponseListOfGames
type ResponseListOfGames struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []*games.GameInfo `json:"data"`
	}
}

// swagger:route GET /games Games RequestListOfGames
//
// Get games list
//
// Get a list of existing games
// Sort values: name, name_desc, created, created_desc
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
//       200: ResponseListOfGames
//       default: ResponseError
func ListHandler(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	items, e := games.NewService().List(sort)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
}
