package main

const (
	MaxCardsOnPage = 69
)

type Config struct {
	SourceDir string
	ResultDir string
	CachePath string
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			SourceDir: "./desc/eng_v1",
			ResultDir: "result_png/",
			CachePath: ".cache/",
		}
	}
	return config
}
