package decks

import (
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/errors"
)

func ItemDeck(gameName, collectionName, deckName string) (result *DeckInfo, e *errors.Error) {
	// Check if collection and collection info exists
	e = collections.FullCollectionCheck(gameName, collectionName)
	if e != nil {
		return
	}

	deckName += ".json"

	// Check deck exist
	exist, e := DeckIsExist(gameName, collectionName, deckName)
	if e != nil {
		return
	}
	if !exist {
		e = errors.DeckNotExists
		return
	}

	// Get info
	result, e = DeckGetInfo(gameName, collectionName, deckName)
	if e != nil {
		return
	}
	return
}
