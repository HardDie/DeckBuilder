package main

import "golang.org/x/exp/constraints"

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func PutDeckToDownloadManager(d *Deck, dm *DownloadManager) {
	dm.AddFile(d.GetBackSideURL(), d.GetBackSideName())
	for _, card := range d.GetCards() {
		dm.AddFile(card.GetFrontSideURL(), card.GetFrontSideName())
		if card.IsUniqueBack() {
			dm.AddFile(card.GetUniqueBackSineURL(), card.GetUniqueBackSideName())
		}
	}
}

func PutDeckToDeckBuilder(d *Deck, db *DeckBuilder) {
	for _, card := range d.GetCards() {
		db.AddCard(card)
	}
}
