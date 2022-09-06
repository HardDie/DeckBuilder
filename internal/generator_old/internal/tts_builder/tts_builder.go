package ttsbuilder

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/generator_old/internal/types"
	"tts_deck_build/internal/tts_entity"
)

type TTSBuilder struct {
	replaces map[string]string
	objects  map[string]*tts_entity.DeckObject

	resObjects []interface{}
}

func NewTTSBuilder() *TTSBuilder {
	data, err := os.ReadFile(filepath.Join(config.GetConfig().Results(), "images.json"))
	if err != nil {
		log.Fatal(err.Error())
	}
	replaces := make(map[string]string)
	err = json.Unmarshal(data, &replaces)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &TTSBuilder{
		objects:  make(map[string]*tts_entity.DeckObject),
		replaces: replaces,
	}
}

func (b *TTSBuilder) generateTTSDeckDescription(deck *types.Deck) tts_entity.DeckDescription {
	face, ok := b.replaces[deck.FileName]
	if !ok {
		log.Fatalf("Can't find URL for image: %s", deck.FileName)
	}
	back, ok := b.replaces[deck.GetBackSideName()]
	if !ok {
		log.Fatalf("Can't find URL for image: %s", deck.GetBackSideName())
	}
	return tts_entity.DeckDescription{
		FaceURL:    face,
		BackURL:    back,
		NumWidth:   deck.Columns,
		NumHeight:  deck.Rows,
		UniqueBack: false,
		Type:       0,
	}
}
func (b *TTSBuilder) generateTTSCard(card *types.Card, cardID int) tts_entity.Card {
	return tts_entity.Card{
		Name:        "Card",
		Nickname:    card.Title,
		Description: new(string),
		CardID:      cardID,
		LuaScript:   card.GetLua(),
	}
}
func (b *TTSBuilder) AddCard(deck *types.Deck, card *types.Card, deckID, cardID int) {
	// Get deck object
	ttsDeck, ok := b.objects[card.Collection]
	if !ok {
		ttsDeck = types.NewTTSDeckObject(deck.Deck.Type, card.Collection)
		b.resObjects = append(b.resObjects, ttsDeck)
		b.objects[card.Collection] = ttsDeck
	}

	// Check if deck exists in list
	if _, ok = ttsDeck.CustomDeck[deckID]; !ok {
		ttsDeck.CustomDeck[deckID] = b.generateTTSDeckDescription(deck)
	}

	// Add card id to deck
	ttsDeck.DeckIDs = append(ttsDeck.DeckIDs, cardID)
	// Add card object to deck
	ttsDeck.ContainedObjects = append(ttsDeck.ContainedObjects, b.generateTTSCard(card, cardID))
}
func (b *TTSBuilder) GetObjects() (result []interface{}) {
	// Sort keys
	keys := make([]string, 0)
	for key := range b.objects {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		object := b.objects[key]
		// If one card in deck, add separated card as object
		if len(object.ContainedObjects) == 1 {
			object.ContainedObjects[0].Transform = &object.Transform
			object.ContainedObjects[0].CustomDeck = object.CustomDeck
			result = append(result, object.ContainedObjects[0])
			continue
		}
		// Add deck as object
		result = append(result, object)
	}
	return
}