package api

import (
	"net/http"

	"github.com/gorilla/mux"

	serversImage "github.com/HardDie/DeckBuilder/internal/servers/image"
)

func RegisterImageServer(route *mux.Router, srv serversImage.Image) {
	GamesRoute := route.PathPrefix("/api/games").Subrouter()
	GamesRoute.HandleFunc("/{game}/image", srv.GameHandler).Methods(http.MethodGet)

	CollectionsRoute := GamesRoute.PathPrefix("/{game}/collections").Subrouter()
	CollectionsRoute.HandleFunc("/{collection}/image", srv.CollectionHandler).Methods(http.MethodGet)

	DecksRoute := CollectionsRoute.PathPrefix("/{collection}/decks").Subrouter()
	DecksRoute.HandleFunc("/{deck}/image", srv.DeckHandler).Methods(http.MethodGet)

	CardsRoute := DecksRoute.PathPrefix("/{deck}/cards").Subrouter()
	CardsRoute.HandleFunc("/{card}/image", srv.CardHandler).Methods(http.MethodGet)
}

type UnimplementedImageServer struct {
}

var (
	// Validation
	_ serversImage.Image = &UnimplementedImageServer{}
)

// Requesting an image of existing card
//
// swagger:parameters RequestCardImage
type RequestCardImage struct {
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

// Card image
//
// swagger:response ResponseCardImage
type ResponseCardImage struct {
	// In: body
	Body []byte
}

// swagger:route GET /api/games/{game}/collections/{collection}/decks/{deck}/cards/{card}/image Images RequestCardImage
//
// # Get card image
//
// Get an image of existing card
//
//	Produces:
//	- application/json
//	- image/png
//	- image/jpeg
//	- image/gif
//
//	Responses:
//	  200: ResponseCardImage
//	  default: ResponseError
func (s *UnimplementedImageServer) CardHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting an image of existing collection
//
// swagger:parameters RequestCollectionImage
type RequestCollectionImage struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
}

// Collection image
//
// swagger:response ResponseCollectionImage
type ResponseCollectionImage struct {
	// In: body
	Body []byte
}

// swagger:route GET /api/games/{game}/collections/{collection}/image Images RequestCollectionImage
//
// # Get collection image
//
// Get an image of existing collection
//
//	Produces:
//	- application/json
//	- image/png
//	- image/jpeg
//	- image/gif
//
//	Responses:
//	  200: ResponseCollectionImage
//	  default: ResponseError
func (s *UnimplementedImageServer) CollectionHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting an image of existing deck
//
// swagger:parameters RequestDeckImage
type RequestDeckImage struct {
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

// Deck image
//
// swagger:response ResponseDeckImage
type ResponseDeckImage struct {
	// In: body
	Body []byte
}

// swagger:route GET /api/games/{game}/collections/{collection}/decks/{deck}/image Images RequestDeckImage
//
// # Get deck image
//
// Get an image of existing deck
//
//	Produces:
//	- application/json
//	- image/png
//	- image/jpeg
//	- image/gif
//
//	Responses:
//	  200: ResponseDeckImage
//	  default: ResponseError
func (s *UnimplementedImageServer) DeckHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting an image of existing game
//
// swagger:parameters RequestGameImage
type RequestGameImage struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// Game image
//
// swagger:response ResponseGameImage
type ResponseGameImage struct {
	// In: body
	Body []byte
}

// swagger:route GET /api/games/{game}/image Images RequestGameImage
//
// # Get game image
//
// Get an image of existing game
//
//	Produces:
//	- application/json
//	- image/png
//	- image/jpeg
//	- image/gif
//
//	Responses:
//	  200: ResponseGameImage
//	  default: ResponseError
func (s *UnimplementedImageServer) GameHandler(w http.ResponseWriter, r *http.Request) {}
