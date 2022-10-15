package server

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/service"
)

type GameServer struct {
	gameService  service.IGameService
	systemServer *SystemServer
}

func NewGameServer(gameService service.IGameService, systemServer *SystemServer) *GameServer {
	return &GameServer{
		gameService:  gameService,
		systemServer: systemServer,
	}
}

func (s *GameServer) CreateHandler(w http.ResponseWriter, r *http.Request) {
	dtoObject := &dto.CreateGameDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := s.gameService.Create(dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
func (s *GameServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	e := s.gameService.Delete(gameID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *GameServer) DuplicateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	dtoObject := &dto.DuplicateGameDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := s.gameService.Duplicate(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
func (s *GameServer) ExportHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	archive, e := s.gameService.Export(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	if _, err := w.Write(archive); err != nil {
		errors.IfErrorLog(err)
	}
}
func (s *GameServer) ImportHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseMultipartForm(0)
	if e != nil {
		errors.IfErrorLog(e)
		e = errors.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	name := r.FormValue("name")

	f, _, e := r.FormFile("file")
	if e != nil {
		errors.IfErrorLog(e)
		e = errors.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	data, e := io.ReadAll(f)
	if e != nil {
		errors.IfErrorLog(e)
		e = errors.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	item, e := s.gameService.Import(data, name)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *GameServer) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	item, e := s.gameService.Item(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *GameServer) ListHandler(w http.ResponseWriter, r *http.Request) {
	s.systemServer.StopQuit()

	sort := r.URL.Query().Get("sort")
	items, e := s.gameService.List(sort)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
}
func (s *GameServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	dtoObject := &dto.UpdateGameDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := s.gameService.Update(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
