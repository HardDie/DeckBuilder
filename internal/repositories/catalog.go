package repositories

import (
	"errors"
	"net/http"

	"github.com/HardDie/fsentry"

	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/network"
)

// FsentrySentinels are the catalog errors for one fsentry call.
// A nil Exist or NotExist skips that case. Message keeps the fsentry text
// on those two, with HTTP 400. ErrBadName is always BadName.
type FsentrySentinels struct {
	Exist    *er.Err
	NotExist *er.Err
	Message  bool
}

// ImageBytes downloads imageURL when it is set, then validates the bytes.
// Otherwise it validates data.
func ImageBytes(imageURL string, data []byte) ([]byte, error) {
	if imageURL != "" {
		downloaded, err := network.DownloadBytes(imageURL)
		if err != nil {
			return nil, err
		}
		data = downloaded
	}
	if _, err := images.ValidateImage(data); err != nil {
		return nil, err
	}
	return data, nil
}

// MapFsentry maps ErrExist, ErrNotExist, and ErrBadName onto the sentinels
// the caller passed. Anything else is an internal error carrying err's text.
func MapFsentry(err error, sentinels FsentrySentinels) error {
	if err == nil {
		return nil
	}
	switch {
	case sentinels.Exist != nil && errors.Is(err, fsentry.ErrExist):
		return sentinel(sentinels.Exist, err, sentinels.Message)
	case sentinels.NotExist != nil && errors.Is(err, fsentry.ErrNotExist):
		return sentinel(sentinels.NotExist, err, sentinels.Message)
	case errors.Is(err, fsentry.ErrBadName):
		return er.BadName
	default:
		return er.InternalError.AddMessage(err.Error())
	}
}

func sentinel(target *er.Err, err error, message bool) error {
	if !message {
		return target
	}
	return target.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
}

// ImageChange is the image part of a create or update.
type ImageChange struct {
	// URL is the image URL to store.
	URL string
	// Data is the new image to write. Nil keeps the current file.
	Data []byte
	// Clear removes the current file.
	Clear bool
}

// ResolveImage decides the image part before anything is written.
// oldURL is the stored URL, "" on create.
// A new file wins over a URL, and its stored URL is "".
// Otherwise a URL different from oldURL is a new image.
// A new image is downloaded and validated here.
// On error the change keeps oldURL and the current file,
// and the error says why the new image was not applied.
func ResolveImage(oldURL, newURL string, newFile []byte) (ImageChange, error) {
	if newFile != nil {
		newURL = ""
	} else if newURL == oldURL {
		return ImageChange{URL: oldURL}, nil
	} else if newURL == "" {
		return ImageChange{Clear: true}, nil
	}
	data, err := ImageBytes(newURL, newFile)
	if err != nil {
		return ImageChange{URL: oldURL}, err
	}
	return ImageChange{URL: newURL, Data: data}, nil
}
