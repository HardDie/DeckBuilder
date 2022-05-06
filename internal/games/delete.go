package games

import (
	"tts_deck_build/internal/errors"
)

func DeleteGame(name string) (e error) {
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
	return GameDelete(name)
}
