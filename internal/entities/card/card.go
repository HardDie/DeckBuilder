package card

import (
	"strings"
	"time"
)

type Card struct {
	ID          int64
	Name        string
	Description string
	Image       string
	Variables   map[string]string
	Count       int
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Dynamic fields

	GameID       string
	CollectionID string
	DeckID       string
}

func (e Card) GetName() string {
	return strings.ToLower(e.Name)
}
func (e Card) GetCreatedAt() time.Time {
	return e.CreatedAt
}
