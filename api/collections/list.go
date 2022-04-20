package collections

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/network"
)

// Requesting a list of existing collections
//
// swagger:parameters RequestListOfCollections
type RequestListOfCollections struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// List of collections
//
// swagger:response ResponseListOfCollections
type ResponseListOfCollections struct {
	// In: body
	Body struct {
		collections.ListOfCollectionsResponse
	}
}

// swagger:route GET /games/{game}/collections Collections RequestListOfCollections
//
// Get collections list
//
// Get a list of existing collections
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
//       200: ResponseListOfCollections
//       default: ResponseError
func ListHandler(w http.ResponseWriter, r *http.Request) {
	gameName := mux.Vars(r)["game"]
	items, e := collections.ListOfCollections(gameName)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
	return
}
