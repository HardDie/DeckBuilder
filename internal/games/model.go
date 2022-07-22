package games

import (
	"path/filepath"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

type GameInfo struct {
	ID          string             `json:"id"`
	Name        utils.QuotedString `json:"name"`
	Description utils.QuotedString `json:"description"`
	Image       string             `json:"image"`
	CreatedAt   *time.Time         `json:"createdAt"`
	UpdatedAt   *time.Time         `json:"updatedAt"`
}

func NewGameInfo(name, desc, image string) *GameInfo {
	return &GameInfo{
		ID:          utils.NameToID(name),
		Name:        utils.NewQuotedString(name),
		Description: utils.NewQuotedString(desc),
		Image:       image,
		CreatedAt:   utils.Allocate(time.Now()),
	}
}

func (i *GameInfo) Path() string {
	return filepath.Join(config.GetConfig().Games(), i.ID)
}

func (i *GameInfo) InfoPath() string {
	return filepath.Join(config.GetConfig().Games(), i.ID, config.GetConfig().InfoFilename)
}

func (i *GameInfo) ImagePath() string {
	return filepath.Join(config.GetConfig().Games(), i.ID, config.GetConfig().ImageFilename)
}

func (i *GameInfo) Compare(val *GameInfo) bool {
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

func (i *GameInfo) GetName() string {
	return i.Name.String()
}

func (i *GameInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}

func (i *GameInfo) SetQuotedOutput() {
	i.Name.SetQuotedOutput()
	i.Description.SetQuotedOutput()
}

func (i *GameInfo) SetRawOutput() {
	i.Name.SetRawOutput()
	i.Description.SetRawOutput()
}
