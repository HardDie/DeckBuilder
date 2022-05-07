package types

import (
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/generator/internal/utils"
)

type Card struct {
	// Link for front image download
	Link *string `json:"link"`
	// Link for background image, if required unique
	Background *string `json:"background,omitempty"`
	// Card title
	Title *string `json:"title"`
	// Description of card, if exists
	Description *string `json:"description,omitempty"`
	// Count of cards in result deck
	Count *int `json:"count,omitempty"`
	// Value for scripts
	Scripts map[string]string `json:"scripts"`

	// Cards in same folder exist in same 'collection'
	Collection string `json:"collection"`
	// Full filename with all prefixes: version, collection, deck type, original name(URL path)
	FileName string `json:"fileName"`
	// Same as FileName but for unique back
	BackFileName *string `json:"backFileName"`
}

func (c *Card) GetFrontSideName() string {
	return c.FileName
}
func (c *Card) GetFrontSideURL() string {
	return *c.Link
}
func (c *Card) IsUniqueBack() bool {
	return c.Background != nil
}
func (c *Card) GetUniqueBackSideName() string {
	return *c.BackFileName
}
func (c *Card) GetUniqueBackSineURL() string {
	return *c.Background
}

func (c *Card) FillWithInfo(version, collection, deckType string) {
	c.Collection = collection
	c.FileName = version + "_" + collection + "_" + deckType + "_" + utils.CleanTitle(*c.Title) + "_" + utils.GetFilenameFromUrl(*c.Link)
	if c.Background != nil {
		name := version + "_" + collection + "_" + deckType + "_" + utils.CleanTitle(*c.Title) + "_" + utils.GetFilenameFromUrl(*c.Background)
		c.BackFileName = &name
	}
}

func (c *Card) GetFileName() string {
	return c.FileName
}
func (c *Card) GetFilePath() string {
	return filepath.Join(config.GetConfig().CachePath, c.GetFileName())
}
func (c *Card) GetLua() string {
	var res string
	for key, value := range c.Scripts {
		if len(res) > 0 {
			res += "\n"
		}
		res += key + "=" + value
	}
	return res
}
