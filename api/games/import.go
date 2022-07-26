package games

import (
	"io"
	"net/http"

	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

// Creating game from archive
//
// swagger:parameters RequestImportGame
type RequestImportGame struct {
	// Specify a name for the imported game
	// In: formData
	// Required: false
	Name string `json:"name"`
	// Binary data of the imported file
	// In: formData
	// Required: true
	File []byte `json:"file"`
}

// Import game
//
// swagger:response ResponseGameImport
type ResponseGameImport struct {
}

// swagger:route POST /games/import Games RequestImportGame
//
// Import game from archive
//
// Creat game from archive
//
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: ResponseGameImport
//       default: ResponseError
func ImportHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseMultipartForm(0)
	if e != nil {
		errors.IfErrorLog(e)
		e = errors.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	name := r.FormValue("name")

	f, _, e := r.FormFile("file")
	if e != nil {
		errors.IfErrorLog(e)
		e = errors.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	data, e := io.ReadAll(f)
	if e != nil {
		errors.IfErrorLog(e)
		e = errors.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	e = games.NewService().Import(data, name)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
}
