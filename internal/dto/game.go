package dto

type CreateGameDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageFile   []byte `json:"imageFile"`
}

type UpdateGameDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageFile   []byte `json:"imageFile"`
}

type DuplicateGameDTO struct {
	Name string `json:"name"`
}
