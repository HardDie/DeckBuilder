package decks

import "tts_deck_build/internal/errors"

func DeleteDeck(gameName, collectionName, deckName string) (e error) {
	// Check if deck exists
	exist, e := DeckIsExist(gameName, collectionName, deckName)
	if e != nil {
		return
	}
	if !exist {
		e = errors.DeckNotExists
		return
	}

	// Try to delete deck
	return DeckDelete(gameName, collectionName, deckName)
}
