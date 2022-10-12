package entity

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/fs"
)

type CardInfo struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	CachedImage string            `json:"cachedImage,omitempty"`
	Variables   map[string]string `json:"variables"`
	Count       int               `json:"count"`
	CreatedAt   *time.Time        `json:"createdAt"`
	UpdatedAt   *time.Time        `json:"updatedAt"`
}

func (i *CardInfo) ImagePath(gameID, collectionID, deckID string, cfg *config.Config) string {
	return filepath.Join(cfg.Games(), gameID, collectionID, deckID, "cards", fs.Int64ToString(i.ID)+".bin")
}
func (i *CardInfo) GetName() string {
	return strings.ToLower(i.Name)
}
func (i *CardInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}
func (i *CardInfo) FillCachedImage(cfg *config.Config, gameID, collectionID, deckID string) {
	i.CachedImage = fmt.Sprintf(cfg.CardImagePath, gameID, collectionID, deckID, i.ID)
}
