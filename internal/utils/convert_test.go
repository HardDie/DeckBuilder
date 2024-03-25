package utils

import (
	"reflect"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/tts_entity"
)

func TestObjectJSONObject(t *testing.T) {
	var in any = tts_entity.Bag{
		Name:             "bag",
		Transform:        tts_entity.Transform{},
		Nickname:         "My bag",
		Description:      "Really my bag",
		ContainedObjects: nil,
	}

	var res tts_entity.Bag
	err := ObjectJSONObject(in, &res)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(in, res) {
		t.Fatal("Objects must be equal")
	}
}
