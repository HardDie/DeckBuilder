package tts_entity

type Card struct {
	Name        string                  `json:"Name"`
	Nickname    *string                 `json:"Nickname"`
	Description *string                 `json:"Description"`
	CardID      int                     `json:"CardID"`
	LuaScript   string                  `json:"LuaScript"`
	Transform   *Transform              `json:"Transform,omitempty"`
	CustomDeck  map[int]DeckDescription `json:"CustomDeck,omitempty"`
}
