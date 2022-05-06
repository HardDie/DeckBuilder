package games

import (
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/network"
)

func CreateImage(imgURL, game string) (e *errors.Error) {
	// Try to download image
	img, e := network.DownloadImage(imgURL)
	if e != nil {
		return
	}

	// Create image file
	gameImageFile := filepath.Join(config.GetConfig().Games(), game, config.GetConfig().ImageFilename)
	e = fs.WriteImageToFile(gameImageFile, img)
	return
}
