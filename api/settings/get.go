package settings

import (
	"net/http"

	"tts_deck_build/internal/network"
	"tts_deck_build/internal/settings"
)

// swagger:parameters RequestSettings
type RequestSettings struct {
}

// Settings
//
// swagger:response ResponseSettings
type ResponseSettings struct {
	// In: body
	Body struct {
		// Required: true
		Data settings.SettingInfo `json:"data"`
	}
}

// swagger:route GET /settings Settings RequestSettings
//
// Get settings
//
// Get default or changed settings
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
//       200: ResponseSettings
//       default: ResponseError
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	setting, e := settings.NewService().Get()
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, setting)
}
