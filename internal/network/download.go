package network

import (
	stderrors "errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/HardDie/DeckBuilder/internal/errors"
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
		errors.IfErrorLog(err)
		return nil, errors.NetworkBadURL.AddMessage(err.Error())
	}

	// GET request for image
	resp, err := client.Get(imageURL.String())
	if err != nil {
		errors.IfErrorLog(err)
		if isTimeout(err) {
			return nil, timeoutError(client)
		}
		return nil, errors.NetworkBadRequest.AddMessage(err.Error())
	}
	defer func() { errors.IfErrorLog(resp.Body.Close()) }()

	// Bad response
	if resp.StatusCode != http.StatusOK {
		return nil, errors.NetworkBadResponse.AddMessage(fmt.Sprintf("server answered %d", resp.StatusCode))
	}

	tooLarge := errors.NetworkBadResponse.AddMessage(fmt.Sprintf("image is larger than %d MB", limit>>20))
	if resp.ContentLength > limit {
		return nil, tooLarge
	}

	// Read one byte past the limit to catch bodies without a length.
	data, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		if isTimeout(err) {
			return nil, timeoutError(client)
		}
		return nil, errors.NetworkBadResponse.AddMessage(err.Error())
	}
	if int64(len(data)) > limit {
		return nil, tooLarge
	}
	return data, nil
}

func isTimeout(err error) bool {
	var netErr net.Error
	return stderrors.As(err, &netErr) && netErr.Timeout()
}

func timeoutError(client *http.Client) error {
	return errors.NetworkTimeout.AddMessage(fmt.Sprintf("download timed out after %d s", int(client.Timeout.Seconds())))
}
