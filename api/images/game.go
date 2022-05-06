package images

import (
	"net/http"

	"github.com/gorilla/mux"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Requesting an image of existing game
//
// swagger:parameters RequestGameImage
type RequestGameImage struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// Game image
//
// swagger:response ResponseGameImage
type ResponseGameImage struct {
	// In: body
	Body []byte
}

// swagger:route GET /games/{game}/image Images RequestGameImage
//
// Get game image
//
// Get an image of existing game
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
//       200: ResponseGameImage
//       default: ResponseError
func GameHandler(w http.ResponseWriter, r *http.Request) {
	gameName := mux.Vars(r)["game"]
	img, imgType, e := games.GetImage(gameName)
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
