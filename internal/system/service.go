package system

import (
	"os"
)

type SystemService struct {
}

func NewService() *SystemService {
	return &SystemService{}
}

func (s *SystemService) Quit() {
	os.Exit(0)
}
