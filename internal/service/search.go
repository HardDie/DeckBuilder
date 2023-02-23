package service

import (
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type ISearchService interface {
	RecursiveSearch(sortField, search, gameID, collectionID string) (*entity.RecursiveSearchItems, *network.Meta, error)
}
type SearchService struct {
	gameService       IGameService
	collectionService ICollectionService
	deckService       IDeckService
	cardService       ICardService
}

func NewSearchService(
	gameService IGameService,
	collectionService ICollectionService,
	deckService IDeckService,
	cardService ICardService,
) *SearchService {
	return &SearchService{
		gameService:       gameService,
		collectionService: collectionService,
		deckService:       deckService,
		cardService:       cardService,
	}
}

func (s *SearchService) RecursiveSearch(sortField, search, gameID, collectionID string) (*entity.RecursiveSearchItems, *network.Meta, error) {
	var err error
	res := &entity.RecursiveSearchItems{}
	meta := &network.Meta{}

	if gameID == "" {
		// Check if game with such search mask exist
		foundGames, _, err := s.gameService.List(sortField, search)
		if err != nil {
			return nil, nil, err
		}
		for _, game := range foundGames {
			res.Games = append(res.Games, game.ID)
			meta.Total += 1
		}
	}

	// Iterate through all games
	var allGames []*entity.GameInfo
	if gameID == "" {
		allGames, _, err = s.gameService.List(sortField, "")
		if err != nil {
			return nil, nil, err
		}
	} else {
		game, err := s.gameService.Item(gameID)
		if err != nil {
			return nil, nil, err
		}
		allGames = append(allGames, game)
	}
	for _, game := range allGames {
		if collectionID == "" {
			// Check if collection with such search mask exist
			foundCollections, _, err := s.collectionService.List(game.ID, sortField, search)
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

		var allCollections []*entity.CollectionInfo
		if collectionID == "" {
			allCollections, _, err = s.collectionService.List(game.ID, sortField, "")
			if err != nil {
				return nil, nil, err
			}
		} else {
			collection, err := s.collectionService.Item(game.ID, collectionID)
			if err != nil {
				return nil, nil, err
			}
			allCollections = append(allCollections, collection)
		}
		for _, collection := range allCollections {
			// Check if deck with such search mask exist
			foundDecks, _, err := s.deckService.List(game.ID, collection.ID, sortField, search)
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
			allDecks, _, err := s.deckService.List(game.ID, collection.ID, sortField, "")
			if err != nil {
				return nil, nil, err
			}
			for _, deck := range allDecks {
				// Check if card with such search mask exist
				foundCards, _, err := s.cardService.List(game.ID, collection.ID, deck.ID, sortField, search)
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
