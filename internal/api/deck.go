package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type IDeckServer interface {
	AllDecksHandler(w http.ResponseWriter, r *http.Request)
	CreateHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
	ItemHandler(w http.ResponseWriter, r *http.Request)
	ListHandler(w http.ResponseWriter, r *http.Request)
	UpdateHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterDeckServer(route *mux.Router, srv IDeckServer) {
	DecksRoute := route.PathPrefix("/api/games/{game}/collections/{collection}/decks").Subrouter()
	DecksRoute.HandleFunc("", srv.ListHandler).Methods(http.MethodGet)
	DecksRoute.HandleFunc("", srv.CreateHandler).Methods(http.MethodPost)
	DecksRoute.HandleFunc("/{deck}", srv.DeleteHandler).Methods(http.MethodDelete)
	DecksRoute.HandleFunc("/{deck}", srv.ItemHandler).Methods(http.MethodGet)
	DecksRoute.HandleFunc("/{deck}", srv.UpdateHandler).Methods(http.MethodPatch)
	route.HandleFunc("/api/games/{game}/decks", srv.AllDecksHandler).Methods(http.MethodGet)
}

type UnimplementedDeckServer struct {
}

var (
	// Validation
	_ IDeckServer = &UnimplementedDeckServer{}
)

// Requesting a list of all decks in game
//
// swagger:parameters RequestListOfAllDecks
type RequestListOfAllDecks struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// List of decks
//
// swagger:response ResponseListOfAllDecks
type ResponseListOfAllDecks struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []*entity.DeckInfo `json:"data"`
	}
}

// swagger:route GET /api/games/{game}/decks Decks RequestListOfAllDecks
//
// # Get list of all decks in game
//
// Get a list of all existing decks in game
//
//	Responses:
//	  200: ResponseListOfAllDecks
//	  default: ResponseError
func (s *UnimplementedDeckServer) AllDecksHandler(w http.ResponseWriter, r *http.Request) {}

// Request to create a deck
//
// swagger:parameters RequestCreateDeck
type RequestCreateDeck struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: formData
	// Required: true
	Name string `json:"name"`
	// In: formData
	// Required: true
	Description string `json:"description"`
	// In: formData
	// Required: false
	Image string `json:"image"`
	// In: formData
	// Required: false
	ImageFile []byte `json:"imageFile"`
}

// Status of deck creation
//
// swagger:response ResponseCreateDeck
type ResponseCreateDeck struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data entity.DeckInfo `json:"data"`
	}
}

// swagger:route POST /api/games/{game}/collections/{collection}/decks Decks RequestCreateDeck
//
// # Create deck
//
// Allows you to create a new deck
//
//	Consumes:
//	- multipart/form-data
//
//	Responses:
//	  200: ResponseCreateDeck
//	  default: ResponseError
func (s *UnimplementedDeckServer) CreateHandler(w http.ResponseWriter, r *http.Request) {}

// Request to delete a deck
//
// swagger:parameters RequestDeleteDeck
type RequestDeleteDeck struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
}

// Deck deletion status
//
// swagger:response ResponseDeleteDeck
type ResponseDeleteDeck struct {
}

// swagger:route DELETE /api/games/{game}/collections/{collection}/decks/{deck} Decks RequestDeleteDeck
//
// # Delete deck
//
// Allows you to delete an existing deck
//
//	Responses:
//	  200: ResponseDeleteDeck
//	  default: ResponseError
func (s *UnimplementedDeckServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting an existing deck
//
// swagger:parameters RequestDeck
type RequestDeck struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
}

// Deck
//
// swagger:response ResponseDeck
type ResponseDeck struct {
	// In: body
	Body struct {
		// Required: true
		Data entity.DeckInfo `json:"data"`
	}
}

// swagger:route GET /api/games/{game}/collections/{collection}/decks/{deck} Decks RequestDeck
//
// # Get deck
//
// Get an existing deck
//
//	Responses:
//	  200: ResponseDeck
//	  default: ResponseError
func (s *UnimplementedDeckServer) ItemHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting a list of existing decks
//
// swagger:parameters RequestListOfDecks
type RequestListOfDecks struct {
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

// List of decks
//
// swagger:response ResponseListOfDecks
type ResponseListOfDecks struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []*entity.DeckInfo `json:"data"`
		// Required: true
		Meta *network.Meta `json:"meta"`
	}
}

// swagger:route GET /api/games/{game}/collections/{collection}/decks Decks RequestListOfDecks
//
// # Get decks list
//
// Get a list of existing decks
// Sort values: name, name_desc, created, created_desc
//
//	Responses:
//	  200: ResponseListOfDecks
//	  default: ResponseError
func (s *UnimplementedDeckServer) ListHandler(w http.ResponseWriter, r *http.Request) {}

// Request to update a deck
//
// swagger:parameters RequestUpdateDeck
type RequestUpdateDeck struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
	// In: formData
	// Required: true
	Name string `json:"name"`
	// In: formData
	// Required: true
	Description string `json:"description"`
	// In: formData
	// Required: false
	Image string `json:"image"`
	// In: formData
	// Required: false
	ImageFile []byte `json:"imageFile"`
}

// Status of deck update
//
// swagger:response ResponseUpdateDeck
type ResponseUpdateDeck struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data entity.DeckInfo `json:"data"`
	}
}

// swagger:route PATCH /api/games/{game}/collections/{collection}/decks/{deck} Decks RequestUpdateDeck
//
// # Update deck
//
// Allows you to update an existing deck
//
//	Consumes:
//	- multipart/form-data
//
//	Responses:
//	  200: ResponseUpdateDeck
//	  default: ResponseError
func (s *UnimplementedDeckServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {}
