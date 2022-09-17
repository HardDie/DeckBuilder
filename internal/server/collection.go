package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/network"
)

type CollectionServer struct {
	cfg *config.Config
}

func NewCollectionServer(cfg *config.Config) *CollectionServer {
	return &CollectionServer{
		cfg: cfg,
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

	item, e := collections.NewService(s.cfg).Create(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
func (s *CollectionServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	e := collections.NewService(s.cfg).Delete(gameID, collectionID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *CollectionServer) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	item, e := collections.NewService(s.cfg).Item(gameID, collectionID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *CollectionServer) ListHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	sort := r.URL.Query().Get("sort")
	items, e := collections.NewService(s.cfg).List(gameID, sort)
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

	item, e := collections.NewService(s.cfg).Update(gameID, collectionID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, item)
}
