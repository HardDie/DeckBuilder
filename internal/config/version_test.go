package config

import (
	"runtime/debug"
	"testing"
)

func TestResolveVersionStamped(t *testing.T) {
	if got := ResolveVersion("v9.9.9"); got != "v9.9.9" {
		t.Fatalf("got %q", got)
	}
}

func TestVersionFromBuildInfo(t *testing.T) {
	tests := []struct {
		name string
		info *debug.BuildInfo
		want string
	}{
		{
			name: "module version",
			info: &debug.BuildInfo{Main: debug.Module{Version: "v1.2.3"}},
			want: "v1.2.3",
		},
		{
			name: "vcs revision",
			info: &debug.BuildInfo{
				Main: debug.Module{Version: "(devel)"},
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "0123456789abcdef"},
				},
			},
			want: "0123456789ab",
		},
		{
			name: "missing",
			info: &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}},
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := versionFromBuildInfo(tt.info); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}
