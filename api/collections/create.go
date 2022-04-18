package collections

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/utils"
)

// Request to create a collection
//
// swagger:parameters RequestCreateCollection
type RequestCreateCollection struct {
	// In: body
	// Required: true
	Body struct {
		collections.CreateCollectionRequest
	}
}

// Status of collection creation
//
// swagger:response ResponseCreateCollection
type ResponseCreateCollection struct {
}

// swagger:route POST /games/{gameName}/collections Collections RequestCreateCollection
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
	gameName := mux.Vars(r)["gameName"]
	req := &collections.CreateCollectionRequest{}
	e := utils.RequestToObject(r.Body, &req)
	if e != nil {
		utils.ResponseError(w, e)
		return
	}

	if e = collections.CreateCollection(gameName, req); e != nil {
		utils.ResponseError(w, e)
	}
	return
}
