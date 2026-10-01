package generator

import (
	"github.com/HardDie/DeckBuilder/internal/render/compose"
)

type Generator struct {
	svc compose.Generator
}

func New(svc compose.Generator) *Generator {
	return &Generator{svc: svc}
}

// Game starts generate and returns. Drawing continues in a goroutine.
// Scale below 1 becomes 1, matching the old HTTP handler.
func (g *Generator) Game(gameID, sortOrder string, scale int) error {
	if scale < 1 {
		scale = 1
	}
	return g.svc.GenerateGame(gameID, compose.GenerateGameRequest{
		SortOrder: sortOrder,
		Scale:     scale,
	})
}
