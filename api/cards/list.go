package cards

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/network"
)

// Requesting a list of existing cards
//
// swagger:parameters RequestListOfCard
type RequestListOfCard struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
	// In: query
	// Required: false
	Sort string `json:"sort"`
}

// List of cards
//
// swagger:response ResponseListOfCard
type ResponseListOfCard struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []*cards.CardInfo `json:"data"`
	}
}

// swagger:route GET /games/{game}/collections/{collection}/decks/{deck}/cards Cards RequestListOfCard
//
// Get cards list
//
// Get a list of existing cards
// Sort values: name, name_desc, created, created_desc
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
//       200: ResponseListOfCard
//       default: ResponseError
func ListHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	sort := r.URL.Query().Get("sort")
	items, e := cards.NewService().List(gameID, collectionID, deckID, sort)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
}
