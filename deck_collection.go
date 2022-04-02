package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type DeckBuilder struct {
	cards []*Card

	deckType     string
	backSideName string
	backSideURL  string
}

func NewDeckBuilder(deck *Deck) *DeckBuilder {
	return &DeckBuilder{
		deckType:     deck.GetType(),
		backSideName: deck.GetBackSideName(),
		backSideURL:  deck.GetBackSideURL(),
	}
}
func (b *DeckBuilder) AddCard(card *Card) {
	b.cards = append(b.cards, card)
}
func (b *DeckBuilder) splitCards() (cards [][]*Card) {
	for leftBorder := 0; leftBorder < len(b.cards); leftBorder += MaxCardsOnPage {
		// Calculate right border for current deck
		rightBorder := min(len(b.cards), leftBorder+MaxCardsOnPage)
		cards = append(cards, b.cards[leftBorder:rightBorder])
	}
	return
}
func (b *DeckBuilder) GetDecks() (decks []*Deck) {
	for index, cards := range b.splitCards() {
		// Calculate optimal count of columns and rows for result image
		columns, rows := b.getImageSize(len(cards) + 1)
		decks = append(decks, &Deck{
			Cards:   cards,
			Columns: columns,
			Rows:    rows,
			FileName: fmt.Sprintf("%s_%d_%d_%dx%d.png", cleanTitle(b.deckType), index+1, len(cards),
				columns, rows),
			BackSide: &b.backSideURL,
		})
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
func (b *DeckBuilder) GetResultImages() []string {
	var images []string
	for _, deck := range b.GetDecks() {
		images = append(images, deck.FileName)
	}
	images = append(images, b.backSideName)
	return images
}
func (b *DeckBuilder) GenerateTTSDeck(replaces map[string]string) []TTSDeckObject {
	var res []TTSDeckObject

	var obj TTSDeckObject

	var lastCollection string
	var lastDeck int

	for i, deck := range b.GetDecks() {
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

type WholeCollection []*DeckBuilder

func (col WholeCollection) GetResultImages() map[string]string {
	res := make(map[string]string)
	for _, dc := range col {
		for _, image := range dc.GetResultImages() {
			res[image] = ""
		}
	}
	return res
}
func (col WholeCollection) GenerateTTSDeck(replaces map[string]string) []byte {
	res := TTSSaveObject{}
	for _, dc := range col {
		res.ObjectStates = append(res.ObjectStates, dc.GenerateTTSDeck(replaces)...)
	}
	data, _ := json.MarshalIndent(res, "", "  ")
	return data
}
