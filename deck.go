package main

type Deck struct {
	// Type of deck
	Type string `json:"type"`
	// List of cards
	Cards []*Card `json:"cards"`
	// Backside image
	Backside *string `json:"backside"`
	// Path prefix
	Prefix string `json:"prefix"`

	FileName string `json:"fileName"`
	Columns  int    `json:"columns"`
	Rows     int    `json:"rows"`
}
