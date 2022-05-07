package config

import "path/filepath"

const (
	MaxCardsOnPage = 69
)

type Config struct {
	SourceDir string `json:"sourceDir"`
	ResultDir string `json:"resultDir"`
	CachePath string `json:"cachePath"`
	Debug     bool   `json:"debug"`

	Data string `json:"data"`
	Game string `json:"game"`

	InfoFilename  string `json:"infoFilename"`
	ImageFilename string `json:"imageFilename"`
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			SourceDir: "./data/desc/eng_v2",
			ResultDir: "./data/result_png/",
			CachePath: "./data/cache",
			Debug:     false,

			Data: "data",
			Game: "games",

			InfoFilename:  ".info.json",
			ImageFilename: "image.bin",
		}
	}
	return config
}

func (c *Config) Games() string {
	return filepath.Join(c.Data, c.Game)
}
