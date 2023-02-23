package entity

type RecursiveCollectionItem struct {
	GameID       string `json:"gameId"`
	CollectionID string `json:"collectionId"`
}
type RecursiveDeckItem struct {
	GameID       string `json:"gameId"`
	CollectionID string `json:"collectionId"`
	DeckID       string `json:"deckId"`
}
type RecursiveCardItem struct {
	GameID       string `json:"gameId"`
	CollectionID string `json:"collectionId"`
	DeckID       string `json:"deckId"`
	CardID       int64  `json:"cardId"`
}
type RecursiveSearchItems struct {
	Games       []string                  `json:"games,omitempty"`
	Collections []RecursiveCollectionItem `json:"collections,omitempty"`
	Decks       []RecursiveDeckItem       `json:"decks,omitempty"`
	Cards       []RecursiveCardItem       `json:"cards,omitempty"`
}
