// Package version is a binary version any program can stamp.
//
// [Describe] reads a git checkout.
// An exact tag of HEAD wins.
// Otherwise the value is the 12-character commit hash.
//
// Stamp that value at link time:
//
//	-X github.com/HardDie/DeckBuilder/pkg/version.Build=<value>
//
// [String] returns the stamp. A blank stamp is "dev".
package version
