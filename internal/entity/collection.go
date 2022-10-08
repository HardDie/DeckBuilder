package entity

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type CollectionInfo struct {
	ID          string             `json:"id"`
	Name        utils.QuotedString `json:"name"`
	Description utils.QuotedString `json:"description"`
	Image       string             `json:"image"`
	CachedImage string             `json:"cachedImage,omitempty"`
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

func (i *CollectionInfo) Path(gameID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, i.ID)
}
func (i *CollectionInfo) InfoPath(gameID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, i.ID, cfg.InfoFilename)
}
func (i *CollectionInfo) ImagePath(gameID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, i.ID, cfg.ImageFilename)
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
	return strings.ToLower(i.Name.String())
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
func (i *CollectionInfo) FillCachedImage(cfg *config.Config, gameID string) {
	i.CachedImage = fmt.Sprintf(cfg.CollectionImagePath, gameID, i.ID)
}
