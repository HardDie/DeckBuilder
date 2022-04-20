package decks

import (
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
)

// Info
func DeckIsExist(gameName, collectionName, deckName string) (isExist bool, e *errors.Error) {
	infoFile := filepath.Join(config.GetConfig().Games(), gameName, collectionName, deckName)
	return fs.FileExist(infoFile)
}
func DeckGetInfo(gameName, collectionName, deckName string) (result *DeckInfo, e *errors.Error) {
	infoFile := filepath.Join(config.GetConfig().Games(), gameName, collectionName, deckName)
	return fs.ReadDataFromFile[DeckInfo](infoFile)
}
