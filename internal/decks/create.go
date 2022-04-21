package decks

import (
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/errors"
)

type CreateDeckRequest struct {
	DeckInfo
}

func CreateDeck(gameName, collectionName string, req *CreateDeckRequest) (e *errors.Error) {
	e = collections.FullCollectionCheck(gameName, collectionName)
	if e != nil {
		return
	}

	// Check if deck already exists
	exist, e := DeckIsExist(gameName, collectionName, req.Type)
	if e != nil {
		return
	}
	if exist {
		e = errors.DeckExist
		return
	}

	// Try to create deck
	e = DeckCreate(gameName, collectionName, req.Type, req.DeckInfo)
	return
}
