// Command version prints the git tag of HEAD, or the short commit hash.
package main

import (
	"fmt"
	"os"

	"github.com/HardDie/DeckBuilder/pkg/version"
)

func main() {
	v, err := version.Describe("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(v)
}
