package server

import (
	"net/http"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/system"
)

type SystemServer struct {
	cfg *config.Config
}

func NewSystemServer(cfg *config.Config) *SystemServer {
	return &SystemServer{
		cfg: cfg,
	}
}

func (s *SystemServer) QuitHandler(w http.ResponseWriter, r *http.Request) {
	system.NewService(s.cfg).Quit()
}
func (s *SystemServer) GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	setting, e := system.NewService(s.cfg).GetSettings()
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

	setting, e := system.NewService(s.cfg).UpdateSettings(dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, setting)
}
