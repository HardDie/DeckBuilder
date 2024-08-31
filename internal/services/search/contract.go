package search

import (
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
)

type Search interface {
	RecursiveSearch(sortField, search, gameID, collectionID string) (*RecursiveSearchResponse, error)
}

type RecursiveSearchResponse struct {
	Games       []*entitiesGame.Game
	Collections []*entitiesCollection.Collection
	Decks       []*entitiesDeck.Deck
	Cards       []*entitiesCard.Card
}
