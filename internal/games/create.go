package games

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

type CreateGameRequest struct {
	GameInfo
}

func CreateGame(req *CreateGameRequest) (e *errors.Error) {
	// Check if game already exists
	dstDir := filepath.Join(config.GetConfig().Games(), req.Name)
	_, err := os.Stat(dstDir)
	if !os.IsNotExist(err) {
		e = errors.FileExists
		return
	}

	// Try to create folder with game
	err = os.Mkdir(dstDir, 0755)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}

	// Create info file
	err = ioutil.WriteFile(filepath.Join(dstDir, GameInfoFilename), utils.ToJson(req), 0644)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
