package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
)

type TTSBuilder struct {
	replaces map[string]string
	objects  map[string]*TTSDeckObject

	resObjects []interface{}
}

func NewTTSBuilder() *TTSBuilder {
	data, err := ioutil.ReadFile(filepath.Join(GetConfig().ResultDir, "images.json"))
	if err != nil {
		log.Fatal(err.Error())
	}
	replaces := make(map[string]string)
	err = json.Unmarshal(data, &replaces)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &TTSBuilder{
		objects:  make(map[string]*TTSDeckObject),
		replaces: replaces,
	}
}

func (b *TTSBuilder) generateTTSDeckDescription(deck *Deck) TTSDeckDescription {
	face, ok := b.replaces[deck.FileName]
	if !ok {
		log.Fatalf("Can't find URL for image: %s", deck.FileName)
	}
	back, ok := b.replaces[deck.GetBackSideName()]
	if !ok {
		log.Fatalf("Can't find URL for image: %s", deck.GetBackSideName())
	}
	return TTSDeckDescription{
		FaceURL:    face,
		BackURL:    back,
		NumWidth:   deck.Columns,
		NumHeight:  deck.Rows,
		UniqueBack: false,
		Type:       0,
	}
}
func (b *TTSBuilder) generateTTSCard(card *Card, cardId int, transform TTSTransform) TTSCard {
	return TTSCard{
		Name:        "Card",
		Nickname:    card.Title,
		Description: new(string),
		CardID:      cardId,
		LuaScript:   card.GetLua(),
		Transform:   transform,
	}
}
func (b *TTSBuilder) AddCard(deck *Deck, card *Card, deckId, cardId int) {
	// Get deck object
	ttsDeck, ok := b.objects[card.Collection]
	if !ok {
		ttsDeck = NewTTSDeckObject(deck.Type, card.Collection)
		b.resObjects = append(b.resObjects, ttsDeck)
		b.objects[card.Collection] = ttsDeck
	}

	// Check if deck exists in list
	if _, ok = ttsDeck.CustomDeck[deckId]; !ok {
		ttsDeck.CustomDeck[deckId] = b.generateTTSDeckDescription(deck)
	}

	// Add card id to deck
	ttsDeck.DeckIDs = append(ttsDeck.DeckIDs, cardId)
	// Add card object to deck
	ttsDeck.ContainedObjects = append(ttsDeck.ContainedObjects, b.generateTTSCard(card, cardId, ttsDeck.Transform))
}
func (b *TTSBuilder) GetObjects() (result []interface{}) {
	// Sort keys
	var keys []string
	for key := range b.objects {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		object := b.objects[key]
		// If one card in deck, add separated card as object
		if len(object.ContainedObjects) == 1 {
			object.ContainedObjects[0].CustomDeck = object.CustomDeck
			result = append(result, object.ContainedObjects[0])
			continue
		}
		// Add deck as object
		result = append(result, object)
	}
	return
}
