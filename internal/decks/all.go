package decks

import (
	"io/ioutil"
	"path/filepath"
	"sort"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
)

func ListOfAllDecks(gameName string) (result *ListOfDecksResponse, e error) {
	result = &ListOfDecksResponse{
		Decks: make([]*DeckInfo, 0),
	}

	e = games.FullGameCheck(gameName)
	if e != nil {
		return
	}

	// Get list of collections
	gamePath := filepath.Join(config.GetConfig().Games(), gameName)
	collections, err := ioutil.ReadDir(gamePath)
	if err != nil {
		e = errors.InternalError.AddMessage(err.Error())
		return
	}

	uniq := make(map[string]struct{})

	for _, collection := range collections {
		if !collection.IsDir() {
			continue
		}

		var decks []*DeckInfo

		// Get list of decks
		collectionPath := filepath.Join(gamePath, collection.Name())
		deckFiles, err := ioutil.ReadDir(collectionPath)
		if err != nil {
			e = errors.InternalError.AddMessage(err.Error())
			return
		}

		e, decks = GetDecksFromCollection(gameName, collection.Name(), deckFiles)
		if e != nil {
			return
		}
		for _, deck := range decks {
			if _, ok := uniq[deck.Type+deck.BacksideImage]; ok {
				continue
			}
			uniq[deck.Type+deck.BacksideImage] = struct{}{}
			result.Decks = append(result.Decks, deck)
		}
	}

	// Sort decks in result
	sort.SliceStable(result.Decks, func(i, j int) bool {
		return result.Decks[i].Type < result.Decks[j].Type
	})
	return
}
