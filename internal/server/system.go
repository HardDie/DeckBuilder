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
	"github.com/HardDie/DeckBuilder/internal/system"
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
	logger.Debug.Println("Start destroy timer")
	quitCtx, quitCancel = context.WithCancel(context.Background())
	quitTimer = time.NewTimer(time.Second * 5)
	go func() {
		select {
		case <-quitTimer.C:
			break
		case <-quitCtx.Done():
			logger.Debug.Println("Cancel destroy timer")
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
func (s *SystemServer) StatusHandler(w http.ResponseWriter, r *http.Request) {
	status := progress.GetProgress().GetStatus()
	network.Response(w, status)
}
