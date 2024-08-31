package search

import (
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
)

type search struct {
	serviceGame       servicesGame.Game
	serviceCollection servicesCollection.Collection
	serviceDeck       servicesDeck.Deck
	serviceCard       servicesCard.Card
}

func New(
	serviceGame servicesGame.Game,
	serviceCollection servicesCollection.Collection,
	serviceDeck servicesDeck.Deck,
	serviceCard servicesCard.Card,
) Search {
	return &search{
		serviceGame:       serviceGame,
		serviceCollection: serviceCollection,
		serviceDeck:       serviceDeck,
		serviceCard:       serviceCard,
	}
}

func (s *search) RecursiveSearch(sortField, search, gameID, collectionID string) (*entity.RecursiveSearchItems, *network.Meta, error) {
	var err error
	res := &entity.RecursiveSearchItems{}
	meta := &network.Meta{}

	if gameID == "" {
		// Check if game with such search mask exist
		foundGames, err := s.serviceGame.List(sortField, search)
		if err != nil {
			return nil, nil, err
		}
		for _, game := range foundGames {
			res.Games = append(res.Games, game.ID)
			meta.Total += 1
		}
	}

	// Iterate through all games
	var allGames []*entitiesGame.Game
	if gameID == "" {
		allGames, err = s.serviceGame.List(sortField, "")
		if err != nil {
			return nil, nil, err
		}
	} else {
		game, err := s.serviceGame.Item(gameID)
		if err != nil {
			return nil, nil, err
		}
		allGames = append(allGames, game)
	}
	for _, game := range allGames {
		if collectionID == "" {
			// Check if collection with such search mask exist
			foundCollections, err := s.serviceCollection.List(game.ID, sortField, search)
			if err != nil {
				return nil, nil, err
			}
			for _, collection := range foundCollections {
				res.Collections = append(res.Collections, entity.RecursiveCollectionItem{
					GameID:       game.ID,
					CollectionID: collection.ID,
				})
				meta.Total += 1
			}
		}

		// Iterate through all collections

		var allCollections []*entitiesCollection.Collection
		if collectionID == "" {
			allCollections, err = s.serviceCollection.List(game.ID, sortField, "")
			if err != nil {
				return nil, nil, err
			}
		} else {
			collection, err := s.serviceCollection.Item(game.ID, collectionID)
			if err != nil {
				return nil, nil, err
			}
			allCollections = append(allCollections, collection)
		}
		for _, collection := range allCollections {
			// Check if deck with such search mask exist
			foundDecks, err := s.serviceDeck.List(game.ID, collection.ID, sortField, search)
			if err != nil {
				return nil, nil, err
			}
			for _, deck := range foundDecks {
				res.Decks = append(res.Decks, entity.RecursiveDeckItem{
					GameID:       game.ID,
					CollectionID: collection.ID,
					DeckID:       deck.ID,
				})
				meta.Total += 1
			}

			// Iterate through all decks
			allDecks, err := s.serviceDeck.List(game.ID, collection.ID, sortField, "")
			if err != nil {
				return nil, nil, err
			}
			for _, deck := range allDecks {
				// Check if card with such search mask exist
				foundCards, _, err := s.serviceCard.List(game.ID, collection.ID, deck.ID, sortField, search)
				if err != nil {
					return nil, nil, err
				}
				for _, card := range foundCards {
					res.Cards = append(res.Cards, entity.RecursiveCardItem{
						GameID:       game.ID,
						CollectionID: collection.ID,
						DeckID:       deck.ID,
						CardID:       card.ID,
					})
					meta.Total += 1
				}
			}
		}
	}
	return res, meta, nil
}
