package decks

import (
	"path/filepath"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

type Deck struct {
	Deck *DeckInfo `json:"deck"`
}

type DeckInfo struct {
	ID            string     `json:"id"`
	Type          string     `json:"type"`
	BacksideImage string     `json:"backside"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}

func NewDeckInfo(deckType, image string) *DeckInfo {
	return &DeckInfo{
		ID:            utils.NameToID(deckType),
		Type:          deckType,
		BacksideImage: image,
		CreatedAt:     utils.Allocate(time.Now()),
	}
}

func (i *DeckInfo) Path(gameID, collectionID string) string {
	return filepath.Join(config.GetConfig().Games(), gameID, collectionID, i.ID+".json")
}

func (i *DeckInfo) ImagePath(gameID, collectionID string) string {
	return filepath.Join(config.GetConfig().Games(), gameID, collectionID, i.ID+".bin")
}

func (i *DeckInfo) Compare(val *DeckInfo) bool {
	if i.ID != val.ID {
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

func (i *DeckInfo) GetName() string {
	return i.Type
}

func (i *DeckInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}
