package catalog

import (
	"fmt"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type WriteRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageFile   []byte `json:"imageFile"`
}

func GameDTO(cfg config.Config, item entitiesGame.Game) dto.Game {
	return dto.Game{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: fmt.Sprintf(cfg.GameImagePath+"?%s", item.ID, utils.HashForTime(&item.UpdatedAt)),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func CollectionDTO(cfg config.Config, gameID string, item entitiesCollection.Collection) dto.Collection {
	return dto.Collection{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: fmt.Sprintf(cfg.CollectionImagePath+"?%s", gameID, item.ID, utils.HashForTime(&item.UpdatedAt)),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
