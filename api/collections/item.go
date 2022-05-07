package collections

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/network"
)

// Requesting an existing collection
//
// swagger:parameters RequestCollection
type RequestCollection struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
}

// Collection
//
// swagger:response ResponseCollection
type ResponseCollection struct {
	// In: body
	Body struct {
		// Required: true
		Data collections.CollectionInfo `json:"data"`
	}
}

// swagger:route GET /games/{game}/collections/{collection} Collections RequestCollection
//
// Get collection
//
// Get an existing collection
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
//       200: ResponseCollection
//       default: ResponseError
func ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameId := mux.Vars(r)["game"]
	collectionId := mux.Vars(r)["collection"]
	item, e := collections.NewService().Item(gameId, collectionId)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
