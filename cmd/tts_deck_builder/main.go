// Package main TTS_deck_builder
//
// Entry point for the application.
//
// Terms Of Service:
//
//     Schemes: http
//     Host: localhost:5000
//     BasePath: /
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - binary
//
// swagger:meta
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"tts_deck_build/api"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/crawl"
	deckBuilder "tts_deck_build/internal/deck_builder"
	downloadManager "tts_deck_build/internal/download_manager"
	"tts_deck_build/internal/helpers"
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

func WebServer() {
	log.Println("Listening on :5000...")

	http.Handle("/", api.GetRoutes())
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	return
}

func createDirIfNotExists(folder string) {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		err = os.Mkdir(folder, 0755)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}
	if err != nil {
		log.Fatal(err.Error())
	}
}
func setup() {
	createDirIfNotExists(config.GetConfig().Data)
	createDirIfNotExists(config.GetConfig().Games())
	createDirIfNotExists(config.GetConfig().CachePath)
	createDirIfNotExists(config.GetConfig().ResultDir)
}

func main() {
	// Setup logs
	log.SetFlags(log.Llongfile | log.Ltime)

	setup()

	// Setup run flags
	genImgMode := flag.Bool("generate_image", false, "Run process of generating deck images")
	genDeckMode := flag.Bool("generate_object", false, "Run process of generating json deck object")
	helpMode := flag.Bool("help", false, "Show help")
	flag.Parse()

	switch {
	case *genImgMode:
		GenerateDeckImages()
	case *genDeckMode:
		GenerateDeckObject()
	case *helpMode:
		fmt.Println("How to use:")
		fmt.Println("1. Build images from ${sourceDir}/*.json descriptions (-generate_image)")
		fmt.Println("2. Upload images on some hosting (steam cloud)")
		fmt.Println("3. Write URL for each image in ${resultDir}/images.json file")
		fmt.Println("4. Build deck object ${resultDir}/deck.json (-generate_object)")
		fmt.Println("5. Put deck object into \"Tabletop Simulator/Saves/Saved Objects\" folder")
		fmt.Println()
		fmt.Println("Choose one of the mode:")
		flag.PrintDefaults()
	default:
		WebServer()
	}
}
