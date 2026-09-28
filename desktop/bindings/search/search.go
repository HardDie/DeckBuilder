package search

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesSearch "github.com/HardDie/DeckBuilder/internal/services/search"
)

type Result struct {
	Data dto.RecursiveSearch `json:"data"`
	Meta *network.Meta       `json:"meta"`
}

type Search struct {
	svc servicesSearch.Search
}

func New(svc servicesSearch.Search) *Search {
	return &Search{svc: svc}
}

func (s *Search) Root(sort, query string) (*Result, error) {
	return s.search(sort, query, "", "")
}

func (s *Search) Game(gameID, sort, query string) (*Result, error) {
	return s.search(sort, query, gameID, "")
}

func (s *Search) Collection(gameID, collectionID, sort, query string) (*Result, error) {
	return s.search(sort, query, gameID, collectionID)
}

func (s *Search) search(sort, query, gameID, collectionID string) (*Result, error) {
	resp, err := s.svc.RecursiveSearch(sort, query, gameID, collectionID)
	if err != nil {
		return nil, err
	}

	response := dto.RecursiveSearch{}
	for _, game := range resp.Games {
		response.Games = append(response.Games, game.ID)
	}
	for _, collection := range resp.Collections {
		response.Collections = append(response.Collections, dto.RecursiveSearchCollection{
			GameID:       collection.GameID,
			CollectionID: collection.ID,
		})
	}
	for _, deck := range resp.Decks {
		response.Decks = append(response.Decks, dto.RecursiveSearchDeck{
			GameID:       deck.GameID,
			CollectionID: deck.CollectionID,
			DeckID:       deck.ID,
		})
	}
	for _, card := range resp.Cards {
		response.Cards = append(response.Cards, dto.RecursiveSearchCard{
			GameID:       card.GameID,
			CollectionID: card.CollectionID,
			DeckID:       card.DeckID,
			CardID:       card.ID,
		})
	}

	return &Result{
		Data: response,
		Meta: &network.Meta{
			Total: len(resp.Games) + len(resp.Collections) + len(resp.Decks) + len(resp.Cards),
		},
	}, nil
}
