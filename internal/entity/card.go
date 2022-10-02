package entity

import (
	"path/filepath"
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
	ID          int64              `json:"id"`
	Name        utils.QuotedString `json:"name"`
	Description utils.QuotedString `json:"description"`
	Image       string             `json:"image"`
	CachedImage string             `json:"cachedImage"`
	Variables   map[string]string  `json:"variables"`
	Count       int                `json:"count"`
	CreatedAt   *time.Time         `json:"createdAt"`
	UpdatedAt   *time.Time         `json:"updatedAt"`
}

func NewCardInfo(name, desc, image string, variables map[string]string, count int) *CardInfo {
	if variables == nil {
		variables = make(map[string]string)
	}
	if count < 1 {
		count = 1
	}
	return &CardInfo{
		ID:          0,
		Name:        utils.NewQuotedString(name),
		Description: utils.NewQuotedString(desc),
		Image:       image,
		Variables:   variables,
		Count:       count,
		CreatedAt:   utils.Allocate(time.Now()),
	}
}

func (i *CardInfo) ImagePath(gameID, collectionID, deckID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, collectionID, deckID, fs.Int64ToString(i.ID)+".bin")
}
func (i *CardInfo) Compare(val *CardInfo) bool {
	if i.Name != val.Name {
		return false
	}
	if i.Description != val.Description {
		return false
	}
	if i.Image != val.Image {
		return false
	}
	if i.Count != val.Count {
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
	return i.Name.String()
}
func (i *CardInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}
func (i *CardInfo) SetQuotedOutput() {
	i.Name.SetQuotedOutput()
	i.Description.SetQuotedOutput()
}
func (i *CardInfo) SetRawOutput() {
	i.Name.SetRawOutput()
	i.Description.SetRawOutput()
}
