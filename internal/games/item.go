package games

import (
	"encoding/json"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

func ItemGame(name string) (result *GameInfo, e *errors.Error) {
	// Check if game exists
	gameDir := filepath.Join(config.GetConfig().Games(), name)
	_, err := os.Stat(gameDir)
	if os.IsNotExist(err) {
		e = errors.GameNotExists
		return
	}

	// Check if info file exist
	infoFile := filepath.Join(gameDir, GameInfoFilename)
	_, err = os.Stat(infoFile)
	if os.IsNotExist(err) {
		e = errors.GameNotExists
		return
	}

	// Open file
	file, err := os.Open(infoFile)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	defer func() { utils.IfErrorLog(file.Close()) }()

	// Parse
	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
