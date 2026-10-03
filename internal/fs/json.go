package fs

import (
	"encoding/json"
	"fmt"
	"io"
)

func JsonToWriter[T any](w io.Writer, data T) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("write json: %w", err)
	}
	return nil
}
