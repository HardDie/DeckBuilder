package decks

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/network"
)

// Request to update a deck
//
// swagger:parameters RequestUpdateDeck
type RequestUpdateDeck struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
	// In: body
	// Required: true
	Body struct {
		// Required: true
		decks.UpdateDeckDTO
	}
}

// Status of deck update
//
// swagger:response ResponseUpdateDeck
type ResponseUpdateDeck struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data decks.DeckInfo `json:"data"`
	}
}

// swagger:route PATCH /games/{game}/collections/{collection}/decks/{deck} Decks RequestUpdateDeck
//
// Update deck
//
// Allows you to update an existing deck
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
//       200: ResponseUpdateDeck
//       default: ResponseError
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameId := mux.Vars(r)["game"]
	collectionId := mux.Vars(r)["collection"]
	deckId := mux.Vars(r)["deck"]
	dto := &decks.UpdateDeckDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := decks.NewService().Update(gameId, collectionId, deckId, dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
