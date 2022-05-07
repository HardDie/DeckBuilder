package games

type GameInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func NewGameInfo(id, name, desc, image string) *GameInfo {
	return &GameInfo{
		Id:          id,
		Name:        name,
		Description: desc,
		Image:       image,
	}
}

func (i *GameInfo) Compare(val *GameInfo) bool {
	if i.Id != val.Id {
		return false
	}
	if i.Name != val.Name {
		return false
	}
	if i.Description != val.Description {
		return false
	}
	if i.Image != val.Image {
		return false
	}
	return true
}
