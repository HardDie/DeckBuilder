package decks

import (
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

type UpdateDeckRequest struct {
	DeckInfoWithoutId
}

func UpdateDeck(gameName, collectionName, deckName string, req *UpdateDeckRequest) (res DeckInfo, e *errors.Error) {
	res = DeckInfo{
		Id:                utils.NameToId(req.DeckInfoWithoutId.Type),
		DeckInfoWithoutId: req.DeckInfoWithoutId,
	}

	// Check if deck id correct
	if len(res.Id) == 0 {
		e = errors.BadName
		return
	}

	// Check if game and collection exists
	e = collections.FullCollectionCheck(gameName, collectionName)
	if e != nil {
		return
	}

	// Check if deck exists
	exist, e := DeckIsExist(gameName, collectionName, deckName)
	if e != nil {
		return
	}
	if !exist {
		e = errors.DeckNotExists
		return
	}

	// Rename deck if name changed
	if res.Id != deckName {
		e = DeckRename(gameName, collectionName, deckName, res.Id)
		if e != nil {
			return
		}
	}

	// Update info file
	e = DeckCreate(gameName, collectionName, res.Id, res)
	return
}
