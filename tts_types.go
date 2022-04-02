package main

// Description of deck
type TTSDeckDescription struct {
	FaceURL      string `json:"FaceURL"`
	BackURL      string `json:"BackURL"`
	NumWidth     int    `json:"NumWidth"`
	NumHeight    int    `json:"NumHeight"`
	BackIsHidden bool   `json:"BackIsHidden"`
	UniqueBack   bool   `json:"UniqueBack"`
	Type         int    `json:"Type"`
}

// Description of card inside deck
type TTSCard struct {
	Name        string       `json:"Name"`
	Nickname    *string      `json:"Nickname"`
	Description *string      `json:"Description"`
	CardID      int          `json:"CardID"`
	LuaScript   string       `json:"LuaScript"`
	Transform   TTSTransform `json:"Transform"`
}

type TTSTransform struct {
	PosX   float64 `json:"posX"`
	PosY   float64 `json:"posY"`
	PosZ   float64 `json:"posZ"`
	ScaleX float64 `json:"scaleX"`
	ScaleY float64 `json:"scaleY"`
	ScaleZ float64 `json:"scaleZ"`
}

// One deck of cards as object
type TTSDeckObject struct {
	Name             string                     `json:"Name"`
	Transform        TTSTransform               `json:"Transform"`
	Nickname         string                     `json:"Nickname"`
	Description      string                     `json:"Description"`
	DeckIDs          []int                      `json:"DeckIDs"`
	CustomDeck       map[int]TTSDeckDescription `json:"CustomDeck"`
	ContainedObjects []TTSCard                  `json:"ContainedObjects"`
}

// Json object with decks
type TTSSaveObject struct {
	ObjectStates []TTSDeckObject `json:"ObjectStates"`
}

var (
	deckVel    = 4
	deckOffset = -deckVel
)

func NewTTSDeckObject(nick, desc string) TTSDeckObject {
	deckOffset += deckVel
	return TTSDeckObject{
		Name:        "Deck",
		Nickname:    nick,
		Description: desc,
		Transform: TTSTransform{
			PosX:   float64(deckOffset),
			PosY:   10,
			PosZ:   0,
			ScaleX: 1.42055011,
			ScaleY: 1,
			ScaleZ: 1.42055011,
		},
		CustomDeck: make(map[int]TTSDeckDescription),
	}
}
