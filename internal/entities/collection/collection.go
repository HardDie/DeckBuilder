package collection

import (
	"strings"
	"time"
)

type Collection struct {
	ID          string
	Name        string
	Description string
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (e Collection) GetName() string {
	return strings.ToLower(e.Name)
}
func (e Collection) GetCreatedAt() time.Time {
	return e.CreatedAt
}
