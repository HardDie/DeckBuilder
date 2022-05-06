package games

import (
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
)

func GetImage(game string) (img []byte, imgType string, e *errors.Error) {
	// Check if game and game info exists
	e = FullGameCheck(game)
	if e != nil {
		return
	}

	// Build path to image
	gameImageFile := filepath.Join(config.GetConfig().Games(), game, config.GetConfig().ImageFilename)

	// Check if image exist
	isExist, e := fs.FileExist(gameImageFile)
	if e != nil {
		return
	}
	if !isExist {
		e = errors.GameImageNotExists
		return
	}

	// Read image from file
	img, imgType, e = fs.ReadImageFromFile(gameImageFile)
	return
}
