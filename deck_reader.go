package main

func (d *Deck) GetBacksideImagePath() DownloadInfo {
	return DownloadInfo{
		FilePath: GetConfig().CachePath + getFilenameFromUrl(*d.Backside),
		FileName: getFilenameFromUrl(*d.Backside),
		URL:      *d.Backside,
	}
}

func (d *Deck) GetUniqueBackImagePath() []DownloadInfo {
	var pairs []DownloadInfo
	for _, card := range d.Cards {
		// Если задняя карта не указана
		if card.GetBackFileName() == nil {
			continue
		}

		pairs = append(pairs, DownloadInfo{
			FilePath: card.GetBackFilePath(),
			FileName: *card.GetBackFileName(),
			URL:      *card.Background,
		})
	}
	return pairs
}

func (d *Deck) GetImagesPath() []DownloadInfo {
	var pairs []DownloadInfo
	for _, card := range d.Cards {
		pairs = append(pairs, DownloadInfo{
			FilePath: card.GetFilePath(),
			FileName: card.GetFileName(),
			URL:      *card.Link,
		})
	}
	return pairs
}

func (d *Deck) GetDownloadList() []DownloadInfo {
	var pairs []DownloadInfo
	if d.Backside != nil {
		pairs = append(pairs, d.GetBacksideImagePath())
	}
	pairs = append(pairs, d.GetImagesPath()...)
	pairs = append(pairs, d.GetUniqueBackImagePath()...)
	return pairs
}
