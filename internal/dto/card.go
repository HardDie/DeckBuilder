package dto

type CreateCardDTO struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Variables   map[string]string `json:"variables"`
	Count       int               `json:"count"`
	ImageFile   []byte            `json:"imageFile"`
}

type UpdateCardDTO struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Variables   map[string]string `json:"variables"`
	Count       int               `json:"count"`
	ImageFile   []byte            `json:"imageFile"`
}
