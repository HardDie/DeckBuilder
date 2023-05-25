package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/utils"
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
	i.CachedImage = fmt.Sprintf(cfg.DeckImagePath+"?%s", gameID, collectionID, i.ID, utils.HashForTime(i.UpdatedAt))
}
