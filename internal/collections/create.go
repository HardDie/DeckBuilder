package collections

import "tts_deck_build/internal/errors"

type CreateCollectionRequest struct {
	CollectionInfo
}

func CreateCollection(gameName string, req *CreateCollectionRequest) (e *errors.Error) {
	// Check if collection already exists
	exist, e := CollectionIsExist(gameName, req.Name)
	if e != nil {
		return
	}
	if exist {
		e = errors.GameExist
		return
	}

	// Try to create folder with game
	e = CollectionCreate(gameName, req.Name)
	if e != nil {
		return
	}

	// Create info file
	e = CollectionAddInfo(gameName, req.Name, req.CollectionInfo)
	return
}
