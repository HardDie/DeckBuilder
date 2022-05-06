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

func NewGameInfo(id, name, desc, image string) *GameInfo {
	return &GameInfo{
		Id: id,
		GameInfoWithoutId: GameInfoWithoutId{
			Name:        name,
			Description: desc,
			Image:       image,
		},
	}
}
