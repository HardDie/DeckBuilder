package games

import (
	"tts_deck_build/internal/errors"
)

type UpdateGameRequest struct {
	GameInfo
}

func UpdateGame(name string, req *UpdateGameRequest) (e *errors.Error) {
	// Validate
	if len(req.Name) == 0 {
		e = errors.DataInvalid.AddMessage("The name of the game cannot be empty")
		return
	}

	// Check if game exists
	exist, e := GameIsExist(name)
	if e != nil {
		return
	}
	if !exist {
		return
	}

	// Update info file
	e = GameAddInfo(name, req.GameInfo)
	if e != nil {
		return
	}

	// Rename folder if name changed
	if req.Name != name {
		e = GameRename(name, req.Name)
	}
	return
}
