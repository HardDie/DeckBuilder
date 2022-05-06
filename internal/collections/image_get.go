package collections

import (
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
)

func GetImage(game, collection string) (img []byte, imgType string, e error) {
	// Check if collection and collection info exists
	e = FullCollectionCheck(game, collection)
	if e != nil {
		return
	}

	// Build path to image
	collectionImageFile := filepath.Join(config.GetConfig().Games(), game, collection, config.GetConfig().ImageFilename)

	// Check if image exist
	isExist, e := fs.FileExist(collectionImageFile)
	if e != nil {
		return
	}
	if !isExist {
		e = errors.CollectionImageNotExists
		return
	}

	// Read image from file
	img, imgType, e = fs.ReadImageFromFile(collectionImageFile)
	return
}
