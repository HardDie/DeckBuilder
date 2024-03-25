package tts_entity

import (
	"fmt"
	"strings"
)

type Card struct {
	GUID        string                  `json:"GUID"`
	Name        string                  `json:"Name"`
	Nickname    string                  `json:"Nickname"`
	Description string                  `json:"Description"`
	CardID      int                     `json:"CardID"`
	LuaScript   string                  `json:"LuaScript"`
	Transform   *Transform              `json:"Transform,omitempty"`
	CustomDeck  map[int]DeckDescription `json:"CustomDeck,omitempty"`
	States      map[string]Card         `json:"States,omitempty"`
}

func NewCard(guid, name, description string, pageId, cardIndex int, variablesMap map[string]string, deckDesc DeckDescription) Card {
	// Converting lua variables into strings
	var variables []string
	for key, value := range variablesMap {
		variables = append(variables, fmt.Sprintf(`%s=%q`, key, value))
	}
	return Card{
		GUID:        guid,
		Name:        "Card",
		Nickname:    name,
		Description: description,
		CardID:      pageId*100 + cardIndex,
		LuaScript:   strings.Join(variables, "\n"),
		CustomDeck: map[int]DeckDescription{
			pageId: deckDesc,
		},
		Transform: &transform,
	}
}
