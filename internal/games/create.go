package games

import (
	"tts_deck_build/internal/errors"
)

type CreateGameRequest struct {
	GameInfo
}

func CreateGame(req *CreateGameRequest) (e *errors.Error) {
	// Check if game already exists
	exist, e := GameIsExist(req.Name)
	if e != nil {
		return
	}
	if exist {
		e = errors.GameExist
		return
	}

	// Try to create folder with game
	e = GameCreate(req.Name)
	if e != nil {
		return
	}

	// Create info file
	e = GameAddInfo(req.Name, req.GameInfo)
	return
}
