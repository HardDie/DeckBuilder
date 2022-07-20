package cards

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/network"
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
	Card int64 `json:"card"`
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
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	item, e := cards.NewService().Item(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
