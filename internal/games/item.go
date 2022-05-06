package games

func ItemGame(name string) (result *GameInfo, e error) {
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
