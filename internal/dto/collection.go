package dto

type CreateCollectionDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type UpdateCollectionDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
