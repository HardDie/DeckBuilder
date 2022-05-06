package games

import (
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

type UpdateGameRequest struct {
	GameInfoWithoutId
}

func UpdateGame(name string, req *UpdateGameRequest) (res GameInfo, e error) {
	res = GameInfo{
		Id:                utils.NameToId(req.GameInfoWithoutId.Name),
		GameInfoWithoutId: req.GameInfoWithoutId,
	}

	// Check if game id correct
	if len(res.Id) == 0 {
		e = errors.BadName
		return
	}

	// Check if game exists
	exist, e := GameIsExist(name)
	if e != nil {
		return
	}
	if !exist {
		e = errors.GameNotExists
		return
	}

	// Rename folder if name changed
	if res.Id != name {
		e = GameRename(name, res.Id)
		if e != nil {
			return
		}
	}

	// Update info file
	e = GameAddInfo(res.Id, res)
	return
}
