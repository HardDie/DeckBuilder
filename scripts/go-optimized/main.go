// Command go-optimized runs the go tool without the debug flags that wails dev adds.
//
// wails dev always passes -gcflags "all=-N -l". That turns off optimization and
// inlining in all Go code, so a render is about 3.4× slower than a release build.
// make dev passes this program to wails dev as -compiler; make dev-debug does not.
// REAL_GO names the go tool to run; the default is "go" from PATH.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const debugGCFlags = "all=-N -l"

func main() {
	goTool := os.Getenv("REAL_GO")
	if goTool == "" {
		goTool = "go"
	}
	cmd := exec.Command(goTool, withoutDebugFlags(os.Args[1:])...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		os.Exit(exit.ExitCode())
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "go-optimized:", err)
		os.Exit(1)
	}
}

// withoutDebugFlags drops every "-gcflags all=-N -l" pair and keeps the rest in order.
func withoutDebugFlags(args []string) []string {
	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if args[i] == "-gcflags" && i+1 < len(args) && args[i+1] == debugGCFlags {
			i++
			continue
		}
		if args[i] == "-gcflags="+debugGCFlags {
			continue
		}
		out = append(out, args[i])
	}
	return out
}
