package games

import (
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

type GameInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func NewGameInfo(name, desc, image string) *GameInfo {
	return &GameInfo{
		Id:          utils.NameToId(name),
		Name:        name,
		Description: desc,
		Image:       image,
	}
}

func (i *GameInfo) Path() string {
	return filepath.Join(config.GetConfig().Games(), i.Id)
}

func (i *GameInfo) InfoPath() string {
	return filepath.Join(config.GetConfig().Games(), i.Id, config.GetConfig().InfoFilename)
}

func (i *GameInfo) ImagePath() string {
	return filepath.Join(config.GetConfig().Games(), i.Id, config.GetConfig().ImageFilename)
}

func (i *GameInfo) Compare(val *GameInfo) bool {
	if i.Id != val.Id {
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
