package game

import "time"

type Game struct {
	ID          string
	Name        string
	Description string
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// ImageError is set by create / update when the new image was not applied.
	// It is never stored.
	ImageError error
}

func (e Game) GetName() string {
	return e.Name
}
func (e Game) GetCreatedAt() time.Time {
	return e.CreatedAt
}
