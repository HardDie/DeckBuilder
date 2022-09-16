package server

import (
	"net/http"

	"tts_deck_build/internal/dto"
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
	dtoObject := &dto.UpdateSettingsDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	setting, e := system.NewService().UpdateSettings(dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, setting)
}
