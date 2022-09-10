package tts_entity

type Bag struct {
	Name             string    `json:"Name"`
	Transform        Transform `json:"Transform"`
	Nickname         string    `json:"Nickname"`
	Description      string    `json:"Description"`
	ContainedObjects []any     `json:"ContainedObjects"`
}
