package main

import (
	"fmt"
	"image"
	"log"
	"strconv"
)

const (
	// DEBUG
	BuildImage = !true
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

func GenerateDeck(imgSlice []image.Image, imgBack image.Image, prefix string) (cols, rows int, fileName string) {
	images := len(imgSlice)
	_, cols, rows = BestSize(len(imgSlice) + 1)
	fileName = fmt.Sprintf("%s_%d_%dx%d.png", prefix, images, cols, rows)

	if BuildImage {
		bound := imgSlice[0].Bounds().Max
		rgba := CreateImage(bound.X, bound.Y, cols, rows)
		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				if len(imgSlice) <= (row*cols + col) {
					continue
				}
				img := imgSlice[row*cols+col]
				rgba.Draw(col, row, img)
			}
		}
		rgba.Draw(cols-1, rows-1, imgBack)
		rgba.SaveImage(ResultPath + fileName)
	}
	return
}

func BuildDeck(deckCol *DeckCollection) {
	var imgSlice []image.Image
	var imgBack image.Image

	imgBack = OpenImage(deckCol.BackFilePath)

	for index, deck := range deckCol.Decks {
		imgSlice = nil
		// Load images
		if BuildImage {
			log.Println("Start loading images...")
		}
		for _, card := range deck.Cards {
			if BuildImage {
				imgSlice = append(imgSlice, OpenImage(card.GetFilePath()))
			} else {
				imgSlice = append(imgSlice, nil)
			}
		}
		// Generate image
		if BuildImage {
			log.Println("Start generating image...")
		}
		deck.Columns, deck.Rows, deck.FileName = GenerateDeck(imgSlice, imgBack, cleanTitle(deck.Type)+"_"+strconv.Itoa(index+1))
	}
}
