package collections

import (
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

type CreateCollectionRequest struct {
	CollectionInfoWithoutId
}

func CreateCollection(gameName string, req *CreateCollectionRequest) (res CollectionInfo, e *errors.Error) {
	// Check if collection already exists
	exist, e := CollectionIsExist(gameName, req.Name)
	if e != nil {
		return
	}
	if exist {
		e = errors.CollectionExist
		return
	}

	res = CollectionInfo{
		Id:                      utils.NameToId(req.CollectionInfoWithoutId.Name),
		CollectionInfoWithoutId: req.CollectionInfoWithoutId,
	}

	// Check if collection id correct
	if len(res.Id) == 0 {
		e = errors.BadName
		return
	}

	// Try to create folder with game
	e = CollectionCreate(gameName, res.Id)
	if e != nil {
		return
	}

	// Create info file
	e = CollectionAddInfo(gameName, res.Id, res)
	return
}
