package search

import (
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type Search interface {
	RecursiveSearch(sortField, search, gameID, collectionID string) (*entity.RecursiveSearchItems, *network.Meta, error)
}
