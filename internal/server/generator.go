package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/service"
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

	if dtoObject.Scale < 1 {
		dtoObject.Scale = 1
	}

	gameID := mux.Vars(r)["game"]
	e = s.generatorService.GenerateGame(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
}
