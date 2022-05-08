package decks

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/network"
)

// Request to delete a deck
//
// swagger:parameters RequestDeleteDeck
type RequestDeleteDeck struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
}

// Deck deletion status
//
// swagger:response ResponseDeleteDeck
type ResponseDeleteDeck struct {
}

// swagger:route DELETE /games/{game}/collections/{collection}/decks/{deck} Decks RequestDeleteDeck
//
// Delete deck
//
// Allows you to delete an existing deck
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
//       200: ResponseDeleteDeck
//       default: ResponseError
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameId := mux.Vars(r)["game"]
	collectionId := mux.Vars(r)["collection"]
	deckId := mux.Vars(r)["deck"]
	e := decks.NewService().Delete(gameId, collectionId, deckId)
	if e != nil {
		network.ResponseError(w, e)
	}
}
