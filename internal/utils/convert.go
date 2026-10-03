package utils

import (
	"encoding/json"
	"fmt"
)

func ObjectJSONObject(in any, out any) error {
	data, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("convert object: %w", err)
	}
	err = json.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("convert object: %w", err)
	}
	return nil
}
