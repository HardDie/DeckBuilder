package collections

import (
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/utils"
)

type UpdateCollectionRequest struct {
	CollectionInfoWithoutId
}

func UpdateCollection(gameName, collectionName string, req *UpdateCollectionRequest) (res CollectionInfo, e *errors.Error) {
	res = CollectionInfo{
		Id:                      utils.NameToId(req.CollectionInfoWithoutId.Name),
		CollectionInfoWithoutId: req.CollectionInfoWithoutId,
	}

	// Check if collection id correct
	if len(res.Id) == 0 {
		e = errors.BadName
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
		e = errors.CollectionNotExists
		return
	}

	// Rename folder if name changed
	if res.Id != collectionName {
		e = CollectionRename(gameName, collectionName, res.Id)
		if e != nil {
			return
		}
	}

	// Update info file
	e = CollectionAddInfo(gameName, res.Id, res)
	return
}
