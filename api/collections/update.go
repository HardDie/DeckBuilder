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
	GameName string `json:"gameName"`
	// In: path
	// Required: true
	CollectionName string `json:"collectionName"`
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
}

// swagger:route PATCH /games/{gameName}/collections/{collectionName} Collections RequestUpdateCollection
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
	gameName := mux.Vars(r)["gameName"]
	collectionName := mux.Vars(r)["collectionName"]
	req := &collections.UpdateCollectionRequest{}
	e := network.RequestToObject(r.Body, &req)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	if e = collections.UpdateCollection(gameName, collectionName, req); e != nil {
		network.ResponseError(w, e)
	}
	return
}
