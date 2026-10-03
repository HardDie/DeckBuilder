package card

import "time"

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

	// HasImage tells whether the image file exists. Rendering needs it.
	// It is never stored.
	HasImage bool

	// ImageError is set by create / update when the new image was not applied.
	// It is never stored.
	ImageError error
}

func (e Card) GetName() string {
	return e.Name
}
func (e Card) GetCreatedAt() time.Time {
	return e.CreatedAt
}
