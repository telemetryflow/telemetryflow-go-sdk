// Package version_test provides unit tests for the version package.
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

package version_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/telemetryflow/telemetryflow-go-sdk/internal/version"
)

func TestVersionConstants(t *testing.T) {
	t.Run("should have valid version format", func(t *testing.T) {
		// Version should be semantic versioning format
		assert.Regexp(t, `^\d+\.\d+\.\d+`, version.Version)
	})

	t.Run("should have git commit value", func(t *testing.T) {
		assert.NotEmpty(t, version.GitCommit)
	})

	t.Run("should have git branch value", func(t *testing.T) {
		assert.NotEmpty(t, version.GitBranch)
	})

	t.Run("should have build time value", func(t *testing.T) {
		assert.NotEmpty(t, version.BuildTime)
	})
}

func TestInfo(t *testing.T) {
	t.Run("should return formatted version info", func(t *testing.T) {
		info := version.Info()

		assert.NotEmpty(t, info)
		assert.Contains(t, info, "Version:")
		assert.Contains(t, info, version.Version)
		assert.Contains(t, info, "Git Commit:")
		assert.Contains(t, info, "Git Branch:")
		assert.Contains(t, info, "Build Time:")
		assert.Contains(t, info, "Go Version:")
		assert.Contains(t, info, "Platform:")
	})

	t.Run("should contain current runtime info", func(t *testing.T) {
		info := version.Info()

		assert.Contains(t, info, runtime.Version())
		assert.Contains(t, info, runtime.GOOS)
		assert.Contains(t, info, runtime.GOARCH)
	})
}

func TestShort(t *testing.T) {
	t.Run("should return version string", func(t *testing.T) {
		short := version.Short()

		assert.NotEmpty(t, short)
		assert.Equal(t, version.Version, short)
		// Version should be semantic versioning format
		assert.Regexp(t, `^\d+\.\d+\.\d+`, short)
	})
}

func TestFull(t *testing.T) {
	t.Run("should return version with commit", func(t *testing.T) {
		full := version.Full()

		assert.NotEmpty(t, full)
		assert.Contains(t, full, version.Version)
		assert.Contains(t, full, version.GitCommit)
	})

	t.Run("should have hyphen separator", func(t *testing.T) {
		full := version.Full()

		parts := strings.Split(full, "-")
		assert.GreaterOrEqual(t, len(parts), 2)
		assert.Equal(t, version.Version, parts[0])
	})
}

func TestGoVersion(t *testing.T) {
	t.Run("should return Go runtime version", func(t *testing.T) {
		goVersion := version.GoVersion()

		assert.NotEmpty(t, goVersion)
		assert.Equal(t, runtime.Version(), goVersion)
	})

	t.Run("should start with 'go'", func(t *testing.T) {
		goVersion := version.GoVersion()

		assert.True(t, strings.HasPrefix(goVersion, "go"))
	})
}

func TestPlatform(t *testing.T) {
	t.Run("should return OS/ARCH format", func(t *testing.T) {
		platform := version.Platform()

		assert.NotEmpty(t, platform)
		assert.Contains(t, platform, "/")

		parts := strings.Split(platform, "/")
		assert.Len(t, parts, 2)
	})

	t.Run("should match runtime values", func(t *testing.T) {
		platform := version.Platform()

		assert.Contains(t, platform, runtime.GOOS)
		assert.Contains(t, platform, runtime.GOARCH)
	})

	t.Run("should have valid OS", func(t *testing.T) {
		platform := version.Platform()
		parts := strings.Split(platform, "/")

		validOS := []string{"linux", "darwin", "windows", "freebsd", "openbsd", "netbsd"}
		assert.Contains(t, validOS, parts[0])
	})

	t.Run("should have valid architecture", func(t *testing.T) {
		platform := version.Platform()
		parts := strings.Split(platform, "/")

		validArch := []string{"amd64", "arm64", "386", "arm", "ppc64", "ppc64le", "s390x", "riscv64"}
		assert.Contains(t, validArch, parts[1])
	})
}

func TestBuildVariables(t *testing.T) {
	t.Run("should have default values when not set via ldflags", func(t *testing.T) {
		// When not built with ldflags, these should be "unknown" or default values
		// This tests the fallback behavior
		assert.NotEmpty(t, version.Version)
		assert.NotEmpty(t, version.GitCommit)
		assert.NotEmpty(t, version.GitBranch)
		assert.NotEmpty(t, version.BuildTime)
	})
}
