package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/constraints"
)

type DeckCollection struct {
	BackURL      string
	BackFileName string
	BackFilePath string

	// List of decks
	Decks []*Deck
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func (dc *DeckCollection) GetLastDeck() *Deck {
	if len(dc.Decks) == 0 {
		return nil
	}
	return dc.Decks[len(dc.Decks)-1]
}

func (dc *DeckCollection) SplitOnDecks(d *Deck) []*Deck {
	count := len(d.Cards) / MaxCardsOnPage
	if len(d.Cards)%MaxCardsOnPage > 0 {
		count++
	}
	fmt.Println("Count of decks:", count, "cards:", len(d.Cards))
	return nil
}

func (dc *DeckCollection) MergeDeck(d *Deck) {
	// If first call, init collection
	if len(dc.Decks) == 0 {
		bs := d.GetBacksideImagePath()
		dc.BackURL = bs.URL
		dc.BackFileName = bs.FileName
		dc.BackFilePath = bs.FilePath

		dc.Decks = []*Deck{
			{
				Type:  d.Type,
				Cards: d.Cards,
			},
		}
		return
	}

	deck := dc.GetLastDeck()
	deck.Cards = append(deck.Cards, d.Cards...)

	if len(deck.Cards) <= MaxCardsOnPage {
		return
	}

	for i := MaxCardsOnPage; i < len(deck.Cards); i += MaxCardsOnPage {
		max := min(i+MaxCardsOnPage, len(deck.Cards))
		dc.Decks = append(dc.Decks, &Deck{
			Type:  deck.Type,
			Cards: deck.Cards[i:max],
		})
	}
	deck.Cards = deck.Cards[0:MaxCardsOnPage]
}

func NewDeckCollection() *DeckCollection {
	return &DeckCollection{}
}

type WholeCollection []*DeckCollection

func (col WholeCollection) GenerateTTSDeck() []byte {
	res := TTSSaveObject{}
	for _, dc := range col {
		res.ObjectStates = append(res.ObjectStates, dc.GenerateTTSDeck()...)
	}
	data, _ := json.MarshalIndent(res, "", "  ")
	return data
}

// DEBUG
var allReplaces []string

func (dc *DeckCollection) GenerateTTSDeck() []TTSDeckObject {
	var res []TTSDeckObject

	var obj TTSDeckObject

	var lastCollection string

	for i, deck := range dc.Decks {
		// DEBUG
		allReplaces = append(allReplaces, deck.FileName, dc.BackFileName)
		// DEBUG
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
				obj.CustomDeck[i+1] = TTSDeckDescription{
					FaceURL:    deck.FileName,
					BackURL:    dc.BackFileName,
					NumWidth:   deck.Columns,
					NumHeight:  deck.Rows,
					UniqueBack: false,
					Type:       0,
				}
			}

			obj.DeckIDs = append(obj.DeckIDs, (i+1)*100+j)
			obj.ContainedObjects = append(obj.ContainedObjects, TTSCard{
				Name:        "Card",
				Nickname:    card.Title,
				Description: new(string),
				LuaScript:   card.GetLua(),
			})
		}
	}

	if len(obj.ContainedObjects) > 0 {
		res = append(res, obj)
	}
	return res
}
