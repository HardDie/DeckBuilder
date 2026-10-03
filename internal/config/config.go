package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

const (
	MaxFilenameLength = 200

	MinWidth  = 2
	MinHeight = 2
	MaxWidth  = 10
	MaxHeight = 7
	MaxCount  = MaxWidth*MaxHeight - 1

	HTTPHost         = "127.0.0.1"
	HTTPPort         = 5000
	HTTPPortAttempts = 20
)

type Config struct {
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

func Get(version string) *Config {
	data := "DeckBuilderData"
	if runtime.GOOS == "darwin" {
		// We cannot create a data folder next to an executable file on the macOS system.
		// So create a data folder in home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Error("no home folder for the data folder", "err", err)
			os.Exit(1)
		}
		data = filepath.Join(home, data)
	}

	return &Config{
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

// Logs is the folder of the app log (see ADR 026).
func (c *Config) Logs() string {
	return filepath.Join(c.Data, "logs")
}

// SetDataPath For tests only!!!
func (c *Config) SetDataPath(dataPath string) {
	c.Data = dataPath
}
