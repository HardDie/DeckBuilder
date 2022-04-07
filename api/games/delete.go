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

// Request to delete a game
//
// swagger:parameters RequestDeleteGame
type RequestDeleteGame struct {
	// In: body
	// Required: true
	Body struct {
		games.DeleteGameRequest
	}
}

// Game deletion status
//
// swagger:response ResponseDeleteGame
type ResponseDeleteGame struct {
	// In: body
	Body struct {
		games.DeleteGameResponse
	}
}

// swagger:route DELETE /games Games RequestDeleteGame
//
// Delete game
//
// Allows you to delete an existing game
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
//       200: ResponseDeleteGame
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	req := &games.DeleteGameRequest{}
	err := dec.Decode(req)
	if err != nil {
		log.Println(err.Error())
		_, err = fmt.Fprintf(w, utils.ToJson(games.DeleteGameResponse{Message: errors.UnknownError}))
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	_, err = fmt.Fprintf(w, utils.ToJson(games.DeleteGame(req)))
	if err != nil {
		log.Println(err.Error())
	}
	return
}
