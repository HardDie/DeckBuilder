package helpers

import (
	db "tts_deck_build/internal/deck_builder"
	dm "tts_deck_build/internal/download_manager"
	"tts_deck_build/internal/types"
)

func PutDeckToDownloadManager(d *types.Deck, dm *dm.DownloadManager) {
	dm.AddFile(d.GetBackSideURL(), d.GetBackSideName())
	for _, card := range d.GetCards() {
		dm.AddFile(card.GetFrontSideURL(), card.GetFrontSideName())
		if card.IsUniqueBack() {
			dm.AddFile(card.GetUniqueBackSineURL(), card.GetUniqueBackSideName())
		}
	}
}

func PutDeckToDeckBuilder(d *types.Deck, db *db.DeckBuilder) {
	for _, card := range d.GetCards() {
		db.AddCard(d, card)
	}
}
