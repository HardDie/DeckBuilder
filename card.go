package main

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

func (c *Card) FillWithInfo(version, collection, deckType string) {
	c.Collection = collection
	c.FileName = version + "_" + collection + "_" + deckType + "_" + cleanTitle(*c.Title) + "_" + getFilenameFromUrl(*c.Link)
	if c.Background != nil {
		name := version + "_" + collection + "_" + deckType + "_" + cleanTitle(*c.Title) + "_" + getFilenameFromUrl(*c.Background)
		c.BackFileName = &name
	}
}

func (c *Card) GetFileName() string {
	return c.FileName
}

func (c *Card) GetBackFileName() *string {
	return c.BackFileName
}

func (c *Card) GetFilePath() string {
	return GetConfig().CachePath + c.GetFileName()
}

func (c *Card) GetBackFilePath() string {
	return GetConfig().CachePath + *c.GetBackFileName()
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
