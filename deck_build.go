package main

import (
	"fmt"
	"image"
)

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
	rgba.SaveImage(GetConfig().ResultDir + deck.FileName)
	return
}

func BuildDeck(deckCol *DeckCollection) {
	var imgSlice []image.Image
	var imgBack image.Image

	imgBack = OpenImage(GetConfig().CachePath + deckCol.BackFileName)

	for _, deck := range deckCol.Decks {
		imgSlice = nil
		// Load images
		for i, card := range deck.Cards {
			imgSlice = append(imgSlice, OpenImage(card.GetFilePath()))
			fmt.Printf("\r[ LOAD ] %s (%d / %d)", deck.FileName, i+1, len(deck.Cards))
		}
		// Generate image
		GenerateDeck(deck, imgSlice, imgBack)
		fmt.Printf("\r[ DONE ] %s\n", deck.FileName)
	}
}
