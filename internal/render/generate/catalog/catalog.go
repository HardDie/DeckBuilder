// Walk a game into decks of cards, sorted by deck name.
package catalog

import (
	"sort"

	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
)

// Deck is one deck identity. The same id, name, and image in two collections share one key.
type Deck struct {
	ID    string
	Name  string
	Image string
}

// Card is one catalog card on that deck. Count is copied from the list.
type Card struct {
	ID           int64
	GameID       string
	CollectionID string
	Count        int
}

// Collect reads every collection, deck, and card.
// Decks with no cards are left out.
// The returned order is deck name, ascending, stable.
func Collect(
	gameID, sortField string,
	collections servicesCollection.Collection,
	decks servicesDeck.Deck,
	cards servicesCard.Card,
) (map[Deck][]Card, []Deck, error) {
	grouped := make(map[Deck][]Card)
	collectionItems, err := collections.List(gameID, sortField, "")
	if err != nil {
		return nil, nil, err
	}
	for _, collectionItem := range collectionItems {
		deckItems, err := decks.List(gameID, collectionItem.ID, sortField, "")
		if err != nil {
			return nil, nil, err
		}
		for _, deckItem := range deckItems {
			deck := Deck{
				ID:    deckItem.ID,
				Name:  deckItem.Name,
				Image: deckItem.Image,
			}
			cardItems, err := cards.List(gameID, collectionItem.ID, deckItem.ID, sortField, "")
			if err != nil {
				return nil, nil, err
			}
			for _, cardItem := range cardItems {
				grouped[deck] = append(grouped[deck], Card{
					ID:           cardItem.ID,
					GameID:       gameID,
					CollectionID: collectionItem.ID,
					Count:        cardItem.Count,
				})
			}
		}
	}

	order := make([]Deck, 0, len(grouped))
	for deck := range grouped {
		order = append(order, deck)
	}
	sort.SliceStable(order, func(i, j int) bool {
		return order[i].Name < order[j].Name
	})
	return grouped, order, nil
}
