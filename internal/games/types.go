package games

const (
	GameInfoFilename = "info.json"
)

type GameInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
