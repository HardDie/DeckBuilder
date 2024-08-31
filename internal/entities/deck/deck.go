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
}

func (e Deck) GetName() string {
	return strings.ToLower(e.Name)
}
func (e Deck) GetCreatedAt() time.Time {
	return e.CreatedAt
}
