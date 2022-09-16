package dto

type CreateDeckDTO struct {
	Type          string `json:"type"`
	BacksideImage string `json:"backside"`
}

type UpdateDeckDTO struct {
	Type          string `json:"type"`
	BacksideImage string `json:"backside"`
}
