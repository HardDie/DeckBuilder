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
type RequestGame struct {
	// In: path
	// Required: true
	GameName string `json:"gameName"`
	// In: path
	// Required: true
	CollectionName string `json:"collectionName"`
}

// Collection
//
// swagger:response ResponseCollection
type ResponseCollection struct {
	// In: body
	Body struct {
		collections.CollectionInfo
	}
}

// swagger:route GET /games/{gameName}/collections/{collectionName} Collections RequestCollection
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
	gameName := mux.Vars(r)["gameName"]
	collectionName := mux.Vars(r)["collectionName"]
	item, e := collections.ItemCollection(gameName, collectionName)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
	return
}
