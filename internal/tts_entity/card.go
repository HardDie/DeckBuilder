package tts_entity

import (
	"fmt"
	"maps"
	"slices"
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

func NewCard(
	guid, name, description string,
	pageId, cardIndex int,
	variablesMap map[string]string,
	deckDesc DeckDescription,
	cardSize Transform,
) Card {
	// One key="value" line per variable, sorted by key so the JSON is stable.
	variables := make([]string, 0, len(variablesMap))
	for _, key := range slices.Sorted(maps.Keys(variablesMap)) {
		variables = append(variables, fmt.Sprintf(`%s=%q`, key, variablesMap[key]))
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
		Transform: &Transform{
			ScaleX: cardSize.ScaleX,
			ScaleY: cardSize.ScaleY,
			ScaleZ: cardSize.ScaleZ,
		},
	}
}

func (c Card) GetName() string {
	return c.Name
}
func (c Card) GetNickname() string {
	return c.Nickname
}
