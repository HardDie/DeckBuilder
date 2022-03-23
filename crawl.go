package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
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
func Crawl(path string, listOfDecks map[string][]*Deck) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		newPath := path + "/" + file.Name()
		if file.IsDir() {
			Crawl(newPath, listOfDecks)
			continue
		}
		log.Println("Parse file:", newPath)

		// Path to file convert to string
		tokens := strings.Split(newPath, "/")
		prefix := strings.Join(tokens[2:len(tokens)-1], "_")

		deck := parseJson(newPath)
		deck.Prefix = prefix
		listOfDecks[deck.Type] = append(listOfDecks[deck.Type], deck)
		// Set for each card
		for _, card := range deck.Cards {
			card.FillWithInfo(tokens[len(tokens)-2], deck.Prefix, deck.Type)
		}
	}
}
