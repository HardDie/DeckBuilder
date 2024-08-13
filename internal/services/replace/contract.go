package replace

import "github.com/HardDie/DeckBuilder/internal/tts_entity"

type Replace interface {
	Prepare(data []byte) ([]Couple, error)
	Replace(data, mapping []byte) (*tts_entity.RootObjects, error)
}

type Couple struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
