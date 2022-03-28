package main

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/exp/constraints"
)

func BestSize(count int) (cols, rows int) {
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
				Type:       d.Type,
				Collection: d.Collection,
				Cards:      d.Cards,
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
			Type:       deck.Type,
			Collection: deck.Collection,
			Cards:      deck.Cards[i:max],
		})
	}
	deck.Cards = deck.Cards[0:MaxCardsOnPage]
}

func NewDeckCollection() *DeckCollection {
	return &DeckCollection{}
}

type WholeCollection []*DeckCollection

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

func (dc *DeckCollection) GetResultImages() []string {
	var images []string
	for _, deck := range dc.Decks {
		images = append(images, deck.FileName)
	}
	images = append(images, dc.BackFileName)
	return images
}

func (dc *DeckCollection) GenerateTTSDeck(replaces map[string]string) []TTSDeckObject {
	var res []TTSDeckObject

	var obj TTSDeckObject

	var lastCollection string
	var lastDeck int

	for i, deck := range dc.Decks {
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
				back, ok := replaces[dc.BackFileName]
				if !ok {
					log.Fatalf("Can't find URL for image: %s", dc.BackFileName)
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
				back, ok := replaces[dc.BackFileName]
				if !ok {
					log.Fatalf("Can't find URL for image: %s", dc.BackFileName)
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

func (dc *DeckCollection) FillAttributes() {
	for index, deck := range dc.Decks {
		deck.Columns, deck.Rows = BestSize(len(deck.Cards) + 1)
		deck.FileName = fmt.Sprintf("%s_%d_%d_%dx%d.png", cleanTitle(deck.Type), index+1, len(deck.Cards),
			deck.Columns, deck.Rows)
	}
}
