package decks

import (
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

type DeckInfo struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	BacksideImage string `json:"backside"`
}

func NewDeckInfo(deckType, image string) *DeckInfo {
	return &DeckInfo{
		Id:            utils.NameToId(deckType),
		Type:          deckType,
		BacksideImage: image,
	}
}

func (i *DeckInfo) Path(gameId, collectionId string) string {
	return filepath.Join(config.GetConfig().Games(), gameId, collectionId, i.Id+".json")
}

func (i *DeckInfo) ImagePath(gameId, collectionId string) string {
	return filepath.Join(config.GetConfig().Games(), gameId, collectionId, i.Id+".bin")
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
