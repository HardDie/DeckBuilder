package games

import (
	"tts_deck_build/internal/errors"
)

func ItemGame(name string) (result *GameInfo, e *errors.Error) {
	// Check if game and game info exists
	e = FullGameCheck(name)
	if e != nil {
		return
	}

	// Get info
	result, e = GameGetInfo(name)
	if e != nil {
		return
	}
	return
}
