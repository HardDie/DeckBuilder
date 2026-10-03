package deck

import "time"

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

	// HasImage tells whether the image file exists. Rendering needs it.
	// It is never stored.
	HasImage bool

	// CardsMissingImage is true when a card in this deck has no image. It is never stored.
	CardsMissingImage bool

	// ImageError is set by create / update when the new image was not applied.
	// It is never stored.
	ImageError error
}

func (e Deck) GetName() string {
	return e.Name
}
func (e Deck) GetCreatedAt() time.Time {
	return e.CreatedAt
}
