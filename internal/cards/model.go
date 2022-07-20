package cards

import (
	"path/filepath"
	"strconv"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/utils"
)

type Card struct {
	Deck  interface{}         `json:"deck"`
	Cards map[int64]*CardInfo `json:"cards"`
}

type CardInfo struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Variables   map[string]string `json:"variables"`
	CreatedAt   *time.Time        `json:"createdAt"`
	UpdatedAt   *time.Time        `json:"updatedAt"`
}

func NewCardInfo(title, desc, image string, variables map[string]string) *CardInfo {
	card := &CardInfo{
		ID:          0,
		Title:       strconv.Quote(title),
		Description: strconv.Quote(desc),
		Image:       image,
		Variables:   make(map[string]string),
		CreatedAt:   utils.Allocate(time.Now()),
	}
	if len(variables) > 0 {
		card.Variables = variables
	}
	return card
}

func (i *CardInfo) ImagePath(gameID, collectionID, deckID string) string {
	filename := deckID + "_" + fs.Int64ToString(i.ID) + ".bin"
	return filepath.Join(config.GetConfig().Games(), gameID, collectionID, filename)
}

func (i *CardInfo) Compare(val *CardInfo) bool {
	if i.Title != val.Title {
		return false
	}
	if i.Description != val.Description {
		return false
	}
	if i.Image != val.Image {
		return false
	}
	if len(i.Variables) != len(val.Variables) {
		return false
	}
	for key, value := range i.Variables {
		value2, ok := val.Variables[key]
		if !ok {
			return false
		}
		if value != value2 {
			return false
		}
	}
	return true
}

func (i *CardInfo) GetName() string {
	return i.Title
}

func (i *CardInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}
