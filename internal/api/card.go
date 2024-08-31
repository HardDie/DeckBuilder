package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	serversCard "github.com/HardDie/DeckBuilder/internal/servers/card"
)

func RegisterCardServer(route *mux.Router, srv serversCard.Card) {
	CardsRoute := route.PathPrefix("/api/games/{game}/collections/{collection}/decks/{deck}/cards").Subrouter()
	CardsRoute.HandleFunc("", srv.ListHandler).Methods(http.MethodGet)
	CardsRoute.HandleFunc("", srv.CreateHandler).Methods(http.MethodPost)
	CardsRoute.HandleFunc("/{card}", srv.DeleteHandler).Methods(http.MethodDelete)
	CardsRoute.HandleFunc("/{card}", srv.ItemHandler).Methods(http.MethodGet)
	CardsRoute.HandleFunc("/{card}", srv.UpdateHandler).Methods(http.MethodPatch)
}

type UnimplementedCardServer struct {
}

var (
	// Validation
	_ serversCard.Card = &UnimplementedCardServer{}
)

// Request to create a card
//
// swagger:parameters RequestCreateCard
type RequestCreateCard struct {
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
	Variables string `json:"variables"`
	// In: formData
	// Required: true
	Count int `json:"count"`
	// In: formData
	// Required: false
	ImageFile []byte `json:"imageFile"`
}

// Status of card creation
//
// swagger:response ResponseCreateCard
type ResponseCreateCard struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data dto.Card `json:"data"`
	}
}

// swagger:route POST /api/games/{game}/collections/{collection}/decks/{deck}/cards Cards RequestCreateCard
//
// # Create card
//
// Allows you to create a new card
//
//	Consumes:
//	- multipart/form-data
//
//	Responses:
//	  200: ResponseCreateCard
//	  default: ResponseError
func (s *UnimplementedCardServer) CreateHandler(w http.ResponseWriter, r *http.Request) {}

// Request to delete a card
//
// swagger:parameters RequestDeleteCard
type RequestDeleteCard struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
	// In: path
	// Required: true
	Card string `json:"card"`
}

// Card deletion status
//
// swagger:response ResponseDeleteCard
type ResponseDeleteCard struct {
}

// swagger:route DELETE /api/games/{game}/collections/{collection}/decks/{deck}/cards/{card} Cards RequestDeleteCard
//
// # Delete card
//
// Allows you to delete an existing card
//
//	Responses:
//	  200: ResponseDeleteCard
//	  default: ResponseError
func (s *UnimplementedCardServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting an existing card
//
// swagger:parameters RequestCard
type RequestCard struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
	// In: path
	// Required: true
	Card int64 `json:"card"`
}

// Card
//
// swagger:response ResponseCard
type ResponseCard struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data dto.Card `json:"data"`
	}
}

// swagger:route GET /api/games/{game}/collections/{collection}/decks/{deck}/cards/{card} Cards RequestCard
//
// # Get card
//
// Get an existing card
//
//	Responses:
//	  200: ResponseCard
//	  default: ResponseError
func (s *UnimplementedCardServer) ItemHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting a list of existing cards
//
// swagger:parameters RequestListOfCard
type RequestListOfCard struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
	// In: query
	// Required: false
	Sort string `json:"sort"`
	// In: query
	// Required: false
	Search string `json:"search"`
}

// List of cards
//
// swagger:response ResponseListOfCard
type ResponseListOfCard struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []*dto.Card `json:"data"`
		// Required: true
		Meta *network.Meta `json:"meta"`
	}
}

// swagger:route GET /api/games/{game}/collections/{collection}/decks/{deck}/cards Cards RequestListOfCard
//
// # Get cards list
//
// Get a list of existing cards
// Sort values: name, name_desc, created, created_desc
//
//	Responses:
//	  200: ResponseListOfCard
//	  default: ResponseError
func (s *UnimplementedCardServer) ListHandler(w http.ResponseWriter, r *http.Request) {}

// Request to update a card
//
// swagger:parameters RequestUpdateCard
type RequestUpdateCard struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: path
	// Required: true
	Deck string `json:"deck"`
	// In: path
	// Required: true
	Card string `json:"card"`
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
	Variables string `json:"variables"`
	// In: formData
	// Required: true
	Count int `json:"count"`
	// In: formData
	// Required: false
	ImageFile []byte `json:"imageFile"`
}

// Status of card update
//
// swagger:response ResponseUpdateCard
type ResponseUpdateCard struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data dto.Card `json:"data"`
	}
}

// swagger:route PATCH /api/games/{game}/collections/{collection}/decks/{deck}/cards/{card} Cards RequestUpdateCard
//
// # Update card
//
// Allows you to update an existing card
//
//	Consumes:
//	- multipart/form-data
//
//	Responses:
//	  200: ResponseUpdateCard
//	  default: ResponseError
func (s *UnimplementedCardServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {}
