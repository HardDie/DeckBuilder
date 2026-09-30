package version

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Build is the version stamped at link time.
// Set it with -X github.com/HardDie/DeckBuilder/pkg/version.Build=…
// The value is an exact git tag, or the short commit hash.
// A blank stamp means "dev".
var Build = "dev"

// String returns the stamped version.
// Empty or blank is "dev".
func String() string {
	v := strings.TrimSpace(Build)
	if v == "" {
		return "dev"
	}
	return v
}

// Describe returns the version of the git checkout in dir.
// An exact tag of HEAD wins.
// Otherwise the value is the 12-character hash of the current commit.
// An empty dir uses the process working directory.
// It returns an error when git is missing or dir has no commit.
func Describe(dir string) (string, error) {
	tag, err := gitOutput(dir, "describe", "--tags", "--exact-match")
	if err == nil && tag != "" {
		return tag, nil
	}
	hash, err := gitOutput(dir, "rev-parse", "--short=12", "HEAD")
	if err != nil {
		return "", fmt.Errorf("version: %w", err)
	}
	if hash == "" {
		return "", fmt.Errorf("version: empty commit")
	}
	return hash, nil
}

func gitOutput(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return "", fmt.Errorf("%w: %s", err, msg)
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
