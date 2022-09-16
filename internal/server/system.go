package server

import (
	"net/http"

	"tts_deck_build/internal/network"
	"tts_deck_build/internal/system"
)

type SystemServer struct {
}

func NewSystemServer() *SystemServer {
	return &SystemServer{}
}

func (s *SystemServer) QuitHandler(w http.ResponseWriter, r *http.Request) {
	system.NewService().Quit()
}
func (s *SystemServer) GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	setting, e := system.NewService().GetSettings()
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, setting)
}
func (s *SystemServer) UpdateSettingsHandler(w http.ResponseWriter, r *http.Request) {
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
