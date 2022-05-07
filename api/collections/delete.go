package collections

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/network"
)

// Request to delete a collection
//
// swagger:parameters RequestDeleteCollection
type RequestDeleteCollection struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
}

// Collection deletion status
//
// swagger:response ResponseDeleteCollection
type ResponseDeleteCollection struct {
}

// swagger:route DELETE /games/{game}/collections/{collection} Collections RequestDeleteCollection
//
// Delete collection
//
// Allows you to delete an existing collection
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
//       200: ResponseDeleteCollection
//       default: ResponseError
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameId := mux.Vars(r)["game"]
	collectionId := mux.Vars(r)["collection"]
	e := collections.NewService().Delete(gameId, collectionId)
	if e != nil {
		network.ResponseError(w, e)
	}
}
