package collections

import (
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

type CollectionInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func NewCollectionInfo(name, desc, image string) *CollectionInfo {
	return &CollectionInfo{
		Id:          utils.NameToId(name),
		Name:        name,
		Description: desc,
		Image:       image,
	}
}

func (i *CollectionInfo) Path(gameId string) string {
	return filepath.Join(config.GetConfig().Games(), gameId, i.Id)
}

func (i *CollectionInfo) InfoPath(gameId string) string {
	return filepath.Join(config.GetConfig().Games(), gameId, i.Id, config.GetConfig().InfoFilename)
}

func (i *CollectionInfo) ImagePath(gameId string) string {
	return filepath.Join(config.GetConfig().Games(), gameId, i.Id, config.GetConfig().ImageFilename)
}

func (i *CollectionInfo) Compare(val *CollectionInfo) bool {
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
