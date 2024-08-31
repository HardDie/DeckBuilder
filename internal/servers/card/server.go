package card

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/network"
	serversSystem "github.com/HardDie/DeckBuilder/internal/servers/system"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type card struct {
	cfg          config.Config
	serviceCard  servicesCard.Card
	serverSystem serversSystem.System
}

func New(
	cfg config.Config,
	serviceCard servicesCard.Card,
	serverSystem serversSystem.System,
) Card {
	return &card{
		cfg:          cfg,
		serviceCard:  serviceCard,
		serverSystem: serverSystem,
	}
}

func (s *card) CreateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]

	e := r.ParseMultipartForm(0)
	if e != nil {
		er.IfErrorLog(e)
		e = er.InternalError.AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	data, e := utils.GetFileFromMultipart("imageFile", r)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	var variables map[string]string

	variablesJson := r.FormValue("variables")
	if variablesJson != "" {
		e = json.Unmarshal([]byte(variablesJson), &variables)
		if e != nil {
			er.IfErrorLog(e)
			e = er.InternalError.HTTP(http.StatusBadRequest).AddMessage("Bad variables json")
			network.ResponseError(w, e)
			return
		}
	}

	item, e := s.serviceCard.Create(gameID, collectionID, deckID, servicesCard.CreateRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		Variables:   variables,
		Count:       fs.StringToInt(r.FormValue("count")),
		ImageFile:   data,
	})
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, dto.Card{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		Variables:   item.Variables,
		Count:       item.Count,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}
func (s *card) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	e = s.serviceCard.Delete(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *card) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	item, e := s.serviceCard.Item(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, dto.Card{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		Variables:   item.Variables,
		Count:       item.Count,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}
func (s *card) ListHandler(w http.ResponseWriter, r *http.Request) {
	s.serverSystem.StopQuit()

	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")
	items, e := s.serviceCard.List(gameID, collectionID, deckID, sort, search)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	respItems := make([]*dto.Card, 0, len(items))
	var cardsTotal int
	for _, item := range items {
		cardsTotal += item.Count
		respItems = append(respItems, &dto.Card{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Image:       item.Image,
			CachedImage: s.calculateCachedImage(*item),
			Variables:   item.Variables,
			Count:       item.Count,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	network.ResponseWithMeta(w, respItems, &network.Meta{
		Total:      len(respItems),
		CardsTotal: cardsTotal,
	})
}
func (s *card) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	e = r.ParseMultipartForm(0)
	if e != nil {
		er.IfErrorLog(e)
		e = er.InternalError.AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	data, e := utils.GetFileFromMultipart("imageFile", r)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	var variables map[string]string

	variablesJson := r.FormValue("variables")
	if variablesJson != "" {
		e = json.Unmarshal([]byte(variablesJson), &variables)
		if e != nil {
			er.IfErrorLog(e)
			e = er.InternalError.HTTP(http.StatusBadRequest).AddMessage("Bad variables json")
			network.ResponseError(w, e)
			return
		}
	}

	item, e := s.serviceCard.Update(gameID, collectionID, deckID, cardID, servicesCard.UpdateRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		Variables:   variables,
		Count:       fs.StringToInt(r.FormValue("count")),
		ImageFile:   data,
	})
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, dto.Card{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		Variables:   item.Variables,
		Count:       item.Count,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}

func (s *card) calculateCachedImage(card entitiesCard.Card) string {
	return fmt.Sprintf(s.cfg.CardImagePath+"?%s", card.GameID, card.CollectionID, card.DeckID, card.ID, utils.HashForTime(&card.UpdatedAt))
}
