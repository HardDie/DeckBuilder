package entity

import (
	"fmt"
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
