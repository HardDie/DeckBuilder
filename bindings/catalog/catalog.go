package catalog

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type WriteRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	ImageFile   []byte            `json:"imageFile"`
	Count       int               `json:"count"`
	Variables   map[string]string `json:"variables"`
}

func (r WriteRequest) ImageBytes() []byte {
	if len(r.ImageFile) == 0 {
		return nil
	}
	return r.ImageFile
}

func GameDTO(cfg config.Config, item entitiesGame.Game) dto.Game {
	return dto.Game{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: fmt.Sprintf(cfg.GameImagePath+"?%s", item.ID, utils.HashForTime(&item.UpdatedAt)),
		CreatedAt:   formatTimestamp(item.CreatedAt),
		UpdatedAt:   formatTimestamp(item.UpdatedAt),
	}
}

func CollectionDTO(cfg config.Config, gameID string, item entitiesCollection.Collection) dto.Collection {
	return dto.Collection{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: fmt.Sprintf(cfg.CollectionImagePath+"?%s", gameID, item.ID, utils.HashForTime(&item.UpdatedAt)),
		CreatedAt:   formatTimestamp(item.CreatedAt),
		UpdatedAt:   formatTimestamp(item.UpdatedAt),
	}
}

func DeckDTO(cfg config.Config, item entitiesDeck.Deck) dto.Deck {
	return dto.Deck{
		ID:                item.ID,
		Name:              item.Name,
		Description:       item.Description,
		Image:             item.Image,
		CachedImage:       fmt.Sprintf(cfg.DeckImagePath+"?%s", item.GameID, item.CollectionID, item.ID, utils.HashForTime(&item.UpdatedAt)),
		HasImage:          item.HasImage,
		CardsMissingImage: item.CardsMissingImage,
		CreatedAt:         formatTimestamp(item.CreatedAt),
		UpdatedAt:         formatTimestamp(item.UpdatedAt),
	}
}

func CardDTO(cfg config.Config, item entitiesCard.Card) dto.Card {
	return dto.Card{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: fmt.Sprintf(cfg.CardImagePath+"?%s", item.GameID, item.CollectionID, item.DeckID, item.ID, utils.HashForTime(&item.UpdatedAt)),
		Variables:   item.Variables,
		Count:       item.Count,
		HasImage:    item.HasImage,
		CreatedAt:   formatTimestamp(item.CreatedAt),
		UpdatedAt:   formatTimestamp(item.UpdatedAt),
	}
}

// ImageWarning tells the user why a new image was not applied.
// It is "" when err is nil. An unexpected cause is logged and shown generically.
func ImageWarning(err error) string {
	if err == nil {
		return ""
	}
	msg := apperr.Message(err)
	if msg == apperr.Unexpected {
		slog.Error("image not saved", "err", err)
	}
	return "Image was not saved: " + msg
}

// Reminders for a deck or card saved without an image. Rendering needs both.
const (
	DeckNoImage = "This deck has no image. Rendering needs a back image for every deck."
	CardNoImage = "This card has no image. Rendering needs an image for every card."
)

// SaveWarning joins why a new image was not applied and, when the entity
// still has no image, the reminder noImage. It is "" when there is nothing to say.
func SaveWarning(imageErr error, hasImage bool, noImage string) string {
	var parts []string
	if w := ImageWarning(imageErr); w != "" {
		parts = append(parts, w)
	}
	if !hasImage {
		parts = append(parts, noImage)
	}
	return strings.Join(parts, " ")
}

// formatTimestamp matches encoding/json for time.Time.
func formatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}
