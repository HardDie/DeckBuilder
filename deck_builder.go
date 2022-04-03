package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type deckBuilderDeck struct {
	cards []*Card

	deckType     string
	backSideName string
	backSideURL  string
}
type DeckBuilder struct {
	decks map[string]*deckBuilderDeck
}

func NewDeckBuilder() *DeckBuilder {
	return &DeckBuilder{
		decks: make(map[string]*deckBuilderDeck),
	}
}

// collect
func (b *DeckBuilder) AddCard(deck *Deck, card *Card) {
	if _, ok := b.decks[deck.GetType()]; !ok {
		b.decks[deck.GetType()] = &deckBuilderDeck{
			deckType:     deck.GetType(),
			backSideName: deck.GetBackSideName(),
			backSideURL:  deck.GetBackSideURL(),
		}
	}
	b.decks[deck.GetType()].cards = append(b.decks[deck.GetType()].cards, card)
}
func (b *DeckBuilder) splitCards(deckType string) (cards [][]*Card) {
	for leftBorder := 0; leftBorder < len(b.decks[deckType].cards); leftBorder += MaxCardsOnPage {
		// Calculate right border for current deck
		rightBorder := min(len(b.decks[deckType].cards), leftBorder+MaxCardsOnPage)
		cards = append(cards, b.decks[deckType].cards[leftBorder:rightBorder])
	}
	return
}
func (b *DeckBuilder) getImageSize(count int) (cols, rows int) {
	cols = 10
	rows = 7
	images := cols * rows
	for r := 2; r <= 7; r++ {
		for c := 2; c <= 10; c++ {
			possible := c * r
			if possible < images && possible >= count {
				images = possible
				cols = c
				rows = r
			}
		}
	}
	return
}
func (b *DeckBuilder) GetDecks(deckType string) (decks []*Deck) {
	for index, cards := range b.splitCards(deckType) {
		// Calculate optimal count of columns and rows for result image
		columns, rows := b.getImageSize(len(cards) + 1)
		decks = append(decks, &Deck{
			Cards:   cards,
			Columns: columns,
			Rows:    rows,
			FileName: fmt.Sprintf("%s_%d_%d_%dx%d.png", cleanTitle(b.decks[deckType].deckType), index+1, len(cards),
				columns, rows),
			BackSide: &b.decks[deckType].backSideURL,
		})
	}
	return
}
func (b *DeckBuilder) GetTypes() (types []string) {
	for deckType := range b.decks {
		types = append(types, deckType)
	}
	return
}

// draw
func (b *DeckBuilder) DrawDecks() map[string]string {
	// List of result files
	res := make(map[string]string)
	for _, deckType := range b.GetTypes() {
		decks := b.GetDecks(deckType)
		for _, deck := range decks {
			NewDeckDrawer(deck).Draw()
			// Add current deck title
			res[deck.FileName] = ""
		}
		// Add back side image title
		res[decks[0].GetBackSideName()] = ""
	}
	return res
}

// tts
func (b *DeckBuilder) generateTTSDeck(replaces map[string]string, deckType string) []TTSDeckObject {
	var res []TTSDeckObject

	var obj TTSDeckObject

	var lastCollection string
	var lastDeck int

	for i, deck := range b.GetDecks(deckType) {
		for j, card := range deck.Cards {
			if lastCollection != card.Collection {
				if lastCollection == "" {
					obj = NewTTSDeckObject(deck.Type, card.Collection)
					lastCollection = card.Collection
				} else {
					lastCollection = card.Collection
					res = append(res, obj)
					obj = NewTTSDeckObject(deck.Type, card.Collection)
				}
				face, ok := replaces[deck.FileName]
				if !ok {
					log.Fatalf("Can't find URL for image: %s", deck.FileName)
				}
				back, ok := replaces[deck.GetBackSideName()]
				if !ok {
					log.Fatalf("Can't find URL for image: %s", deck.GetBackSideName())
				}
				obj.CustomDeck[i+1] = TTSDeckDescription{
					FaceURL:    face,
					BackURL:    back,
					NumWidth:   deck.Columns,
					NumHeight:  deck.Rows,
					UniqueBack: false,
					Type:       0,
				}
				lastDeck = i
			}

			if lastDeck != i {
				lastDeck = i
				face, ok := replaces[deck.FileName]
				if !ok {
					log.Fatalf("Can't find URL for image: %s", deck.FileName)
				}
				back, ok := replaces[deck.GetBackSideName()]
				if !ok {
					log.Fatalf("Can't find URL for image: %s", deck.GetBackSideName())
				}
				obj.CustomDeck[i+1] = TTSDeckDescription{
					FaceURL:    face,
					BackURL:    back,
					NumWidth:   deck.Columns,
					NumHeight:  deck.Rows,
					UniqueBack: false,
					Type:       0,
				}
			}

			cardId := (i+1)*100 + j
			obj.DeckIDs = append(obj.DeckIDs, cardId)
			obj.ContainedObjects = append(obj.ContainedObjects, TTSCard{
				Name:        "Card",
				Nickname:    card.Title,
				Description: new(string),
				CardID:      cardId,
				LuaScript:   card.GetLua(),
				Transform:   obj.Transform,
			})
		}
	}

	if len(obj.ContainedObjects) > 0 {
		res = append(res, obj)
	}
	return res
}
func (b *DeckBuilder) GenerateTTSDeck(replaces map[string]string) []byte {
	res := TTSSaveObject{}
	for _, deckType := range b.GetTypes() {
		res.ObjectStates = append(res.ObjectStates, b.generateTTSDeck(replaces, deckType)...)
	}
	data, _ := json.MarshalIndent(res, "", "  ")
	return data
}
