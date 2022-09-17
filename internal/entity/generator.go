package entity

import (
	"tts_deck_build/internal/config"
)

type PageInfo struct {
	Columns, Rows int
	Width, Height int
	Count         int
	Name          string
}

// DeckType Unique description for deck
type DeckType struct {
	DeckID string
	Image  string
}

// CardDescription Full description for single card
type CardDescription struct {
	GameID       string
	CollectionID string
	CardID       int64
}

// CardPage Image page for cards with max size 10x7
type CardPage []CardDescription

// CardArray Array of pages cards in same deck type
type CardArray struct {
	Pages []CardPage
}

func NewCardArray() *CardArray {
	return &CardArray{
		Pages: []CardPage{make(CardPage, 0)},
	}
}
func (cArr *CardArray) AddCard(gameID, collectionID string, cardID int64) {
	// Get the index of the last element in the array
	lastArray := len(cArr.Pages) - 1
	if len(cArr.Pages[lastArray]) == config.MaxCount {
		// If the card array has reached the limit of the number of cards per page,
		// create a new page.
		cArr.Pages = append(cArr.Pages, make(CardPage, 0))
		lastArray++
	}
	// Add card to page
	cArr.Pages[lastArray] = append(cArr.Pages[lastArray], CardDescription{
		GameID:       gameID,
		CollectionID: collectionID,
		CardID:       cardID,
	})
}

// DeckArray Full collection of different deck types split by pages
type DeckArray struct {
	Decks        map[DeckType]*CardArray
	selectedDeck DeckType
}

func NewDeckArray() *DeckArray {
	return &DeckArray{
		Decks: make(map[DeckType]*CardArray),
	}
}

// SelectDeck Allows you to select the type of deck to which cards will be added
// in the following calls to the AddCard() method
func (dArr *DeckArray) SelectDeck(title, image string) {
	dArr.selectedDeck = DeckType{
		DeckID: title,
		Image:  image,
	}
	if _, ok := dArr.Decks[dArr.selectedDeck]; !ok {
		dArr.Decks[dArr.selectedDeck] = NewCardArray()
	}
}

// AddCard Allows you to add a card to the selected deck
func (dArr *DeckArray) AddCard(gameID, collectionID string, cardID int64) {
	dArr.Decks[dArr.selectedDeck].AddCard(gameID, collectionID, cardID)
}
