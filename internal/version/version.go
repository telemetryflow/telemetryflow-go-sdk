// Package version provides build version information for TelemetryFlow Go SDK.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform
// Copyright (c) 2024-2026 Telemetri Data Indonesia. All rights reserved.
// Open Source Software built by Telemetri Data Indonesia.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version

import (
	"fmt"
	"runtime"
)

// Build information - set via ldflags at build time
var (
	Version   = "1.2.0"
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
