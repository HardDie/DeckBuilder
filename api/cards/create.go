package cards

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/network"
)

// Request to create a card
//
// swagger:parameters RequestCreateCard
type RequestCreateCard struct {
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
		cards.CreateCardDTO
	}
}

// Status of card creation
//
// swagger:response ResponseCreateCard
type ResponseCreateCard struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data cards.CardInfo `json:"data"`
	}
}

// swagger:route POST /games/{game}/collections/{collection}/decks/{deck}/cards Cards RequestCreateCard
//
// Create card
//
// Allows you to create a new card
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
//       200: ResponseCreateCard
//       default: ResponseError
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	dto := &cards.CreateCardDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	_, _, _ = gameID, collectionID, deckID
	// network.Response(w, item)
}
