package dto

import "time"

type Card struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	CachedImage string            `json:"cachedImage,omitempty"`
	Variables   map[string]string `json:"variables"`
	Count       int               `json:"count"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}
