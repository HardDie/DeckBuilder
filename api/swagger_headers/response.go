package swagger_headers

import (
	"tts_deck_build/internal/errors"
)

// Default error response
//
// swagger:response ResponseError
type ResponseGame struct {
	// In: body
	Body struct {
		errors.Error
	}
}
