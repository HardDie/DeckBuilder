package entity

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type GameInfo struct {
	ID          string             `json:"id"`
	Name        utils.QuotedString `json:"name"`
	Description utils.QuotedString `json:"description"`
	Image       string             `json:"image"`
	CachedImage string             `json:"cachedImage,omitempty"`
	CreatedAt   *time.Time         `json:"createdAt"`
	UpdatedAt   *time.Time         `json:"updatedAt"`
}

func NewGameInfo(name, desc, image string) *GameInfo {
	return &GameInfo{
		ID:          utils.NameToID(name),
		Name:        utils.NewQuotedString(name),
		Description: utils.NewQuotedString(desc),
		Image:       image,
		CreatedAt:   utils.Allocate(time.Now()),
	}
}

func (i *GameInfo) Path(cfg *config.Config) string {
	return filepath.Join(cfg.Games(), i.ID)
}
func (i *GameInfo) InfoPath(cfg *config.Config) string {
	return filepath.Join(cfg.Games(), i.ID, cfg.InfoFilename)
}
func (i *GameInfo) ImagePath(cfg *config.Config) string {
	return filepath.Join(cfg.Games(), i.ID, cfg.ImageFilename)
}
func (i *GameInfo) Compare(val *GameInfo) bool {
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
func (i *GameInfo) GetName() string {
	return strings.ToLower(i.Name.String())
}
func (i *GameInfo) GetCreatedAt() time.Time {
	if i.CreatedAt != nil {
		return *i.CreatedAt
	}
	return time.Time{}
}
func (i *GameInfo) SetQuotedOutput() {
	i.Name.SetQuotedOutput()
	i.Description.SetQuotedOutput()
}
func (i *GameInfo) SetRawOutput() {
	i.Name.SetRawOutput()
	i.Description.SetRawOutput()
}
func (i *GameInfo) FillCachedImage(cfg *config.Config) {
	i.CachedImage = fmt.Sprintf(cfg.GameImagePath, i.ID)
}
