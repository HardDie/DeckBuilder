package dto

type RecursiveSearchCollection struct {
	GameID       string `json:"gameId"`
	CollectionID string `json:"collectionId"`
}
type RecursiveSearchDeck struct {
	GameID       string `json:"gameId"`
	CollectionID string `json:"collectionId"`
	DeckID       string `json:"deckId"`
}
type RecursiveSearchCard struct {
	GameID       string `json:"gameId"`
	CollectionID string `json:"collectionId"`
	DeckID       string `json:"deckId"`
	CardID       int64  `json:"cardId"`
}
type RecursiveSearch struct {
	Games       []string                    `json:"games,omitempty"`
	Collections []RecursiveSearchCollection `json:"collections,omitempty"`
	Decks       []RecursiveSearchDeck       `json:"decks,omitempty"`
	Cards       []RecursiveSearchCard       `json:"cards,omitempty"`
}
