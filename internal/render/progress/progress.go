// Sheet progress for one generation. The window polls this value.
package progress

import "sync"

const (
	Empty      = "empty"
	InProgress = "in_progress"
	Done       = "done"
	Error      = "error"
)

// State is one reading of the current run.
type State struct {
	Done    int
	Total   int
	Percent float32
	Status  string
}

var (
	mu    sync.Mutex
	state = State{Status: Empty}
)

// Reset clears the run.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	state = State{Status: Empty}
}

// Begin marks the run in progress and zeroes the bar.
func Begin() {
	mu.Lock()
	defer mu.Unlock()
	state = State{Status: InProgress}
}

// Sheets records finished sheets and the sheet count.
// The percent is done/total*100. A total of 0 keeps the percent at 0.
func Sheets(done, total int) {
	mu.Lock()
	defer mu.Unlock()
	state.Done = done
	state.Total = total
	state.Percent = 0
	if total > 0 {
		state.Percent = float32(done) / float32(total) * 100
	}
}

// Finish marks the run done. The last percent stays.
func Finish() {
	mu.Lock()
	defer mu.Unlock()
	state.Status = Done
}

// Fail marks the run failed. The last percent stays.
func Fail() {
	mu.Lock()
	defer mu.Unlock()
	state.Status = Error
}

// Get returns the current reading.
func Get() State {
	mu.Lock()
	defer mu.Unlock()
	return state
}
