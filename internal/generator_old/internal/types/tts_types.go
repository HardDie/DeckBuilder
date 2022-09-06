package types

import (
	"tts_deck_build/internal/tts_entity"
)

var (
	deckVel    = 4
	deckOffset = -deckVel
)

func NewTTSDeckObject(nick, desc string) *tts_entity.DeckObject {
	deckOffset += deckVel
	return &tts_entity.DeckObject{
		Name:        "Deck",
		Nickname:    nick,
		Description: desc,
		Transform: tts_entity.Transform{
			PosX:   float64(deckOffset),
			PosY:   10,
			PosZ:   0,
			ScaleX: 1.42055011,
			ScaleY: 1,
			ScaleZ: 1.42055011,
		},
		CustomDeck: make(map[int]tts_entity.DeckDescription),
	}
}
