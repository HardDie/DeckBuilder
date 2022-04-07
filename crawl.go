package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Parse json file to deck
func parseJson(path string) *Deck {
	desc := &Deck{}

	// Open file
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Decode
	dec := json.NewDecoder(f)
	if err = dec.Decode(desc); err != nil {
		log.Fatal(err.Error())
	}

	return desc
}
func cleanTitle(in string) string {
	res := strings.ReplaceAll(in, " / ", "_")
	res = strings.ReplaceAll(res, "/", "_")
	res = strings.ReplaceAll(res, "!", "")
	res = strings.ReplaceAll(res, "'", "")
	res = strings.ReplaceAll(res, ".", "")
	return strings.ReplaceAll(res, " ", "_")
}
func getFilenameFromUrl(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, filename := path.Split(u.Path)
	return filename
}

// Check every folder and get cards information
func crawl(path string) (result []*Deck) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		newPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			// Skip 'files' folder
			if file.Name() == "files" {
				continue
			}
			result = append(result, crawl(newPath)...)
			continue
		}
		log.Println("Parse file:", newPath)

		// Parse only json files
		if filepath.Ext(newPath) != ".json" {
			continue
		}
		deck := parseJson(newPath)

		result = append(result, deck)
		// Set for each card
		for _, card := range deck.Cards {
			card.FillWithInfo(deck.Version, deck.Collection, deck.Type)
		}
	}
	return
}

// Separate decks by type
// Top level map[string] - split by types (ex.: Loot, Monster)
// Next level []*Deck - split by collection (ex.: Base, DLC)
func Crawl(path string) map[string][]*Deck {
	result := make(map[string][]*Deck)
	// Get all decks
	decks := crawl(path)
	// Split decks by type
	for _, deck := range decks {
		result[deck.Type] = append(result[deck.Type], deck)
	}
	return result
}
