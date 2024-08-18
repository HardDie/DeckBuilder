package game

import (
	"strings"
	"time"
)

type Game struct {
	ID          string
	Name        string
	Description string
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (e Game) GetName() string {
	return strings.ToLower(e.Name)
}
func (e Game) GetCreatedAt() time.Time {
	return e.CreatedAt
}
