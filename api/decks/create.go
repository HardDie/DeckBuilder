package decks

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/network"
)

// Request to create a deck
//
// swagger:parameters RequestCreateDeck
type RequestCreateDeck struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: body
	// Required: true
	Body struct {
		decks.CreateDeckRequest
	}
}

// Status of deck creation
//
// swagger:response ResponseCreateDeck
type ResponseCreateDeck struct {
	decks.DeckInfo
}

// swagger:route POST /games/{game}/collections/{collection}/decks Decks RequestCreateDeck
//
// Create deck
//
// Allows you to create a new deck
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
//       200: ResponseCreateDeck
//       default: ResponseError
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	gameName := mux.Vars(r)["game"]
	collectionName := mux.Vars(r)["collection"]
	req := &decks.CreateDeckRequest{}
	e := network.RequestToObject(r.Body, &req)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := decks.CreateDeck(gameName, collectionName, req)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
