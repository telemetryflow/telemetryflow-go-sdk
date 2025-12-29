// Package version provides build version information for TelemetryFlow Go SDK.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
package version

import (
	"fmt"
	"runtime"
)

// Build information - set via ldflags at build time
var (
	Version   = "1.1.1"
	GitCommit = "unknown"
	GitBranch = "unknown"
	BuildTime = "unknown"
)

// Info returns version information as a formatted string
func Info() string {
	return fmt.Sprintf(
		"Version: %s\nGit Commit: %s\nGit Branch: %s\nBuild Time: %s\nGo Version: %s\nPlatform: %s/%s",
		Version,
		GitCommit,
		GitBranch,
		BuildTime,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// Short returns a short version string
func Short() string {
	return Version
}

// Full returns the full version with commit
func Full() string {
	return fmt.Sprintf("%s-%s", Version, GitCommit)
}

// GoVersion returns the Go runtime version
func GoVersion() string {
	return runtime.Version()
}

// Platform returns the current platform
func Platform() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}
