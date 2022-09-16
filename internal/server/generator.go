package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/generator"
	"tts_deck_build/internal/network"
)

type GeneratorServer struct {
}

func NewGeneratorServer() *GeneratorServer {
	return &GeneratorServer{}
}

func (s *GeneratorServer) GameHandler(w http.ResponseWriter, r *http.Request) {
	dtoObject := &dto.GenerateGameDTO{}
	e := network.RequestToObject(r.Body, &dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	gameID := mux.Vars(r)["game"]
	e = generator.NewService().GenerateGame(gameID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
}
