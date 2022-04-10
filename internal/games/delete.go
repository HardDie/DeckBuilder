package games

import (
	"tts_deck_build/internal/errors"
)

func DeleteGame(name string) (e *errors.Error) {
	// Check if game exists
	exist, e := GameIsExist(name)
	if e != nil {
		return
	}
	if !exist {
		e = errors.GameNotExists
		return
	}

	// Try to delete game
	e = GameDelete(name)
	return
}
