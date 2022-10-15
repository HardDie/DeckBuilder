package entity

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
)

type GameInfo struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	CachedImage string     `json:"cachedImage,omitempty"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

func (i *GameInfo) Path(cfg *config.Config) string {
	return filepath.Join(cfg.Games(), i.ID)
}
func (i *GameInfo) GetName() string {
	return strings.ToLower(i.Name)
}
func (i *GameInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}
func (i *GameInfo) FillCachedImage(cfg *config.Config) {
	i.CachedImage = fmt.Sprintf(cfg.GameImagePath, i.ID)
}
