package config

import "path/filepath"

const (
	MaxCardsOnPage    = 69
	MaxFilenameLength = 200
)

type Config struct {
	Debug bool `json:"debug"`

	Data   string `json:"data"`
	Game   string `json:"game"`
	Cache  string `json:"cache"`
	Result string `json:"result"`

	Setting       string `json:"settings"`
	InfoFilename  string `json:"infoFilename"`
	ImageFilename string `json:"imageFilename"`
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			Debug: false,

			Data:   "data",
			Game:   "games",
			Cache:  "cache",
			Result: "result",

			Setting:       "settings.json",
			InfoFilename:  ".info.json",
			ImageFilename: ".image.bin",
		}
	}
	return config
}

func (c *Config) Settings() string {
	return filepath.Join(c.Data, c.Setting)
}
func (c *Config) Games() string {
	return filepath.Join(c.Data, c.Game)
}
func (c *Config) Caches() string {
	return filepath.Join(c.Data, c.Cache)
}
func (c *Config) Results() string {
	return filepath.Join(c.Data, c.Result)
}

// SetDataPath For tests only!!!
func (c *Config) SetDataPath(dataPath string) {
	c.Data = dataPath
}
