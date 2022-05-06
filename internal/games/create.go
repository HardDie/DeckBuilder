package games

import (
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

type CreateGameRequest struct {
	GameInfoWithoutId
}

func CreateGame(req *CreateGameRequest) (res GameInfo, e *errors.Error) {
	// Check if game already exists
	exist, e := GameIsExist(req.Name)
	if e != nil {
		return
	}
	if exist {
		e = errors.GameExist
		return
	}

	res = GameInfo{
		Id:                utils.NameToId(req.GameInfoWithoutId.Name),
		GameInfoWithoutId: req.GameInfoWithoutId,
	}

	// Check if game id correct
	if len(res.Id) == 0 {
		e = errors.BadName
		return
	}

	// Try to create folder with game
	e = GameCreate(res.Id)
	if e != nil {
		return
	}

	// Create info file
	e = GameAddInfo(res.Id, res)
	if e != nil {
		return
	}

	// Download image if set
	if len(res.Image) > 0 {
		e = CreateImage(res.Image, res.Id)
	}
	return
}
