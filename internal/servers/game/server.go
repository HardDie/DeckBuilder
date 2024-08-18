package game

import (
	"net/http"

	"github.com/gorilla/mux"

	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/network"
	serversSystem "github.com/HardDie/DeckBuilder/internal/servers/system"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type game struct {
	serviceGame  servicesGame.Game
	serverSystem serversSystem.System
}

func New(serviceGame servicesGame.Game, serverSystem serversSystem.System) Game {
	return &game{
		serviceGame:  serviceGame,
		serverSystem: serverSystem,
	}
}

func (s *game) CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	item, e := s.serviceGame.Create(servicesGame.CreateRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		ImageFile:   data,
	})
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
func (s *game) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	e := s.serviceGame.Delete(gameID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *game) DuplicateHandler(w http.ResponseWriter, r *http.Request) {
	type duplicateRequest struct {
		Name string `json:"name"`
	}
	gameID := mux.Vars(r)["game"]
	dtoObject := &duplicateRequest{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := s.serviceGame.Duplicate(gameID, servicesGame.DuplicateRequest{
		Name: dtoObject.Name,
	})
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
func (s *game) ExportHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	archive, e := s.serviceGame.Export(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	if _, err := w.Write(archive); err != nil {
		er.IfErrorLog(err)
	}
}
func (s *game) ImportHandler(w http.ResponseWriter, r *http.Request) {
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

	item, e := s.serviceGame.Import(data, name)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *game) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	item, e := s.serviceGame.Item(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *game) ListHandler(w http.ResponseWriter, r *http.Request) {
	s.serverSystem.StopQuit()

	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")
	items, meta, e := s.serviceGame.List(sort, search)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.ResponseWithMeta(w, items, meta)
}
func (s *game) UpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	item, e := s.serviceGame.Update(gameID, servicesGame.UpdateRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		ImageFile:   data,
	})
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
