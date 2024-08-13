package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/network"
	servicesSearch "github.com/HardDie/DeckBuilder/internal/services/search"
)

type SearchServer struct {
	serviceSearch servicesSearch.Search
}

func NewSearchServer(serviceSearch servicesSearch.Search) *SearchServer {
	return &SearchServer{
		serviceSearch: serviceSearch,
	}
}

func (s *SearchServer) RootHandler(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")

	resp, meta, err := s.serviceSearch.RecursiveSearch(sort, search, "", "")
	if err != nil {
		network.ResponseError(w, err)
		return
	}
	network.ResponseWithMeta(w, resp, meta)
}
func (s *SearchServer) GameHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")

	resp, meta, err := s.serviceSearch.RecursiveSearch(sort, search, gameID, "")
	if err != nil {
		network.ResponseError(w, err)
		return
	}
	network.ResponseWithMeta(w, resp, meta)
}
func (s *SearchServer) CollectionHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")

	resp, meta, err := s.serviceSearch.RecursiveSearch(sort, search, gameID, collectionID)
	if err != nil {
		network.ResponseError(w, err)
		return
	}
	network.ResponseWithMeta(w, resp, meta)
}
