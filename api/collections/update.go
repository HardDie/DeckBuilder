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
		// Required: true
		collections.UpdateCollectionDTO
	}
}

// Status of collection update
//
// swagger:response ResponseUpdateCollection
type ResponseUpdateCollection struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data collections.CollectionInfo `json:"data"`
	}
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
	gameId := mux.Vars(r)["game"]
	collectionId := mux.Vars(r)["collection"]
	dto := &collections.UpdateCollectionDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := collections.NewService().Update(gameId, collectionId, dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
