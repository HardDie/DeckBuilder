package collections

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/utils"
)

// Requesting a list of existing collections
//
// swagger:parameters RequestListOfCollections
type RequestListOfCollections struct {
	// In: path
	// Required: true
	GameName string `json:"gameName"`
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

// swagger:route GET /games/{gameName}/collections Collections RequestListOfCollections
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
	gameName := mux.Vars(r)["gameName"]
	items, e := collections.ListOfCollections(gameName)
	if e != nil {
		utils.ResponseError(w, e)
		return
	}

	_, err := w.Write(utils.ToJson(items))
	utils.IfErrorLog(err)
	return
}
