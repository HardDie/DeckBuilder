package network

import (
	"image"
	"net/http"
	"net/url"

	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
)

func DownloadImage(path string) (img image.Image, e error) {
	// Parse URL
	imageUrl, err := (&url.URL{}).Parse(path)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.BarURL.AddMessage(err.Error())
		return
	}

	// GET requst for image
	resp, err := http.Get(imageUrl.String())
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.BadHTTPRequest.AddMessage(err.Error())
		return
	}
	defer func() { errors.IfErrorLog(resp.Body.Close()) }()

	// Parse and validate image
	img, e = fs.BytesToImage(resp.Body)
	if e != nil {
		return
	}
	return
}
