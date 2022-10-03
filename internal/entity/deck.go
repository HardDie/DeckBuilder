package entity

import (
	"path/filepath"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

type Deck struct {
	Deck  *DeckInfo   `json:"deck"`
	Cards interface{} `json:"cards"`
}

type DeckInfo struct {
	ID          string             `json:"id"`
	Name        utils.QuotedString `json:"name"`
	Description utils.QuotedString `json:"description"`
	Image       string             `json:"image"`
	CachedImage string             `json:"cachedImage"`
	CreatedAt   *time.Time         `json:"createdAt"`
	UpdatedAt   *time.Time         `json:"updatedAt"`
}

func NewDeckInfo(name, desc, image string) *DeckInfo {
	return &DeckInfo{
		ID:          utils.NameToID(name),
		Name:        utils.NewQuotedString(name),
		Description: utils.NewQuotedString(desc),
		Image:       image,
		CreatedAt:   utils.Allocate(time.Now()),
	}
}

func (i *DeckInfo) Path(gameID, collectionID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, collectionID, i.ID+".json")
}
func (i *DeckInfo) CardImagesPath(gameID, collectionID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, collectionID, i.ID)
}
func (i *DeckInfo) ImagePath(gameID, collectionID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, collectionID, i.ID+".bin")
}
func (i *DeckInfo) Compare(val *DeckInfo) bool {
	if i.ID != val.ID {
		return false
	}
	if i.Name != val.Name {
		return false
	}
	if i.Description != val.Description {
		return false
	}
	if i.Image != val.Image {
		return false
	}
	return true
}
func (i *DeckInfo) GetName() string {
	return i.Name.String()
}
func (i *DeckInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}
func (i *DeckInfo) SetQuotedOutput() {
	i.Name.SetQuotedOutput()
	i.Description.SetQuotedOutput()
}
func (i *DeckInfo) SetRawOutput() {
	i.Name.SetRawOutput()
	i.Description.SetRawOutput()
}
