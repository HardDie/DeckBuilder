package entity

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
)

type DeckInfo struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	CachedImage string     `json:"cachedImage,omitempty"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

func (i *DeckInfo) Path(gameID, collectionID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, collectionID, i.ID+".json")
}
func (i *DeckInfo) ImagePath(gameID, collectionID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, collectionID, i.ID+".bin")
}
func (i *DeckInfo) GetName() string {
	return strings.ToLower(i.Name)
}
func (i *DeckInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}
func (i *DeckInfo) FillCachedImage(cfg *config.Config, gameID, collectionID string) {
	i.CachedImage = fmt.Sprintf(cfg.DeckImagePath, gameID, collectionID, i.ID)
}
