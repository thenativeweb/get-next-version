package version

import (
	"log"
	"runtime/debug"
)

var (
	Version    = "(version unavailable)"
	GitVersion = "(version unavailable)"
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		log.Fatal("failed to read build info")
	}

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			GitVersion = setting.Value
			break
		}
	}
}
