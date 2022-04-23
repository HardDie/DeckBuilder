package decks

type DeckInfoWithoutId struct {
	Type          string `json:"type"`
	BacksideImage string `json:"backside"`
}

type DeckInfo struct {
	Id string `json:"id"`
	DeckInfoWithoutId
}
