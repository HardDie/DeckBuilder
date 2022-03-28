package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

// Read configurations, download images, build deck image files
func GenerateDeckImages() {
	// Read all decks
	listOfDecks := Crawl(GetConfig().SourceDir)

	// Get download list
	var pairs []DownloadInfo
	for _, decks := range listOfDecks {
		for _, deck := range decks {
			pairs = append(pairs, deck.GetDownloadList()...)
		}
	}

	// Download all images
	DownloadFiles(pairs)

	// Build
	collection := make(map[string]*DeckCollection)
	for deckType, decks := range listOfDecks {
		deckCol := NewDeckCollection()
		for _, deck := range decks {
			deckCol.MergeDeck(deck)
		}
		deckCol.FillAttributes()
		collection[deckType] = deckCol
	}

	// Generate images
	var wc WholeCollection
	for _, deckCol := range collection {
		BuildDeck(deckCol)
		wc = append(wc, deckCol)
	}

	// Write all created files
	data, _ := json.MarshalIndent(wc.GetResultImages(), "", "	")
	err := ioutil.WriteFile(GetConfig().ResultDir+"/images.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Read configurations, generate TTS json object with description
func GenerateDeckObject() {
	// Read all decks
	listOfDecks := Crawl(GetConfig().SourceDir)

	// Build
	collection := make(map[string]*DeckCollection)
	for deckType, decks := range listOfDecks {
		deckCol := NewDeckCollection()
		for _, deck := range decks {
			deckCol.MergeDeck(deck)
		}
		deckCol.FillAttributes()
		collection[deckType] = deckCol
	}

	// Generate images
	var wc WholeCollection
	for _, deckCol := range collection {
		wc = append(wc, deckCol)
	}

	data, err := ioutil.ReadFile(GetConfig().ResultDir + "/images.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	replaces := make(map[string]string)
	err = json.Unmarshal(data, &replaces)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Generate TTS object
	res := wc.GenerateTTSDeck(replaces)

	// Write deck json to file
	err = ioutil.WriteFile(GetConfig().ResultDir+"/deck.json", res, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Setup logs
	log.SetFlags(log.Lshortfile | log.Ltime)

	// Setup run flags
	genImgMode := flag.Bool("generate_image", false, "Run process of generating deck images")
	genDeckMode := flag.Bool("generate_object", false, "Run process of generating json deck object")
	flag.Parse()

	// One of the modes must be selected
	if *genImgMode == *genDeckMode {
		fmt.Println("How to use:")
		fmt.Println("1. Build images from ${sourceDir}/*.json descriptions (-generate_image)")
		fmt.Println("2. Upload images on some hosting (steam cloud)")
		fmt.Println("3. Write URL for each image in ${resultDir}/images.json file")
		fmt.Println("4. Build deck object ${resultDir}/deck.json (-generate_object)")
		fmt.Println("5. Put deck object into \"Tabletop Simulator/Saves/Saved Objects\" folder")
		fmt.Println()
		fmt.Println("Choose one of the mode:")
		flag.PrintDefaults()
		return
	}

	switch {
	case *genImgMode:
		GenerateDeckImages()
	case *genDeckMode:
		GenerateDeckObject()
	}
}
