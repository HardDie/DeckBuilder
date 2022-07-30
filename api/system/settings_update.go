package system

import (
	"net/http"

	"tts_deck_build/internal/network"
	"tts_deck_build/internal/system"
)

// Request to update a settings
//
// swagger:parameters RequestUpdateSettings
type RequestUpdateSettings struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		system.UpdateSettingsDTO
	}
}

// Settings
//
// swagger:response ResponseUpdateSettings
type ResponseUpdateSettings struct {
	// In: body
	Body struct {
		// Required: true
		Data system.SettingInfo `json:"data"`
	}
}

// swagger:route PATCH /system/settings System RequestUpdateSettings
//
// Update settings
//
// API to update settings
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
//       200: ResponseUpdateSettings
//       default: ResponseError
func UpdateSettingsHandler(w http.ResponseWriter, r *http.Request) {
	dto := &system.UpdateSettingsDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	setting, e := system.NewService().UpdateSettings(dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, setting)
}
