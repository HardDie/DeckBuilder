package main

type Deck struct {
	// Type of deck
	Type string `json:"type"`
	// List of cards
	Cards []*Card `json:"cards"`
	// Backside image
	Backside *string `json:"backside"`

	// Version of deck
	Version string `json:"version"`
	// Type of collection (example: Base, DLC)
	Collection string `json:"collection"`

	FileName string `json:"fileName"`
	Columns  int    `json:"columns"`
	Rows     int    `json:"rows"`
}
