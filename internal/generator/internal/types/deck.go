package types

import (
	"log"

	"tts_deck_build/internal/utils"
)

type Deck struct {
	// Type of deck
	Type string `json:"type"`
	// List of cards
	Cards []*Card `json:"cards"`
	// BackSide image
	BackSide *string `json:"backside"`

	// Version of deck
	Version string `json:"version"`
	// Type of collection (example: Base, DLC)
	Collection string `json:"collection"`

	FileName string `json:"fileName"`
	Columns  int    `json:"columns"`
	Rows     int    `json:"rows"`
}

func (d *Deck) GetBackSideName() string {
	if d.BackSide == nil {
		log.Fatalf("Can't get back side image")
	}
	return utils.GetFilenameFromUrl(*d.BackSide)
}
func (d *Deck) GetBackSideURL() string {
	if d.BackSide == nil {
		log.Fatalf("Can't get back side url")
	}
	return *d.BackSide
}
func (d *Deck) GetCards() []*Card {
	return d.Cards
}
func (d *Deck) GetType() string {
	return d.Type
}
func (d *Deck) GetCollection() string {
	return d.Collection
}
