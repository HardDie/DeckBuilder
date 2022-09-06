package generator_old

import (
	"log"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/generator_old/internal/crawl"
	"tts_deck_build/internal/generator_old/internal/deck_builder"
	"tts_deck_build/internal/generator_old/internal/helpers"
)

// Read configurations, generate TTS json object with description
func GenerateDeckObject(gameID string) {
	gamePath := filepath.Join(config.GetConfig().Games(), gameID)
	// Read all decks
	listOfDecks, err := crawl.Crawl(gamePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Build
	db := deckbuilder.NewDeckBuilder()
	for _, decks := range listOfDecks {
		for _, deck := range decks {
			helpers.PutDeckToDeckBuilder(deck, db)
		}
	}

	// Generate TTS object
	res := db.GenerateTTSDeck()

	// Write deck json to file
	err = os.WriteFile(filepath.Join(config.GetConfig().Results(), "deck.json"), res, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
