package network

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

const (
	// downloadTimeout covers connect, headers, and body.
	downloadTimeout = 120 * time.Second
	// maxImageBytes caps one downloaded image.
	maxImageBytes = 100 << 20
)

var downloadClient = &http.Client{Timeout: downloadTimeout}

// DownloadBytes fetches source with a timeout and a size cap.
func DownloadBytes(source string) ([]byte, error) {
	return download(downloadClient, source, maxImageBytes)
}

func download(client *http.Client, source string, limit int64) ([]byte, error) {
	// Parse URL
	imageURL, err := (&url.URL{}).Parse(source)
	if err != nil {
		logger.IfError(err)
		return nil, apperr.ErrDownloadBadURL
	}

	// GET request for image
	resp, err := client.Get(imageURL.String())
	if err != nil {
		logger.IfError(err)
		if isTimeout(err) {
			return nil, timeoutError(client)
		}
		return nil, apperr.ErrDownloadFailed
	}
	defer func() { logger.IfError(resp.Body.Close()) }()

	// Bad response
	if resp.StatusCode != http.StatusOK {
		return nil, apperr.Withf(apperr.ErrDownloadFailed, "the image could not be downloaded: the server answered %d", resp.StatusCode)
	}

	tooLarge := apperr.Withf(apperr.ErrImageTooLarge, "image is larger than %d MB", limit>>20)
	if resp.ContentLength > limit {
		return nil, tooLarge
	}

	beginDownload(resp.ContentLength)
	defer endDownload()

	// Read one byte past the limit to catch bodies without a length.
	data, err := io.ReadAll(io.LimitReader(progressReader{r: resp.Body}, limit+1))
	if err != nil {
		if isTimeout(err) {
			return nil, timeoutError(client)
		}
		logger.IfError(err)
		return nil, apperr.ErrDownloadFailed
	}
	if int64(len(data)) > limit {
		return nil, tooLarge
	}
	return data, nil
}

func isTimeout(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func timeoutError(client *http.Client) error {
	return apperr.Withf(apperr.ErrDownloadTimeout, "the image download timed out after %d s", int(client.Timeout.Seconds()))
}
