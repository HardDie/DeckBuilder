package images

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/network"
)

// Requesting an image of existing collection
//
// swagger:parameters RequestCollectionImage
type RequestCollectionImage struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
}

// Collection image
//
// swagger:response ResponseCollectionImage
type ResponseCollectionImage struct {
	// In: body
	Body []byte
}

// swagger:route GET /games/{game}/collections/{collection}/image Images RequestCollectionImage
//
// Get collection image
//
// Get an image of existing collection
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - image/png
//     - image/jpeg
//     - image/gif
//
//     Schemes: http
//
//     Responses:
//       200: ResponseCollectionImage
//       default: ResponseError
func CollectionHandler(w http.ResponseWriter, r *http.Request) {
	gameName := mux.Vars(r)["game"]
	collectionName := mux.Vars(r)["collection"]
	img, imgType, e := collections.GetImage(gameName, collectionName)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
	return
}
