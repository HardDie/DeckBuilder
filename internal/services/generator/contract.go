package generator

import "github.com/HardDie/DeckBuilder/internal/dto"

type Generator interface {
	GenerateGame(gameID string, dtoObject *dto.GenerateGameDTO) error
}
