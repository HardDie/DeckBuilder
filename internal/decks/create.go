package decks

import (
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

type CreateDeckRequest struct {
	DeckInfoWithoutId
}

func CreateDeck(gameName, collectionName string, req *CreateDeckRequest) (res DeckInfo, e *errors.Error) {
	e = collections.FullCollectionCheck(gameName, collectionName)
	if e != nil {
		return
	}

	res = DeckInfo{
		Id:                utils.NameToId(req.DeckInfoWithoutId.Type),
		DeckInfoWithoutId: req.DeckInfoWithoutId,
	}

	// Check if deck id correct
	if len(res.Id) == 0 {
		e = errors.BadName
		return
	}

	// Check if deck already exists
	exist, e := DeckIsExist(gameName, collectionName, res.Id)
	if e != nil {
		return
	}
	if exist {
		e = errors.DeckExist
		return
	}

	// Try to create deck
	e = DeckCreate(gameName, collectionName, res.Id, res)
	return
}
