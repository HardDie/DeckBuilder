package tts_entity

type DeckDescription struct {
	FaceURL      string `json:"FaceURL"`
	BackURL      string `json:"BackURL"`
	NumWidth     int    `json:"NumWidth"`
	NumHeight    int    `json:"NumHeight"`
	BackIsHidden bool   `json:"BackIsHidden"`
	UniqueBack   bool   `json:"UniqueBack"`
	Type         int    `json:"Type"`
}

type DeckObject struct {
	Name             string                  `json:"Name"`
	Transform        Transform               `json:"Transform"`
	Nickname         string                  `json:"Nickname"`
	Description      string                  `json:"Description"`
	DeckIDs          []int                   `json:"DeckIDs"`
	CustomDeck       map[int]DeckDescription `json:"CustomDeck"`
	ContainedObjects []Card                  `json:"ContainedObjects"`
}
