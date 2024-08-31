package dto

type Status struct {
	Type     string  `json:"type"`
	Message  string  `json:"message"`
	Progress float32 `json:"progress"`
	Status   string  `json:"status"`
}
