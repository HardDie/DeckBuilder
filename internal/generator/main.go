package generator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/generator/internal/crawl"
	deckBuilder "tts_deck_build/internal/generator/internal/deck_builder"
	downloadManager "tts_deck_build/internal/generator/internal/download_manager"

	// downloadManager "tts_deck_build/internal/download_manager"
	"tts_deck_build/internal/generator/internal/helpers"
	// "tts_deck_build/internal/helpers"
)

// Read configurations, download images, build deck image files
func GenerateDeckImages() {
	// Read all decks
	listOfDecks := crawl.Crawl(config.GetConfig().SourceDir)

	dm := downloadManager.NewDownloadManager(config.GetConfig().CachePath)
	// Fill download list
	for _, decks := range listOfDecks {
		for _, deck := range decks {
			helpers.PutDeckToDownloadManager(deck, dm)
		}
	}
	// Download all images
	dm.Download()

	// Build
	db := deckBuilder.NewDeckBuilder()
	for _, decks := range listOfDecks {
		for _, deck := range decks {
			helpers.PutDeckToDeckBuilder(deck, db)
		}
	}

	// Generate images
	images := db.DrawDecks()

	// Write all created files
	data, _ := json.MarshalIndent(images, "", "	")
	err := ioutil.WriteFile(filepath.Join(config.GetConfig().ResultDir, "images.json"), data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Read configurations, generate TTS json object with description
func GenerateDeckObject() {
	// Read all decks
	listOfDecks := crawl.Crawl(config.GetConfig().SourceDir)

	// Build
	db := deckBuilder.NewDeckBuilder()
	for _, decks := range listOfDecks {
		for _, deck := range decks {
			helpers.PutDeckToDeckBuilder(deck, db)
		}
	}

	// Generate TTS object
	res := db.GenerateTTSDeck()

	// Write deck json to file
	err := ioutil.WriteFile(filepath.Join(config.GetConfig().ResultDir, "deck.json"), res, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
