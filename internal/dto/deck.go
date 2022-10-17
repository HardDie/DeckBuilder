package dto

type CreateDeckDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageFile   []byte `json:"imageFile"`
}

type UpdateDeckDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageFile   []byte `json:"imageFile"`
}
