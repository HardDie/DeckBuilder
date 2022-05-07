package images

import (
	"bytes"
	"image"

	"tts_deck_build/internal/errors"
)

func ValidateImage(input []byte) (string, error) {
	_, imgType, err := image.Decode(bytes.NewBuffer(input))
	if err != nil {
		return "", errors.UnknownImageType.AddMessage(err.Error())
	}
	return imgType, nil
}
