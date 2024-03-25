package utils

import (
	"encoding/json"

	er "github.com/HardDie/DeckBuilder/internal/errors"
)

func ObjectJSONObject(in any, out any) error {
	data, err := json.Marshal(in)
	if err != nil {
		err = er.InternalError.AddMessage(err.Error())
		return err
	}
	err = json.Unmarshal(data, out)
	if err != nil {
		err = er.InternalError.AddMessage(err.Error())
		return err
	}
	return nil
}
