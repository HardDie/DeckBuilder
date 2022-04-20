package decks

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
)

type ListOfDecksResponse struct {
	Decks []*DeckInfo `json:"decks"`
}

func ListOfDecks(gameName, collectionName string) (result *ListOfDecksResponse, e *errors.Error) {
	e = collections.FullCollectionCheck(gameName, collectionName)
	if e != nil {
		return
	}
	result = &ListOfDecksResponse{
		Decks: make([]*DeckInfo, 0),
	}

	collectionPath := filepath.Join(config.GetConfig().Games(), gameName, collectionName)
	files, err := ioutil.ReadDir(collectionPath)
	if err != nil {
		e = errors.InternalError.AddMessage(err.Error())
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		if file.Name() == config.GetConfig().InfoFilename {
			continue
		}

		var item *DeckInfo

		// Get info
		item, e = DeckGetInfo(gameName, collectionName, file.Name())
		if e != nil {
			log.Println("Bad info:", file.Name())
			continue
		}

		// Append collection info to list
		result.Decks = append(result.Decks, item)
	}
	return
}
