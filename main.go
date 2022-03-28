package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

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
		collection[deckType] = deckCol
	}

	// Generate images
	var wc WholeCollection
	for _, deckCol := range collection {
		BuildDeck(deckCol)
		wc = append(wc, deckCol)
	}

	// Generate TTS object
	res := wc.GenerateTTSDeck()

	// Write deck json to file
	err := ioutil.WriteFile(GetConfig().ResultDir+"/deck.json", res, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Script for replacing links
	script := `
#!/bin/bash
declare -A arr
arr+=(
`
	for _, val := range allReplaces {
		script += fmt.Sprintf("[\"%s\"]=\"\"\n", val)
	}
	script += `)

if [[ ! -f deck_backup.json ]]; then
	cp deck.json deck_backup.json
fi
for key in ${!arr[@]}; do
	value=$(echo ${arr[$key]} | sed 's/\//\\\//g')
	sed -i "s/$key/$value/g" deck.json
done
`
	err = ioutil.WriteFile("./result_png/replace.sh", []byte(script), 0744)
	if err != nil {
		log.Fatal(err)
	}
}
