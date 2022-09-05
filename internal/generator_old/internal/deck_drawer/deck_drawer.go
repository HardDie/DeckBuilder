package deckdrawer

import (
	"fmt"
	"image"
	"log"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/generator_old/internal/types"
	"tts_deck_build/internal/images"
)

type DeckDrawer struct {
	cards        []*types.Card
	backSideName string
	columns      int
	rows         int
	deckFileName string

	images        []image.Image
	imageBackSide image.Image
}

func NewDeckDrawer(deck *types.Deck) *DeckDrawer {
	return &DeckDrawer{
		cards:        deck.Cards,
		backSideName: deck.GetBackSideName(),
		columns:      deck.Columns,
		rows:         deck.Rows,
		deckFileName: deck.FileName,
	}
}

func (d *DeckDrawer) log(logType string, progress, total int) {
	if total > 0 {
		fmt.Printf("\r[ %s ] %s %d / %d", logType, d.deckFileName, progress, total)
	} else {
		fmt.Printf("\r[ %s ] %s          ", logType, d.deckFileName)
	}
	if logType == "DONE" {
		fmt.Println()
	}
}
func (d *DeckDrawer) loadCards() {
	for i, card := range d.cards {
		img, err := fs.OpenAndProcess(card.GetFilePath(), images.ImageFromReader)
		if err != nil {
			log.Fatalln(err.Error())
		}
		d.images = append(d.images, img)
		d.log("LOAD", i+1, len(d.cards))
	}
	img, err := fs.OpenAndProcess(filepath.Join(config.GetConfig().Caches(), d.backSideName), images.ImageFromReader)
	if err != nil {
		log.Fatalln(err.Error())
	}
	d.imageBackSide = img
}
func (d *DeckDrawer) collectDeckImage() *image.RGBA {
	bound := d.images[0].Bounds().Max
	deckImage := images.CreateImage(bound.X*d.columns, bound.Y*d.rows)
	// Draw front side images
	for row := 0; row < d.rows; row++ {
		for col := 0; col < d.columns; col++ {
			if len(d.images) <= (row*d.columns + col) {
				continue
			}
			curImg := d.images[row*d.columns+col]
			images.Draw(deckImage, col, row, curImg)
			d.log("DRAW", row*d.columns+col+1, len(d.images))
		}
	}
	// On bottom right place draw back side image
	images.Draw(deckImage, d.columns-1, d.rows-1, d.imageBackSide)
	return deckImage
}
func (d *DeckDrawer) Draw() {
	d.loadCards()
	deckImage := d.collectDeckImage()
	d.log("SAVE", 0, 0)
	err := fs.CreateAndProcess[image.Image](filepath.Join(config.GetConfig().Results(), d.deckFileName), deckImage, images.SaveToWriter)
	if err != nil {
		log.Fatal(err.Error())
	}
	d.log("DONE", 0, 0)
}
