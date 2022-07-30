package system

import (
	"net/http"

	"tts_deck_build/internal/system"
)

// swagger:parameters RequestQuit
type RequestQuit struct {
}

// Quit
//
// swagger:response ResponseQuit
type ResponseQuit struct {
}

// swagger:route DELETE /system/quit System RequestQuit
//
// Close application
//
// Close app on back side
//
//    Consumes:
//    - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: ResponseQuit
//       default: ResponseError
func QuitHandler(w http.ResponseWriter, r *http.Request) {
	system.NewService().Quit()
}
