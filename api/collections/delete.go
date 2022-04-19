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
	GameName string `json:"gameName"`
	// In: path
	// Required: true
	CollectionName string `json:"collectionName"`
}

// Collection deletion status
//
// swagger:response ResponseDeleteCollection
type ResponseDeleteCollection struct {
}

// swagger:route DELETE /games/{gameName}/collections/{collectionName} Collections RequestDeleteCollection
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
	gameName := mux.Vars(r)["gameName"]
	collectionName := mux.Vars(r)["collectionName"]
	e := collections.DeleteCollection(gameName, collectionName)
	if e != nil {
		network.ResponseError(w, e)
	}
	return
}
