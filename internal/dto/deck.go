package dto

type CreateDeckDTO struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type UpdateDeckDTO struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
