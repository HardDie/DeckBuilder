package collections

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

type CollectionInfo struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

func NewCollectionInfo(name, desc, image string) *CollectionInfo {
	return &CollectionInfo{
		Id:          utils.NameToId(name),
		Name:        name,
		Description: desc,
		Image:       image,
		CreatedAt:   utils.Allocate(time.Now()),
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

func (i *CollectionInfo) GetName() string {
	return i.Name
}

func (i *CollectionInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}

func Sort(items *[]*CollectionInfo, field string) {
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
