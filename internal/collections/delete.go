package collections

import "tts_deck_build/internal/errors"

func DeleteCollection(gameName, collectionName string) (e *errors.Error) {
	// Check if collection exists
	exist, e := CollectionIsExist(gameName, collectionName)
	if e != nil {
		return
	}
	if !exist {
		e = errors.GameNotExists
		return
	}

	// Try to delete collection
	return CollectionDelete(gameName, collectionName)
}
