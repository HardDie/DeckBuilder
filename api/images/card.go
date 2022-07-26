package images

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/network"
)

// Requesting an image of existing card
//
// swagger:parameters RequestCardImage
type RequestCardImage struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
	// In: path
	// Required: true
	Card string `json:"card"`
}

// Card image
//
// swagger:response ResponseCardImage
type ResponseCardImage struct {
	// In: body
	Body []byte
}

// swagger:route GET /games/{game}/collections/{collection}/decks/{deck}/cards/{card}/image Images RequestCardImage
//
// Get card image
//
// Get an image of existing card
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
//       200: ResponseCardImage
//       default: ResponseError
func CardHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	img, imgType, e := cards.NewService().GetImage(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
