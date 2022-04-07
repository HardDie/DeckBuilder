package games

import (
	"log"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
)

type CreateGameRequest struct {
	Name string `json:"name"`
}

type CreateGameResponse struct {
	Message string `json:"message"`
}

func CreateGame(req *CreateGameRequest) (response CreateGameResponse) {
	dstDir := filepath.Join(config.GetConfig().Games(), req.Name)
	_, err := os.Stat(dstDir)
	if !os.IsNotExist(err) {
		response.Message = errors.FileExists
		return
	}
	err = os.Mkdir(dstDir, 0755)
	if err != nil {
		log.Println(err.Error())
		response.Message = errors.UnknownError
		return
	}
	response.Message = errors.Done
	return
}
