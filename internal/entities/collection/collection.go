package collection

import "time"

type Collection struct {
	ID          string
	Name        string
	Description string
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Dynamic fields

	GameID string

	// ImageError is set by create / update when the new image was not applied.
	// It is never stored.
	ImageError error
}

func (e Collection) GetName() string {
	return e.Name
}
func (e Collection) GetCreatedAt() time.Time {
	return e.CreatedAt
}
