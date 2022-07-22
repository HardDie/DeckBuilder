package utils

import (
	"encoding/json"
	"testing"
)

func TestQuoted(t *testing.T) {
	type Test struct {
		Str QuotedString `json:"str"`
	}

	t.Run("first", func(t *testing.T) {
		str := &Test{
			Str: NewQuotedString("hello world"),
		}

		wait := "hello world"
		if str.Str.String() != wait {
			t.Fatalf("Wait %s [got] %s\n", wait, str.Str.String())
		}

		str.Str.SetQuotedOutput()
		wait = "\"hello world\""
		if str.Str.String() != wait {
			t.Fatalf("Wait %s [got] %s", wait, str.Str)
		}

		str.Str.SetRawOutput()
		wait = "hello world"
		if str.Str.String() != wait {
			t.Fatalf("Wait %s [got] %s", wait, str.Str.String())
		}
	})

	t.Run("second", func(t *testing.T) {
		in := []byte("{\"str\": \"good job\"}")
		str := &Test{}

		err := json.Unmarshal(in, str)
		if err != nil {
			t.Fatal(err)
		}

		wait := "good job"
		if str.Str.String() != wait {
			t.Fatalf("Wait %s [got] %s\n", wait, str.Str.String())
		}

		str.Str.SetQuotedOutput()
		wait = "\"good job\""
		if str.Str.String() != wait {
			t.Fatalf("Wait %s [got] %s", wait, str.Str)
		}

		str.Str.SetRawOutput()
		wait = "good job"
		if str.Str.String() != wait {
			t.Fatalf("Wait %s [got] %s", wait, str.Str.String())
		}
	})

	t.Run("third", func(t *testing.T) {
		in := []byte("{\"str\": \"\\\"have a nice day\\\"\"}")
		str := &Test{}

		err := json.Unmarshal(in, str)
		if err != nil {
			t.Fatal(err)
		}

		wait := "have a nice day"
		if str.Str.String() != wait {
			t.Fatalf("Wait %s [got] %s\n", wait, str.Str.String())
		}

		str.Str.SetQuotedOutput()
		wait = "\"have a nice day\""
		if str.Str.String() != wait {
			t.Fatalf("Wait %s [got] %s", wait, str.Str)
		}

		str.Str.SetRawOutput()
		wait = "have a nice day"
		if str.Str.String() != wait {
			t.Fatalf("Wait %s [got] %s", wait, str.Str.String())
		}
	})
}
