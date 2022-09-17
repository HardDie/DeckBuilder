package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/service"
)

type CardServer struct {
	cardService service.ICardService
}

func NewCardServer(cardService service.ICardService) *CardServer {
	return &CardServer{
		cardService: cardService,
	}
}

func (s *CardServer) CreateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	dtoObject := &dto.CreateCardDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	item, e := s.cardService.Create(gameID, collectionID, deckID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *CardServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	e = s.cardService.Delete(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *CardServer) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	item, e := s.cardService.Item(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *CardServer) ListHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	sort := r.URL.Query().Get("sort")
	items, e := s.cardService.List(gameID, collectionID, deckID, sort)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
}
func (s *CardServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	dtoObject := &dto.UpdateCardDTO{}
	e = network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	item, e := s.cardService.Update(gameID, collectionID, deckID, cardID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
