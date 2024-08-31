package collection

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/network"
	serversSystem "github.com/HardDie/DeckBuilder/internal/servers/system"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type collection struct {
	cfg               config.Config
	serviceCollection servicesCollection.Collection
	serverSystem      serversSystem.System
}

func New(
	cfg config.Config,
	serviceCollection servicesCollection.Collection,
	serverSystem serversSystem.System,
) Collection {
	return &collection{
		cfg:               cfg,
		serviceCollection: serviceCollection,
		serverSystem:      serverSystem,
	}
}

func (s *collection) CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	item, e := s.serviceCollection.Create(gameID, servicesCollection.CreateRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		ImageFile:   data,
	})
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, dto.Collection{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(gameID, *item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}
func (s *collection) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	e := s.serviceCollection.Delete(gameID, collectionID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *collection) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	item, e := s.serviceCollection.Item(gameID, collectionID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, dto.Collection{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(gameID, *item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}
func (s *collection) ListHandler(w http.ResponseWriter, r *http.Request) {
	s.serverSystem.StopQuit()

	gameID := mux.Vars(r)["game"]
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")
	items, e := s.serviceCollection.List(gameID, sort, search)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	respItems := make([]*dto.Collection, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, &dto.Collection{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Image:       item.Image,
			CachedImage: s.calculateCachedImage(gameID, *item),
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	network.ResponseWithMeta(w, respItems, &network.Meta{
		Total: len(respItems),
	})
}
func (s *collection) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]

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

	item, e := s.serviceCollection.Update(gameID, collectionID, servicesCollection.UpdateRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		ImageFile:   data,
	})
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, dto.Collection{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(gameID, *item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}

func (s *collection) calculateCachedImage(gameID string, collection entitiesCollection.Collection) string {
	return fmt.Sprintf(s.cfg.CollectionImagePath+"?%s", gameID, collection.ID, utils.HashForTime(&collection.UpdatedAt))
}
