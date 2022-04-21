package decks

import (
	"log"
	"os"
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

func GetDecksFromCollection(gameName, collectionName string, files []os.FileInfo) (e *errors.Error, decks []*DeckInfo) {
	for _, file := range files {
		if file.IsDir() {
			// Skip folders
			continue
		}
		if filepath.Ext(file.Name()) != ".json" {
			// Skip non json files
			continue
		}
		if file.Name() == config.GetConfig().InfoFilename {
			// Skip info collection files
			continue
		}

		var item *DeckInfo

		// Get info
		item, e = DeckGetInfo(gameName, collectionName, file.Name())
		if e != nil {
			log.Println("Bad deck:", file.Name())
			continue
		}

		// Append collection info to list
		decks = append(decks, item)
	}
	return
}
