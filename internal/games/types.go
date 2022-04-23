package games

type GameInfoWithoutId struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type GameInfo struct {
	Id string `json:"id"`
	GameInfoWithoutId
}
