package main

const (
	MaxCardsOnPage = 69
)

type Config struct {
	SourceDir string `json:"sourceDir"`
	ResultDir string `json:"resultDir"`
	CachePath string `json:"cachePath"`
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			SourceDir: "desc/eng_v1",
			ResultDir: "result_png/",
			CachePath: ".cache/",
		}
	}
	return config
}
