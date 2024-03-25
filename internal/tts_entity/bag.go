package tts_entity

type Bag struct {
	Name             string    `json:"Name"`
	Transform        Transform `json:"Transform"`
	Nickname         string    `json:"Nickname"`
	Description      string    `json:"Description"`
	ContainedObjects []any     `json:"ContainedObjects"`
}

func NewBag(nickname string) Bag {
	return Bag{
		Name:      "Bag",
		Nickname:  nickname,
		Transform: transform,
	}
}

func (b Bag) GetName() string {
	return b.Name
}
func (b Bag) GetNickname() string {
	return b.Nickname
}
