package network

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"tts_deck_build/internal/errors"
)

func DownloadBytes(source string) ([]byte, error) {
	// Parse URL
	imageURL, err := (&url.URL{}).Parse(source)
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.NetworkBadURL.AddMessage(err.Error())
	}

	// GET request for image
	resp, err := http.Get(imageURL.String())
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
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.InternalError.AddMessage(err.Error())
	}
	return data, nil
}
