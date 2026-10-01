package generator

type Generator interface {
	GenerateGame(gameID string, req GenerateGameRequest) error
}

type GenerateGameRequest struct {
	SortOrder string
	Scale     int
}
