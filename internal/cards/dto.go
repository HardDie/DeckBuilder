package cards

type CreateCardDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Variables   map[string]string `json:"variables"`
}

type UpdateCardDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Variables   map[string]string `json:"variables"`
}
