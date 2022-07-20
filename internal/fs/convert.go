package fs

import (
	"fmt"
	"strconv"

	"tts_deck_build/internal/errors"
)

func StringToInt64(in string) (int64, error) {
	val, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		errors.IfErrorLog(err)
		return 0, errors.BadId.AddMessage(err.Error())
	}
	return val, nil
}

func Int64ToString(in int64) string {
	return fmt.Sprintf("%d", in)
}
