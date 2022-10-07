package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
)

type IGameServer interface {
	CreateHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
	DuplicateHandler(w http.ResponseWriter, r *http.Request)
	ExportHandler(w http.ResponseWriter, r *http.Request)
	ImportHandler(w http.ResponseWriter, r *http.Request)
	ItemHandler(w http.ResponseWriter, r *http.Request)
	ListHandler(w http.ResponseWriter, r *http.Request)
	UpdateHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterGameServer(route *mux.Router, srv IGameServer) {
	GamesRoute := route.PathPrefix("/api/games").Subrouter()
	GamesRoute.HandleFunc("", srv.ListHandler).Methods(http.MethodGet)
	GamesRoute.HandleFunc("", srv.CreateHandler).Methods(http.MethodPost)
	GamesRoute.HandleFunc("/{game}", srv.DeleteHandler).Methods(http.MethodDelete)
	GamesRoute.HandleFunc("/{game}", srv.ItemHandler).Methods(http.MethodGet)
	GamesRoute.HandleFunc("/{game}", srv.UpdateHandler).Methods(http.MethodPatch)
	GamesRoute.HandleFunc("/{game}/duplicate", srv.DuplicateHandler).Methods(http.MethodPost)
	GamesRoute.HandleFunc("/{game}/export", srv.ExportHandler).Methods(http.MethodGet)
	GamesRoute.HandleFunc("/import", srv.ImportHandler).Methods(http.MethodPost)
}

type UnimplementedGameServer struct {
}

var (
	// Validation
	_ IGameServer = &UnimplementedGameServer{}
)

// Request to create a game
//
// swagger:parameters RequestCreateGame
type RequestCreateGame struct {
	// In: body
	// Required: true
	Body struct {
		dto.CreateGameDTO
	}
}

// Status of game creation
//
// swagger:response ResponseCreateGame
type ResponseCreateGame struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data entity.GameInfo `json:"data"`
	}
}

// swagger:route POST /api/games Games RequestCreateGame
//
// Create game
//
// Allows you to create a new game
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: ResponseCreateGame
//       default: ResponseError
func (s *UnimplementedGameServer) CreateHandler(w http.ResponseWriter, r *http.Request) {}

// Request to delete a game
//
// swagger:parameters RequestDeleteGame
type RequestDeleteGame struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// Game deletion status
//
// swagger:response ResponseDeleteGame
type ResponseDeleteGame struct {
}

// swagger:route DELETE /api/games/{game} Games RequestDeleteGame
//
// Delete game
//
// Allows you to delete an existing game
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: ResponseDeleteGame
//       default: ResponseError
func (s *UnimplementedGameServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {}

// Request to duplicate a game
//
// swagger:parameters RequestDuplicateGame
type RequestDuplicateGame struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: body
	// Required: true
	Body struct {
		// Required: true
		dto.DuplicateGameDTO
	}
}

// Status of game duplicate
//
// swagger:response ResponseDuplicateGame
type ResponseDuplicateGame struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data entity.GameInfo `json:"data"`
	}
}

// swagger:route POST /api/games/{game}/duplicate Games RequestDuplicateGame
//
// Duplicate game
//
// Allows you to create a copy of an existing game
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: ResponseDuplicateGame
//       default: ResponseError
func (s *UnimplementedGameServer) DuplicateHandler(w http.ResponseWriter, r *http.Request) {}

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
// Export game to archive
//
// Get an existing game archive
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - application/zip
//
//     Schemes: http
//
//     Responses:
//       200: ResponseGameArchive
//       default: ResponseError
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
		Data entity.GameInfo `json:"data"`
	}
}

// swagger:route POST /api/games/import Games RequestImportGame
//
// Import game from archive
//
// Creat game from archive
//
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: ResponseGameImport
//       default: ResponseError
func (s *UnimplementedGameServer) ImportHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting an existing game
//
// swagger:parameters RequestGame
type RequestGame struct {
	// In: path
	// Required: true
	Game string `json:"game"`
}

// Game
//
// swagger:response ResponseGame
type ResponseGame struct {
	// In: body
	Body struct {
		// Required: true
		Data entity.GameInfo `json:"data"`
	}
}

// swagger:route GET /api/games/{game} Games RequestGame
//
// Get game
//
// Get an existing game
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: ResponseGame
//       default: ResponseError
func (s *UnimplementedGameServer) ItemHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting a list of existing games
//
// swagger:parameters RequestListOfGames
type RequestListOfGames struct {
	// In: query
	// Required: false
	Sort string `json:"sort"`
}

// List of games
//
// swagger:response ResponseListOfGames
type ResponseListOfGames struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []*entity.GameInfo `json:"data"`
	}
}

// swagger:route GET /api/games Games RequestListOfGames
//
// Get games list
//
// Get a list of existing games
// Sort values: name, name_desc, created, created_desc
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: ResponseListOfGames
//       default: ResponseError
func (s *UnimplementedGameServer) ListHandler(w http.ResponseWriter, r *http.Request) {}

// Request to update a game
//
// swagger:parameters RequestUpdateGame
type RequestUpdateGame struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: body
	// Required: true
	Body struct {
		// Required: true
		dto.UpdateGameDTO
	}
}

// Status of game update
//
// swagger:response ResponseUpdateGame
type ResponseUpdateGame struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data entity.GameInfo `json:"data"`
	}
}

// swagger:route PATCH /api/games/{game} Games RequestUpdateGame
//
// Update game
//
// Allows you to update an existing game
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: ResponseUpdateGame
//       default: ResponseError
func (s *UnimplementedGameServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {}
