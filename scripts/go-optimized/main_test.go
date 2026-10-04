package main

import (
	"slices"
	"testing"
)

func TestWithoutDebugFlags(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"wails dev build",
			[]string{"build", "-buildvcs=false", "-gcflags", "all=-N -l", "-tags", "dev", "-o", "app"},
			[]string{"build", "-buildvcs=false", "-tags", "dev", "-o", "app"}},
		{"equals form", []string{"build", "-gcflags=all=-N -l", "."}, []string{"build", "."}},
		{"other gcflags kept", []string{"build", "-gcflags", "all=-m", "."}, []string{"build", "-gcflags", "all=-m", "."}},
		{"trailing -gcflags kept", []string{"build", "-gcflags"}, []string{"build", "-gcflags"}},
		{"mod tidy unchanged", []string{"mod", "tidy"}, []string{"mod", "tidy"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := withoutDebugFlags(tt.in); !slices.Equal(got, tt.want) {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
