package game

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/network"
	serversSystem "github.com/HardDie/DeckBuilder/internal/servers/system"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type game struct {
	cfg          config.Config
	serviceGame  servicesGame.Game
	serverSystem serversSystem.System
}

func New(
	cfg config.Config,
	serviceGame servicesGame.Game,
	serverSystem serversSystem.System,
) Game {
	return &game{
		cfg:          cfg,
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

	type createResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Image       string    `json:"image"`
		CachedImage string    `json:"cachedImage,omitempty"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	network.Response(w, createResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
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

	type duplicateResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Image       string    `json:"image"`
		CachedImage string    `json:"cachedImage,omitempty"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	network.Response(w, duplicateResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
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

	type importResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Image       string    `json:"image"`
		CachedImage string    `json:"cachedImage,omitempty"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	network.Response(w, importResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}
func (s *game) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	item, e := s.serviceGame.Item(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	type itemResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Image       string    `json:"image"`
		CachedImage string    `json:"cachedImage,omitempty"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	network.Response(w, itemResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}
func (s *game) ListHandler(w http.ResponseWriter, r *http.Request) {
	s.serverSystem.StopQuit()

	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")
	items, e := s.serviceGame.List(sort, search)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	type listItemResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Image       string    `json:"image"`
		CachedImage string    `json:"cachedImage,omitempty"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	respItems := make([]*listItemResponse, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, &listItemResponse{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Image:       item.Image,
			CachedImage: s.calculateCachedImage(*item),
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	network.ResponseWithMeta(w, respItems, &network.Meta{
		Total: len(respItems),
	})
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

	type updateResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Image       string    `json:"image"`
		CachedImage string    `json:"cachedImage,omitempty"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	network.Response(w, updateResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}

func (s *game) calculateCachedImage(game entitiesGame.Game) string {
	return fmt.Sprintf(s.cfg.GameImagePath+"?%s", game.ID, utils.HashForTime(&game.UpdatedAt))
}
