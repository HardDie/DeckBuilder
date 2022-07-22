package utils

import (
	"encoding/json"
	"strconv"
)

type QuotedString struct {
	data        string
	quoteOutput bool
}

func NewQuotedString(str string) QuotedString {
	return QuotedString{
		data: str,
	}
}

func (s QuotedString) String() string {
	if s.quoteOutput {
		return strconv.Quote(s.data)
	}
	return s.data
}

func (s *QuotedString) MarshalJSON() ([]byte, error) {
	if s.quoteOutput {
		return json.Marshal(strconv.Quote(s.data))
	}
	return json.Marshal(s.data)
}

func (s *QuotedString) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.data)
	if err != nil {
		return err
	}
	// If data quoted, remove quotes
	if len(s.data) > 0 && []rune(s.data)[0] == '"' {
		s.data, err = strconv.Unquote(s.data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *QuotedString) SetQuotedOutput() {
	s.quoteOutput = true
}

func (s *QuotedString) SetRawOutput() {
	s.quoteOutput = false
}
