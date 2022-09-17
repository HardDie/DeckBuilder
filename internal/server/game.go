package server

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

type GameServer struct {
	cfg *config.Config
}

func NewGameServer(cfg *config.Config) *GameServer {
	return &GameServer{
		cfg: cfg,
	}
}

func (s *GameServer) CreateHandler(w http.ResponseWriter, r *http.Request) {
	dtoObject := &dto.CreateGameDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := games.NewService(s.cfg).Create(dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
func (s *GameServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	e := games.NewService(s.cfg).Delete(gameID)
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

	item, e := games.NewService(s.cfg).Duplicate(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
func (s *GameServer) ExportHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	archive, e := games.NewService(s.cfg).Export(gameID)
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

	e = games.NewService(s.cfg).Import(data, name)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
}
func (s *GameServer) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	item, e := games.NewService(s.cfg).Item(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *GameServer) ListHandler(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	items, e := games.NewService(s.cfg).List(sort)
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

	item, e := games.NewService(s.cfg).Update(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
