package helpers

import (
	"tts_deck_build/internal/generator_old/internal/deck_builder"
	"tts_deck_build/internal/generator_old/internal/types"
)

func PutDeckToDeckBuilder(d *types.Deck, db *deckbuilder.DeckBuilder) {
	for _, card := range d.GetCards() {
		db.AddCard(d, card)
	}
}
