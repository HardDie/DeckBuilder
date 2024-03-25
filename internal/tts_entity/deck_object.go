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

func NewDeck(nickname string) DeckObject {
	return DeckObject{
		Name:       "Deck",
		Nickname:   nickname,
		CustomDeck: make(map[int]DeckDescription),
		Transform:  transform,
	}
}

func (d *DeckObject) AddCard(card Card) {
	// Place the card ID in the list of cards inside the deck object
	d.DeckIDs = append(d.DeckIDs, card.CardID)
	// Place card in the list of cards inside the deck
	d.ContainedObjects = append(d.ContainedObjects, card)
}

func (d DeckObject) GetName() string {
	return d.Name
}
func (d DeckObject) GetNickname() string {
	return d.Nickname
}
