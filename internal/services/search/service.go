package search

import (
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
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

func (s *search) RecursiveSearch(sortField, search, gameID, collectionID string) (*RecursiveSearchResponse, error) {
	var err error
	res := &RecursiveSearchResponse{}

	if gameID == "" {
		// Check if game with such search mask exist
		foundGames, err := s.serviceGame.List(sortField, search)
		if err != nil {
			return nil, err
		}
		for _, game := range foundGames {
			res.Games = append(res.Games, game)
		}
	}

	// Iterate through all games
	var allGames []*entitiesGame.Game
	if gameID == "" {
		allGames, err = s.serviceGame.List(sortField, "")
		if err != nil {
			return nil, err
		}
	} else {
		game, err := s.serviceGame.Item(gameID)
		if err != nil {
			return nil, err
		}
		allGames = append(allGames, game)
	}
	for _, game := range allGames {
		if collectionID == "" {
			// Check if collection with such search mask exist
			foundCollections, err := s.serviceCollection.List(game.ID, sortField, search)
			if err != nil {
				return nil, err
			}
			for _, collection := range foundCollections {
				res.Collections = append(res.Collections, collection)
			}
		}

		// Iterate through all collections

		var allCollections []*entitiesCollection.Collection
		if collectionID == "" {
			allCollections, err = s.serviceCollection.List(game.ID, sortField, "")
			if err != nil {
				return nil, err
			}
		} else {
			collection, err := s.serviceCollection.Item(game.ID, collectionID)
			if err != nil {
				return nil, err
			}
			allCollections = append(allCollections, collection)
		}
		for _, collection := range allCollections {
			// Check if deck with such search mask exist
			foundDecks, err := s.serviceDeck.List(game.ID, collection.ID, sortField, search)
			if err != nil {
				return nil, err
			}
			for _, deck := range foundDecks {
				res.Decks = append(res.Decks, deck)
			}

			// Iterate through all decks
			allDecks, err := s.serviceDeck.List(game.ID, collection.ID, sortField, "")
			if err != nil {
				return nil, err
			}
			for _, deck := range allDecks {
				// Check if card with such search mask exist
				foundCards, err := s.serviceCard.List(game.ID, collection.ID, deck.ID, sortField, search)
				if err != nil {
					return nil, err
				}
				for _, card := range foundCards {
					res.Cards = append(res.Cards, card)
				}
			}
		}
	}
	return res, nil
}
