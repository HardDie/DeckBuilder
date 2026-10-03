package dto

type Deck struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	CachedImage string `json:"cachedImage,omitempty"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`

	// Image status for the list; rendering needs every image.
	// HasImage is false when the deck has no back image.
	// CardsMissingImage is true when a card in the deck has no image.
	HasImage          bool `json:"hasImage"`
	CardsMissingImage bool `json:"cardsMissingImage"`
}
