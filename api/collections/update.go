package collections

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/network"
)

// Request to update a collection
//
// swagger:parameters RequestUpdateCollection
type RequestUpdateCollection struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: body
	// Required: true
	Body struct {
		collections.UpdateCollectionRequest
	}
}

// Status of collection update
//
// swagger:response ResponseUpdateCollection
type ResponseUpdateCollection struct {
	collections.CollectionInfo
}

// swagger:route PATCH /games/{game}/collections/{collection} Collections RequestUpdateCollection
//
// Update collection
//
// Allows you to update an existing collection
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
//       200: ResponseUpdateCollection
//       default: ResponseError
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameName := mux.Vars(r)["game"]
	collectionName := mux.Vars(r)["collection"]
	req := &collections.UpdateCollectionRequest{}
	e := network.RequestToObject(r.Body, &req)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := collections.UpdateCollection(gameName, collectionName, req)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
