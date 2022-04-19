package collections

import (
	"tts_deck_build/internal/errors"
)

func ItemCollection(gameName, collectionName string) (result *CollectionInfo, e *errors.Error) {
	// Check if collection and collection info exists
	e = FullCollectionCheck(gameName, collectionName)
	if e != nil {
		return
	}

	// Get info
	result, e = CollectionGetInfo(gameName, collectionName)
	if e != nil {
		return
	}
	return
}
