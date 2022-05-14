package generator

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/generator/internal/crawl"
	"tts_deck_build/internal/generator/internal/deck_builder"
	"tts_deck_build/internal/generator/internal/download_manager"
	"tts_deck_build/internal/generator/internal/helpers"
)

// Read configurations, download images, build deck image files
func GenerateDeckImages(gameID string) error {
	gamePath := filepath.Join(config.GetConfig().Games(), gameID)
	// Read all decks
	listOfDecks := crawl.Crawl(gamePath)

	dm := downloadmanager.NewDownloadManager(config.GetConfig().Caches())
	// Fill download list
	for _, decks := range listOfDecks {
		for _, deck := range decks {
			helpers.PutDeckToDownloadManager(deck, dm)
		}
	}
	// Download all images
	err := fs.CreateFolderIfNotExist(config.GetConfig().Caches())
	if err != nil {
		return err
	}
	dm.Download()

	// Build
	db := deckbuilder.NewDeckBuilder()
	for _, decks := range listOfDecks {
		for _, deck := range decks {
			helpers.PutDeckToDeckBuilder(deck, db)
		}
	}

	// Generate images
	err = fs.CreateFolderIfNotExist(config.GetConfig().Results())
	if err != nil {
		return err
	}
	images := db.DrawDecks()

	// Write all created files
	data, err := json.MarshalIndent(images, "", "	")
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(config.GetConfig().Results(), "images.json"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Read configurations, generate TTS json object with description
func GenerateDeckObject(gameID string) {
	gamePath := filepath.Join(config.GetConfig().Games(), gameID)
	// Read all decks
	listOfDecks := crawl.Crawl(gamePath)

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
	err := os.WriteFile(filepath.Join(config.GetConfig().Results(), "deck.json"), res, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
