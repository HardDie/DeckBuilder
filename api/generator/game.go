package generator

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/generator"
	"tts_deck_build/internal/network"
)

// Request to start generating result objects
//
// swagger:parameters RequestGameGenerate
type RequestGameImage struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: body
	// Required: false
	Body struct {
		generator.GenerateGameDTO
	}
}

// Generating game objects
//
// swagger:response ResponseGameGenerate
type ResponseGameGenerate struct {
}

// swagger:route POST /games/{game}/generate Generator RequestGameGenerate
//
// Start generating items for TTS
//
// Allow to run the background process of generating images and json item for the game
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
//       200: ResponseGameGenerate
//       default: ResponseError
func GameHandler(w http.ResponseWriter, r *http.Request) {
	dto := &generator.GenerateGameDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	gameID := mux.Vars(r)["game"]
	e = generator.NewService().GenerateGame(gameID, dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
}
