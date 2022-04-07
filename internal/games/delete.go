package games

import (
	"log"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
)

type DeleteGameRequest struct {
	Name string `json:"name"`
}

type DeleteGameResponse struct {
	Message string `json:"message"`
}

func DeleteGame(req *DeleteGameRequest) (response DeleteGameResponse) {
	dstDir := filepath.Join(config.GetConfig().Games(), req.Name)
	_, err := os.Stat(dstDir)
	if os.IsNotExist(err) {
		response.Message = errors.FileNotExist
		return
	}
	err = os.RemoveAll(dstDir)
	if err != nil {
		log.Println(err.Error())
		response.Message = errors.UnknownError
		return
	}
	response.Message = errors.Done
	return
}
