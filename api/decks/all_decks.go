package decks

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/network"
)

// Requesting a list of all decks in game
//
// swagger:parameters RequestListOfAllDecks
type RequestListOfAllDecks struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// List of decks
//
// swagger:response ResponseListOfAllDecks
type ResponseListOfAllDecks struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []*decks.DeckInfo `json:"data"`
	}
}

// swagger:route GET /games/{game}/decks Decks RequestListOfAllDecks
//
// Get list of all decks in game
//
// Get a list of all existing decks in game
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
//       200: ResponseListOfAllDecks
//       default: ResponseError
func AllDecksHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	items, e := decks.NewService().ListAllUnique(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
}
