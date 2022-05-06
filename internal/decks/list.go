package decks

import (
	"io/ioutil"
	"path/filepath"
	"sort"

	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
)

type ListOfDecksResponse struct {
	Decks []*DeckInfo `json:"decks"`
}

func ListOfDecks(gameName, collectionName string) (result *ListOfDecksResponse, e error) {
	result = &ListOfDecksResponse{
		Decks: make([]*DeckInfo, 0),
	}

	e = collections.FullCollectionCheck(gameName, collectionName)
	if e != nil {
		return
	}

	// Get decks from collection
	collectionPath := filepath.Join(config.GetConfig().Games(), gameName, collectionName)
	files, err := ioutil.ReadDir(collectionPath)
	if err != nil {
		e = errors.InternalError.AddMessage(err.Error())
		return
	}

	// Get decks info
	e, result.Decks = GetDecksFromCollection(gameName, collectionName, files)

	// Sort decks in result
	sort.SliceStable(result.Decks, func(i, j int) bool {
		return result.Decks[i].Type < result.Decks[j].Type
	})
	return
}
