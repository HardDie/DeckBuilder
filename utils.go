package main

func PutDeckToDownloadManager(d *Deck, dm *DownloadManager) {
	dm.AddFile(d.GetBackSideURL(), d.GetBackSideName())
	for _, card := range d.GetCards() {
		dm.AddFile(card.GetFrontSideURL(), card.GetFrontSideName())
		if card.IsUniqueBack() {
			dm.AddFile(card.GetUniqueBackSineURL(), card.GetUniqueBackSideName())
		}
	}
}
