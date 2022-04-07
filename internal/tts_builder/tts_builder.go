package tts_builder

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/types"
)

type TTSBuilder struct {
	replaces map[string]string
	objects  map[string]*types.TTSDeckObject

	resObjects []interface{}
}

func NewTTSBuilder() *TTSBuilder {
	data, err := ioutil.ReadFile(filepath.Join(config.GetConfig().ResultDir, "images.json"))
	if err != nil {
		log.Fatal(err.Error())
	}
	replaces := make(map[string]string)
	err = json.Unmarshal(data, &replaces)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &TTSBuilder{
		objects:  make(map[string]*types.TTSDeckObject),
		replaces: replaces,
	}
}

func (b *TTSBuilder) generateTTSDeckDescription(deck *types.Deck) types.TTSDeckDescription {
	face, ok := b.replaces[deck.FileName]
	if !ok {
		log.Fatalf("Can't find URL for image: %s", deck.FileName)
	}
	back, ok := b.replaces[deck.GetBackSideName()]
	if !ok {
		log.Fatalf("Can't find URL for image: %s", deck.GetBackSideName())
	}
	return types.TTSDeckDescription{
		FaceURL:    face,
		BackURL:    back,
		NumWidth:   deck.Columns,
		NumHeight:  deck.Rows,
		UniqueBack: false,
		Type:       0,
	}
}
func (b *TTSBuilder) generateTTSCard(card *types.Card, cardId int, transform types.TTSTransform) types.TTSCard {
	return types.TTSCard{
		Name:        "Card",
		Nickname:    card.Title,
		Description: new(string),
		CardID:      cardId,
		LuaScript:   card.GetLua(),
		Transform:   transform,
	}
}
func (b *TTSBuilder) AddCard(deck *types.Deck, card *types.Card, deckId, cardId int) {
	// Get deck object
	ttsDeck, ok := b.objects[card.Collection]
	if !ok {
		ttsDeck = types.NewTTSDeckObject(deck.Type, card.Collection)
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
