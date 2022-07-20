package cards

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/network"
)

// Request to delete a card
//
// swagger:parameters RequestDeleteCard
type RequestDeleteCard struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
	// In: path
	// Required: true
	Card string `json:"card"`
}

// Card deletion status
//
// swagger:response ResponseDeleteCard
type ResponseDeleteCard struct {
}

// swagger:route DELETE /games/{game}/collections/{collection}/decks/{deck}/cards/{card} Cards RequestDeleteCard
//
// Delete card
//
// Allows you to delete an existing card
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
//       200: ResponseDeleteCard
//       default: ResponseError
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	e = cards.NewService().Delete(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
