package swagger_headers

import "tts_deck_build/internal/errors"

// Default error response
//
// swagger:response ResponseError
type ResponseError struct {
	// In: body
	Body struct {
		// Сообщение ошибки
		// Required: true
		Error errors.Err `json:"error"`
	}
}
