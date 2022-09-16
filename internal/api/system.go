package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/system"
)

type ISystemServer interface {
	QuitHandler(w http.ResponseWriter, r *http.Request)
	GetSettingsHandler(w http.ResponseWriter, r *http.Request)
	UpdateSettingsHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterSystemServer(route *mux.Router, srv ISystemServer) {
	SettingsRoute := route.PathPrefix("/api/system").Subrouter()
	SettingsRoute.HandleFunc("/quit", srv.QuitHandler).Methods(http.MethodDelete)
	SettingsRoute.HandleFunc("/settings", srv.GetSettingsHandler).Methods(http.MethodGet)
	SettingsRoute.HandleFunc("/settings", srv.UpdateSettingsHandler).Methods(http.MethodPatch)
}

type UnimplementedSystemServer struct {
}

var (
	// Validation
	_ ISystemServer = &UnimplementedSystemServer{}
)

// swagger:parameters RequestQuit
type RequestQuit struct {
}

// Quit
//
// swagger:response ResponseQuit
type ResponseQuit struct {
}

// swagger:route DELETE /api/system/quit System RequestQuit
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
func (s *UnimplementedSystemServer) QuitHandler(w http.ResponseWriter, r *http.Request) {}

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

// swagger:route GET /api/system/settings System RequestSettings
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
func (s *UnimplementedSystemServer) GetSettingsHandler(w http.ResponseWriter, r *http.Request) {}

// Request to update a settings
//
// swagger:parameters RequestUpdateSettings
type RequestUpdateSettings struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		dto.UpdateSettingsDTO
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

// swagger:route PATCH /api/system/settings System RequestUpdateSettings
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
func (s *UnimplementedSystemServer) UpdateSettingsHandler(w http.ResponseWriter, r *http.Request) {}
