package cards

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/cards"
)

// Requesting an existing card
//
// swagger:parameters RequestCard
type RequestCard struct {
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

// Card
//
// swagger:response ResponseCard
type ResponseCard struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data cards.CardInfo `json:"data"`
	}
}

// swagger:route GET /games/{game}/collections/{collection}/decks/{deck}/cards/{card} Cards RequestCard
//
// Get card
//
// Get an existing card
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
//       200: ResponseCard
//       default: ResponseError
func ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID := mux.Vars(r)["card"]

	_, _, _, _ = gameID, collectionID, deckID, cardID
	// network.Response(w, item)
}
