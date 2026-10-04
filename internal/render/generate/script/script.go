// Build the TTS Saved Object from pages that are already measured.
package script

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/render/generate/catalog"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

// Page is one measured sheet. Paths are absolute.
type Page struct {
	DeckID string
	Index  int
	Image  string
	Back   string
	Cols   int
	Rows   int
}

// Document is the Saved Object file plus the inner bag sent to TTS.
type Document struct {
	Path string
	Root tts_entity.RootObjects
	Bag  tts_entity.Bag
}

// Build walks the cards with a counter so page index and slot match the sheets.
func Build(
	dir string,
	gameItem *entitiesGame.Game,
	decks map[catalog.Deck][]catalog.Card,
	order []catalog.Deck,
	pages []Page,
	cfg *entitiesSettings.Settings,
	cards servicesCard.Card,
) (Document, error) {
	byKey := make(map[string]Page, len(pages))
	for _, page := range pages {
		byKey[page.DeckID+"_"+strconv.Itoa(page.Index)] = page
	}

	bag := tts_entity.NewBag(gameItem.Name)
	collectionBags := make(map[string]*tts_entity.Bag)

	var deckIDOffset int
	var commonIndex int
	for _, deckInfo := range order {
		cardsInDeck := decks[deckInfo]
		commonIndex++

		deck := newDeck(deckInfo.Name, cfg)
		index := 1
		filled := 0
		pageInfo := byKey[deckInfo.ID+"_"+strconv.Itoa(index)]
		deckDescription := description(pageInfo)
		deck.CustomDeck[index+deckIDOffset] = deckDescription

		var prevCollection string
		var prevCollectionDeck string

		for _, card := range cardsInDeck {
			if filled == 0 {
				prevCollection = card.CollectionID
				prevCollectionDeck = card.CollectionID + deckInfo.ID
			}
			if filled >= config.MaxCount {
				index++
				commonIndex++
				filled = 0
				pageInfo = byKey[deckInfo.ID+"_"+strconv.Itoa(index)]
				deckDescription = description(pageInfo)
				deck.CustomDeck[index+deckIDOffset] = deckDescription
			}
			if card.CollectionID+deckInfo.ID != prevCollectionDeck {
				prevCollectionDeck = card.CollectionID + deckInfo.ID
				putDeck(collectionBags, prevCollection, deck)
				prevCollection = card.CollectionID
				deck = newDeck(deckInfo.Name, cfg)
				deck.CustomDeck[index+deckIDOffset] = deckDescription
			}

			filled++
			cardItem, err := cards.Item(card.GameID, card.CollectionID, deckInfo.ID, card.ID)
			if err != nil {
				return Document{}, err
			}

			cardGUID := fmt.Sprintf("%06d", commonIndex)
			commonIndex++
			cardObject := tts_entity.NewCard(
				cardGUID,
				cardItem.Name,
				cardItem.Description,
				index+deckIDOffset,
				filled-1,
				cardItem.Variables,
				deckDescription,
				cardTransform(cfg),
			)
			for i := 0; i < cardItem.Count; i++ {
				deck.AddCard(cardObject)
			}
		}

		if filled != 0 {
			putDeck(collectionBags, prevCollection, deck)
		}
		deckIDOffset += index
	}

	// Collection bags in collection-id order, so the same catalog gives the same JSON.
	for _, collection := range slices.Sorted(maps.Keys(collectionBags)) {
		bag.ContainedObjects = append(bag.ContainedObjects, collectionBags[collection])
	}
	bag.Description = fmt.Sprintf("Created at: %v", time.Now().Format("2006-01-02 15:04:05"))
	return Document{
		Path: filepath.Join(dir, gameItem.ID+".json"),
		Root: tts_entity.RootObjects{ObjectStates: []tts_entity.Bag{bag}},
		Bag:  bag,
	}, nil
}

func newDeck(name string, cfg *entitiesSettings.Settings) tts_entity.DeckObject {
	return tts_entity.NewDeck(name, cardTransform(cfg))
}

// cardTransform scales width (X) and length (Z) like TTS does; thickness (Y) stays 1.
func cardTransform(cfg *entitiesSettings.Settings) tts_entity.Transform {
	return tts_entity.Transform{
		ScaleX: cfg.CardScale,
		ScaleY: 1,
		ScaleZ: cfg.CardScale,
	}
}

func description(info Page) tts_entity.DeckDescription {
	return tts_entity.DeckDescription{
		FaceURL:   "file:///" + info.Image,
		BackURL:   "file:///" + info.Back,
		NumWidth:  info.Cols,
		NumHeight: info.Rows,
	}
}

func putDeck(bags map[string]*tts_entity.Bag, collection string, deck tts_entity.DeckObject) {
	if _, ok := bags[collection]; !ok {
		bags[collection] = utils.Allocate(tts_entity.NewBag(collection))
	}
	switch {
	case len(deck.ContainedObjects) == 1:
		bags[collection].ContainedObjects = append(bags[collection].ContainedObjects, deck.ContainedObjects[0])
	case len(deck.ContainedObjects) > 1:
		bags[collection].ContainedObjects = append(bags[collection].ContainedObjects, deck)
	}
}
