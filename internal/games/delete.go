package games

import (
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

func DeleteGame(name string) (e *errors.Error) {
	// Check if game exists
	dstDir := filepath.Join(config.GetConfig().Games(), name)
	_, err := os.Stat(dstDir)
	if os.IsNotExist(err) {
		e = errors.GameNotExists
		return
	}

	// Try to delete game
	err = os.RemoveAll(dstDir)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
