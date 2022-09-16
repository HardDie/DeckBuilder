package dto

type CreateCardDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Variables   map[string]string `json:"variables"`
	Count       int               `json:"count"`
}

type UpdateCardDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Variables   map[string]string `json:"variables"`
	Count       int               `json:"count"`
}
