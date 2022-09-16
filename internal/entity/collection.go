package entity

import (
	"path/filepath"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

type CollectionInfo struct {
	ID          string             `json:"id"`
	Name        utils.QuotedString `json:"name"`
	Description utils.QuotedString `json:"description"`
	Image       string             `json:"image"`
	CreatedAt   *time.Time         `json:"createdAt"`
	UpdatedAt   *time.Time         `json:"updatedAt"`
}

func NewCollectionInfo(name, desc, image string) *CollectionInfo {
	return &CollectionInfo{
		ID:          utils.NameToID(name),
		Name:        utils.NewQuotedString(name),
		Description: utils.NewQuotedString(desc),
		Image:       image,
		CreatedAt:   utils.Allocate(time.Now()),
	}
}

func (i *CollectionInfo) Path(gameID string) string {
	return filepath.Join(config.GetConfig().Games(), gameID, i.ID)
}

func (i *CollectionInfo) InfoPath(gameID string) string {
	return filepath.Join(config.GetConfig().Games(), gameID, i.ID, config.GetConfig().InfoFilename)
}

func (i *CollectionInfo) ImagePath(gameID string) string {
	return filepath.Join(config.GetConfig().Games(), gameID, i.ID, config.GetConfig().ImageFilename)
}

func (i *CollectionInfo) Compare(val *CollectionInfo) bool {
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

func (i *CollectionInfo) GetName() string {
	return i.Name.String()
}

func (i *CollectionInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}

func (i *CollectionInfo) SetQuotedOutput() {
	i.Name.SetQuotedOutput()
	i.Description.SetQuotedOutput()
}

func (i *CollectionInfo) SetRawOutput() {
	i.Name.SetRawOutput()
	i.Description.SetRawOutput()
}
