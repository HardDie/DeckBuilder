package collections

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/network"
)

// Request to create a collection
//
// swagger:parameters RequestCreateCollection
type RequestCreateCollection struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: body
	// Required: true
	Body struct {
		// Required: true
		collections.CreateCollectionDTO
	}
}

// Status of collection creation
//
// swagger:response ResponseCreateCollection
type ResponseCreateCollection struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data collections.CollectionInfo `json:"data"`
	}
}

// swagger:route POST /games/{game}/collections Collections RequestCreateCollection
//
// Create collection
//
// Allows you to create a new collection
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
//       200: ResponseCreateCollection
//       default: ResponseError
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	dto := &collections.CreateCollectionDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := collections.NewService().Create(gameID, dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
