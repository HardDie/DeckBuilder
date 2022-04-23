package decks

import (
	"log"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
)

// Deck
func DeckIsExist(gameName, collectionName, deckName string) (isExist bool, e *errors.Error) {
	infoFile := filepath.Join(config.GetConfig().Games(), gameName, collectionName, deckName+".json")
	return fs.FileExist(infoFile)
}
func DeckCreate(gameName, collectionName, deckName string, info DeckInfo) (e *errors.Error) {
	deckPath := filepath.Join(config.GetConfig().Games(), gameName, collectionName, deckName+".json")
	return fs.WriteDataToFile(deckPath, info)
}
func DeckRename(gameName, collectionName, oldName, newName string) (e *errors.Error) {
	oldDeckPath := filepath.Join(config.GetConfig().Games(), gameName, collectionName, oldName+".json")
	newDeckPath := filepath.Join(config.GetConfig().Games(), gameName, collectionName, newName+".json")
	err := os.Rename(oldDeckPath, newDeckPath)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func DeckDelete(gameName, collectionName, deckName string) (e *errors.Error) {
	deckPath := filepath.Join(config.GetConfig().Games(), gameName, collectionName, deckName+".json")
	return fs.RemoveDir(deckPath)
}

// Info
func DeckGetInfo(gameName, collectionName, deckName string) (result *DeckInfo, e *errors.Error) {
	infoFile := filepath.Join(config.GetConfig().Games(), gameName, collectionName, deckName)
	return fs.ReadDataFromFile[DeckInfo](infoFile)
}

func GetDecksFromCollection(gameName, collectionName string, files []os.FileInfo) (e *errors.Error, decks []*DeckInfo) {
	decks = make([]*DeckInfo, 0)

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
