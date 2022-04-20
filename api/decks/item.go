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
	GameName string `json:"gameName"`
	// In: path
	// Required: true
	CollectionName string `json:"collectionName"`
	// In: path
	// Required: true
	DeckName string `json:"deckName"`
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

// swagger:route GET /games/{gameName}/collections/{collectionName}/decks/{deckName} Decks RequestDeck
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
	gameName := mux.Vars(r)["gameName"]
	collectionName := mux.Vars(r)["collectionName"]
	deckName := mux.Vars(r)["deckName"]
	item, e := decks.ItemDeck(gameName, collectionName, deckName)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
	return
}
