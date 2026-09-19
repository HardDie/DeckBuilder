package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	serversGame "github.com/HardDie/DeckBuilder/internal/servers/game"
)

func RegisterGameServer(route *mux.Router, srv serversGame.Game) {
	GamesRoute := route.PathPrefix("/api/games").Subrouter()
	GamesRoute.HandleFunc("/import", srv.ImportHandler).Methods(http.MethodPost)
	GamesRoute.HandleFunc("/{game}/export", srv.ExportHandler).Methods(http.MethodGet)
}

type UnimplementedGameServer struct {
}

var (
	// Validation
	_ serversGame.Game = &UnimplementedGameServer{}
)

// Requesting an existing game archive
//
// swagger:parameters RequestArchiveGame
type RequestArchiveGame struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// Game archive
//
// swagger:response ResponseGameArchive
type ResponseGameArchive struct {
	// In: body
	Body []byte
}

// swagger:route GET /api/games/{game}/export Games RequestArchiveGame
//
// # Export game to archive
//
// Get an existing game archive
//
//	Produces:
//	- application/json
//	- application/zip
//
//	Responses:
//	  200: ResponseGameArchive
//	  default: ResponseError
func (s *UnimplementedGameServer) ExportHandler(w http.ResponseWriter, r *http.Request) {}

// Creating game from archive
//
// swagger:parameters RequestImportGame
type RequestImportGame struct {
	// Specify a name for the imported game
	// In: formData
	// Required: false
	Name string `json:"name"`
	// Binary data of the imported file
	// In: formData
	// Required: true
	File []byte `json:"file"`
}

// Import game
//
// swagger:response ResponseGameImport
type ResponseGameImport struct {
	// In: body
	Body struct {
		// Required: true
		Data dto.Game `json:"data"`
	}
}

// swagger:route POST /api/games/import Games RequestImportGame
//
// # Import game from archive
//
// Creat game from archive
//
//	Consumes:
//	- multipart/form-data
//
//	Responses:
//	  200: ResponseGameImport
//	  default: ResponseError
func (s *UnimplementedGameServer) ImportHandler(w http.ResponseWriter, r *http.Request) {}
