package fsentry_error

import "fmt"

var (
	ErrorBadName  = fmt.Errorf("bad name")
	ErrorBadPath  = fmt.Errorf("bad path")
	ErrorExist    = fmt.Errorf("object exist")
	ErrorNotExist = fmt.Errorf("object not exist")
	ErrorInternal = fmt.Errorf("internal error")
)

func Wrap(err, localErr error) error {
	return fmt.Errorf("%s: %w", err, localErr)
}
