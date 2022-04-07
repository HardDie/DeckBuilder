package config

const (
	MaxCardsOnPage = 69
)

type Config struct {
	SourceDir string `json:"sourceDir"`
	ResultDir string `json:"resultDir"`
	CachePath string `json:"cachePath"`
	Debug     bool   `json:"debug"`
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			SourceDir: "desc/eng_v2",
			ResultDir: "result_png/",
			CachePath: ".cache/",
			Debug:     false,
		}
	}
	return config
}
