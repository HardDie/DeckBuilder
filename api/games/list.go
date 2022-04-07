package games

import (
	"fmt"
	"log"
	"net/http"

	"tts_deck_build/internal/games"
	"tts_deck_build/internal/utils"
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
// Get games
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
func ListHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, utils.ToJson(games.ListOfGames()))
	if err != nil {
		log.Println(err.Error())
	}
	return
}
