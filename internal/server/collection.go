package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/service"
)

type CollectionServer struct {
	collectionService service.ICollectionService
}

func NewCollectionServer(collectionService service.ICollectionService) *CollectionServer {
	return &CollectionServer{
		collectionService: collectionService,
	}
}

func (s *CollectionServer) CreateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	dtoObject := &dto.CreateCollectionDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := s.collectionService.Create(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
func (s *CollectionServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	e := s.collectionService.Delete(gameID, collectionID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *CollectionServer) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	item, e := s.collectionService.Item(gameID, collectionID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *CollectionServer) ListHandler(w http.ResponseWriter, r *http.Request) {
	NewSystemServer(nil).StopQuit()

	gameID := mux.Vars(r)["game"]
	sort := r.URL.Query().Get("sort")
	items, e := s.collectionService.List(gameID, sort)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, items)
}
func (s *CollectionServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	dtoObject := &dto.UpdateCollectionDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	item, e := s.collectionService.Update(gameID, collectionID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
