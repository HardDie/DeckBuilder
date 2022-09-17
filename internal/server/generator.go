package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/service"
)

type GeneratorServer struct {
	generatorService service.IGeneratorService
}

func NewGeneratorServer(generatorService service.IGeneratorService) *GeneratorServer {
	return &GeneratorServer{
		generatorService: generatorService,
	}
}

func (s *GeneratorServer) GameHandler(w http.ResponseWriter, r *http.Request) {
	dtoObject := &dto.GenerateGameDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	gameID := mux.Vars(r)["game"]
	e = s.generatorService.GenerateGame(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
}
