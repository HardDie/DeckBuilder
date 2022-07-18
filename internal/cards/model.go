package cards

import (
	"strconv"
	"time"

	"tts_deck_build/internal/utils"
)

type CardInfo struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Variables   map[string]string `json:"variables"`
	CreatedAt   *time.Time        `json:"createdAt"`
	UpdatedAt   *time.Time        `json:"updatedAt"`
}

func NewCardInfo(title, desc, image string) *CardInfo {
	return &CardInfo{
		ID:          utils.NameToID(title),
		Title:       strconv.Quote(title),
		Description: strconv.Quote(desc),
		Image:       image,
		Scripts:     make(map[string]string),
		CreatedAt:   utils.Allocate(time.Now()),
	}
}
