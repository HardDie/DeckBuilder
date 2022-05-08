package images

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/network"
)

// Requesting an image of existing deck
//
// swagger:parameters RequestDeckImage
type RequestDeckImage struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
}

// Deck image
//
// swagger:response ResponseDeckImage
type ResponseDeckImage struct {
	// In: body
	Body []byte
}

// swagger:route GET /games/{game}/collections/{collection}/decks/{deck}/image Images RequestDeckImage
//
// Get deck image
//
// Get an image of existing deck
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
//       200: ResponseDeckImage
//       default: ResponseError
func DeckHandler(w http.ResponseWriter, r *http.Request) {
	gameId := mux.Vars(r)["game"]
	collectionId := mux.Vars(r)["collection"]
	deckId := mux.Vars(r)["deck"]
	img, imgType, e := decks.NewService().GetImage(gameId, collectionId, deckId)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
