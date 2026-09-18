package fsentry_types

import (
	"encoding/json"
	"strconv"
)

type QuotedString string

func QS(val string) QuotedString {
	return QuotedString(val)
}

func (s QuotedString) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.Quote(string(s)))
}
func (s *QuotedString) UnmarshalJSON(data []byte) error {
	if s == nil {
		return nil
	}
	var val string
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}
	// Remove quotes
	val, err = strconv.Unquote(val)
	if err != nil {
		return err
	}
	*s = QuotedString(val)
	return nil
}

func (s QuotedString) String() string {
	return string(s)
}
