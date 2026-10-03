package fs

import (
	"strconv"

	"github.com/HardDie/DeckBuilder/internal/apperr"
)

func StringToInt64(in string) (int64, error) {
	val, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return 0, apperr.ErrBadCardID
	}
	return val, nil
}
