package envcheck

import (
	"runtime"
	"strings"
)

// parseGoVersion extracts the semver number from the output of a "go version" command.
// This code is extracted from Fyne's cmd/fyne/internal/mobile/env.go
func parseGoVersion(out string) string {
	fields := strings.Split(string(out), " ")
	trimGo := func(in string) string {
		return strings.TrimPrefix(in, "go")
	}
	if len(fields) < 3 {
		return trimGo(runtime.Version()) // failed to parse go version
	}

	goVer := trimGo(fields[2])

	// If a go command is a development version, the version
	// information may only appears in the third elements.
	// For instance:
	// go version devel go1.18-527609d47b Wed Aug 25 17:07:58 2021 +0200 darwin/arm64
	if goVer == "devel" && len(fields) >= 4 {
		prefix := strings.Split(fields[3], "-")
		// a go command may miss version information. If that happens
		// we just use the environment version.
		if len(prefix) > 0 {
			return trimGo(prefix[0])
		} else {
			return trimGo(runtime.Version())
		}
	}

	return goVer
}
