package tts_entity

type Transform struct {
	PosX   float64 `json:"posX"`
	PosY   float64 `json:"posY"`
	PosZ   float64 `json:"posZ"`
	ScaleX float64 `json:"scaleX"`
	ScaleY float64 `json:"scaleY"`
	ScaleZ float64 `json:"scaleZ"`
}

type DeckDescription struct {
	FaceURL      string `json:"FaceURL"`
	BackURL      string `json:"BackURL"`
	NumWidth     int    `json:"NumWidth"`
	NumHeight    int    `json:"NumHeight"`
	BackIsHidden bool   `json:"BackIsHidden"`
	UniqueBack   bool   `json:"UniqueBack"`
	Type         int    `json:"Type"`
}
