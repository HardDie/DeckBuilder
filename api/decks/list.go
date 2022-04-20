package decks

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/network"
)

// Requesting a list of existing decks
//
// swagger:parameters RequestListOfDecks
type RequestListOfDecks struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
}

// List of decks
//
// swagger:response ResponseListOfDecks
type ResponseListOfDecks struct {
	// In: body
	Body struct {
		decks.ListOfDecksResponse
	}
}

// swagger:route GET /games/{game}/collections/{collection}/decks Decks RequestListOfDecks
//
// Get decks list
//
// Get a list of existing decks
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
//       200: ResponseListOfDecks
//       default: ResponseError
func ListHandler(w http.ResponseWriter, r *http.Request) {
	gameName := mux.Vars(r)["game"]
	collectionName := mux.Vars(r)["collection"]
	items, e := decks.ListOfDecks(gameName, collectionName)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
	return
}
