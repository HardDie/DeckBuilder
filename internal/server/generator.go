package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesGenerator "github.com/HardDie/DeckBuilder/internal/services/generator"
)

type GeneratorServer struct {
	serviceGenerator servicesGenerator.Generator
}

func NewGeneratorServer(serviceGenerator servicesGenerator.Generator) *GeneratorServer {
	return &GeneratorServer{
		serviceGenerator: serviceGenerator,
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
	e = s.serviceGenerator.GenerateGame(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
}
