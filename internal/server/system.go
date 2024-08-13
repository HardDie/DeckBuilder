package server

import (
	"context"
	"net/http"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/progress"
	servicesSystem "github.com/HardDie/DeckBuilder/internal/services/system"
)

const (
	DestroyTimer = time.Second * 60
)

var quitTimer *time.Timer
var quitCtx context.Context
var quitCancel func()

type SystemServer struct {
	serviceSystem servicesSystem.System
	cfg           *config.Config
}

func NewSystemServer(cfg *config.Config, serviceSystem servicesSystem.System) *SystemServer {
	return &SystemServer{
		cfg:           cfg,
		serviceSystem: serviceSystem,
	}
}

func (s *SystemServer) QuitHandler(w http.ResponseWriter, _ *http.Request) {
	if s.cfg.Debug {
		return
	}

	if quitTimer != nil {
		return
	}
	logger.Debug.Println("Start destroy timer:", DestroyTimer.String())
	quitCtx, quitCancel = context.WithCancel(context.Background())
	quitTimer = time.NewTimer(DestroyTimer)
	go func() {
		select {
		case <-quitTimer.C:
			break
		case <-quitCtx.Done():
			logger.Debug.Println("Cancel destroy timer")
			quitTimer = nil
			return
		}
		s.serviceSystem.Quit()
	}()
}
func (s *SystemServer) StopQuit() {
	if quitTimer == nil {
		return
	}
	quitCancel()
}
func (s *SystemServer) GetSettingsHandler(w http.ResponseWriter, _ *http.Request) {
	setting, e := s.serviceSystem.GetSettings()
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

	setting, e := s.serviceSystem.UpdateSettings(dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, setting)
}
func (s *SystemServer) StatusHandler(w http.ResponseWriter, _ *http.Request) {
	status := progress.GetProgress().GetStatus()
	if status.Status == progress.StatusError || status.Status == progress.StatusDone {
		progress.GetProgress().Flush()
	}
	network.Response(w, status)
}
func (s *SystemServer) GetVersionHandler(w http.ResponseWriter, _ *http.Request) {
	network.Response(w, s.cfg.Version)
}
