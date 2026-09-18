package card

import (
	"time"

	"github.com/HardDie/fsentry"
)

type model struct {
	ID          int64                           `json:"id"`
	Name        fsentry.QuotedString            `json:"name"`
	Description fsentry.QuotedString            `json:"description"`
	Image       fsentry.QuotedString            `json:"image"`
	Variables   map[string]fsentry.QuotedString `json:"variables"`
	Count       int                             `json:"count"`
	CreatedAt   time.Time                       `json:"createdAt"`
	UpdatedAt   time.Time                       `json:"updatedAt"`
}
