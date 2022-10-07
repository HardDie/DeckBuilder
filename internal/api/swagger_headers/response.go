package swaggerheaders

import "github.com/HardDie/DeckBuilder/internal/errors"

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
