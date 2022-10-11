package entity

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
)

type CollectionInfo struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	CachedImage string     `json:"cachedImage,omitempty"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

func (i *CollectionInfo) Path(gameID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, i.ID)
}
func (i *CollectionInfo) ImagePath(gameID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, i.ID, cfg.ImageFilename)
}
func (i *CollectionInfo) GetName() string {
	return strings.ToLower(i.Name)
}
func (i *CollectionInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}
func (i *CollectionInfo) FillCachedImage(cfg *config.Config, gameID string) {
	i.CachedImage = fmt.Sprintf(cfg.CollectionImagePath, gameID, i.ID)
}
