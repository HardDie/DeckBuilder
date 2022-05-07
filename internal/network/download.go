package network

import (
	"fmt"
	"image"
	"io/ioutil"
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
		e = errors.NetworkBadURL.AddMessage(err.Error())
		return
	}

	// GET requst for image
	resp, err := http.Get(imageUrl.String())
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.NetworkBadRequest.AddMessage(err.Error())
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

func DownloadBytes(source string) ([]byte, error) {
	// Parse URL
	imageUrl, err := (&url.URL{}).Parse(source)
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.NetworkBadURL.AddMessage(err.Error())
	}

	// GET request for image
	resp, err := http.Get(imageUrl.String())
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.NetworkBadRequest.AddMessage(err.Error())
	}
	defer func() { errors.IfErrorLog(resp.Body.Close()) }()

	// Bad response
	if resp.StatusCode != http.StatusOK {
		return nil, errors.NetworkBadResponse.AddMessage(fmt.Sprintf("code: %d", resp.StatusCode))
	}

	// Read response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.InternalError.AddMessage(err.Error())
	}
	return data, nil
}
