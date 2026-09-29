package config

import "runtime/debug"

// ResolveVersion returns the version of this repository.
// A stamped ldflag wins. Otherwise the module version, then the VCS revision.
func ResolveVersion(stamped string) string {
	if stamped != "" {
		return stamped
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	return versionFromBuildInfo(info)
}

func versionFromBuildInfo(info *debug.BuildInfo) string {
	switch info.Main.Version {
	case "", "(devel)":
	default:
		return info.Main.Version
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" && setting.Value != "" {
			rev := setting.Value
			if len(rev) > 12 {
				rev = rev[:12]
			}
			return rev
		}
	}
	return "unknown"
}
