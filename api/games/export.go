package games

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Requesting an existing game archive
//
// swagger:parameters RequestArchiveGame
type RequestArchiveGame struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// Game archive
//
// swagger:response ResponseGameArchive
type ResponseGameArchive struct {
	// In: body
	Body []byte
}

// swagger:route GET /games/{game}/export Games RequestArchiveGame
//
// Export game to archive
//
// Get an existing game archive
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - application/zip
//
//     Schemes: http
//
//     Responses:
//       200: ResponseGameArchive
//       default: ResponseError
func ExportHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	archive, e := games.NewService().Export(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	if _, err := w.Write(archive); err != nil {
		errors.IfErrorLog(err)
	}
}
