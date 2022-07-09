package settings

import (
	"net/http"

	"tts_deck_build/internal/network"
	"tts_deck_build/internal/settings"
)

// Request to update a settings
//
// swagger:parameters RequestUpdateSettings
type RequestUpdateSettings struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		settings.UpdateSettingsDTO
	}
}

// Settings
//
// swagger:response ResponseUpdateSettings
type ResponseUpdateSettings struct {
	// In: body
	Body struct {
		// Required: true
		Data settings.SettingInfo `json:"data"`
	}
}

// swagger:route PATCH /settings Settings RequestUpdateSettings
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
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	dto := &settings.UpdateSettingsDTO{}
	e := network.RequestToObject(r.Body, &dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	setting, e := settings.NewService().Update(dto)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, setting)
}
