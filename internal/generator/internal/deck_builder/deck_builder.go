package deckbuilder

import (
	"encoding/json"
	"fmt"
	"sort"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/generator/internal/deck_drawer"
	"tts_deck_build/internal/generator/internal/tts_builder"
	"tts_deck_build/internal/generator/internal/types"
	"tts_deck_build/internal/generator/internal/utils"
)

type deckBuilderDeck struct {
	cards []*types.Card

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
func (b *DeckBuilder) AddCard(deck *types.Deck, card *types.Card) {
	if _, ok := b.decks[deck.GetType()]; !ok {
		b.decks[deck.GetType()] = &deckBuilderDeck{
			deckType:     deck.GetType(),
			backSideName: deck.GetBackSideName(),
			backSideURL:  deck.GetBackSideURL(),
		}
	}
	b.decks[deck.GetType()].cards = append(b.decks[deck.GetType()].cards, card)
}
func (b *DeckBuilder) splitCards(deckType string) (cards [][]*types.Card) {
	for leftBorder := 0; leftBorder < len(b.decks[deckType].cards); leftBorder += config.MaxCardsOnPage {
		// Calculate right border for current deck
		rightBorder := utils.Min[int](len(b.decks[deckType].cards), leftBorder+config.MaxCardsOnPage)
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
func (b *DeckBuilder) GetDecks(deckType string) (decks []*types.Deck) {
	for index, cards := range b.splitCards(deckType) {
		// Calculate optimal count of columns and rows for result image
		columns, rows := b.getImageSize(len(cards) + 1)
		decks = append(decks, &types.Deck{
			Cards:   cards,
			Columns: columns,
			Rows:    rows,
			FileName: fmt.Sprintf("%s_%d_%d_%dx%d.png", utils.CleanTitle(b.decks[deckType].deckType), index+1, len(cards),
				columns, rows),
			BackSide: &b.decks[deckType].backSideURL,
			Type:     deckType,
		})
	}
	return
}
func (b *DeckBuilder) GetTypes() (types []string) {
	for deckType := range b.decks {
		types = append(types, deckType)
	}
	sort.SliceStable(types, func(i, j int) bool {
		return types[i] < types[j]
	})
	return
}

// draw
func (b *DeckBuilder) DrawDecks() map[string]string {
	// List of result files
	res := make(map[string]string)
	for _, deckType := range b.GetTypes() {
		decks := b.GetDecks(deckType)
		for _, deck := range decks {
			deckdrawer.NewDeckDrawer(deck).Draw()
			// Add current deck title
			res[deck.FileName] = ""
		}
		// Add back side image title
		res[decks[0].GetBackSideName()] = ""
	}
	return res
}

// tts
func (b *DeckBuilder) GenerateTTSDeck() []byte {
	res := types.TTSSaveObject{}
	for _, deckType := range b.GetTypes() {
		tts := ttsbuilder.NewTTSBuilder()
		decks := b.GetDecks(deckType)
		for deckID, deck := range decks {
			for j, card := range deck.Cards {
				cardID := (deckID+1)*100 + j
				tts.AddCard(deck, card, deckID+1, cardID)
			}
		}
		res.ObjectStates = append(res.ObjectStates, tts.GetObjects()...)
	}
	data, _ := json.MarshalIndent(res, "", "  ")
	return data
}
