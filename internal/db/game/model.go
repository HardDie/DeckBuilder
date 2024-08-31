package game

import "github.com/HardDie/fsentry/pkg/fsentry_types"

type model struct {
	Description fsentry_types.QuotedString `json:"description"`
	Image       fsentry_types.QuotedString `json:"image"`
}
