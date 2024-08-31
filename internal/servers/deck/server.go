package deck

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/network"
	serversSystem "github.com/HardDie/DeckBuilder/internal/servers/system"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type deck struct {
	cfg          config.Config
	serviceDeck  servicesDeck.Deck
	serverSystem serversSystem.System
}

func New(
	cfg config.Config,
	serviceDeck servicesDeck.Deck,
	serverSystem serversSystem.System,
) Deck {
	return &deck{
		cfg:          cfg,
		serviceDeck:  serviceDeck,
		serverSystem: serverSystem,
	}
}

func (s *deck) AllDecksHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	items, e := s.serviceDeck.ListAllUnique(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	respItems := make([]*dto.Deck, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, &dto.Deck{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Image:       item.Image,
			CachedImage: s.calculateCachedImage(*item),
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	network.Response(w, respItems)
}
func (s *deck) CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	item, e := s.serviceDeck.Create(gameID, collectionID, servicesDeck.CreateRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		ImageFile:   data,
	})
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, dto.Deck{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}
func (s *deck) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	e := s.serviceDeck.Delete(gameID, collectionID, deckID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *deck) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	item, e := s.serviceDeck.Item(gameID, collectionID, deckID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, dto.Deck{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}
func (s *deck) ListHandler(w http.ResponseWriter, r *http.Request) {
	s.serverSystem.StopQuit()

	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")
	items, e := s.serviceDeck.List(gameID, collectionID, sort, search)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	respItems := make([]*dto.Deck, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, &dto.Deck{
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
func (s *deck) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]

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

	item, e := s.serviceDeck.Update(gameID, collectionID, deckID, servicesDeck.UpdateRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		ImageFile:   data,
	})
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, dto.Deck{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}

func (s *deck) calculateCachedImage(deck entitiesDeck.Deck) string {
	return fmt.Sprintf(s.cfg.DeckImagePath+"?%s", deck.GameID, deck.CollectionID, deck.ID, utils.HashForTime(&deck.UpdatedAt))
}
