package cards

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/network"
)

// Request to update a card
//
// swagger:parameters RequestUpdateCard
type RequestUpdateCard struct {
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
	// In: body
	// Required: true
	Body struct {
		// Required: true
		cards.UpdateCardDTO
	}
}

// Status of card update
//
// swagger:response ResponseUpdateCard
type ResponseUpdateCard struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data cards.CardInfo `json:"data"`
	}
}

// swagger:route PATCH /games/{game}/collections/{collection}/decks/{deck}/cards/{card} Cards RequestUpdateCard
//
// Update card
//
// Allows you to update an existing card
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
//       200: ResponseUpdateCard
//       default: ResponseError
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID := mux.Vars(r)["card"]
	dto := &cards.CreateCardDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	_, _, _, _ = gameID, collectionID, deckID, cardID
	// network.Response(w, item)
}
