package tts_entity

type DeckObject struct {
	Name             string                  `json:"Name"`
	Transform        Transform               `json:"Transform"`
	Nickname         string                  `json:"Nickname"`
	Description      string                  `json:"Description"`
	DeckIDs          []int                   `json:"DeckIDs"`
	CustomDeck       map[int]DeckDescription `json:"CustomDeck"`
	ContainedObjects []Card                  `json:"ContainedObjects"`
}
