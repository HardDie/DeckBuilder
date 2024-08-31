package search

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesSearch "github.com/HardDie/DeckBuilder/internal/services/search"
)

type search struct {
	serviceSearch servicesSearch.Search
}

func New(serviceSearch servicesSearch.Search) Search {
	return &search{
		serviceSearch: serviceSearch,
	}
}

func (s *search) RootHandler(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")

	resp, err := s.serviceSearch.RecursiveSearch(sort, search, "", "")
	if err != nil {
		network.ResponseError(w, err)
		return
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

	network.ResponseWithMeta(w, response, &network.Meta{
		Total: len(resp.Games) + len(resp.Collections) + len(resp.Decks) + len(resp.Cards),
	})
}
func (s *search) GameHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")

	resp, err := s.serviceSearch.RecursiveSearch(sort, search, gameID, "")
	if err != nil {
		network.ResponseError(w, err)
		return
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

	network.ResponseWithMeta(w, response, &network.Meta{
		Total: len(resp.Games) + len(resp.Collections) + len(resp.Decks) + len(resp.Cards),
	})
}
func (s *search) CollectionHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")

	resp, err := s.serviceSearch.RecursiveSearch(sort, search, gameID, collectionID)
	if err != nil {
		network.ResponseError(w, err)
		return
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

	network.ResponseWithMeta(w, response, &network.Meta{
		Total: len(resp.Games) + len(resp.Collections) + len(resp.Decks) + len(resp.Cards),
	})
}
