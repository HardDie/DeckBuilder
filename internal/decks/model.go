package decks

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

type DeckInfo struct {
	Id            string     `json:"id"`
	Type          string     `json:"type"`
	BacksideImage string     `json:"backside"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}

func NewDeckInfo(deckType, image string) *DeckInfo {
	return &DeckInfo{
		Id:            utils.NameToId(deckType),
		Type:          deckType,
		BacksideImage: image,
		CreatedAt:     utils.Allocate(time.Now()),
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

func (i *DeckInfo) GetName() string {
	return i.Type
}

func (i *DeckInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}

func Sort(items *[]*DeckInfo, field string) {
	field = strings.ToLower(field)
	sort.SliceStable(*items, func(i, j int) bool {
		switch field {
		default: // "name"
			return (*items)[i].GetName() < (*items)[j].GetName()
		case "name_desc":
			return (*items)[i].GetName() > (*items)[j].GetName()
		case "created":
			return (*items)[i].GetCreatedAt().Before((*items)[j].GetCreatedAt())
		case "created_desc":
			return (*items)[i].GetCreatedAt().After((*items)[j].GetCreatedAt())
		}
	})
}
