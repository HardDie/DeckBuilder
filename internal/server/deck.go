package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/network"
)

type DeckServer struct {
}

func NewDeckServer() *DeckServer {
	return &DeckServer{}
}

func (s *DeckServer) AllDecksHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	items, e := decks.NewService().ListAllUnique(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
}
func (s *DeckServer) CreateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	dtoObject := &dto.CreateDeckDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := decks.NewService().Create(gameID, collectionID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
func (s *DeckServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	e := decks.NewService().Delete(gameID, collectionID, deckID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *DeckServer) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	item, e := decks.NewService().Item(gameID, collectionID, deckID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *DeckServer) ListHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	sort := r.URL.Query().Get("sort")
	items, e := decks.NewService().List(gameID, collectionID, sort)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
}
func (s *DeckServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	dtoObject := &dto.UpdateDeckDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := decks.NewService().Update(gameID, collectionID, deckID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
