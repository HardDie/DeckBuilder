package system

import (
	"net/http"

	"tts_deck_build/internal/network"
	"tts_deck_build/internal/system"
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
		Data system.SettingInfo `json:"data"`
	}
}

// swagger:route GET /system/settings System RequestSettings
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
func GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	setting, e := system.NewService().GetSettings()
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, setting)
}
