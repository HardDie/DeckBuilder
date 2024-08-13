package deck

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/network"
	serversSystem "github.com/HardDie/DeckBuilder/internal/servers/system"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type deck struct {
	serviceDeck  servicesDeck.Deck
	serverSystem serversSystem.System
}

func New(serviceDeck servicesDeck.Deck, serverSystem serversSystem.System) Deck {
	return &deck{
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
	network.Response(w, items)
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

	dtoObject := &dto.CreateDeckDTO{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		ImageFile:   data,
	}

	item, e := s.serviceDeck.Create(gameID, collectionID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
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
	network.Response(w, item)
}
func (s *deck) ListHandler(w http.ResponseWriter, r *http.Request) {
	s.serverSystem.StopQuit()

	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")
	items, meta, e := s.serviceDeck.List(gameID, collectionID, sort, search)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.ResponseWithMeta(w, items, meta)
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

	dtoObject := &dto.UpdateDeckDTO{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		ImageFile:   data,
	}

	item, e := s.serviceDeck.Update(gameID, collectionID, deckID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
