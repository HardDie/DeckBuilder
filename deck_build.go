package main

import (
	"fmt"
	"image"
)

const (
	// DEBUG
	BuildImage = true
)

func BestSize(count int) (images, cols, rows int) {
	cols = 10
	rows = 7
	images = cols * rows
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
	images = count
	return
}

func GenerateDeck(deck *Deck, imgSlice []image.Image, imgBack image.Image) {
	bound := imgSlice[0].Bounds().Max
	rgba := CreateImage(bound.X, bound.Y, deck.Columns, deck.Rows)
	for row := 0; row < deck.Rows; row++ {
		for col := 0; col < deck.Columns; col++ {
			if len(imgSlice) <= (row*deck.Columns + col) {
				continue
			}
			img := imgSlice[row*deck.Columns+col]
			rgba.Draw(col, row, img)
			fmt.Printf("\r[ DRAW ] %s (%d / %d)", deck.FileName, row*deck.Columns+col+1, len(imgSlice))
		}
	}
	rgba.Draw(deck.Columns-1, deck.Rows-1, imgBack)
	fmt.Printf("\r[ SAVE ] %s            ", deck.FileName)
	rgba.SaveImage(ResultPath + deck.FileName)
	return
}

func BuildDeck(deckCol *DeckCollection) {
	var imgSlice []image.Image
	var imgBack image.Image

	imgBack = OpenImage(deckCol.BackFilePath)

	for index, deck := range deckCol.Decks {
		_, deck.Columns, deck.Rows = BestSize(len(deck.Cards) + 1)
		deck.FileName = fmt.Sprintf("%s_%d_%d_%dx%d.png", cleanTitle(deck.Type), index+1, len(deck.Cards)+1,
			deck.Columns, deck.Rows)

		if BuildImage {
			imgSlice = nil
			// Load images
			for i, card := range deck.Cards {
				imgSlice = append(imgSlice, OpenImage(card.GetFilePath()))
				fmt.Printf("\r[ LOAD ] %s (%d / %d)", deck.FileName, i+1, len(imgSlice))
			}
			// Generate image
			GenerateDeck(deck, imgSlice, imgBack)
			fmt.Printf("\r[ DONE ] %s\n", deck.FileName)
		}
	}
}
