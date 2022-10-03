package server

import (
	"context"
	"net/http"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/system"
)

var quitTimer *time.Timer
var quitCtx context.Context
var quitCancel func()

type SystemServer struct {
	cfg *config.Config
}

func NewSystemServer(cfg *config.Config) *SystemServer {
	return &SystemServer{
		cfg: cfg,
	}
}

func (s *SystemServer) QuitHandler(w http.ResponseWriter, r *http.Request) {
	if quitTimer != nil {
		return
	}
	quitCtx, quitCancel = context.WithCancel(context.Background())
	quitTimer = time.NewTimer(time.Second * 5)
	go func() {
		select {
		case <-quitTimer.C:
			break
		case <-quitCtx.Done():
			quitTimer = nil
			return
		}
		system.NewService(s.cfg).Quit()
	}()
}
func (s *SystemServer) StopQuit() {
	if quitTimer == nil {
		return
	}
	quitCancel()
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
