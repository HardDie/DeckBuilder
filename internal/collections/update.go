package collections

import (
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
)

type UpdateCollectionRequest struct {
	CollectionInfo
}

func UpdateCollection(gameName, collectionName string, req *UpdateCollectionRequest) (e *errors.Error) {
	// Validate
	if len(req.Name) == 0 {
		e = errors.DataInvalid.AddMessage("The name of the collection cannot be empty")
		return
	}

	// Check if game exists
	e = games.FullGameCheck(gameName)
	if e != nil {
		return
	}

	// Check if collection exists
	exist, e := CollectionIsExist(gameName, collectionName)
	if e != nil {
		return
	}
	if !exist {
		return
	}

	// Update info file
	e = CollectionAddInfo(gameName, collectionName, req.CollectionInfo)
	if e != nil {
		return
	}

	// Rename folder if name changed
	if req.Name != collectionName {
		e = CollectionRename(gameName, collectionName, req.Name)
	}
	return
}
