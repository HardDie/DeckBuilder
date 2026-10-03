package generate

import (
	"fmt"

	"github.com/HardDie/DeckBuilder/internal/render/generate/catalog"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
)

// MissingImages lists, in render order, the decks without a back image
// and the cards without a face. Rendering needs both, so the caller
// should not start while the list is not empty.
// No image is read: decks and cards report HasImage from a folder listing.
func MissingImages(decks map[catalog.Deck][]catalog.Card, order []catalog.Deck, deckSvc servicesDeck.Deck) ([]string, error) {
	var missing []string
	for _, deck := range order {
		cards := decks[deck]
		if len(cards) == 0 {
			continue
		}
		// Prepare takes the back from the deck in the first card's collection.
		item, err := deckSvc.Item(cards[0].GameID, cards[0].CollectionID, deck.ID)
		if err != nil {
			return nil, err
		}
		if !item.HasImage {
			missing = append(missing, fmt.Sprintf("deck %q (back image)", deck.Name))
		}
		for _, card := range cards {
			if !card.HasImage {
				missing = append(missing, fmt.Sprintf("card %q in deck %q", card.Name, deck.Name))
			}
		}
	}
	return missing, nil
}
