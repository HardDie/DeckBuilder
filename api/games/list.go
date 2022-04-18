package games

import (
	"log"
	"net/http"

	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Requesting a list of existing games
//
// swagger:parameters RequestListOfGames
type RequestListOfGames struct {
}

// List of games
//
// swagger:response ResponseListOfGames
type ResponseListOfGames struct {
	// In: body
	Body struct {
		games.ListOfGamesResponse
	}
}

// swagger:route GET /games Games RequestListOfGames
//
// Get games list
//
// Get a list of existing games
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
func ListHandler(w http.ResponseWriter, _ *http.Request) {
	items, e := games.ListOfGames()
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	_, err := w.Write(network.ToJson(items))
	if err != nil {
		log.Println(err.Error())
	}
	return
}
