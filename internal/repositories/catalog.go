package repositories

import (
	"errors"
	"fmt"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/network"
)

// FsentrySentinels are the catalog errors for one fsentry call.
// A nil Exist or NotExist skips that case.
type FsentrySentinels struct {
	Exist    *apperr.Error
	NotExist *apperr.Error
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

// MapFsentry turns an fsentry error into one the user can read.
// A missing game, collection, or deck on the path comes first.
// Then ErrExist and ErrNotExist become the caller's sentinels, and ErrBadName becomes ErrBadName.
// Anything else is unexpected and is wrapped for the log.
func MapFsentry(err error, sentinels FsentrySentinels) error {
	if err == nil {
		return nil
	}
	if missing := missingParent(err); missing != nil {
		return missing
	}
	switch {
	case sentinels.Exist != nil && errors.Is(err, fsentry.ErrExist):
		return sentinels.Exist
	case sentinels.NotExist != nil && errors.Is(err, fsentry.ErrNotExist):
		return sentinels.NotExist
	case errors.Is(err, fsentry.ErrBadName):
		return apperr.ErrBadName
	default:
		return fmt.Errorf("catalog store: %w", err)
	}
}

// missingParent maps a missing parent folder from fsentry (*fsentry.BadPathError).
// Its Path runs from "games" through the missing segment:
// two segments is the game, three the collection, four the deck.
// Any other error, including the bare ErrBadPath, gives nil.
func missingParent(err error) error {
	var bad *fsentry.BadPathError
	if !errors.As(err, &bad) || bad == nil || len(bad.Path) < 2 || bad.Path[0] != gamesPath {
		return nil
	}
	switch len(bad.Path) {
	case 2:
		return apperr.ErrGameNotFound
	case 3:
		return apperr.ErrCollectionNotFound
	case 4:
		return apperr.ErrDeckNotFound
	default:
		return nil
	}
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
