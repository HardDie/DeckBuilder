package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/service"
	"github.com/HardDie/DeckBuilder/internal/utils"
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
	e := r.ParseMultipartForm(0)
	if e != nil {
		er.IfErrorLog(e)
		e = er.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	data, e := utils.GetFileFromMultipart("imageFile", r)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	dtoObject := &dto.CreateGameDTO{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		ImageFile:   data,
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
		er.IfErrorLog(err)
	}
}
func (s *GameServer) ImportHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseMultipartForm(0)
	if e != nil {
		er.IfErrorLog(e)
		e = er.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	name := r.FormValue("name")

	data, e := utils.GetFileFromMultipart("file", r)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	if data == nil {
		e = er.BadArchive.AddMessage("The file must be passed as an argument")
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
	search := r.URL.Query().Get("search")
	items, e := s.gameService.List(sort, search)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
}
func (s *GameServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]

	e := r.ParseMultipartForm(0)
	if e != nil {
		er.IfErrorLog(e)
		e = er.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	data, e := utils.GetFileFromMultipart("imageFile", r)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	dtoObject := &dto.UpdateGameDTO{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		ImageFile:   data,
	}

	item, e := s.gameService.Update(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
