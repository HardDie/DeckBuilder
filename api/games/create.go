package games

import (
	"encoding/json"
	"fmt"
	"log"
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
	// In: body
	Body struct {
		games.CreateGameResponse
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
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	req := &games.CreateGameRequest{}
	err := dec.Decode(req)
	if err != nil {
		log.Println(err.Error())
		_, err = fmt.Fprintf(w, utils.ToJson(games.CreateGameResponse{Message: errors.UnknownError}))
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	_, err = fmt.Fprintf(w, utils.ToJson(games.CreateGame(req)))
	if err != nil {
		log.Println(err.Error())
	}
	return
}
