package deck

import (
	"strings"
	"time"
)

type Deck struct {
	ID          string
	Name        string
	Description string
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Dynamic fields

	GameID       string
	CollectionID string

	// ImageError is set by create / update when the new image was not applied.
	// It is never stored.
	ImageError error
}

func (e Deck) GetName() string {
	return strings.ToLower(e.Name)
}
func (e Deck) GetCreatedAt() time.Time {
	return e.CreatedAt
}
