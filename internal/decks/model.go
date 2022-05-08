package decks

type DeckInfo struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	BacksideImage string `json:"backside"`
}

func NewDeckInfo(id, deckType, image string) *DeckInfo {
	return &DeckInfo{
		Id:            id,
		Type:          deckType,
		BacksideImage: image,
	}
}

func (i *DeckInfo) Compare(val *DeckInfo) bool {
	if i.Id != val.Id {
		return false
	}
	if i.Type != val.Type {
		return false
	}
	if i.BacksideImage != val.BacksideImage {
		return false
	}
	return true
}
