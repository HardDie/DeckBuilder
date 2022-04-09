package games

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

type UpdateGameRequest struct {
	GameInfo
}

func UpdateGame(name string, req *UpdateGameRequest) (e *errors.Error) {
	// Validate
	if len(req.Name) == 0 {
		e = errors.DataInvalid.AddMessage("The name of the game cannot be empty")
	}

	// Check if game exists
	gameDir := filepath.Join(config.GetConfig().Games(), name)
	_, err := os.Stat(gameDir)
	if os.IsNotExist(err) {
		e = errors.GameNotExists
		return
	}

	// Update info file
	infoFile := filepath.Join(gameDir, GameInfoFilename)
	err = ioutil.WriteFile(infoFile, utils.ToJson(req), 0644)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}

	// Rename folder if name changed
	if req.Name != name {
		newGameDir := filepath.Join(config.GetConfig().Games(), req.Name)
		err = os.Rename(gameDir, newGameDir)
		if err != nil {
			utils.IfErrorLog(err)
			e = errors.InternalError.AddMessage(err.Error())
		}
	}
	return
}
