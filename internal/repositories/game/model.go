package game

import "github.com/HardDie/fsentry"

type model struct {
	Description fsentry.QuotedString `json:"description"`
	Image       fsentry.QuotedString `json:"image"`
}
