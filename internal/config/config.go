package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/HardDie/DeckBuilder/internal/logger"
)

const (
	MaxFilenameLength = 200

	MinWidth  = 2
	MinHeight = 2
	MaxWidth  = 10
	MaxHeight = 7
	MaxCount  = MaxWidth*MaxHeight - 1
)

type Config struct {
	Debug   bool   `json:"debug"`
	Version string `json:"version"`

	Data   string `json:"data"`
	Game   string `json:"game"`
	Cache  string `json:"cache"`
	Result string `json:"result"`

	CardImagePath       string `json:"cardImagePath"`
	DeckImagePath       string `json:"deckImagePath"`
	CollectionImagePath string `json:"collectionImagePath"`
	GameImagePath       string `json:"gameImagePath"`
}

func Get(debugFlag bool, version string) *Config {
	data := "DeckBuilderData"
	if runtime.GOOS == "darwin" {
		// We cannot create a data folder next to an executable file on the macOS system.
		// So create a data folder in home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Error.Fatal("Unable to define a user's home directory")
		}
		data = filepath.Join(home, data)
	}

	return &Config{
		Debug:   debugFlag,
		Version: version,

		Data:   data,
		Game:   "games",
		Cache:  "cache",
		Result: "result",

		CardImagePath:       "/api/games/%s/collections/%s/decks/%s/cards/%d/image",
		DeckImagePath:       "/api/games/%s/collections/%s/decks/%s/image",
		CollectionImagePath: "/api/games/%s/collections/%s/image",
		GameImagePath:       "/api/games/%s/image",
	}
}

func (c *Config) Games() string {
	return filepath.Join(c.Data, c.Game)
}
func (c *Config) Results() string {
	return filepath.Join(c.Data, c.Result)
}

// SetDataPath For tests only!!!
func (c *Config) SetDataPath(dataPath string) {
	c.Data = dataPath
}
