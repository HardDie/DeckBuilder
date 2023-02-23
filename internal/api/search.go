package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/network"
)

type ISearchServer interface {
	RootHandler(w http.ResponseWriter, r *http.Request)
	GameHandler(w http.ResponseWriter, r *http.Request)
	CollectionHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterSearchServer(route *mux.Router, srv ISearchServer) {
	SearchRoute := route.PathPrefix("/api/search").Subrouter()
	SearchRoute.HandleFunc("", srv.RootHandler).Methods(http.MethodGet)
	SearchRoute.HandleFunc("/games/{game}", srv.GameHandler).Methods(http.MethodGet)
	SearchRoute.HandleFunc("/games/{game}/collections/{collection}", srv.CollectionHandler).Methods(http.MethodGet)
}

type UnimplementedSearchServer struct {
}

var (
	// Validation
	_ ISearchServer = &UnimplementedSearchServer{}
)

// Recursive search for all object types
//
// swagger:parameters RequestRootSearch
type RequestRootSearch struct {
	// In: query
	// Required: false
	Sort string `json:"sort"`
	// In: query
	// Required: false
	Search string `json:"search"`
}

// List of found objects
//
// swagger:response ResponseRootSearch
type ResponseRootSearch struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []byte `json:"data"`
		// Required: true
		Meta *network.Meta `json:"meta"`
	}
}

// swagger:route GET /api/search Search RequestRootSearch
//
// # Recursive search for all games
//
// Get a list of all objects that matched the search mask
// Sort values: name, name_desc, created, created_desc
//
//	Responses:
//	  200: ResponseRootSearch
//	  default: ResponseError
func (s *UnimplementedSearchServer) RootHandler(w http.ResponseWriter, r *http.Request) {}

// Recursive search for all types of objects in the specified game
//
// swagger:parameters RequestGameSearch
type RequestGameSearch struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: query
	// Required: false
	Sort string `json:"sort"`
	// In: query
	// Required: false
	Search string `json:"search"`
}

// List of found objects
//
// swagger:response ResponseGameSearch
type ResponseGameSearch struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []byte `json:"data"`
		// Required: true
		Meta *network.Meta `json:"meta"`
	}
}

// swagger:route GET /api/search/games/{game} Search RequestGameSearch
//
// # Recursive search in a particular game
//
// Get a list of all objects that match the search mask and are in a particular game
// Sort values: name, name_desc, created, created_desc
//
//	Responses:
//	  200: ResponseGameSearch
//	  default: ResponseError
func (s *UnimplementedSearchServer) GameHandler(w http.ResponseWriter, r *http.Request) {}

// Recursive search for all types of objects in the specified game and collection
//
// swagger:parameters RequestCollectionSearch
type RequestCollectionSearch struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: query
	// Required: false
	Sort string `json:"sort"`
	// In: query
	// Required: false
	Search string `json:"search"`
}

// List of found objects
//
// swagger:response ResponseCollectionSearch
type ResponseCollectionSearch struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []byte `json:"data"`
		// Required: true
		Meta *network.Meta `json:"meta"`
	}
}

// swagger:route GET /api/search/games/{game}/collections/{collection} Search RequestCollectionSearch
//
// # Recursive search in a particular game and collection
//
// Get a list of all objects that match the search mask and are in a particular game and collection
// Sort values: name, name_desc, created, created_desc
//
//	Responses:
//	  200: ResponseCollectionSearch
//	  default: ResponseError
func (s *UnimplementedSearchServer) CollectionHandler(w http.ResponseWriter, r *http.Request) {}
