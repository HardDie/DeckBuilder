package crawl

import (
	"path/filepath"

	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/generator_old/internal/types"
)

// Separate decks by type
// Top level map[string] - split by types (ex.: Loot, Monster)
// Next level []*Deck - split by collection (ex.: Base, DLC)
func Crawl(gamePath string) (map[string][]*types.Deck, error) {
	result := make(map[string][]*types.Deck)
	// Get all collections
	collections, err := fs.ListOfFolders(gamePath)
	if err != nil {
		return nil, err
	}
	for _, col := range collections {
		collectionPath := filepath.Join(gamePath, col)
		// Get all decks for current collection
		decks, err := fs.ListOfFiles(collectionPath)
		if err != nil {
			return nil, err
		}
		for _, deck := range decks {
			deckPath := filepath.Join(collectionPath, deck)

			// Parse json
			deckObj, err := fs.OpenAndProcess(deckPath, fs.JsonFromReader[types.Deck])
			if err != nil {
				return nil, err
			}

			// Add deck to result objects separated by types
			if len(deckObj.Cards) > 0 {
				result[deckObj.Deck.Type] = append(result[deckObj.Deck.Type], deckObj)
			}

			// Set info for each card
			for _, card := range deckObj.Cards {
				card.FillWithInfo(deckObj.Version, deckObj.Collection, deckObj.Deck.Type)
			}
		}
	}
	return result, nil
}
