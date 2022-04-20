package decks

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/network"
)

// Requesting an existing deck
//
// swagger:parameters RequestDeck
type RequestDeck struct {
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

// Deck
//
// swagger:response ResponseDeck
type ResponseDeck struct {
	// In: body
	Body struct {
		decks.DeckInfo
	}
}

// swagger:route GET /games/{game}/collections/{collection}/decks/{deck} Decks RequestDeck
//
// Get deck
//
// Get an existing deck
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
//       200: ResponseDeck
//       default: ResponseError
func ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameName := mux.Vars(r)["game"]
	collectionName := mux.Vars(r)["collection"]
	deckName := mux.Vars(r)["deck"]
	item, e := decks.ItemDeck(gameName, collectionName, deckName)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
	return
}
