package api

import (
	"net/http"

	"github.com/gorilla/mux"

	serversGenerator "github.com/HardDie/DeckBuilder/internal/servers/generator"
)

func RegisterGeneratorServer(route *mux.Router, srv serversGenerator.Generator) {
	GeneratorsRoute := route.PathPrefix("/api/games/{game}").Subrouter()
	GeneratorsRoute.HandleFunc("/generate", srv.GameHandler).Methods(http.MethodPost)
}

type UnimplementedGeneratorServer struct {
}

var (
	// Validation
	_ serversGenerator.Generator = &UnimplementedGeneratorServer{}
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
		SortOrder string `json:"sortOrder"`
		Scale     int    `json:"scale"`
	}
}

// Generating game objects
//
// swagger:response ResponseGameGenerate
type ResponseGameGenerate struct {
}

// swagger:route POST /api/games/{game}/generate Generator RequestGameGenerate
//
// # Start generating items for TTS
//
// Allow to run the background process of generating images and json item for the game
//
//	Responses:
//	  200: ResponseGameGenerate
//	  default: ResponseError
func (s *UnimplementedGeneratorServer) GameHandler(w http.ResponseWriter, r *http.Request) {}
