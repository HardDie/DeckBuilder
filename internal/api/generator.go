package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/dto"
)

type IGeneratorServer interface {
	GameHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterGeneratorServer(route *mux.Router, srv IGeneratorServer) {
	GeneratorsRoute := route.PathPrefix("/api/games/{game}").Subrouter()
	GeneratorsRoute.HandleFunc("/generate", srv.GameHandler).Methods(http.MethodGet)
}

type UnimplementedGeneratorServer struct {
}

var (
	// Validation
	_ IGeneratorServer = &UnimplementedGeneratorServer{}
)

// Request to start generating result objects
//
// swagger:parameters RequestGameGenerate
type RequestGameGenerate struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: body
	// Required: false
	Body struct {
		dto.GenerateGameDTO
	}
}

// Generating game objects
//
// swagger:response ResponseGameGenerate
type ResponseGameGenerate struct {
}

// swagger:route POST /api/games/{game}/generate Generator RequestGameGenerate
//
// Start generating items for TTS
//
// Allow to run the background process of generating images and json item for the game
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
//       200: ResponseGameGenerate
//       default: ResponseError
func (s *UnimplementedGeneratorServer) GameHandler(w http.ResponseWriter, r *http.Request) {}
