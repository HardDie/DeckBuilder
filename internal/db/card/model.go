package card

import (
	"time"

	"github.com/HardDie/fsentry/pkg/fsentry_types"
)

type model struct {
	ID          int64                                 `json:"id"`
	Name        fsentry_types.QuotedString            `json:"name"`
	Description fsentry_types.QuotedString            `json:"description"`
	Image       fsentry_types.QuotedString            `json:"image"`
	Variables   map[string]fsentry_types.QuotedString `json:"variables"`
	Count       int                                   `json:"count"`
	CreatedAt   *time.Time                            `json:"createdAt"`
	UpdatedAt   *time.Time                            `json:"updatedAt"`
}
